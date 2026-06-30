package main

import (
	"fmt"

	"github.com/DangeL187/erax"
)

func main() {
	// New creates a regular error.
	//
	// It behaves just like errors.New, so use it whenever you need
	// to create a new root error.
	err := erax.New("user not found")

	// Error() returns only the message.
	fmt.Println(err)

	fmt.Println()

	// erax.Format prints the full error trace.
	fmt.Println(erax.Format(err))
}
