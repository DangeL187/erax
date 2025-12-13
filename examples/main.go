package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/DangeL187/erax"
)

func jsonPrint(data map[string]any) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes))
}

func CreateUser() error {
	err := errors.New("email is already in use")

	err = erax.Wrap(err, "failed to create user")
	err = erax.WithMeta(err, "code", "503")
	err = erax.WithMeta(err, "info", "This is a really\nreally long information.")
	err = erax.WithMeta(err, "user_error", "An account with this email already exists.")
	return err
}

func Register() error {
	err := CreateUser()
	return erax.WrapWithError(err, errors.New("random error"), "failed to register\nbecause of ducks!")
}

func main() {
	err := Register()
	if err == nil {
		return
	}

	fmt.Println("Default Logs:")
	fmt.Println(erax.Format(err))
	fmt.Println()

	fmt.Println("JSON Logs:")
	errJSON, _ := erax.FormatToJSONString(err)
	fmt.Println(errJSON)
	fmt.Println()

	fmt.Println("From-JSON Logs:")
	fmt.Println(erax.Format(erax.FromJSONMap(erax.FormatToJSONMap(err))))
	fmt.Println()

	fmt.Println("Response from server:")
	errCode, _ := erax.GetMeta(err, "code")
	errUserError, _ := erax.GetMeta(err, "user_error")
	jsonPrint(map[string]any{
		"code":       errCode,
		"user_error": errUserError,
	})
}
