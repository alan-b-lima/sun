package parser

import (
	"fmt"
	"slices"
)

func (s Tree) String() string {
	b := builder{tab: "    "}

	b.WriteString("SpellSource")
	b.Indent()
	b.WriteString("\n")
	stringify(&b, s.Spell)
	b.Dedent()

	return b.String()
}

func (s SpellSource) string(b *builder) {
	b.WriteString("Name(")
	stringify(b, s.Name)
	b.WriteString(")")

	if len(s.Decls) == 0 {
		b.WriteString("\nDeclList(nil)")
		return
	}

	b.WriteString("\nDeclList")
	b.Indent()
	for _, decl := range s.Decls {
		switch decl.(type) {
		case AtomGroupDecl:
			b.WriteString("\nAtomGroupDecl")
		case CellarDecl:
			b.WriteString("\nCellarDecl")
		case RuneDecl:
			b.WriteString("\nRuneDecl")
		case StateDecl:
			b.WriteString("\nStateDecl")
		}

		b.Indent()
		b.WriteString("\n")
		stringify(b, decl)
		b.Dedent()
	}
	b.Dedent()
}

func (s Atom) string(b *builder) {
	b.WriteString("~")
	b.WriteString(string(s))
	b.WriteString("~")
}

func (s Name) string(b *builder) {
	b.WriteString(string(s))
}

func (s AtomGroupDecl) string(b *builder) {
	b.WriteString("Name(")
	stringify(b, s.Name)
	b.WriteString(")")

	if len(s.Atoms) == 0 {
		b.WriteString("\nAtomList(nil)")
		return
	}

	b.WriteString("\nAtomList")
	b.Indent()
	for _, atom := range s.Atoms {
		switch atom.(type) {
		case Atom:
			b.WriteString("\nAtom(")
		case Name:
			b.WriteString("\nName(")
		}

		stringify(b, atom)
		b.WriteString(")")
	}
	b.Dedent()
}

func (s CellarDecl) string(b *builder) {
	b.WriteString("Name(")
	stringify(b, s.Name)
	b.WriteString(")")
}

func (s RuneDecl) string(b *builder) {
	b.WriteString("Name(")
	stringify(b, s.Name)
	b.WriteString(")")
}

func (s StateDecl) string(b *builder) {
	b.WriteString("State")
	b.Indent()
	b.WriteString("\n")
	stringify(b, s.State)
	b.Dedent()

	if len(s.Transitions) == 0 {
		b.WriteString("\nTransitionList(nil)")
		return
	}

	b.WriteString("\nTransitionList")
	b.Indent()
	for _, transition := range s.Transitions {
		b.WriteString("\nTransition")
		b.Indent()
		b.WriteString("\n")
		stringify(b, transition)
		b.Dedent()
	}
	b.Dedent()
}

func (s State) string(b *builder) {
	b.WriteString("Name(")
	stringify(b, s.Name)
	b.WriteString(")")

	if len(s.Params) == 0 {
		b.WriteString("\nParamList(nil)")
		return
	}

	b.WriteString("\nParamList")
	b.Indent()
	for _, param := range s.Params {
		b.WriteString("\nParam")
		b.Indent()
		b.WriteString("\n")
		stringify(b, param)
		b.Dedent()
	}
	b.Dedent()
}

func (s Transition) string(b *builder) {
	b.WriteString("AtomCond(")
	stringify(b, s.AtomCond)
	b.WriteString(")")

	if len(s.CellarConds) != 0 {
		b.WriteString("\nCellarCondList")
		b.Indent()
		for _, cond := range s.CellarConds {
			b.WriteString("\nCellarCond")
			b.Indent()
			b.WriteString("\n")
			stringify(b, cond)
			b.Dedent()
		}
		b.Dedent()
	} else {
		b.WriteString("\nCellarCondList(nil)")
	}

	b.WriteString("\nBehavior")
	b.Indent()
	b.WriteString("\n")
	stringify(b, s.Behavior)
	b.Dedent()

	if len(s.Moves) != 0 {
		b.WriteString("\nMoveList")
		b.Indent()
		for _, move := range s.Moves {
			b.WriteString("\nMove(")
			stringify(b, move)
			b.WriteString(")")
		}
		b.Dedent()
	} else {
		b.WriteString("\nMoveList(nil)")
	}

	if !s.Halt {
		b.WriteString("\nFinal")
		b.Indent()
		b.WriteString("\n")
		stringify(b, s.Final)
		b.Dedent()
	} else {
		b.WriteString("\nFinal(nil)")
	}
}

func (s AtomCond) string(b *builder) {
	switch s.Tag {
	case TagAny:
		b.WriteString("nil")
	case TagAtom:
		stringify(b, s.Atom)
	case TagName:
		stringify(b, s.Group)
	}
}

func (s CellarCond) string(b *builder) {
	b.WriteString("Cellar(")
	stringify(b, s.Name)
	b.WriteString(")")

	b.WriteString("\nAtom(")
	switch s.Tag {
	case TagAtom:
		stringify(b, s.Atom)
	case TagName:
		stringify(b, s.Name)
	}
	b.WriteString(")")
}

func (s Behavior) string(b *builder) {
	switch s.Action {
	case ActionNil:
		b.WriteString("Nil")
		return
	case ActionAbsorb:
		b.WriteString("Absorb")
	case ActionRelease:
		b.WriteString("Release")
	case ActionWrite:
		b.WriteString("Write")
	}

	b.WriteString("\n")
	stringify(b, s.Cellar)

	if s.Action == ActionWrite {
		b.WriteString("\n")
		stringify(b, s.Rune)
	}
}

func (s Move) string(b *builder) {
	switch s {
	case MoveNil:
		b.WriteString("nil")
	case MoveUp:
		b.WriteString("up")
	case MoveDown:
		b.WriteString("down")
	case MoveLeft:
		b.WriteString("left")
	case MoveRight:
		b.WriteString("right")
	case MoveFace:
		b.WriteString("face")
	case MoveBack:
		b.WriteString("back")
	}
}

func stringify(b *builder, v any) {
	switch v := v.(type) {
	case interface{ string(*builder) }:
		v.string(b)

	case fmt.Stringer:
		b.WriteString(v.String())

	default:
		fmt.Fprint(b, v)
	}
}

type builder struct {
	buf     []byte
	advance int
	tab     string
}

func (b *builder) Indent() {
	b.advance++
}

func (b *builder) Dedent() {
	b.advance--
}

func (b *builder) Write(buf []byte) (int, error) {
	return write(b, buf)
}

func (b *builder) WriteString(buf string) (int, error) {
	return write(b, buf)
}

func write[T ~[]byte | ~string](b *builder, buf T) (int, error) {
	var nls int
	for i := range len(buf) {
		if buf[i] == '\n' {
			nls++
		}
	}
	if nls == 0 {
		b.buf = append(b.buf, buf...)
		return len(buf), nil
	}

	n := len(b.buf)
	m := n + len(buf) + b.advance*len(b.tab)*nls

	b.buf = slices.Grow(b.buf, m)

	var last int
	for i := range len(buf) {
		if buf[i] == '\n' {
			b.buf = append(b.buf, buf[last:i+1]...)
			last = i + 1

			for range b.advance {
				b.buf = append(b.buf, b.tab...)
			}
		}
	}
	if last < len(buf) {
		b.buf = append(b.buf, buf[last:]...)
	}

	return len(b.buf) - n, nil
}

func (b *builder) String() string {
	return string(b.buf)
}
