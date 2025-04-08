package goop

import (
	"fmt"
	"io"
	"os"
)

func makeDisplayFunc(stdout io.Writer) GoopFunc {
	if stdout == nil {
		stdout = os.Stdout
	}
	return func(args []*Val) (*Val, error) {
		anys := make([]any, 0)
		for _, arg := range args {
			anys = append(anys, arg)
		}
		fmt.Fprint(stdout, anys...)
		fmt.Fprintln(stdout)
		return nil, nil
	}
}
