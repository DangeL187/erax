package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/DangeL187/erax/erax"
)

func jsonPrint(data map[string]interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes))
}

func log(text string) {
	fmt.Println(text)
}

func CreateUser() erax.Error {
	err := errors.New("email already in use")
	return erax.New(err, "failed to create user").
		WithMeta("code", 503).
		WithMeta("user_error", "An account with this email already exists.")
}

func Register() erax.Error {
	err := CreateUser()
	return erax.New(err, "failed to register\nbecause of ducks!")
}

func main() {
	err := Register()
	if err != nil {
		fmt.Println("Logs:")
		log(erax.Trace(err))

		fmt.Println()

		fmt.Println("Response from server:")
		errUserError, _ := err.Meta("user_error")
		errCode, _ := err.Meta("code")
		jsonPrint(map[string]interface{}{
			"data":  errCode,
			"error": errUserError,
		})
	}
}
