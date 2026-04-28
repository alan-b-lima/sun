package spell

func (s *CastingSpell) move(move Move) bool {
	move = move.Colapse(s.facing)
	if !move.Valid(s.facing) {
		return false
	}

	s.X += move.X
	s.Y += move.Y
	return true
}

type Move struct {
	X, Y, Xf int
}

func (move Move) Valid(facing Facing) bool {
	move = move.Colapse(facing)
	return -1 <= move.X && move.X <= 1 && -1 <= move.Y && move.Y <= 1
}

func (move Move) Colapse(facing Facing) Move {
	switch facing {
	case FacingLeft:
		move.X -= move.Xf
	case FacingRight:
		move.X += move.Xf
	}

	move.Xf = 0
	return move
}

func MoveCost(move Move, facing Facing) Energy {
	move = move.Colapse(facing)
	return Energy(abs(move.X) + abs(move.Y))
}

type Facing bool

const (
	FacingLeft  Facing = true
	FacingRight Facing = false
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
