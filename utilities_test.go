package main

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	// Test
	have, _ := Encrypt("r")
	want := "4dc7c9ec434ed06502767136789763ec11d2c4b7"
	if have != want {
		t.Error("TestEncrypt: incorrect encryption")
	}

	// Test
	_, err := Encrypt("")
	if err.Error() != "Empty string" {
		t.Error("TestEncrypt: empty string")
	}
}
