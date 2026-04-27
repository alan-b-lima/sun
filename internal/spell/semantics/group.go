package semantics

import (
	"slices"

	"github.com/alan-b-lima/sun/internal/atoms"
	"github.com/alan-b-lima/sun/internal/spell/parser"
)

type Group []atoms.Atom

type (
	group_solver map[Symbol]group

	group struct {
		Atoms []atoms.Atom

		groups []Symbol
		mark   mark
	}
)

func (solver group_solver) Add(scope *Scope, decl parser.AtomGroupDecl) error {
	symbol, err := scope.LookupKind(decl.Name.String(), KindGroup)
	if err != nil {
		return err
	}

	var group group
	for _, node := range decl.Atoms {
		switch node := node.(type) {
		case parser.Atom:
			atom := atoms.FromString(node.String())
			if !atom.Valid() {
				return ErrUnknownAtom
			}

			group.Atoms = append(group.Atoms, atom)

		case parser.Name:
			symbol, err := scope.LookupKind(node.String(), KindGroup)
			if err != nil {
				return err
			}

			group.groups = append(group.groups, symbol)
		}
	}

	solver[symbol] = group
	return nil
}

func (solver group_solver) Solve() (map[Symbol]Group, error) {
	groups := make(map[Symbol]Group, len(solver))

	for symbol := range solver {
		if err := visit(solver, symbol); err != nil {
			return nil, err
		}

		groups[symbol] = solver[symbol].Atoms
	}

	return groups, nil
}

func (solver group_solver) yield(symbol Symbol) {
	group := solver[symbol]

	atoms := group.Atoms
	for _, symbol := range group.groups {
		atoms = append(atoms, solver[symbol].Atoms...)
	}

	slices.Sort(atoms)
	group.Atoms = slices.Compact(atoms)
	group.groups = nil

	solver[symbol] = group
}

type mark int8

const (
	white mark = iota
	gray
	black
)

func visit(graph group_solver, symbol Symbol) error {
	switch mark := graph[symbol].mark; mark {
	case black:
		return nil
	case gray:
		return ErrGroupCycleFound
	case white:
	}

	graph.mark(symbol, gray)

	for _, symbol := range graph[symbol].groups {
		if err := visit(graph, symbol); err != nil {
			return err
		}
	}

	graph.mark(symbol, black)
	graph.yield(symbol)
	return nil
}

func (solver group_solver) mark(symbol Symbol, mark mark) {
	group := solver[symbol]
	group.mark = mark
	solver[symbol] = group
}
