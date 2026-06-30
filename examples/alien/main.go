package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/DangeL187/erax"
	pkgerrors "github.com/pkg/errors"
)

func alienShowcase() {
	// erax is designed to work with "alien" errors
	// like github.com/pkg/errors without losing trace visibility.

	err := pkgerrors.New("database connection failed")

	err = pkgerrors.Wrap(err, "failed to load user")
	err = pkgerrors.WithMessage(err, "service layer error")

	err = erax.Wrap(err, "erax error message")

	// Format still produces a readable trace.
	fmt.Println(erax.Format(err))
}

func castShowcase() {
	// Cast explicitly converts standard Go wrapped errors
	// into erax structured error trees.

	// WrapCast and WrapWithErrorsCast do it implicitly so you won't need to.

	err := fmt.Errorf("repository: %w", io.EOF)
	err = fmt.Errorf("service: %w", err)

	err = erax.Cast(err)

	fmt.Println(erax.Format(err))
}

func isAsShowcase() {
	base := erax.New("database error")

	err := erax.Wrap(base, "failed to load user")

	// errors.Is still works
	fmt.Println("Is base:", errors.Is(err, base))

	// errors.As also works for erax errors
	var target error
	if errors.As(err, &target) {
		fmt.Println("As matched:", target)
	}

	// Standard errors also work inside erax
	fmt.Println("Is EOF:", errors.Is(err, io.EOF)) // false
	fmt.Println("Is EOF:", errors.Is(err, base))   // true
}

func main() {
	fmt.Println()

	alienShowcase()

	fmt.Println()
	fmt.Println("=============================")
	fmt.Println()

	castShowcase()

	fmt.Println()
	fmt.Println("=============================")
	fmt.Println()

	isAsShowcase()

	fmt.Println()
}
