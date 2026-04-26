package lexer

import (
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Token int

const (
	Illegal Token = iota
	EOF           // EOF
	Comment       // // ...

	literal_begin

	Identifier // Identifier
	Atom       // Atom

	literal_end

	keyword_begin

	Absorb  // absorb
	Back    // back
	Cellar  // cellar
	Down    // down
	Face    // face
	Group   // group
	Left    // left
	Nil     // nil
	Release // release
	Right   // right
	Rune    // rune
	Spell   // spell
	State   // state
	Up      // up
	Write   // write

	keyword_end

	punctuation_begin

	LParen // (
	LBrack // [
	LBrace // {

	RParen // )
	RBrack // ]
	RBrace // }

	Semicolon // ;
	Comma     // ,
	Dot       // .

	punctuation_end
)

func (t Token) String() string {
	if 0 <= t && t < Token(len(tokens)) {
		str := tokens[t]
		if str != "" {
			return str
		}
	}

	return "token(" + strconv.Itoa(int(t)) + ")"
}

func (t Token) IsLiteral() bool {
	return literal_begin < t && t < literal_end
}

func (t Token) IsPunctuation() bool {
	return punctuation_begin < t && t < punctuation_end
}

func (t Token) IsKeyword() bool {
	return keyword_begin < t && t < keyword_end
}

func is_identifier_char(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_'
}

var tokens = [...]string{
	Illegal: "%illegal%",

	EOF:     "EOF",
	Comment: "Comment",

	Identifier: "Identifier",
	Atom:       "Atom",

	Absorb:  "absorb",
	Back:    "back",
	Cellar:  "cellar",
	Down:    "down",
	Face:    "face",
	Group:   "group",
	Left:    "left",
	Nil:     "nil",
	Release: "release",
	Right:   "right",
	Rune:    "rune",
	Spell:   "spell",
	State:   "state",
	Up:      "up",
	Write:   "write",

	LParen: "(",
	LBrack: "[",
	LBrace: "{",

	RParen: ")",
	RBrack: "]",
	RBrace: "}",

	Semicolon: ";",
	Comma:     ",",
	Dot:       ".",
}

var (
	punctuation map[rune]Token
	keywords    map[string]Token
)

func init() {
	punctuation = make(map[rune]Token, punctuation_end-punctuation_begin-1)
	for i := punctuation_begin + 1; i < punctuation_end; i++ {
		rune, _ := utf8.DecodeRuneInString(tokens[i])
		punctuation[rune] = i
	}

	keywords = make(map[string]Token, keyword_end-keyword_begin-1)
	for i := keyword_begin + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}
