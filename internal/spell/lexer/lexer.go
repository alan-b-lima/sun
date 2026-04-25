package lexer

import (
	"strings"

	"github.com/alan-b-lima/sun/internal/spell/scanner"
)

type Stream struct {
	src *scanner.Source

	pos     scanner.Position
	token   Token
	content string

	await  bool
	insert bool
}

func NewStream(src *scanner.Source) *Stream {
	return &Stream{
		src:   src,
		token: Illegal,
	}
}

func (s *Stream) Next() {
	for {
		s.pos = s.src.Position()
		switch rune := s.next_rune(); {
		case rune == scanner.EOF:
			if s.insert_semicolon_after() {
				return
			}

			s.token = EOF
			return

		case is_identifier_char(rune):
			s.bubble()

			ident := s.next_identifier()
			if ident == "" {
				// impossible
				s.token = Illegal
				return
			}

			if keyword, in := keywords[ident]; in {
				s.token = keyword
				return
			}

			s.token = Identifier
			s.content = ident
			return

		case rune == '~':
			s.bubble()
			ident := s.next_atom()
			if ident == "" {
				s.token = Illegal
				return
			}

			s.token = Atom
			s.content = ident
			return

		case rune == '/':
			if rune := s.next_rune(); rune != '/' {
				s.pos = s.src.Position()
				s.token = Illegal
				return
			}

			for s.next_rune() != '\n' {
			}

			if s.insert_semicolon_after() {
				return
			}

		case rune == '}' && s.token != Semicolon:
			s.token = Semicolon
			s.bubble()
			return

		case rune == ';' && s.token == Semicolon:

		case rune == '\n':
			if s.insert_semicolon_after() {
				return
			}

		default:
			if token, in := punctuation[rune]; in {
				s.token = token
				return
			}
		}
	}
}

func (s *Stream) Token() Token {
	return s.token
}

func (s *Stream) Position() scanner.Position {
	return s.pos
}

func (s *Stream) Content() string {
	return s.content
}

func (s *Stream) next_rune() rune {
	if s.await {
		s.await = false
	} else {
		s.src.Next()
	}

	return s.src.Rune()
}

func (s *Stream) bubble() {
	s.await = true
}

func (s *Stream) next_identifier() string {
	var b strings.Builder

	for {
		rune := s.next_rune()
		if rune == scanner.EOF {
			break
		}

		if !is_identifier_char(rune) {
			s.bubble()
			break
		}

		b.WriteRune(rune)
	}

	return b.String()
}

func (s *Stream) next_atom() string {
	if rune := s.next_rune(); rune != '~' {
		return ""
	}

	ident := s.next_identifier()
	if ident == "" {
		return ""
	}

	if rune := s.next_rune(); rune != '~' {
		s.bubble()
		return ""
	}

	return ident
}

func (s *Stream) insert_semicolon_after() bool {
	switch s.token {
	case Identifier, Atom, RParen, RBrace:
		s.token = Semicolon
		return true
	}

	return false
}
