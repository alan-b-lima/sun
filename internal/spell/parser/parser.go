package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alan-b-lima/sun/internal/spell/lexer"
	"github.com/alan-b-lima/sun/internal/spell/scanner"
)

func Parse(stream *lexer.Stream) (*Tree, error) {
	parser := parser{Stream: stream}

	var root SpellSource
	root.parse(&parser)

	if err := errors.Join(parser.errors...); err != nil {
		return nil, err
	}
	return &Tree{Spell: root}, nil
}

type parser struct {
	*lexer.Stream
	errors []error

	await bool
}

func (s *parser) Next() {
	if s.await {
		s.await = false
	} else {
		s.Stream.Next()
	}
}

func (p *parser) LookFor(t lexer.Token) bool {
	p.Next()
	if token := p.Token(); token != t {
		p.ErrExpected(t)
		return false
	}

	return true
}

func (s *parser) bubble() {
	s.await = true
}

func (p parser) TokenString() string {
	switch token := p.Token(); {
	case token.IsLiteral():
		return token.String() + "(" + p.Content() + ")"
	default:
		return token.String()
	}
}

func (p *parser) Error(message string) {
	p.errors = append(p.errors, &Error{
		Position: p.Position(),
		Message:  message,
	})
}

func (p *parser) Errorf(format string, v ...any) {
	p.Error(fmt.Sprintf(format, v...))
}

func (p *parser) ErrExpected(tokens ...lexer.Token) {
	switch len(tokens) {
	case 0:
		p.Errorf("expected nothing, got '%v'", p.TokenString())

	case 1:
		p.Errorf("expected '%v', got '%v'", tokens[0], p.TokenString())

	default:
		var b strings.Builder

		b.WriteString("expected ")
		for _, token := range tokens[:len(tokens)-1] {
			fmt.Fprintf(&b, "'%v', ", token)
		}
		fmt.Fprintf(&b, "or '%v', got '%v'", tokens[len(tokens)-1], p.TokenString())

		p.Error(b.String())
	}
}

type Error struct {
	Position scanner.Position
	Message  string
}

func (err *Error) Error() string {
	return err.Position.String() + ": " + err.Message
}

func (s *SpellSource) parse(parser *parser) {
	if !parser.LookFor(lexer.Spell) {
		return
	}

	if !parser.LookFor(lexer.Identifier) {
		return
	}

	s.Name = Name(parser.Content())

	if !parser.LookFor(lexer.Semicolon) {
		return
	}

	for {
		parser.Next()
		switch token := parser.Token(); token {
		case lexer.Group:
			var decl AtomGroupDecl
			decl.parse(parser)
			s.Decls = append(s.Decls, decl)

		case lexer.Cellar:
			var decl CellarDecl
			decl.parse(parser)
			s.Decls = append(s.Decls, decl)

		case lexer.Rune:
			var decl RuneDecl
			decl.parse(parser)
			s.Decls = append(s.Decls, decl)

		case lexer.State:
			var decl StateDecl
			decl.parse(parser)
			s.Decls = append(s.Decls, decl)

		case lexer.EOF:
			return

		default:
			parser.ErrExpected(lexer.Group, lexer.Cellar, lexer.Rune, lexer.State)
			return
		}

		if !parser.LookFor(lexer.Semicolon) {
			return
		}
	}
}

func (s *AtomGroupDecl) parse(parser *parser) {
	if !parser.LookFor(lexer.Identifier) {
		return
	}
	s.Name = Name(parser.Content())

	if !parser.LookFor(lexer.LBrace) {
		return
	}

	for {
		parser.Next()
		switch token := parser.Token(); token {
		case lexer.Atom:
			s.Atoms = append(s.Atoms, Atom(parser.Content()))

		case lexer.Identifier:
			s.Atoms = append(s.Atoms, Name(parser.Content()))

		case lexer.RBrace:
			return

		default:
			parser.ErrExpected(lexer.Atom, lexer.Identifier, lexer.RBrace)
		}

		if !parser.LookFor(lexer.Semicolon) {
			return
		}
	}
}

func (s *CellarDecl) parse(parser *parser) {
	if !parser.LookFor(lexer.Identifier) {
		return
	}
	s.Name = Name(parser.Content())
}

func (s *RuneDecl) parse(parser *parser) {
	if !parser.LookFor(lexer.Identifier) {
		return
	}
	s.Name = Name(parser.Content())
}

func (s *StateDecl) parse(parser *parser) {
	s.State.parse(parser)

	if !parser.LookFor(lexer.LBrace) {
		return
	}

	for {
		parser.Next()
		if parser.Token() == lexer.RBrace {
			return
		}
		parser.bubble()

		var transition Transition
		transition.parse(parser)

		s.Transitions = append(s.Transitions, transition)

		if !parser.LookFor(lexer.Semicolon) {
			return
		}
	}
}

