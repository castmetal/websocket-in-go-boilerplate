package user

import (
	"bytes"
	"context"
	"io"
	"log"
	"net"
	_ "net/http/pprof"

	_errors "websocket-in-go-boilerplate/src/domains/common"
	_use_cases "websocket-in-go-boilerplate/src/domains/user/use-cases"
	_epoll "websocket-in-go-boilerplate/src/infra/epoll"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type (
	UserWebSocketService interface {
		SimpleSocket(ctx context.Context, userId string, conn net.Conn) (bool, error)
		WriteToAllClients(ctx context.Context, userId string, conn net.Conn) (bool, error)
		WriteToAnUser(ctx context.Context, userId string, conn net.Conn) (bool, error)
	}
	userWebSocketService struct {
		Epoll *_epoll.Epoll
	}
)

func NewUserSocketService(epoll *_epoll.Epoll) UserWebSocketService {
	return &userWebSocketService{Epoll: epoll}
}

func (cfg *userWebSocketService) SimpleSocket(ctx context.Context, userId string, conn net.Conn) (bool, error) {
	if err := cfg.Epoll.Add(conn, userId); err != nil {
		log.Printf("Failed to add connection %v", err)
		conn.Close()

		return false, _errors.InvalidConnectionError()
	}

	go func() {
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

				break
			}

			var msgBytes io.Reader = bytes.NewReader(msg)
			useCase, err := _use_cases.NewReceiveUserSocketMessage(
				&msgBytes,
				userId,
			)
			if err != nil {
				log.Println("Msg with err: ", err)
				break
			}

			writeMessage, err := useCase.Execute()
			if err != nil {
				log.Printf("Failed to execute ReceiveUserSocketMessage: %v", err)
				continue
			}

			err = wsutil.WriteServerMessage(conn, op, writeMessage)
			if err != nil {
				break
			}
		}

	}()

	return true, nil
}

func (cfg *userWebSocketService) WriteToAllClients(ctx context.Context, userId string, conn net.Conn) (bool, error) {
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

func (cfg *userWebSocketService) WriteToAnUser(ctx context.Context, userId string, conn net.Conn) (bool, error) {
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

func (cfg *userWebSocketService) removeConn(conn net.Conn, userId string) {
	if err := cfg.Epoll.Remove(conn, userId); err != nil {
		log.Printf("Failed to remove %v", err)
	}

	conn.Close()
}

func (cfg *userWebSocketService) writeMessageToAllConnections(
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
