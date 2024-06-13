package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }() // os.Args is a "global variable", so keep the state from before the test, and restore it after.

	os.Args = []string{"", "../api.go", "../api_handlers.go"}
	main()
}
