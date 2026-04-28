package spell

import "github.com/alan-b-lima/sun/internal/atoms"

func (s *CastingSpell) do(behavior Behavior) bool {
	switch behavior.Action {
	case BehaviorNil:

	case BehaviorAbsorb:
		cellar, ok := s.Cellar(behavior.Cellar)
		if !ok {
			return false
		}

		a := s.world.At(s.X, s.Y)
		s.world.Set(s.X, s.Y, atoms.Air)
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
			res, ok := mix(s.world.At(s.X, s.Y), cell.Atom)
			if !ok {
				return false
			}

			s.world.Set(s.X, s.Y, res)
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

func mix(base, over atoms.Atom) (atoms.Atom, bool) {
	switch base {
	case atoms.Nil, atoms.Air:
		return over, true
	}

	return atoms.Nil, false
}
