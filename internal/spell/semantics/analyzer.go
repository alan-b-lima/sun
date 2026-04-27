package semantics

import (
	"errors"

	"github.com/alan-b-lima/sun/internal/spell/parser"
)

type Analysis struct {
	Outcome

	groups group_solver
	states state_solver
}

type Outcome struct {
	Groups  map[Symbol]Group
	Cellars []Symbol
	Runes   []Symbol
	States  map[Symbol]State
	Initial Symbol

	Table *Scope
}

func Analyze(ast *parser.Tree) (Outcome, error) {
	asis, err := Make(ast.Spell.Decls)
	if err != nil {
		return Outcome{}, err
	}

	if err := asis.Mount(ast.Spell.Decls); err != nil {
		return Outcome{}, err
	}

	return asis.Outcome, nil
}

const InitialState = "start"

var (
	ErrNameConflict   = errors.New("name conflict")
	ErrSymbolNotFound = errors.New("symbol not found")

	ErrUnknownAtom = errors.New("unknown atom")

	ErrNotGroupSymbol  = errors.New("not a group symbol")
	ErrGroupCycleFound = errors.New("group dependency cycle found")

	ErrNotCellarSymbol = errors.New("not a cellar symbol")

	ErrNotRuneSymbol = errors.New("not a rune symbol")

	ErrNotStateSymbol         = errors.New("not a state symbol")
	ErrInitialStateNotFound   = errors.New("initial state not found")
	ErrInitialStateWithParams = errors.New("initial state must not have parameters")

	ErrImproperTag = errors.New("improper tag was generated")

	ErrMoveIllegal = errors.New("illegal move")

	ErrNotParamSymbol  = errors.New("not a state parameter symbol")
	ErrBadArity        = errors.New("expected a different number of parameters")
	ErrParamWithParams = errors.New("parameter must not have parameters")
)

func Make(decls []parser.Decl) (Analysis, error) {
	var groups, cellars, runes, states int
	table := NewScope(nil)

	for _, decl := range decls {
		var ok bool
		switch decl := decl.(type) {
		case parser.AtomGroupDecl:
			_, ok = table.Add(decl.Name.String(), KindGroup)
			groups++
		case parser.CellarDecl:
			_, ok = table.Add(decl.Name.String(), KindCellar)
			cellars++
		case parser.RuneDecl:
			_, ok = table.Add(decl.Name.String(), KindRune)
			runes++
		case parser.StateDecl:
			_, ok = table.Add(decl.State.Name.String(), KindState)
			states++
		}

		if !ok {
			return Analysis{}, ErrNameConflict
		}
	}

	return Analysis{
		Outcome: Outcome{
			Cellars: make([]Symbol, 0, cellars),
			Runes:   make([]Symbol, 0, runes),
			Table:   table,
		},

		groups: make(group_solver, groups),
		states: make(state_solver, states),
	}, nil
}

func (asis *Analysis) Mount(decls []parser.Decl) error {
	for _, decl := range decls {
		switch decl := decl.(type) {
		case parser.AtomGroupDecl:
			if err := asis.groups.Add(asis.Table, decl); err != nil {
				return err
			}

		case parser.CellarDecl:
			symbol, _ := asis.Table.Lookup(decl.Name.String())
			asis.Cellars = append(asis.Cellars, symbol)

		case parser.RuneDecl:
			symbol, _ := asis.Table.Lookup(decl.Name.String())
			asis.Runes = append(asis.Runes, symbol)

		case parser.StateDecl:
			if err := asis.states.Add(asis.Table, decl); err != nil {
				return err
			}
		}
	}

	var err error

	asis.Groups, err = asis.groups.Solve()
	if err != nil {
		return err
	}

	asis.Initial, asis.States, err = asis.states.Solve(asis.Table)
	if err != nil {
		return err
	}

	return nil
}
