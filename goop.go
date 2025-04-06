package goop

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Val struct {
	Type string
	Val  any
}

type Env struct {
	vars   map[string]*Val
	parent *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{
		vars:   make(map[string]*Val),
		parent: parent,
	}
}

func (e *Env) Set(name string, value *Val) {
	e.vars[name] = value
}

func (e *Env) Get(name string) (*Val, bool) {
	if val, ok := e.vars[name]; ok {
		return val, true
	}
	p := e.parent
	for p != nil {
		if val, ok := p.vars[name]; ok {
			return val, true
		}
		p = p.parent
	}
	return nil, false
}

type Program struct {
	env *Env
}

func NewProgram(io.Reader) (*Program, error) {
	return &Program{env: NewEnv(nil)}, nil
}

func LoadProgram(file string) (*Program, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewProgram(f)
}

func (p *Program) Run() error {
	fmt.Println("Running program...")
	return nil
}

type AtomType int

const (
	Identifier AtomType = iota
	Number
	String
	Bool
)

type Atom struct {
	Type AtomType
	Val  any
}

type Sexpr struct {
	Atom *Atom
	List []*Sexpr
}

func (s *Sexpr) String() string {
	if s.IsAtom() {
		return fmt.Sprintf("%s", s.Atom.Val)
	}
	if s.IsList() {
		var sb strings.Builder
		sb.WriteString("(")
		for i, e := range s.List {
			if i > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(e.String())
		}
		sb.WriteString(")")
		return sb.String()
	}
	return ""
}

// NewAtom creates a new Sexpr with an Atom value
func NewAtom(value string, ty AtomType) *Sexpr {
	return &Sexpr{Atom: &Atom{Type: ty, Val: value}}
}

// NewList creates a new Sexpr with a List value
func NewList(items []*Sexpr) *Sexpr {
	return &Sexpr{List: items}
}

// IsAtom returns true if the Sexpr is an Atom
func (s *Sexpr) IsAtom() bool {
	return s.Atom != nil
}

// IsList returns true if the Sexpr is a List
func (s *Sexpr) IsList() bool {
	return s.List != nil
}

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
		s, err := p.parseStringLiteral()
		if err != nil {
			return nil, err
		}
		return &Sexpr{Atom: s}, nil
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

func (p *Parser) parseStringLiteral() (*Atom, error) {
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
	return &Atom{Type: String, Val: sb.String()}, nil
}

func (p *Parser) parseIdentifier() (*Sexpr, error) {
	start := p.pos
out:
	for p.pos < len(p.input) {
		c := p.input[p.pos]
		switch c {
		case '(', ')', ' ', '\n', '\t':
			break out
		default:
			p.pos++
		}
	}
	atom := p.input[start:p.pos]
	if atom == "" {
		return nil, fmt.Errorf("empty atom")
	}
	return NewAtom(atom, Identifier), nil
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

const (
	ExprDefine = "define"
)

type ExprHandler func(*Env, []*Sexpr) (*Val, error)

var exprHandlers = map[string]ExprHandler{
	ExprDefine: nil,
}
