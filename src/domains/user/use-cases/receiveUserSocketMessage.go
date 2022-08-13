package use_cases

import (
	"bytes"
	"fmt"
	"io"

	_errors "websocket-in-go-boilerplate/src/domains/common"
	_interfaces "websocket-in-go-boilerplate/src/domains/common"
)

type (
	UserWebSocketService interface {
		Execute()
	}
	receiveUserSocketMessage struct {
		Message io.Reader
		UserId  string
	}
)

type ReceiveUserSocketMessage struct {
	Message io.Reader
	UserId  string
}

func NewReceiveUserSocketMessage(message *io.Reader, userId string) (_interfaces.IUseCase, error) {
	messageBuffer := &bytes.Buffer{}
	messageBuffer.ReadFrom(*message)

	if _errors.IsNullOrEmptyByte(messageBuffer.Bytes()) {
		return nil, _errors.InvalidMessageError(string(messageBuffer.Bytes()))
	}

	if _errors.IsNullOrEmpty(userId) {
		return nil, _errors.InvalidUserIdError(userId)
	}

	// TODO - validate your JSON params here

	var uc _interfaces.IUseCase = &receiveUserSocketMessage{
		Message: bytes.NewReader(messageBuffer.Bytes()),
		UserId:  userId,
	}

	defer func() { message = nil }()

	return uc, nil
}

// Put here your code to read from a queue or other async message like SQS, RabbitMQ, Kafka or connect this in your cluster API
func (msg *receiveUserSocketMessage) Execute() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.ReadFrom(msg.Message)

	fmt.Printf("Msg received: %v \n", string(buf.Bytes()))

	// To do your code - example
	return []byte("success"), nil
}
