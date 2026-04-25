package spell

func (s *Spell) move(move Move) bool {
	if !move.Valid(s.facing) {
		return false
	}

	s.x += move.X
	s.y += move.Y

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

	return move
}

func MoveCost(move Move, facing Facing) int {
	move = move.Colapse(facing)
	return abs(move.X) + abs(move.Y)
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
