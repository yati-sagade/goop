package goop

import (
	"fmt"
	"strings"
)

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
		return s.Atom.String()
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

func (a *Atom) String() string {
	switch a.Type {
	case String:
		val := a.Val.(string)
		val = strings.ReplaceAll(val, "\\", "\\\\")
		val = strings.ReplaceAll(val, "\n", "\\n")
		val = strings.ReplaceAll(val, "\t", "\\t")
		val = strings.ReplaceAll(val, "\"", "\\\"")
		return fmt.Sprintf("\"%s\"", val)
	case Identifier:
		return fmt.Sprintf("%s", a.Val)
	case Bool:
		val := a.Val.(bool)
		if val {
			return "#t"
		}
		return "#f"
	case Number:
		val := a.Val.(float64)
		if val == float64(int(val)) {
			return fmt.Sprintf("%d", int(val))
		}
		return fmt.Sprintf("%f", val)
	default:
		panic("unknown atom type: " + fmt.Sprint(a.Type))
	}
}

// NewAtom creates a new Sexpr with an Atom value
func NewAtom(ty AtomType, value string) *Sexpr {
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
