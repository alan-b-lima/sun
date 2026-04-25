package scanner

import (
	"fmt"
	"io"
	"unicode/utf8"
)

const (
	EOF rune = -1
)

type Source struct {
	source    []byte
	character rune

	offset       int
	line, column int
}

func New(src []byte) *Source {
	return &Source{
		source:    src,
		character: ' ',
		offset:    0,
		line:      1,
		column:    1,
	}
}

func NewFromReader(r io.Reader) (*Source, error) {
	src, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return New(src), nil
}

func (s *Source) Next() {
	if s.offset >= len(s.source) {
		s.character = EOF
		return
	}

	rune, size := utf8.DecodeRune(s.source[s.offset:])
	s.offset += size

	if rune == '\n' {
		s.line++
		s.column = 0
	}
	s.column++

	s.character = rune
}

func (s *Source) Position() Position {
	return Position{
		Offset: s.offset,
		Line:   s.line,
		Column: s.column,
	}
}

func (s *Source) Rune() rune {
	return s.character
}

type Position struct {
	Offset int
	Line   int
	Column int
}

func (p Position) String() string {
	if p.Line > 0 {
		if p.Column > 0 {
			return fmt.Sprintf("%d:%d", p.Line, p.Column)
		} else {
			return fmt.Sprintf("%d", p.Line)
		}
	}

	return "-"
}
