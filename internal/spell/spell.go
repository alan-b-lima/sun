package spell

import "github.com/alan-b-lima/sun/internal/atom"

type Spell struct {
	Cellars []Cellar
	States  []State
	Initial int

	// runtime info

	world  World
	x, y   int
	facing Facing

	energy Energy
	halted bool
	state  state
}

type World interface {
	Dim() (w, h int)

	At(x, y int) atom.Atom
	Set(x, y int, atom atom.Atom)
}

type Energy int64

func Make(cellars []Cellar, states []State, initial int) Spell {
	return Spell{
		Cellars: cellars,
		States:  states,
		Initial: initial,
	}
}

func (s Spell) Cast(world World, energy Energy, x, y int, facing Facing) Spell {
	s.world = world
	s.x = x
	s.y = y
	s.facing = facing

	s.energy = energy - CellarCost*Energy(len(s.Cellars))
	s.halted = s.energy < 0
	s.state = state{State: s.Initial}

	return s
}

func (s *Spell) Zero() {
	*s = Spell{
		Cellars: s.Cellars,
		States:  s.States,
		Initial: s.Initial,
	}
}

func (s *Spell) Energy() int  { return int(s.energy) }
func (s *Spell) Halted() bool { return s.halted }

func (s *Spell) State(index int) (State, bool) {
	if index >= len(s.States) {
		return State{}, false
	}
	return s.States[index], true
}

func (s *Spell) Cellar(index int) (Cellar, bool) {
	if index >= len(s.Cellars) {
		return Cellar{}, false
	}
	return s.Cellars[index], true
}

func (s *Spell) Perform() {
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

func (s *Spell) find(state State) (Transition, bool) {
	for _, transition := range state.Transitions {
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

func (s *Spell) exec(transition Transition) bool {
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

func (s *Spell) next(final Template) bool {
	next, ok := final.resolve(s.state)
	if !ok {
		return false
	}

	s.state = next
	return true
}
