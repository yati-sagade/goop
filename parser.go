package goop

import (
	"fmt"
	"strings"
)

type Parser struct {
	input string
	pos   int
}

func NewParser(input string) *Parser {
	return &Parser{
		input: input,
		pos:   0,
	}
}

func (p *Parser) Next() (*Sexpr, error) {
	if p.atEnd() {
		return nil, nil
	}
	p.consumeWhitespace()
	if p.atEnd() {
		return nil, nil
	}
	switch c := p.input[p.pos]; c {
	case '"':
		return p.parseStringLiteral()
	case '(':
		return p.parseList()
	default:
		return p.parseIdentifier()
	}
}

func (p *Parser) atEnd() bool {
	return p.pos >= len(p.input)
}

func (p *Parser) consume(prefix string) error {
	if strings.HasPrefix(p.input[p.pos:], prefix) {
		p.pos += len(prefix)
		return nil
	}
	return fmt.Errorf("expected '%s', got '%s'", prefix, p.input[p.pos:])
}

func (p *Parser) parseStringLiteral() (*Sexpr, error) {
	if p.atEnd() {
		return nil, fmt.Errorf("unexpected end of input")
	}
	if err := p.consume("\""); err != nil {
		return nil, err
	}
	if p.atEnd() {
		return nil, fmt.Errorf("unexpected end of input while parsing string literal")
	}
	var sb strings.Builder
	done := false
	for !p.atEnd() && !done {
		switch c := p.input[p.pos]; c {
		case '\\':
			p.pos++
			if p.atEnd() {
				return nil, fmt.Errorf("unexpected end of input while parsing string literal")
			}
			switch d := p.input[p.pos]; d {
			case '"':
				sb.WriteByte('"')
			case '\\':
				sb.WriteByte('\\')
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			default:
				return nil, fmt.Errorf("invalid escape sequence: \\%c", d)
			}
		case '"':
			done = true
		default:
			sb.WriteByte(c)
		}
		p.pos++
	}
	return NewAtom(String, sb.String()), nil
}

func (p *Parser) parseIdentifier() (*Sexpr, error) {
	start := p.pos
	done := false
	for !p.atEnd() && !done {
		c := p.input[p.pos]
		switch c {
		case '(', ')', ' ', '\n', '\t':
			done = true
		default:
			p.pos++
		}
	}
	ident := p.input[start:p.pos]
	if ident == "" {
		return nil, fmt.Errorf("empty identifier")
	}
	return NewAtom(Identifier, ident), nil
}

func (p *Parser) parseList() (*Sexpr, error) {
	p.consumeWhitespace()
	if p.pos >= len(p.input) || p.input[p.pos] != '(' {
		return nil, fmt.Errorf("expected '(', got '%c'", p.input[p.pos])
	}
	p.pos++
	p.consumeWhitespace()
	var items []*Sexpr
	for p.pos < len(p.input) && p.input[p.pos] != ')' {
		item, err := p.Next()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		p.consumeWhitespace()
	}
	if p.pos >= len(p.input) {
		return nil, fmt.Errorf("expected ')', got EOF")
	}
	if p.input[p.pos] == ')' {
		p.pos++
		return NewList(items), nil
	}
	return nil, fmt.Errorf("expected ')', got '%c'", p.input[p.pos])
}

func (p *Parser) consumeWhitespace() {
	for p.pos < len(p.input) && (p.input[p.pos] == ' ' || p.input[p.pos] == '\n' || p.input[p.pos] == '\t') {
		p.pos++
	}
}
