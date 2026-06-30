package main

import (
	"fmt"

	"github.com/DangeL187/erax"
)

func main() {
	// You can fully customize the erax visual theme.

	erax.SetBranchColor("#c6a0f6")
	erax.SetErrorColor("#40a02b")
	erax.SetKeyColor("#209fb5")
	erax.SetValueColor("#dd7878")

	err := erax.New("db timeout")
	err = erax.Wrap(err, "failed to load user")
	err = erax.WithMeta(
		err,
		"service error",
		erax.F("code", "500"),
		erax.F("env", "production"),
	)

	fmt.Println(erax.Format(err))
}
