package main

import (
	"log"
	"strings"

	"github.com/yati-sagade/goop"
)

func main() {

	prog := `(display "Hello,      \nworld!" foo)`
	/*
		t := goop.NewTokenizer(prog)
		tokens, err := t.Go()
		if err != nil {
			log.Fatalf("Error tokenizing: %v", err)
		}
		for _, token := range tokens {
			log.Printf("Token: %v", token.String())
		}
	*/
	p, err := goop.NewProgram(strings.NewReader(prog))
	if err != nil {
		log.Fatalf("Error loading program: %v", err)
	}
	if err := p.Run(goop.RunOptions{}); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
