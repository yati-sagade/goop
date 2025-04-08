package goop

import (
	"fmt"
	"io"
	"os"
)

type Program struct {
	src io.Reader
	env *Env
}

func NewProgram(src io.Reader) (*Program, error) {
	return &Program{src: src, env: NewEnv(nil)}, nil
}

func LoadProgram(file string) (*Program, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewProgram(f)
}

func (p *Program) eval(s *Sexpr) (*Val, error) {
	if s.IsAtom() {
		switch s.Atom.Type {
		case String:
			return &Val{Type: ValTypeString, Val: s.Atom.Val}, nil
		case Identifier:
			val, ok := p.env.Get(s.Atom.Val.(string))
			if !ok {
				return nil, fmt.Errorf("undefined variable: '%s'", s.Atom.Val)
			}
			return val, nil
		}
	}
	return nil, nil
}

type RunOptions struct {
	Stdout io.Writer
}

func (p *Program) maybeRunSpecialForm(s *Sexpr) (bool, error) {
	if s.IsAtom() {
		return false, nil
	}
	if len(s.List) == 0 {
		return false, nil
	}
	if !s.List[0].IsAtom() {
		return false, nil
	}
	if s.List[0].Atom.Type != Identifier {
		return false, nil
	}
	switch s.List[0].Atom.Val.(string) {
	case "define":
		return true, p.define(s.List[1:])
	}
	return false, nil
}

func (p *Program) Run(opts RunOptions) error {
	// Load builtins
	if err := p.loadBuiltins(opts); err != nil {
		return err
	}
	s, err := io.ReadAll(p.src)
	if err != nil {
		return err
	}
	parser := NewParser(string(s))
	for !parser.atEnd() {
		s, err := parser.Next()
		if err != nil {
			return err
		}
		ran, err := p.maybeRunSpecialForm(s)
		if ran && err != nil {
			return err
		} else if !ran && s.IsList() {
			f, err := p.eval(s.List[0])
			if err != nil {
				return err
			}
			if f.Type != ValTypeFunction {
				return fmt.Errorf("not a function: %s", s.List[0].String())
			}
			args := make([]*Val, 0)
			for _, arg := range s.List[1:] {
				v, err := p.eval(arg)
				if err != nil {
					return err
				}
				args = append(args, v)
			}
			gf := f.Val.(GoopFunc)
			ret, err := gf(args)
			if err != nil {
				return err
			}
			if ret != nil {
				fmt.Println(ret.String())
			}
		}
	}
	return nil
}

func (p *Program) loadBuiltins(opts RunOptions) error {
	p.env.Set("display", NewFuncVal(makeDisplayFunc(opts.Stdout)))
	p.env.Set("foo", NewStringVal("haha"))
	return nil
}

func (p *Program) define(args []*Sexpr) error {
	if len(args) != 2 {
		return fmt.Errorf("define: expected 2 arguments, got %d", len(args))
	}
	if !args[0].IsAtom() || args[0].Atom.Type != Identifier {
		return fmt.Errorf("define: expected identifier, got %s", args[0].String())
	}
	v, err := p.eval(args[1])
	if err != nil {
		return fmt.Errorf("define: %v", err)
	}
	p.env.Set(args[0].Atom.Val.(string), v)
	return nil
}
