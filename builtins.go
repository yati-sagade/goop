package goop

import (
	"fmt"
)

func displayFunc(args []*Val) (*Val, error) {
	anys := make([]any, 0)
	for _, arg := range args {
		anys = append(anys, arg)
	}
	fmt.Print(anys...)
	fmt.Println()
	return nil, nil
}
