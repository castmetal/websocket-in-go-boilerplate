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
	go func(e chan error) {
		defer cfg.removeConn(conn, userId)

		for {
			_, err := cfg.Epoll.Wait()
			if err != nil {
				continue
			}

			_, _, err = wsutil.ReadClientData(conn)
			if err != nil {
				log.Printf("Error read message %v", err)

				if err := cfg.Epoll.Remove(conn, userId); err != nil {
					log.Printf("Failed to remove %v", err)
				}

				e <- err
				break
			}

		}
	}(err)

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
	go func(e chan error) {
		defer cfg.removeConn(conn, userId)

		for {
			_, err := cfg.Epoll.Wait()
			if err != nil {
				continue
			}

			msg, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				log.Printf("Error read message %v", err)

				if err := cfg.Epoll.Remove(conn, userId); err != nil {
					log.Printf("Failed to remove %v", err)
				}

				e <- err
				break
			}

			var msgBytes io.Reader = bytes.NewReader(msg)
			writeMessage, err := useCase.Execute(msgBytes)
			if err != nil {
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

	}(err)

	select {
	case _err := <-err:
		return false, _err
	default:
		return true, nil
	}

}

func (cfg *WebSocketRequest) WriteToAllClients(ctx context.Context, userId string, conn net.Conn) (bool, error) {
	go func() {
		defer conn.Close()

		for {
			receivedMessage, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				break
			}

			connections := cfg.Epoll.Connections
			cfg.writeMessageToAllConnections(connections, op, receivedMessage, conn, userId)

			break
		}

	}()

	return true, nil
}

func (cfg *WebSocketRequest) WriteToAnUser(ctx context.Context, userId string, conn net.Conn) (bool, error) {
	go func() {
		defer conn.Close()

		for {
			receivedMessage, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				break
			}

			connections := cfg.Epoll.UserConnections.UserConn[userId]
			cfg.writeMessageToAllConnections(connections, op, receivedMessage, conn, userId)

			break
		}

	}()

	return true, nil
}

func (cfg *WebSocketRequest) removeConn(conn net.Conn, userId string) {
	if err := cfg.Epoll.Remove(conn, userId); err != nil {
		log.Printf("Failed to remove %v", err)
	}

	conn.Close()
}

func (cfg *WebSocketRequest) writeMessageToAllConnections(
	connections map[int]net.Conn,
	op ws.OpCode, receivedMessage []byte,
	conn net.Conn,
	userId string,
) {
	for _, epConn := range connections {
		if epConn == nil {
			break
		}

		err := wsutil.WriteServerMessage(epConn, op, receivedMessage)
		if err != nil {
			cfg.removeConn(conn, userId)
			epConn.Close()
			continue
		}

	}
}
