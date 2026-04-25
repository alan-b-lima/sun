package spell

import "github.com/alan-b-lima/sun/internal/atom"

func (s *Spell) do(behavior Behavior) bool {
	switch behavior.Action {
	case BehaviorNil:

	case BehaviorAbsorb:
		cellar, ok := s.Cellar(behavior.Cellar)
		if !ok {
			return false
		}

		a := s.world.At(s.x, s.y)
		s.world.Set(s.x, s.y, atom.Air)
		cellar.PushAtom(a)

	case BehaviorRelease:
		cellar, ok := s.Cellar(behavior.Cellar)
		if !ok {
			return false
		}

		cell, ok := cellar.Pop()
		if !ok {
			return false
		}

		if cell.Tag == TagAtom {
			res, ok := mix(s.world.At(s.x, s.y), cell.Atom)
			if !ok {
				return false
			}

			s.world.Set(s.x, s.y, res)
		}

	case BehaviorWrite:
		cellar, ok := s.Cellar(behavior.Cellar)
		if !ok {
			return false
		}

		cellar.PushRune(behavior.Rune)
	}

	return true
}

type Behavior struct {
	Action Action
	Cellar int
	Rune   Rune
}

type Action int

const (
	BehaviorNil Action = iota
	BehaviorAbsorb
	BehaviorRelease
	BehaviorWrite
)

var BehaviorCost = [...]Energy{
	BehaviorNil:     0,
	BehaviorAbsorb:  3,
	BehaviorRelease: 2,
	BehaviorWrite:   1,
}

func mix(base, over atom.Atom) (atom.Atom, bool) {
	switch base {
	case atom.Nil, atom.Air:
		return over, true
	}

	return atom.Nil, false
}
