package atoms

import "strconv"

type Atom int16

const (
	atoms_begin Atom = iota - 1

	Nil

	Air
	Dirt

	atoms_end

	invalid
)

const Number = atoms_end - atoms_begin - 1

func (a Atom) Valid() bool {
	return atoms_begin < a && a < atoms_end
}

func (a Atom) String() string {
	if a.Valid() {
		s := atoms[a]
		if s != "" {
			return s
		}
	}

	return "atom(" + strconv.Itoa(int(a)) + ")"
}

func FromString(s string) Atom {
	if len(s) >= 2 && s[0] == '~' && s[len(s)-1] == '~' {
		s = s[1 : len(s)-1]
	}

	for i, atom := range atoms {
		if atom == s {
			return Atom(i)
		}
	}

	return invalid
}

var atoms = [...]string{
	Nil:  "nil",
	Air:  "air",
	Dirt: "dirt",
}

var tilde_atoms [len(atoms)]string

func init() {
	for i, atom := range atoms {
		if atom != "" {
			tilde_atoms[i] = "~" + atom + "~"
		}
	}
}
