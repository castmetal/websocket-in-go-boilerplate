package application_services

import (
	"bytes"
	"context"
	"io"
	"log"
	"net"
	_ "net/http/pprof"
	"strconv"

	_core "websocket-in-go-boilerplate/src/core"
	_epoll "websocket-in-go-boilerplate/src/infra/epoll"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type (
	WebSocketService interface {
		_core.IWebsocket
	}
	WebSocketRequest struct {
		WebSocketService
		Epoll *_epoll.Epoll
	}
)

func NewWebSocketService(epoll *_epoll.Epoll) *WebSocketRequest {
	return &WebSocketRequest{Epoll: epoll}
}

func (cfg *WebSocketRequest) EstabilishConnection(ctx context.Context, userId string, conn net.Conn) (bool, error) {
	if err := cfg.Epoll.Add(conn, userId); err != nil {
		log.Printf("Failed to add connection %v", err)
		conn.Close()

		return false, _core.InvalidConnectionError()
	}

	log.Printf("Connection Estabilished")

	err := make(chan error)
	done := make(chan bool)
	go func(e chan error, done chan bool) {
		defer cfg.removeConn(conn, userId, done)

		for {
			_, err := cfg.Epoll.Wait()
			if err != nil {
				continue
			}

			_, _, err = wsutil.ReadClientData(conn)
			if err != nil {
				log.Printf("Closing connection, reason: %v", err)

				if err := cfg.Epoll.Remove(conn, userId); err != nil {
					log.Printf("Failed to remove connection: %v", err)
				}

				e <- err
				break
			}

		}
	}(err, done)

	select {
	case _err := <-err:
		return false, _err
	default:
		return true, nil
	}
}

func (cfg *WebSocketRequest) ExecuteUseCase(ctx context.Context, useCase _core.IUseCase, userId string, conn net.Conn) (bool, error) {
	if err := cfg.Epoll.Add(conn, userId); err != nil {
		log.Printf("Failed to add connection %v", err)
		conn.Close()

		return false, _core.InvalidConnectionError()
	}

	err := make(chan error)
	done := make(chan bool)
	go func(e chan error, done chan bool) {
		defer cfg.removeConn(conn, userId, done)

		for {
			_, err := cfg.Epoll.Wait()
			if err != nil {
				continue
			}

			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				log.Printf("Closing connection, reason: %v", err)

				if err := cfg.Epoll.Remove(conn, userId); err != nil {
					log.Printf("Failed to remove connection: %v", err)
				}

				e <- err
				break
			}

			var msgBytes io.Reader = bytes.NewReader(msg)
			writeMessage, err := useCase.Execute(ctx, msgBytes)
			if err != nil {
				_ = wsutil.WriteServerMessage(conn, op, []byte(strconv.FormatBool(writeMessage)))
				log.Printf("Failed to execute use case: %v", err)
				e <- err
				break
			}

			err = wsutil.WriteServerMessage(conn, op, []byte(strconv.FormatBool(writeMessage)))
			if err != nil {
				e <- err
				break
			}
		}

	}(err, done)

	select {
	case _err := <-err:
		return false, _err
	case <-done:
		return true, nil
	}

}

func (cfg *WebSocketRequest) WriteToAllClients(ctx context.Context, userId string, conn net.Conn) (bool, error) {
	err := make(chan error)
	done := make(chan bool)
	go func(e chan error, done chan bool) {
		defer conn.Close()

		for {
			receivedMessage, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				e <- err
				break
			}

			connections := cfg.Epoll.GetConnections()
			cfg.writeMessageToAllConnections(connections, op, receivedMessage, conn, userId, done)

			break
		}
	}(err, done)

	select {
	case _err := <-err:
		return false, _err
	case <-done:
		return true, nil
	}
}

func (cfg *WebSocketRequest) WriteToAnUser(ctx context.Context, userId string, conn net.Conn) (bool, error) {
	err := make(chan error)
	done := make(chan bool)
	go func(e chan error, done chan bool) {
		defer conn.Close()

		for {
			receivedMessage, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				break
			}

			connections := cfg.Epoll.GetUserConnections(userId)
			cfg.writeMessageToAllConnections(connections, op, receivedMessage, conn, userId, done)

			break
		}

	}(err, done)

	select {
	case _err := <-err:
		return false, _err
	case <-done:
		return true, nil
	}
}

func (cfg *WebSocketRequest) removeConn(conn net.Conn, userId string, done chan bool) {
	if err := cfg.Epoll.Remove(conn, userId); err != nil {
		log.Printf("Failed to remove connection: %v", err)
	}

	conn.Close()

	done <- true
}

func (cfg *WebSocketRequest) writeMessageToAllConnections(
	connections map[int]net.Conn,
	op ws.OpCode, receivedMessage []byte,
	conn net.Conn,
	userId string,
	done chan bool,
) {
	for _, epConn := range connections {
		if epConn == nil {
			break
		}

		err := wsutil.WriteServerMessage(epConn, op, receivedMessage)
		if err != nil {
			if err := cfg.Epoll.Remove(conn, userId); err != nil {
				log.Printf("Failed to remove connection: %v", err)
			}
			epConn.Close()
			continue
		}

	}

	defer conn.Close()
	done <- true
}
