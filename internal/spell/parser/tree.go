package parser

type Tree struct {
	Spell SpellSource
}

type (
	SpellSource struct {
		Name  string
		Decls DeclList
	}

	DeclList struct {
		Groups  []AtomGroupDecl
		Cellars []CellarDecl
		Runes   []RuneDecl
		States  []StateDecl
	}

	AtomGroupDecl struct {
		Name  string
		Atoms AtomList
	}

	AtomList struct {
		Atoms  []Atom
		Groups []AtomGroup
	}

	Atom string
	Name string

	AtomGroup struct {
		Name string
	}

	CellarDecl struct {
		Name string
	}

	RuneDecl struct {
		Name string
	}

	StateDecl struct {
		State       State
		Transitions []Transition
	}

	Transition struct {
		AtomCond    AtomCond
		CellarConds []CellarCond
		Behavior    Behavior
		Moves       []Move
		Final       State
		AlwaysHalt  bool
	}

	AtomCond struct {
		Tag   Tag
		Atom  Atom
		Group AtomGroup
	}

	CellarCond struct {
		Cellar Name
		Tag    Tag
		Atom   Atom
		Name   Name
	}

	AtomGroupOrRune struct {
		Name string
	}

	Behavior struct {
		Action Action
		Cellar Name
		Rune   Name
	}

	State struct {
		Name  string
		Param []State
	}
)

type Tag int

const (
	TagAny Tag = iota
	TagAtom
	TagName
)

type Action int

const (
	ActionNil Action = iota
	ActionAbsorb
	ActionRelease
	ActionWrite
)

type Move int

const (
	MoveUp Move = iota
	MoveDown
	MoveLeft
	MoveRight
	MoveFace
	MoveBack
)
