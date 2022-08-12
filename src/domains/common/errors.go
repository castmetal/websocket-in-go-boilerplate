package common

import (
	"errors"
	"fmt"
)

const PREFIX = "Domain - "
const (
	isNullOrEmptyText  string = PREFIX + "%s is null or empty !"
	invalidUserIdText  string = PREFIX + "The User '%s' is invalid!"
	invalidMessageText string = PREFIX + "Message '%s' is invalid!"
)

func IsNullOrEmptyError(name string) error {
	return errors.New(fmt.Sprintf(isNullOrEmptyText, name))
}

func InvalidUserIdError(name string) error {
	return errors.New(fmt.Sprintf(invalidUserIdText, name))
}

func InvalidMessageError(name string) error {
	return errors.New(fmt.Sprintf(invalidMessageText, name))
}
