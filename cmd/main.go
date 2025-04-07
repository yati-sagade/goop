package main

import (
	"log"
	"strings"

	"github.com/yati-sagade/goop"
)

func main() {

	prog := `(display "Hello,      \nworld!" foo)`
	p, err := goop.NewProgram(strings.NewReader(prog))
	if err != nil {
		log.Fatalf("Error loading program: %v", err)
	}
	if err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
