package user

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	_use_cases "websocket-in-go-boilerplate/src/domains/user/use-cases"
	_epoll "websocket-in-go-boilerplate/src/epoll"
	_utils "websocket-in-go-boilerplate/src/utils"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type UserSocketController struct {
	Response http.ResponseWriter
	Request  *http.Request
	Epoll    *_epoll.Epoll
}

func (cfg *UserSocketController) removeConn(conn net.Conn, userId string) {
	if err := cfg.Epoll.Remove(conn, userId); err != nil {
		log.Printf("Failed to remove %v", err)
	}

	conn.Close()
}

func (cfg *UserSocketController) SimpleSocket() {
	// Upgrade connection
	userId := cfg.Request.Header.Get(_utils.SystemParams.AUTH_HEADER)

	// TODO - for more security rules, use a middleware before and validates the auth_header with JWT, Oauth or you service provider

	conn, _, _, err := ws.UpgradeHTTP(cfg.Request, cfg.Response)
	if err != nil {
		return
	}

	if err := cfg.Epoll.Add(conn, userId); err != nil {
		log.Printf("Failed to add connection %v", err)
		conn.Close()

		return
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

			useCase, err := _use_cases.NewReceiveUserSocketMessage(
				msg,
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
}

func (cfg *UserSocketController) WriteToAllClients() {
	// Upgrade connection
	userId := cfg.Request.Header.Get(_utils.SystemParams.AUTH_HEADER)

	// TODO - for more security rules, use a middleware before and validates the auth_header with JWT, Oauth or you service provider

	conn, _, _, err := ws.UpgradeHTTP(cfg.Request, cfg.Response)
	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		defer conn.Close()

		for {
			receivedMessage, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				break
			}

			connections := cfg.Epoll.Connections

			for _, epConn := range connections {
				if epConn == nil {
					break
				}

				err = wsutil.WriteServerMessage(epConn, op, receivedMessage)
				if err != nil {
					cfg.removeConn(conn, userId)
					epConn.Close()
					continue
				}

			}

			break
		}

	}()
}

func (cfg *UserSocketController) WriteToAnUser() {
	// Upgrade connection
	userId := cfg.Request.Header.Get(_utils.SystemParams.AUTH_HEADER)

	// TODO - for more security rules, use a middleware before and validates the auth_header with JWT, Oauth or you service provider

	conn, _, _, err := ws.UpgradeHTTP(cfg.Request, cfg.Response)
	if err != nil {
		log.Fatal(err)
		return
	}

	go func() {
		defer conn.Close()

		for {
			receivedMessage, op, err := wsutil.ReadClientData(conn)
			if err != nil {
				break
			}

			connections := cfg.Epoll.UserConnections.UserConn[userId]

			for _, epConn := range connections {
				if epConn == nil {
					break
				}

				err = wsutil.WriteServerMessage(epConn, op, receivedMessage)
				if err != nil {
					cfg.removeConn(conn, userId)
					epConn.Close()
					continue
				}

			}

			break
		}

	}()
}
