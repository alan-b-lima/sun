package semantics

import (
	"errors"

	"github.com/alan-b-lima/sun/internal/atoms"
	"github.com/alan-b-lima/sun/internal/spell/parser"
)

type (
	State struct {
		Arity       int
		Transitions []Transition
	}

	Transition struct {
		AtomCond    AtomCond
		CellarConds CellarConds
		Behavior    Behavior
		Move        Move
		Final       Final
	}
)

type (
	state_solver map[Symbol]state

	state struct {
		Arity       int
		Transitions []parser.Transition
		Table       *Scope
	}
)

func (solver state_solver) Add(scope *Scope, decl parser.StateDecl) error {
	symbol, err := scope.LookupKind(decl.State.Name.String(), KindState)
	if err != nil {
		return err
	}

	scope = NewScope(scope)

	state := state{
		Arity:       len(decl.State.Params),
		Transitions: decl.Transitions,
		Table:       scope,
	}

	for _, param := range decl.State.Params {
		if _, ok := scope.Add(param.Name.String(), KindParam); !ok {
			return ErrNameConflict
		}
	}

	solver[symbol] = state
	return nil
}

func (solver state_solver) Solve(scope *Scope) (Symbol, map[Symbol]State, error) {
	initial, err := scope.LookupKind(InitialState, KindState)
	if err != nil {
		return NoSymbol, nil, ErrInitialStateNotFound
	}

	if solver[initial].Arity > 0 {
		return NoSymbol, nil, ErrInitialStateWithParams
	}

	states := make(map[Symbol]State)

	for symbol, state := range solver {
		transitions := make([]Transition, 0, len(state.Transitions))

		for _, stmt := range state.Transitions {
			transition, err := solve_transition(solver, state.Table, stmt)
			if err != nil {
				return NoSymbol, nil, err
			}

			transitions = append(transitions, transition)
		}

		states[symbol] = State{
			Arity:       state.Arity,
			Transitions: transitions,
		}
	}

	return initial, states, nil
}

func solve_transition(solver state_solver, scope *Scope, stmt parser.Transition) (Transition, error) {
	atom_cond, err := solve_atom_cond(scope, stmt.AtomCond)
	if err != nil {
		return Transition{}, err
	}

	cellar_conds, err := solve_cellar_conds(scope, stmt.CellarConds)
	if err != nil {
		return Transition{}, err
	}

	behavior, err := solve_behavior(scope, stmt.Behavior)
	if err != nil {
		return Transition{}, err
	}

	move, err := solve_move(stmt.Moves)
	if err != nil {
		return Transition{}, err
	}

	final, err := solve_final(solver, scope, stmt.Final, stmt.Halt)
	if err != nil {
		return Transition{}, err
	}

	return Transition{
		AtomCond:    atom_cond,
		CellarConds: cellar_conds,
		Behavior:    behavior,
		Move:        move,
		Final:       final,
	}, nil
}

type AtomCond struct {
	Atom    atoms.Atom
	Group   Symbol
	IsGroup bool
	IsAny   bool
}

func solve_atom_cond(scope *Scope, cond parser.AtomCond) (AtomCond, error) {
	switch cond.Tag {
	case parser.TagAny:
		return AtomCond{IsAny: true}, nil

	case parser.TagAtom:
		atom := atoms.FromString(cond.Atom.String())
		if !atom.Valid() {
			return AtomCond{}, ErrUnknownAtom
		}

		return AtomCond{Atom: atom}, nil

	case parser.TagName:
		symbol, err := scope.LookupKind(cond.Group.String(), KindGroup)
		if err != nil {
			return AtomCond{}, err
		}

		return AtomCond{Group: symbol, IsGroup: true}, nil
	}

	// unreachable
	return AtomCond{}, ErrImproperTag
}

type CellarConds map[Symbol]CellarCond

type CellarCond struct {
	Atom    atoms.Atom
	Group   Symbol
	Rune    Symbol
	IsGroup bool
	IsRune  bool
}

func solve_cellar_conds(scope *Scope, stmt []parser.CellarCond) (CellarConds, error) {
	conds := make(CellarConds, len(stmt))
	for _, expr := range stmt {
		cellar, cond, err := solve_cellar_cond(scope, expr)
		if err != nil {
			return nil, err
		}

		conds[cellar] = cond
	}

	return conds, nil
}

