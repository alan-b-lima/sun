package spell

import "github.com/alan-b-lima/sun/internal/atoms"

type Cellar struct {
	Stack []Cell
}

type Rune uint16

func (c *Cellar) Push(cell Cell) {
	c.Stack = append(c.Stack, cell)
}

func (c *Cellar) PushAtom(atom atoms.Atom) {
	c.Push(Cell{Atom: atom, Tag: TagAtom})
}

func (c *Cellar) PushRune(rune Rune) {
	c.Push(Cell{Rune: rune, Tag: TagRune})
}

func (c *Cellar) Pop() (Cell, bool) {
	if len(c.Stack) == 0 {
		return Cell{}, false
	}

	cell := c.Stack[len(c.Stack)-1]
	c.Stack = c.Stack[:len(c.Stack)-1]
	return cell, true
}

func (c *Cellar) Peek() (Cell, bool) {
	if len(c.Stack) == 0 {
		return Cell{}, false
	}

	cell := c.Stack[len(c.Stack)-1]
	return cell, true
}

const CellarCost = 5

type Cell struct {
	Rune Rune
	Atom atoms.Atom
	Tag  Tag
}

type Tag bool

const (
	TagAtom Tag = true
	TagRune Tag = false
)
