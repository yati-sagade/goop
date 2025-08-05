package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/yati-sagade/goop"
)

func main() {
	var prog string
	if len(os.Args) > 1 {
		if os.Args[1] == "-" {
			input, err := io.ReadAll(os.Stdin)
			if err != nil {
				log.Fatalf("Error reading stdin: %v", err)
			}
			prog = string(input)
		} else {
			input, err := os.ReadFile(os.Args[1])
			if err != nil {
				log.Fatalf("Error reading file %s: %v", os.Args[1], err)
			}
			prog = string(input)
		}
	} else {
		fmt.Println("Usage: goop [filename|-]")
	}
	p, err := goop.NewProgram(strings.NewReader(prog))
	if err != nil {
		log.Fatalf("Error loading program: %v", err)
	}
	if err := p.Run(goop.RunOptions{}); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
