package goop

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
