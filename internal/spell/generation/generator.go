package generation

import (
	"maps"
	"slices"

	"github.com/alan-b-lima/sun/internal/atoms"
	"github.com/alan-b-lima/sun/internal/spell"
	"github.com/alan-b-lima/sun/internal/spell/semantics"
)

type generator struct {
	group_index  index[int]
	cellar_index index[int]
	rune_index   index[spell.Rune]
	states_index index[int]

	groups map[semantics.Symbol]semantics.Group
}

func Generate(outcome semantics.Outcome) spell.Spell {
	gen := generator{
		group_index:  make_index[int](maps.Keys(outcome.Groups), len(outcome.Groups)),
		cellar_index: make_index[int](slices.Values(outcome.Cellars), len(outcome.Cellars)),
		rune_index:   make_index[spell.Rune](slices.Values(outcome.Runes), len(outcome.Runes)),
		states_index: make_index[int](maps.Keys(outcome.States), len(outcome.States)),

		groups: outcome.Groups,
	}

	states := make([]spell.State, len(outcome.States))
	for symbol, state := range outcome.States {
		index := gen.states_index[symbol]
		states[index] = gen.make_state(state)
	}

	return spell.Spell{
		Cellars: len(outcome.Cellars),
		States:  states,
		Initial: gen.states_index[outcome.Initial],
	}
}

var (
	group_none spell.AtomCond
	group_any  spell.AtomCond
)

func init() {
	for i := range atoms.Number {
		group_none[i] = false
		group_any[i] = true
	}
}

func (gen *generator) make_state(state semantics.State) spell.State {
	lines := make([]spell.Line, 0, len(state.Lines))
	for _, line := range state.Lines {
		lines = append(lines, gen.make_line(line))
	}

	return lines
}

func (gen *generator) make_line(line semantics.Line) spell.Line {
	formed := spell.Line{
		AtomCond:    gen.make_atom_cond(line.AtomCond),
		CellarConds: gen.make_cellar_conds(line.CellarConds),
		Behavior:    gen.make_behavior(line.Behavior),
		Move:        spell.Move(line.Move),
		Final:       gen.make_final(line.Final),
	}

	return formed
}

func (gen *generator) make_atom_cond(cond semantics.AtomCond) spell.AtomCond {
	switch {
	case cond.IsAny:
		return group_any

	case cond.IsGroup:
		formed := group_none
		atoms := gen.groups[cond.Group]
		for _, atom := range atoms {
			formed[atom] = true
		}

		return formed

	default:
		formed := group_none
		formed[cond.Atom] = true
		return formed
	}
}

func (gen *generator) make_cellar_conds(conds semantics.CellarConds) []spell.CellarCond {
	formed := make([]spell.CellarCond, 0, len(conds))
	for _, cond := range conds {
		formed = append(formed, gen.make_cellar_cond(cond))
	}

	return formed
}

func (gen *generator) make_cellar_cond(cond semantics.CellarCond) spell.CellarCond {
	switch {
	case cond.IsRune:
		return spell.CellarCond{
			Rune: gen.rune_index[cond.Rune],
		}

	case cond.IsGroup:
		formed := spell.CellarCond{Atom: group_none}

		atoms := gen.groups[cond.Group]
		for _, atom := range atoms {
			formed.Atom[atom] = true
		}

		return formed

	default:
		formed := spell.CellarCond{Atom: group_none}
		formed.Atom[cond.Atom] = true
		return formed
	}
}

func (gen *generator) make_behavior(behavior semantics.Behavior) spell.Behavior {
	switch behavior.Action {
	case semantics.ActionNil:
		return spell.Behavior{Action: spell.BehaviorNil}

	case semantics.ActionAbsorb:
		return spell.Behavior{
			Action: spell.BehaviorAbsorb,
			Cellar: gen.cellar_index[behavior.Cellar],
		}

	case semantics.ActionRelease:
		return spell.Behavior{
			Action: spell.BehaviorRelease,
			Cellar: gen.cellar_index[behavior.Cellar],
		}

	case semantics.ActionWrite:
		return spell.Behavior{
			Action: spell.BehaviorRelease,
			Cellar: gen.cellar_index[behavior.Cellar],
			Rune:   gen.rune_index[behavior.Rune],
		}
	}

	return spell.Behavior{}
}

func (gen *generator) make_final(final semantics.Final) spell.Template {
	formed := spell.Template{
		State:   gen.states_index[final.Symbol],
		Params:  make([]spell.Template, 0, len(final.Params)),
		Ref:     final.IsVar,
		Halting: final.Halting,
	}

	for _, param := range final.Params {
		formed.Params = append(formed.Params, gen.make_final(param))
	}

	return formed
}
