package goop

import (
	"fmt"
	"strings"
	"unicode"
)

type Tokenizer struct {
	input []rune
	pos   int
	line  int
	col   int
}

func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{
		input: []rune(input),
		pos:   0,
		line:  1,
		col:   1,
	}
}

func (t *Tokenizer) atEnd() bool {
	return t.pos >= len(t.input)
}

func (t *Tokenizer) advance() {
	if !t.atEnd() {
		if t.curr() == '\n' {
			t.line++
			t.col = 1
		} else {
			t.col++
		}
		t.pos++
	}
}

func (t *Tokenizer) curr() rune {
	if t.atEnd() {
		panic("current position is at the end of input")
	}
	return t.input[t.pos]
}

func (t *Tokenizer) consumeWhitespace() {
	for !t.atEnd() && unicode.IsSpace(t.curr()) {
		t.advance()
	}
}

func (t *Tokenizer) matchWord(s string) bool {
	if t.atEnd() {
		return false
	}
	if !t.match(s) {
		return false
	}
	if t.pos+len(s) < len(t.input) {
		return unicode.IsSpace(t.input[t.pos+len(s)])
	}
	return true // Appears at the end of input
}

func (t *Tokenizer) match(s string) bool {
	srunes := []rune(s)
	if t.pos+len(srunes) > len(t.input) {
		return false
	}
	for i, r := range srunes {
		if t.input[t.pos+i] != r {
			return false
		}
	}
	return true
}

func (t *Tokenizer) consume(s string) {
	if !t.match(s) {
		panic(fmt.Sprintf("expected '%s', got '%s'", s, string(t.input[t.pos:])))
	}
	t.pos += len([]rune(s))
}

func (t *Tokenizer) Go() ([]Token, error) {
	var toks []Token
	for {
		t.consumeWhitespace()
		if t.atEnd() {
			break
		}
		switch c := t.curr(); c {
		case '(':
			toks = append(toks, Token{Type: TokenTypeListStart, StringValue: "(", Location: Location{Line: 0, Column: t.pos}})
			t.advance()
		case ')':
			toks = append(toks, Token{Type: TokenTypeListEnd, StringValue: ")", Location: Location{Line: 0, Column: t.pos}})
			t.advance()
		case '"':
			strTok, err := t.parseString()
			if err != nil {
				return nil, err
			}
			toks = append(toks, strTok)
		case '#':
			boolTok, err := t.parseBool()
			if err != nil {
				return nil, err
			}
			toks = append(toks, boolTok)
		default:
			idTok, err := t.parseIdentifier()
			if err != nil {
				return nil, err
			}
			toks = append(toks, idTok)
		}
	}
	return toks, nil
}

func (t *Tokenizer) parseString() (Token, error) {
	t.consume("\"")
	loc := Location{Line: t.line, Column: t.col}
	var sb strings.Builder
	for !t.atEnd() && t.curr() != '"' {
		switch c := t.curr(); c {
		case '\\':
			t.advance()
			if t.atEnd() {
				return Token{}, fmt.Errorf("unexpected end of input while parsing string literal at line %d, column %d", t.line, t.col)
			}
			switch d := t.curr(); d {
			case '"':
				sb.WriteByte('"')
			case '\\':
				sb.WriteByte('\\')
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			default:
				return Token{}, fmt.Errorf("invalid escape sequence: \\%c", d)
			}
		case '"':
			t.advance()
			return Token{Type: TokenTypeStrLiteral, StringValue: sb.String(), Location: Location{Line: 0, Column: t.pos}}, nil
		default:
			sb.WriteRune(c)
		}
		t.advance()
	}
	if t.atEnd() {
		return Token{}, fmt.Errorf("unexpected end of input while parsing string literal")
	}
	t.consume("\"")
	return Token{Type: TokenTypeStrLiteral, StringValue: sb.String(), Location: loc}, nil
}

func (t *Tokenizer) parseBool() (Token, error) {
	loc := Location{Line: t.line, Column: t.col}
	if t.matchWord("#t") {
		t.consume("#t")
		return Token{Type: TokenTypeBoolLiteral, BoolValue: true, Location: loc}, nil
	}
	if t.matchWord("#f") {
		t.consume("#f")
		return Token{Type: TokenTypeBoolLiteral, BoolValue: false, Location: loc}, nil
	}
	return Token{}, fmt.Errorf("expected '#t' or '#f', got '%s'", string(t.input[t.pos:]))
}

