package main

import (
	"fmt"

	"github.com/DangeL187/erax"
)

func formatToJSONStringShowcase() {
	// FormatToJSONString serializes an error into a structured JSON string.
	//
	// It preserves:
	// - message
	// - metadata
	// - error tree (cause / errs)

	err := erax.New("db timeout")
	err = erax.Wrap(err, "failed to load user")

	err = erax.WithMeta(
		err,
		"service error",
		erax.F("code", "500"),
	)

	json := erax.FormatToJSONString(err)

	fmt.Println(json)
}

func formatToJSONMapShowcase() {
	err := erax.New("duplicate email")

	err = erax.WithMeta(
		err,
		"validation failed",
		erax.F("code", "409"),
	)

	// Convert error into structured Go map.
	// Useful for logs, tracing systems, or custom serialization.
	m := erax.FormatToJSONMap(err)

	fmt.Println("message:", m["message"])
	fmt.Println("meta:", m["meta"])
}

func fromJSONMapShowcase() {
	err := erax.New("db timeout")
	err = erax.Wrap(err, "failed to load user")

	err = erax.WithMeta(
		err,
		"service error",
		erax.F("code", "500"),
	)

	// Serialize
	m := erax.FormatToJSONMap(err)

	// Deserialize back into erax error
	reconstructed := erax.FromJSONMap(m)

	fmt.Println(erax.Format(reconstructed))
}

func main() {
	fmt.Println()

	formatToJSONStringShowcase()

	fmt.Println()
	fmt.Println("=============================")
	fmt.Println()

	formatToJSONMapShowcase()

	fmt.Println()
	fmt.Println("=============================")
	fmt.Println()

	fromJSONMapShowcase()

	fmt.Println()
}
