package spell

import "github.com/alan-b-lima/sun/internal/atom"

type State struct {
	Transitions []Transition
}

type Transition struct {
	AtomCond    AtomCond
	CellarConds []CellarCond
	Behavior    Behavior
	Move        Move
	Final       Template
}

type AtomCond [atom.Number]bool

func (c *AtomCond) For(atom atom.Atom) bool {
	if int(atom) >= len(c) {
		return false
	}

	return c[atom]
}

type CellarCond struct {
	Atom AtomCond
	Rune Rune
}

func (c *CellarCond) For(cell Cell) bool {
	switch cell.Tag {
	case TagAtom:
		return c.Atom.For(cell.Atom)
	case TagRune:
		return c.Rune == cell.Rune
	}

	return false
}

type state struct {
	State  int
	Params []state
}

type Template struct {
	State  int
	Params []Template
	Ref    bool
}

func (t Template) resolve(config state) (state, bool) {
	if t.Ref {
		if t.State >= len(config.Params) {
			return state{}, false
		}

		return config.Params[t.State], true
	}

	final := state{
		State:  t.State,
		Params: make([]state, 0, len(t.Params)),
	}

	for _, tmpl := range t.Params {
		param, ok := tmpl.resolve(config)
		if !ok {
			return state{}, false
		}

		final.Params = append(final.Params, param)
	}

	return final, true
}
