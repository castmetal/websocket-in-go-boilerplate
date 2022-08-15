package tests

import (
	"bytes"
	"context"
	"io"
	"testing"

	_use_cases "websocket-in-go-boilerplate/src/application/use-cases"
	_in_memory_tests "websocket-in-go-boilerplate/src/application/use-cases/tests/in-memory-tests"
)

type CreateUserTestStruct struct {
	message []byte
	expect  bool
}

var TestData = []CreateUserTestStruct{
	{[]byte(`{"first_name":"Castmetal","last_name":"Metal","user_name":"castmetal","email":"email@gmail.com","password":"password"}`), true},
	{[]byte(`{"first_name":"Castmetal","last_name":"Metal","user_name":"castmetal","email":"email@gmail.com","password":"password"}`), false},
	{[]byte(`{"first_name":"Castmetal","last_name":"Metal","user_name":"castmetal","email":"wrongemail","password":"password"}`), false},
}

// Testing CreateUser Use Case
func Test(t *testing.T) {
	var message io.Reader
	userRepository := _in_memory_tests.NewUserRepositoryFromConfig()

	var dto = &_in_memory_tests.InMemoryCreateUserDTO{}
	useCase, _ := _use_cases.NewCreateUser(userRepository, dto)

	for _, testItem := range TestData {
		message = bytes.NewReader(testItem.message)
		result, err := useCase.Execute(context.Background(), message)
		if result != testItem.expect {
			t.Error(err.Error())
		}
	}

}
