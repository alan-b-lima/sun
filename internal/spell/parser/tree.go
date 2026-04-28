package parser

type Tree struct {
	Spell SpellSource
}

type (
	Decl     interface{ decl() }
	AtomNode interface{ atom() }
)

type (
	SpellSource struct {
		Name  Name
		Decls []Decl
	}

	Name string
	Atom string

	AtomGroupDecl struct {
		Name  Name
		Atoms []AtomNode
	}

	CellarDecl struct {
		Name Name
	}

	RuneDecl struct {
		Name Name
	}

	StateDecl struct {
		State State
		Lines []Line
	}

	Line struct {
		AtomCond    AtomCond
		CellarConds []CellarCond
		Behavior    Behavior
		Moves       []Move
		Final       State
		Halt        bool
	}

	AtomCond struct {
		Tag   Tag
		Atom  Atom
		Group Name
	}

	CellarCond struct {
		Cellar Name
		Tag    Tag
		Atom   Atom
		Name   Name
	}

	Behavior struct {
		Action Action
		Cellar Name
		Rune   Name
	}

	State struct {
		Name   Name
		Params []State
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
	MoveNil Move = iota
	MoveUp
	MoveDown
	MoveLeft
	MoveRight
	MoveFace
	MoveBack
)

func (AtomGroupDecl) decl() {}
func (CellarDecl) decl()    {}
func (RuneDecl) decl()      {}
func (StateDecl) decl()     {}

func (Atom) atom() {}
func (Name) atom() {}