func solve_cellar_cond(scope *Scope, cond parser.CellarCond) (Symbol, CellarCond, error) {
	cellar, err := scope.LookupKind(cond.Cellar.String(), KindCellar)
	if err != nil {
		return NoSymbol, CellarCond{}, err
	}

	switch cond.Tag {
	case parser.TagAtom:
		atom := atoms.FromString(cond.Atom.String())
		if !atom.Valid() {
			return NoSymbol, CellarCond{}, ErrUnknownAtom
		}

		return cellar, CellarCond{Atom: atom}, nil

	default:
		// unreacheble
		return NoSymbol, CellarCond{}, ErrImproperTag

	case parser.TagName:
	}

	symbol, in := scope.Lookup(cond.Name.String())
	if !in {
		return NoSymbol, CellarCond{}, ErrSymbolNotFound
	}

	switch symbol {
	case KindGroup:
		return cellar, CellarCond{Group: symbol, IsGroup: true}, nil
	case KindRune:
		return cellar, CellarCond{Rune: symbol, IsRune: true}, nil
	}

	return NoSymbol, CellarCond{}, errors.Join(ErrNotGroupSymbol, ErrNotRuneSymbol)
}

type Behavior struct {
	Action Action
	Cellar Symbol
	Rune   Symbol
}

type Action int

const (
	ActionNil Action = iota
	ActionAbsorb
	ActionRelease
	ActionWrite
)

var actions = [...]Action{
	parser.ActionNil:     ActionNil,
	parser.ActionAbsorb:  ActionAbsorb,
	parser.ActionRelease: ActionRelease,
	parser.ActionWrite:   ActionWrite,
}

func solve_behavior(scope *Scope, behavior parser.Behavior) (Behavior, error) {
	switch behavior.Action {
	case parser.ActionNil:
		return Behavior{Action: ActionNil}, nil
	default:
		// unreachable
		return Behavior{}, ErrImproperTag
	case parser.ActionAbsorb, parser.ActionRelease, parser.ActionWrite:
	}

	cellar, err := scope.LookupKind(behavior.Cellar.String(), KindCellar)
	if err != nil {
		return Behavior{}, err
	}

	if behavior.Action != parser.ActionWrite {
		return Behavior{
			Action: actions[behavior.Action],
			Cellar: cellar,
		}, nil
	}

	rune, err := scope.LookupKind(behavior.Rune.String(), KindRune)
	if err != nil {
		return Behavior{}, err
	}

	return Behavior{
		Action: actions[behavior.Action],
		Cellar: cellar,
		Rune:   rune,
	}, nil
}

type Move struct {
	X, Y, Xf int
}

func solve_move(moves []parser.Move) (Move, error) {
	var move Move

	for _, m := range moves {
		switch m {
		case parser.MoveNil:
			// nothing
		case parser.MoveUp:
			move.Y -= 1
		case parser.MoveDown:
			move.Y += 1
		case parser.MoveLeft:
			move.X -= 1
		case parser.MoveRight:
			move.X += 1
		case parser.MoveFace:
			move.Xf += 1
		case parser.MoveBack:
			move.Xf -= 1
		}
	}

	var (
		boundX  = -1 <= move.X && move.X <= 1
		boundY  = -1 <= move.Y && move.Y <= 1
		boundXf = -1 <= move.Xf && move.Xf <= 1
	)
	if !boundX || !boundY || !boundXf {
		return Move{}, ErrMoveIllegal
	}

	return move, nil
}

type Final struct {
	Symbol  Symbol
	Params  []Final
	IsVar   bool
	Halting bool
}

func solve_final(solver state_solver, scope *Scope, stmt parser.State, halts bool) (Final, error) {
	if halts {
		return Final{Halting: true}, nil
	}

	return solve_final_params(solver, scope, stmt)
}

func solve_final_params(solver state_solver, scope *Scope, stmt parser.State) (Final, error) {
	symbol, in := scope.Lookup(stmt.Name.String())
	if !in {
		return Final{}, ErrSymbolNotFound
	}

	switch symbol.Kind() {
	case KindState:
		if solver[symbol].Arity != len(stmt.Params) {
			return Final{}, ErrBadArity
		}

		final := Final{
			Symbol: symbol,
			Params: make([]Final, 0, len(stmt.Params)),
		}

		for _, param := range stmt.Params {
			state, err := solve_final_params(solver, scope, param)
			if err != nil {
				return Final{}, err
			}

			final.Params = append(final.Params, state)
		}

		return final, nil

	case KindParam:
		if len(stmt.Params) > 0 {
			return Final{}, ErrParamWithParams
		}

		return Final{
			Symbol: symbol,
			IsVar:  true,
		}, nil

	default:
		return Final{}, ErrNotStateSymbol
	}
}
