package use_cases

import (
	"bytes"
	"fmt"
	"io"

	_errors "websocket-in-go-boilerplate/src/domains/common"
	_interfaces "websocket-in-go-boilerplate/src/domains/common"
)

type ReceiveUserSocketMessage struct {
	Message io.Reader
	UserId  string
}

func NewReceiveUserSocketMessage(message []byte, userId string) (_interfaces.UseCase, error) {
	if _errors.IsNullOrEmptyByte(message) {
		return nil, _errors.InvalidMessageError(string(message))
	}

	if _errors.IsNullOrEmpty(userId) {
		return nil, _errors.InvalidUserIdError(string(message))
	}

	var uc _interfaces.UseCase = &ReceiveUserSocketMessage{
		Message: bytes.NewReader(message),
		UserId:  userId,
	}

	return uc, nil
}

// Put here your code to read from a queue or other async message like SQS, RabbitMQ, Kafka or connect this in your cluster API
func (msg *ReceiveUserSocketMessage) Execute() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.ReadFrom(msg.Message)

	fmt.Printf("Msg received: %v \n", string(buf.Bytes()))

	// To do your code - example
	return []byte("success"), nil
}