func (t *Tokenizer) parseIdentifier() (Token, error) {
	loc := Location{Line: t.line, Column: t.col}
	start := t.pos
	for !t.atEnd() && !unicode.IsSpace(t.curr()) && t.curr() != '(' && t.curr() != ')' {
		t.advance()
	}
	id := string(t.input[start:t.pos])
	if id == "" {
		return Token{}, fmt.Errorf("empty identifier at line %d, column %d", loc.Line, loc.Column)
	}
	return Token{Type: TokenTypeIdent, StringValue: id, Location: loc}, nil
}

type Location struct {
	Line      int
	Column    int
	RawOffset int
}

type Token struct {
	Type     TokenType
	Location Location

	StringValue string
	BoolValue   bool
	NumberValue float64
}

func (token *Token) String() string {
	var val string
	switch token.Type {
	case TokenTypeIdent:
		val = token.StringValue
	case TokenTypeStrLiteral:
		val = token.StringValue
	case TokenTypeBoolLiteral:
		if token.BoolValue {
			val = "#t"
		} else {
			val = "#f"
		}
	default:
		val = token.StringValue
	}
	return fmt.Sprintf("%s(%s)", token.Type, val)
}

type TokenType int

const (
	TokenTypeIdent TokenType = iota
	TokenTypeStrLiteral
	TokenTypeBoolLiteral
	TokenTypeListStart
	TokenTypeListEnd
	TokenTypeIf
	TokenTypeLambda
	TokenTypeDefine
)

func (t TokenType) String() string {
	switch t {
	case TokenTypeIdent:
		return "Identifier"
	case TokenTypeStrLiteral:
		return "StringLiteral"
	case TokenTypeBoolLiteral:
		return "BoolLiteral"
	case TokenTypeListStart:
		return "ListStart"
	case TokenTypeListEnd:
		return "ListEnd"
	case TokenTypeIf:
		return "If"
	case TokenTypeLambda:
		return "Lambda"
	case TokenTypeDefine:
		return "Define"
	default:
		return "Unknown"
	}
}

type Parser struct {
	input []Token
	pos   int
}

func NewParser(input []Token) *Parser {
	return &Parser{
		input: input,
		pos:   0,
	}
}

func (p *Parser) Next() (*Sexpr, error) {
	if p.atEnd() {
		return nil, nil
	}
	if p.atEnd() {
		return nil, nil
	}
	switch c := p.input[p.pos]; c.Type {
	case TokenTypeListStart:
		return p.parseList()
	case TokenTypeListEnd:
		return nil, fmt.Errorf("unexpected ')'")
	case TokenTypeStrLiteral:
		return NewStringAtom(c.StringValue), nil
	case TokenTypeBoolLiteral:
		return NewBoolAtom(c.BoolValue), nil
	case TokenTypeIdent:
		return NewAtom(Identifier, c.StringValue), nil
	}
	return nil, nil
}

func (p *Parser) atEnd() bool {
	return p.pos >= len(p.input)
}

func (p *Parser) consume(tt TokenType) error {
	if p.atEnd() {
		return fmt.Errorf("unexpected end of input")
	}
	if p.input[p.pos].Type != tt {
		return fmt.Errorf("expected %s, got %s", tt, p.input[p.pos].Type)
	}
	p.pos++
	return nil
}

func (p *Parser) advance() {
	if !p.atEnd() {
		p.pos++
	}
}

func (p *Parser) curr() *Token {
	if p.atEnd() {
		panic("current position is at the end of input")
	}
	return &p.input[p.pos]
}

func (p *Parser) parseList() (*Sexpr, error) {
	p.consume(TokenTypeListStart)
	var items []*Sexpr
	for !p.atEnd() && p.curr().Type != TokenTypeListEnd {
		item, err := p.Next()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
		p.advance()
	}
	if p.atEnd() {
		return nil, fmt.Errorf("unexpected end of input while parsing list: missing ')'")
	}
	p.consume(TokenTypeListEnd)
	return NewList(items), nil
}
