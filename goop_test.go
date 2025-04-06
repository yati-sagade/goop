package goop

import (
	"fmt"
	"os"
	"testing"
)

func writeTemp(s string) (string, error) {
	tmpFile, err := os.CreateTemp("", "test_program")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()
	if _, err := tmpFile.WriteString(s); err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func TestDisplay(t *testing.T) {
	// Create a temporary file for the program
	fname, err := writeTemp(`(display "Hello, world!")`)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(fname) // Clean up the temp file after the test
	p, err := LoadProgram(fname)
	if err != nil {
		t.Fatalf("Failed to load program: %v", err)
	}
	fmt.Println("Loaded program:", p)
}
