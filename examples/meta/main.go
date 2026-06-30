package main

import (
	"fmt"

	"github.com/DangeL187/erax"
)

func withMetaShowcase() {
	// WithMeta attaches structured metadata to an error.
	//
	// Metadata is useful when you want to keep machine-readable context
	// alongside a human-readable message.

	// F is a helper for building metadata fields.
	//
	// It is just a small convenience wrapper around:
	//   erax.MetaField{Key: k, Value: v}

	err := erax.WithMeta(
		erax.New("root error"),
		"meta error message",
		erax.F("field-1", "value-1"),
		erax.F("field-2", "value-2"),
	)

	fmt.Println(erax.Format(err))
}

func addMetaShowcase() {
	err := erax.New("root error")

	// AddMeta appends a single metadata field to an error.
	//
	// ⚠️ Prefer WithMeta when possible.
	// Repeated AddMeta calls may cause extra allocations,
	// since metadata is extended incrementally.
	err = erax.AddMeta(err, "meta error message", "field-1", "value-1")

	// ⚠️ Message is ignored if err is already an erax error. Leave it empty.
	err = erax.AddMeta(err, "", "field-2", "value-2")

	fmt.Println(erax.Format(err))
}

func getMetaShowcase() {
	err := erax.WithMeta(
		erax.New("duplicate email"),
		"failed to create user",
		erax.F("code", "409"),
		erax.F("field", "email"),
	)

	// GetMeta searches through the entire error chain
	// and returns the first matching key.
	code, ok := erax.GetMeta(err, "code")
	if !ok {
		fmt.Println("no code found")
		return
	}

	fmt.Println("code:", code)
}

func main() {
	fmt.Println()

	withMetaShowcase()

	fmt.Println()
	fmt.Println("=============================")
	fmt.Println()

	addMetaShowcase()

	fmt.Println()
	fmt.Println("=============================")
	fmt.Println()

	getMetaShowcase()

	fmt.Println()
}
