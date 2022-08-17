package epoll

import (
	"log"
	"net"
	"reflect"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

type UserConnections struct {
	UserConn map[string]map[int]net.Conn
}

type Epoll struct {
	fd              int
	Connections     map[int]net.Conn
	lock            *sync.RWMutex
	UserConnections UserConnections
}

func MkEpoll() (*Epoll, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	return &Epoll{
		fd:          fd,
		lock:        &sync.RWMutex{},
		Connections: make(map[int]net.Conn),
		UserConnections: UserConnections{
			UserConn: make(map[string]map[int]net.Conn),
		},
	}, nil
}

func (e *Epoll) Add(conn net.Conn, userId string) error {
	// Extract file descriptor associated with the connection
	fd := websocketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	e.Connections[fd] = conn

	if len(e.UserConnections.UserConn[userId]) <= 0 {
		e.UserConnections.UserConn[userId] = make(map[int]net.Conn)
	}

	e.UserConnections.UserConn[userId][fd] = conn

	if len(e.Connections)%100 == 0 {
		log.Printf("Total number of connections: %v", len(e.Connections))
	}

	return nil
}

func (e *Epoll) Remove(conn net.Conn, userId string) error {
	fd := websocketFD(conn)
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}

	e.lock.Lock()
	defer e.lock.Unlock()

	delete(e.Connections, fd)
	delete(e.UserConnections.UserConn[userId], fd)

	if len(e.UserConnections.UserConn[userId]) == 0 {
		delete(e.UserConnections.UserConn, userId)
	}

	if len(e.Connections)%100 == 0 {
		log.Printf("Total number of connections: %v", len(e.Connections))
	}

	return nil
}

func (e *Epoll) Wait() ([]net.Conn, error) {
	events := make([]unix.EpollEvent, 100)
	n, err := unix.EpollWait(e.fd, events, 100)
	if err != nil {
		return nil, err
	}

	e.lock.RLock()
	defer e.lock.RUnlock()

	var connections []net.Conn

	for i := 0; i < n; i++ {
		conn := e.Connections[int(events[i].Fd)]
		connections = append(connections, conn)
	}

	return connections, nil
}

func (e *Epoll) GetConnections() map[int]net.Conn {

	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.Connections
}

func (e *Epoll) GetUserConnections(userId string) map[int]net.Conn {

	e.lock.RLock()
	defer e.lock.RUnlock()

	return e.UserConnections.UserConn[userId]
}

func websocketFD(conn net.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("Sysfd").Int())
}
