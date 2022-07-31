package use_cases

import (
	"bytes"
	"fmt"
	"io"
)

type ReceiveUserSocketMessage struct {
	Message io.Reader
	UserId  string
}

// Put here your code to read from a queue or other async message like SQS, RabbitMQ, Kafka or connect this in your cluster API
func (msg ReceiveUserSocketMessage) Execute() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.ReadFrom(msg.Message)

	fmt.Printf("Msg received: %v \n", string(buf.Bytes()))

	// To do your code - example
	return []byte("success"), nil
}
