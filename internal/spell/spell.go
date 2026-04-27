package spell

import "github.com/alan-b-lima/sun/internal/atoms"

type CastingSpell struct {
	*Spell

	// runtime info

	world  World
	x, y   int
	facing Facing

	cellars []Cellar

	energy Energy
	halted bool
	state  state
}

type Spell struct {
	Cellars int
	States  []State
	Initial int
}

type World interface {
	Dim() (w, h int)

	At(x, y int) atoms.Atom
	Set(x, y int, atom atoms.Atom)
}

type Energy int64

func Make(cellars int, states []State, initial int) Spell {
	return Spell{
		Cellars: cellars,
		States:  states,
		Initial: initial,
	}
}

func (spell *Spell) Cast(world World, energy Energy, x, y int, facing Facing) CastingSpell {
	s := CastingSpell{
		Spell: spell,
	}

	s.world = world
	s.x = x
	s.y = y
	s.facing = facing

	s.cellars = make([]Cellar, spell.Cellars)
	s.energy = energy - CellarCost*Energy(s.Cellars)
	s.halted = s.energy < 0
	s.state = state{State: spell.Initial}

	return s
}

func (s *CastingSpell) Energy() int  { return int(s.energy) }
func (s *CastingSpell) Halted() bool { return s.halted }

func (s *CastingSpell) State(index int) (State, bool) {
	if index >= len(s.States) {
		return State{}, false
	}
	return s.States[index], true
}

func (s *CastingSpell) Cellar(index int) (Cellar, bool) {
	if index >= len(s.cellars) {
		return Cellar{}, false
	}
	return s.cellars[index], true
}

func (s *CastingSpell) Perform() {
	if s.halted {
		return
	}

	state, ok := s.State(s.state.State)
	if !ok {
		s.halted = true
		return
	}

	transition, ok := s.find(state)
	if !ok {
		s.halted = true
		return
	}

	if !s.exec(transition) {
		s.halted = true
		return
	}

	if !s.next(transition.Final) {
		s.halted = true
	}
}

func (s *CastingSpell) find(state State) (Transition, bool) {
	for _, transition := range state {
		atom := s.world.At(s.x, s.y)
		if !transition.AtomCond.For(atom) {
			continue
		}

		for i, cond := range transition.CellarConds {
			cellar, ok := s.Cellar(i)
			if !ok {
				continue
			}

			top, ok := cellar.Peek()
			if !ok {
				continue
			}

			if cond.For(top) {
				return transition, true
			}
		}
	}

	return Transition{}, false
}

func (s *CastingSpell) exec(transition Transition) bool {
	cost := 1 +
		BehaviorCost[transition.Behavior.Action] +
		MoveCost(transition.Move, s.facing)

	s.energy -= cost
	if s.energy < 0 {
		return false
	}

	if !s.do(transition.Behavior) {
		return false
	}

	if !s.move(transition.Move) {
		return false
	}

	return true
}

func (s *CastingSpell) next(final Template) bool {
	next, ok := final.resolve(s.state)
	if !ok {
		return false
	}

	s.state = next
	return true
}
