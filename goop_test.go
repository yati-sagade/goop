package goop

import (
	"fmt"
	"strings"
	"testing"
)

func TestDisplay(t *testing.T) {
	tests := []struct {
		args           []string
		expectedOutput string
	}{
		{args: []string{`"Hello, world!"`}, expectedOutput: "Hello, world!\n"},
		{args: []string{`""`}, expectedOutput: "\n"},
		{args: []string{}, expectedOutput: "\n"},
		{args: []string{`"Hello\n     spaces!"`}, expectedOutput: "Hello\n     spaces!\n"},
	}
	for _, test := range tests {
		progname := fmt.Sprintf("(display %s)", strings.Join(test.args, " "))
		t.Run(progname, func(t *testing.T) {
			prog, err := NewProgram(strings.NewReader(progname))
			if err != nil {
				t.Fatalf("Failed to create program: %v", err)
			}
			output := &strings.Builder{}
			if err := prog.Run(RunOptions{Stdout: output}); err != nil {
				t.Errorf("Failed to run program: %v", err)
			}
			if output.String() != test.expectedOutput {
				t.Errorf("Expected output: %q, got: %q", test.expectedOutput, output.String())
			}
		})
	}
}
