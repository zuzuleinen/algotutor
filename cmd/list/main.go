// Command list prints every known course with its slug and enrollment status.
//
//	go run ./cmd/list
//	make list
package main

import (
	"errors"
	"fmt"
	"os"

	"algotutor/internal/courses"
)

func main() {
	state, err := courses.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, c := range courses.Known {
		marker := "  "
		if state != nil && state.IsEnrolled(c.Slug) {
			marker = "* "
		}
		fmt.Printf("%s%-8s %s\n", marker, c.Slug, c.Name)
	}
}
