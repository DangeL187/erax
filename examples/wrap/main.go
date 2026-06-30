package main

import (
	"fmt"

	"github.com/DangeL187/erax"
)

func wrapShowcase() {
	err := erax.New("0. root error")

	// Wrap adds context to an existing error.
	//
	// This is useful when you want to preserve the original error
	// while explaining where or why it occurred.
	err = erax.Wrap(err, "1. next error")

	// Message can be multiline
	err = erax.Wrap(err, "2. next-next error\nsecond line\nthird line")

	// The original error remains at the bottom of the trace.
	// Each Wrap adds a new layer above it.
	fmt.Println(erax.Format(err))
}

func wrapWithErrorsShowcase() {
	err := erax.New("0. root error")

	// WrapWithErrors combines multiple errors into a single one.
	//
	// Unlike Wrap, it can represent several independent causes.
	err = erax.WrapWithErrors(
		err,                 // Goes to the bottom of the list
		"validation failed", // Use the message to describe what the enclosed errors have in common
		erax.New("3. username is already taken"),
		erax.New("2. password is too short"),
		erax.New("1. email is invalid"),
	)

	fmt.Println(erax.Format(err))
}

func wrapWithErrorsNestedShowcase() {
	// WrapWithErrors can be nested to build arbitrarily complex
	// error trees.

	wrap0 := erax.WrapWithErrors(
		nil,
		"0. security",
		erax.New("password was found in a data breach"),
		erax.New("password is too common"),
	)

	wrap1 := erax.WrapWithErrors(
		nil,
		"1. password",
		erax.New("must contain at least one digit"),
		erax.New("must contain at least one uppercase letter"),
		wrap0,
	)

	wrap2 := erax.WrapWithErrors(
		nil,
		"2. user",
		erax.New("username is already taken"),
		erax.New("email is invalid"),
	)

	err := erax.WrapWithErrors(
		nil,
		"request validation failed",
		wrap2,
		wrap1,
	)

	fmt.Println(erax.Format(err))
}

func main() {
	fmt.Println()

	wrapShowcase()

	fmt.Println()
	fmt.Println("=============================")
	fmt.Println()

	wrapWithErrorsShowcase()

	fmt.Println()
	fmt.Println("=============================")
	fmt.Println()

	wrapWithErrorsNestedShowcase()

	fmt.Println()
}
