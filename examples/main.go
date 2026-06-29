package main

import (
	"encoding/json"
	"fmt"

	"github.com/DangeL187/erax"
	"github.com/TylerBrock/colorjson"
)

func FormatJSON(rawJSON string) string {
	formatter := colorjson.NewFormatter()
	formatter.Indent = 2

	var obj interface{}
	err := json.Unmarshal([]byte(rawJSON), &obj)
	if err != nil {
		return ""
	}

	formattedBody, err := formatter.Marshal(obj)
	if err != nil {
		return ""
	}

	return string(formattedBody)
}

func jsonPrint(data map[string]any) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes))
}

func CreateUser() error {
	err := erax.New("email is already in use")

	err = erax.WithMeta(err, "failed to create user",
		erax.F("code", "503"),
		erax.F("info", "This is a really\nreally long information."),
		erax.F("user_error", "An account with this email already exists."),
	)
	return err
}

func Register() error {
	err := CreateUser()
	return erax.WrapWithErrors(err, "failed to register\nbecause of ducks!",
		erax.New("random error"),
	)
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
	errJSON := erax.FormatToJSONString(err)
	fmt.Println(FormatJSON(errJSON))
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
