package main

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	// <setup code>
	// <test-code>
	t.Run(fmt.Sprintf("Test correct encryption"), func(t *testing.T) {
		have, _ := Encrypt("r")
		want := "4dc7c9ec434ed06502767136789763ec11d2c4b7"
		if have != want {
			t.Error("Failed correct encryption")
		}
	})

	t.Run(fmt.Sprintf("Test encrypting empty string"), func(t *testing.T) {
		_, err := Encrypt("")
		if err.Error() != "Empty string" {
			t.Error("TestEncrypt: empty string")
		}
	})

	// <tear-down code>
}
