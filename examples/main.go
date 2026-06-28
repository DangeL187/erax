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
	err := erax.New("email is already in use\nemail-LONG")

	err = erax.WithMeta(err, "failed to create user\nfailed-LONG",
		erax.F("code", "503"),
		erax.F("info", "This is a really\nreally long information."),
		erax.F("user_error", "An account with this email already exists."),
	)
	err = erax.WrapWithErrors(err, "UNITY\nUNITY-LONG", err)
	return err
}

func Register() error {
	err := CreateUser()
	return erax.WrapWithErrors(err, "failed to register\nbecause of ducks!",
		erax.New("random error\nrandom-LONG"),
	)
}

func main() {
	err := Register()
	if err == nil {
		return
	}

	err = erax.WrapWithErrors(err, "[---]\n[---]-LONG", err)
	err = erax.WithMeta(err, "AHAH\nAHAH-LONG",
		erax.F("A", "503"),
		erax.F("B", "This is a really\nreally long information."),
		erax.F("C", "An account with this email already exists."),
	)
	err = erax.WithMeta(err, "AHAH\nAHAH-LONG",
		erax.F("A", "503"),
		erax.F("B", "This is a really\nreally long information."),
		erax.F("C", "An account with this email already exists."),
	)

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

/*

 ▼ [ERROR TRACE]
 ╰─ [---] // writeIndents(nil); cause==nil(YES) -> " ╰─ "; append(T); ... cause==nil(YES) -> writeIndents(T); " ├╮";
     ├╮
     │╰─ [0] failed to register // writeIndents(T); cause==nil(YES) && isLast(NO) -> " │" -> this.cause==nil -> "╰─"; append(F);
     │   because of ducks!
     │    ├╮
     │    │╰─ [0] random error   // writeIndents(TF); cause==nil(YES) && isLast(NO) -> "├─";
     │    ╰╮
     │     ╰─ [1] UNITY          // writeIndents(TF); cause==nil(YES) && isLast -> "├─";
    F│FF       ├╮
TTTT │  TTTT   │├─ [0] failed to create user // ? writeIndents(TF); cause==nil && isLast -> "├─";
     │         ││   ├─ code: 503
1234 │         ││   ├─ info:
     │  1234   ││   │   This is a really
    1│34    123││   │   really long information.
     │         ││   ╰─ user_error: An account with this email already exists.
     │         │╰─ email is already in use
     │         ╰╮
     │          ├─ [1] failed to create user
     │          │   ├─ code: 503
     │          │   ├─ info:
     │          │   │   This is a really
     │          │   │   really long information.
     │          │   ╰─ user_error: An account with this email already exists.
     ╰╮         ╰─ email is already in use
      ╰─ [1] failed to register
         because of ducks!
          ├╮
          │╰─ [0] random error   // writeIndents(TF); cause==nil(YES) && isLast(NO) -> "├─";
          ╰╮
           ╰─ [1] UNITY          // writeIndents(TF); cause==nil(YES) && isLast -> "├─";
1234    1234   ├╮
    1234    123│├─ [0] failed to create user // ? writeIndents(TF); cause==nil && isLast -> "├─";
               ││   ├─ code: 503
               ││   ├─ info:
               ││   │   This is a really
               ││   │   really long information.
               ││   ╰─ user_error: An account with this email already exists.
               │╰─ email is already in use
               ╰╮
                ├─ [1] failed to create user
                │   ├─ code: 503
                │   ├─ info:
                │   │   This is a really
                │   │   really long information.
                │   ╰─ user_error: An account with this email already exists.
                ╰─ email is already in use

*/
