package common

import (
	"net/http"
)

func OnlySameRequestMethod(r *http.Request, method string) error {
	if r.Method != method {
		return OnlySameMethodError()
	}

	return nil
}

func ValidateSameMethod(w http.ResponseWriter, r *http.Request, method string) error {
	err := OnlySameRequestMethod(r, method)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(err.Error()))

		return OnlySameMethodError()
	}

	return nil
}
