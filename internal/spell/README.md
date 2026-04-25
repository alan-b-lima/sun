# Sun Spell Specification

Sun Spell Grammar is a grammar that describes a multi-stack pushdown automaton.

## Notation

We use a variant of WSN, which is itself a variant of EBNF, bootstraptly defined by itself as:

```
Syntax     = { Production } .
Production = identifier "=" Expression "." .
Expression = Term { "|" Term } .
Term       = Factor { Factor } .
Factor     = identifier | literal [ "..." literal ] | Group | Optional | Repetition .
Group      = "(" Expression ")" .
Optional   = "[" Expression "]" .
Repetition = "{" Expression "}" .

identifier = letter { letter | "-" } .
literal    = "`" { code-point } "`" | `"` { code-point } `"` .
```

Productions are expressions constructed from terms and the following operators, in increasing precedence:

```
|   alternation
()  grouping
[]  optional (0 or 1 times)
{}  repetition (0 to n times)
```

There also the non standard explanation definition, the identifiers `letter` and `code-point` weren't defined, they might be as such:

```
letter     = /* an Unicode code point categorized as "Letter" */
code-point = /* an Unicode code point */
```

Technically `code-point` has to exclude its quotes, be them `"` (U+0022) or <code>&#96;</code> (U+0060), escapes like `\n` for feed form (U+000A), `\r` for carriage return (U+000D), and `\t` for tab (U+0009) may be used as we see fit.

## Source representation

Spell sources are Unicode text encoded in UTF-8. No effort is given to canonize the text, i.e, accented code points are distinct from the same character constructed using combining accents and another code point.

The source is interpreted as case-sensitive. 

### Letters and digits

```
letter = /* an Unicode code point categorized as "Letter" */ .
digit  = /* an Unicode code point categorized as "Number, decimal digit" */ .
```

## Lexical elements

### Comments

Comments are a way of documenting the source. There is one form of comment, well known as line comment, starts with the character sequence `//`, two forward slashes `/` (U+002F), and stop at the end of the line.

Comments have no intrinsic meaning and are ignored.

### Semicolons

The formal grammar uses semicolons `;` (U+003B) as terminators for some constructs, however those need not be insert by the sourcerer, as the semicolons are inserted by the lexer:

1. at a line break following:
    - an identifier
    - an atom
    - a closing paren `)` (U+0028)
    - a closing brace `}` (U+007D)

1. before a closing brace `}` (U+007D)

A sequence of semicolon tokens are squashed into a single semicolon token.

### Tokens

Tokens belong to one of three classes: _keywords_, _identifiers_, _punctuation_, and _atoms_.

### Keywords

Keywords are reserved words and cannot be used as identifiers:

```
absorb     back     cellar    down
face       group    left      nil
release    right    rune      spell
state      up       write
```

### Identifiers

Identifiers are used to give name to objects, e.g, cellars, atom groups, and states. An identifer, excluding keywords, is defined as:

```
identifier = { letter | digit | "_" } .
```

### Punctuation

The following character are considered punctuation:

```
(    )    ;    :
[    ]    ,    .
{    }
```

### Atoms

Alike keywords, atoms have special meaning, however they do not overlap with identifiers, defined as:

```
atom = "~" identifier "~" .
```

Atoms are predefined by the interpreter, they may accept values like `~air~` or `~water~`, per se.

## Spell Source

A spell is build from a single source, called spell source, defined as below:

```
SpellSource = SpellClause ";" { Decl ";" } .
```

### Spell clause

A spell clause lives at the start of every single spell source and gives the name of the spell. 

```
SpellClause = "spell" SpellName .
SpellName   = identifier .
```

### Declarations

A declaration is one of three thing, a atom group declaration, a cellar declaration, state declaration.

```
Decl
    = AtomGroupDecl
    | CellarDecl
    | RuneDecl
    | StateDecl
    .
```

### Atom group declaration

An atom group declaration gives name to a set of atoms, which can be used in state transitions.

```
AtomGroupDecl = "group" AtomGroupName "{" { AtomSpec ";" } "}" .
AtomGroupName = identifier .
AtomSpec      = atom | AtomGroupName .
```

### Cellar declaration

A cellar declaration gives name to a cellar, which is a stack of atoms.

```
CellarDecl = "cellar" CellarName .
CellarName = identifier .
```

### Rune declaration

A rune declaration gives name to a rune, which is a special atom that can be used in state transitions, but cannot be absorbed and its release does not alter the world.

```
RuneDecl = "rune" RuneName .
RuneName  = identifier .
```

### State declaration

A state declaration gives name to a state, which defines a node in the automaton and all of its possible transitions.

```
StateDecl = "state" Signature "{" { Transition ";" } "}" .
Signature = identifier [ "(" identifier { "," identifier } ")" ]
```

A declared state may be a plain state or function state, function states have parameters, which are themselves states.

### State

A state is a node in the automaton, it has a name and zero or more parameters, which are themselves states.

```
State     = identifier [ "(" ParamList ")" ] .
ParamList = ParamSpec { "," ParamSpec } .
ParamSpec = State .
```

### Transitions

A transition in an edge in a graph connection the state of the spell automaton. They are defined as:

```
Transition = AtomCond CellarCondList Behavior MoveList [ FinalState ]
```

### Atom conditions

The atom condition matches the atom in the world in the very location the spell is over, the wildcard `.` (U+002E) may also be used to indicate any atom.

```
AtomCond = "." | AtomSpec .
```

### Cellar conditions

The cellar condition list matches the top of the cellars, it is a list of cellar conditions, which are defined as below, a wildcard `.` (U+002E) may also be used to indicate to match any state on the top of the cellars.

```
CellarCondList = "." | CellarCond { "," CellarCond } .
CellarCond     = CellarName "[" CellarCondSpec "]" .
CellarCondSpec = atom | AtomGroupName | RuneName .
```

### Behavior

A behavior is a action performed by the spell, it is defined as below, the keyword `nil` may also be used to indicate no behavior.

```
Behavior
    = "nil"
    | "absorb" "[" CellarName "]"
    | "release" "[" CellarName "]"
    | "write" "[" CellarName "," RuneName "]"
    .
```

The `absorb` behavior pushes the atom from the world into the cellar, the `release` behavior pops an atom from the cellar and release it into the world, and the `write` behavior writes a rune on the top of the cellar.

### Movement

Move list is a list of moves, which are defined as below, the keyword `nil` may also be used to indicate no move.

```
MoveList = "nil" | Move { "," Move } .
Move     = "up" | "down" | "left" | "right" | "face" | "back" .
```

The directions `up`, `down`, `left`, and `right` move the spell in the corresponding direction, while the direction `face` makes the spell move in the direction from which the spell was casted, horizontally, and `back` makes the spell move in the reverse direction of casting.

A spell can only move to an adjacent location, including diagonals, each direction can be thought as a vector (x, y), `nil` = (0, 0); `up` = (0, 1); `down` = (0, -1); `left` = (-1, 0); `right` = (1, 0); `face` = (0, f); `back` = (0, -f). If the sum of all moves yields a value outside of &PlusMinus;1 for any coordinate, the movement is considered illegal.

### Final state

The final state is the state to which the transition leads, if not specified, the spell will terminate.

```
FinalState = State .
```

A final state may depend on the parameters of the current state, for example:

```
state foo1(x) { . . . . bar(x) }
state foo2(x) { . . . . x }
```
