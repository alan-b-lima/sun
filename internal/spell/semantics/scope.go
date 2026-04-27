package semantics

import (
	"math/bits"
	"strconv"
)

type Scope struct {
	Parent   *Scope
	Children []*Scope

	Names   map[Symbol]Name
	Symbols map[Name]Symbol
}

func NewScope(parent *Scope) *Scope {
	scope := &Scope{
		Parent:  parent,
		Names:   make(map[Symbol]Name),
		Symbols: make(map[Name]Symbol),
	}

	if parent != nil {
		parent.Children = append(parent.Children, scope)
	}

	return scope
}

var base int

func (s *Scope) Add(name Name, kind Symbol) (Symbol, bool) {
	if _, in := s.Symbols[name]; in {
		return NoSymbol, false
	}

	base++
	symbol := NewSymbol(base, kind)

	s.Names[symbol] = name
	s.Symbols[name] = symbol

	return symbol, true
}

func (s *Scope) Lookup(name Name) (sym Symbol, found bool) {
	for s != nil {
		if symbol, in := s.Symbols[name]; in {
			return symbol, true
		}
		s = s.Parent
	}

	return NoSymbol, false
}

func (s *Scope) LookupKind(name Name, kind Symbol) (Symbol, error) {
	symbol, found := s.Lookup(name)
	if !found {
		return NoSymbol, ErrSymbolNotFound
	}
	if symbol.Kind() != kind {
		switch kind {
		case KindGroup:
			return NoSymbol, ErrNotGroupSymbol
		case KindCellar:
			return NoSymbol, ErrNotCellarSymbol
		case KindRune:
			return NoSymbol, ErrNotRuneSymbol
		case KindState:
			return NoSymbol, ErrNotStateSymbol
		case KindParam:
			return NoSymbol, ErrNotParamSymbol
		}

		return NoSymbol, ErrImproperTag
	}

	return symbol, nil
}

type Name = string

type Symbol int64

const (
	KindGroup Symbol = iota
	KindCellar
	KindRune
	KindState
	KindParam

	kinds_count

	NoSymbol Symbol = -1
)

var (
	kind_bits = Symbol(bits.Len64(uint64(kinds_count - 1)))
	kind_mask = Symbol(1<<kind_bits) - 1
)

func NewSymbol(base int, kind Symbol) Symbol {
	return (Symbol(base) << kind_bits) | kind
}

func (k Symbol) Kind() Symbol {
	return k & kind_mask
}

func (k Symbol) String() string {
	var s string

	kind := k & kind_mask
	if kind < kinds_count {
		s = kinds[kind]
	}
	if s == "" {
		s = "type(" + strconv.FormatInt(int64(kind), 10) + ")"
	}

	if base := k >> kind_bits; base > 0 {
		return s + "[" + strconv.FormatInt(int64(base), 10) + "]"
	}
	return s
}

var kinds = [...]string{
	KindGroup:  "Group",
	KindCellar: "Cellar",
	KindRune:   "Rune",
	KindState:  "State",
	KindParam:  "Param",
}
