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
		if s.IsList() {
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