func (s *State) parse(parser *parser) {
	if !parser.LookFor(lexer.Identifier) {
		return
	}
	s.Name = Name(parser.Content())

	parser.Next()
	if parser.Token() != lexer.LParen {
		parser.bubble()
		return
	}

	parser.Next()
	if parser.Token() == lexer.RParen {
		return
	}
	parser.bubble()

	for {
		var state State
		state.parse(parser)

		s.Params = append(s.Params, state)

		parser.Next()
		switch parser.Token() {
		case lexer.Comma:
		case lexer.RParen:
			return
		default:
			parser.ErrExpected(lexer.Comma, lexer.RParen)
			return
		}
	}
}

func (s *Transition) parse(parser *parser) {
	s.AtomCond.parse(parser)

	parser.Next()
	if parser.Token() != lexer.Dot {
		parser.bubble()
		for {
			var cond CellarCond
			cond.parse(parser)

			s.CellarConds = append(s.CellarConds, cond)

			parser.Next()
			if parser.Token() != lexer.Comma {
				parser.bubble()
				break
			}
		}
	}

	s.Behavior.parse(parser)

	for {
		var move Move
		move.parse(parser)

		s.Moves = append(s.Moves, move)

		parser.Next()
		if parser.Token() != lexer.Comma {
			parser.bubble()
			break
		}
	}

	parser.Next()
	parser.bubble()

	if parser.Token() == lexer.Semicolon {
		s.Halt = true
		return
	}
	s.Final.parse(parser)
}

func (s *AtomCond) parse(parser *parser) {
	parser.Next()
	switch parser.Token() {
	case lexer.Identifier:
		s.Tag = TagName
		s.Group = Name(parser.Content())

	case lexer.Atom:
		s.Tag = TagAtom
		s.Atom = Atom(parser.Content())

	case lexer.Dot:
		s.Tag = TagAny

	default:
		parser.ErrExpected(lexer.Identifier, lexer.Atom)
	}
}

func (s *CellarCond) parse(parser *parser) {
	if !parser.LookFor(lexer.Identifier) {
		return
	}
	s.Cellar = Name(parser.Content())

	if !parser.LookFor(lexer.LBrack) {
		return
	}

	parser.Next()
	switch parser.Token() {
	case lexer.Identifier:
		s.Tag = TagName
		s.Name = Name(parser.Content())

	case lexer.Atom:
		s.Tag = TagAtom
		s.Atom = Atom(parser.Content())

	default:
		parser.ErrExpected(lexer.Identifier, lexer.Atom)
	}

	if !parser.LookFor(lexer.RBrack) {
		return
	}
}

func (s *Behavior) parse(parser *parser) {
	parser.Next()
	switch token := parser.Token(); token {
	case lexer.Nil:
		s.Action = ActionNil

	case lexer.Absorb, lexer.Release:
		s.Action = ActionAbsorb
		if token == lexer.Release {
			s.Action = ActionRelease
		}

		if !parser.LookFor(lexer.LBrack) {
			return
		}

		if !parser.LookFor(lexer.Identifier) {
			return
		}
		s.Cellar = Name(parser.Content())

		if !parser.LookFor(lexer.RBrack) {
			return
		}

	case lexer.Write:
		s.Action = ActionWrite

		if !parser.LookFor(lexer.LBrack) {
			return
		}

		if !parser.LookFor(lexer.Identifier) {
			return
		}
		s.Cellar = Name(parser.Content())

		if !parser.LookFor(lexer.Comma) {
			return
		}

		if !parser.LookFor(lexer.Identifier) {
			return
		}
		s.Rune = Name(parser.Content())

		if !parser.LookFor(lexer.RBrack) {
			return
		}

	default:
		parser.ErrExpected(lexer.Nil, lexer.Absorb, lexer.Release, lexer.Write)
	}
}

func (s *Move) parse(parser *parser) {
	parser.Next()
	switch parser.Token() {
	case lexer.Nil:
		*s = MoveNone
	case lexer.Up:
		*s = MoveUp
	case lexer.Down:
		*s = MoveDown
	case lexer.Left:
		*s = MoveLeft
	case lexer.Right:
		*s = MoveRight
	case lexer.Face:
		*s = MoveFace
	case lexer.Back:
		*s = MoveBack

	default:
		parser.ErrExpected(lexer.Nil, lexer.Down, lexer.Left, lexer.Right, lexer.Face, lexer.Back)
	}
}
