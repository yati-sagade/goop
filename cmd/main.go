package main

import (
	"fmt"

	"github.com/yati-sagade/goop"
)

func main() {
	p := goop.NewParser(`(display "Hello,      \nworld!")`)
	for {
		s, err := p.Next()
		if err != nil {
			fmt.Println("Error parsing:", err)
		}
		if s == nil {
			break
		}
		fmt.Printf("%+v", s)
	}
}
