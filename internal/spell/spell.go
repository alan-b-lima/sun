package spell

import "github.com/alan-b-lima/sun/internal/atoms"

type CastingSpell struct {
	*Spell

	// runtime info

	world  World
	X, Y   int
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
	s.X = x
	s.Y = y
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

func (s *CastingSpell) Cellar(index int) (*Cellar, bool) {
	if index >= len(s.cellars) {
		return &Cellar{}, false
	}
	return &s.cellars[index], true
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

	line, ok := s.find(state)
	if !ok {
		s.halted = true
		return
	}

	if !s.exec(line) {
		s.halted = true
		return
	}

	if !s.next(line.Final) {
		s.halted = true
	}
}

func (s *CastingSpell) find(state State) (Line, bool) {
Lines:
	for _, line := range state {
		atom := s.world.At(s.X, s.Y)
		if !line.AtomCond.For(atom) {
			continue
		}

		for _, cond := range line.CellarConds {
			cellar, ok := s.Cellar(cond.Cellar)
			if !ok {
				continue Lines
			}

			top, ok := cellar.Peek()
			if !ok {
				continue Lines
			}

			if !cond.For(top) {
				continue Lines
			}
		}

		return line, true
	}

	return Line{}, false
}

func (s *CastingSpell) exec(line Line) bool {
	cost := 1 +
		BehaviorCost[line.Behavior.Action] +
		MoveCost(line.Move, s.facing)

	s.energy -= cost
	if s.energy < 0 {
		return false
	}

	if !s.do(line.Behavior) {
		return false
	}

	if !s.move(line.Move) {
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
