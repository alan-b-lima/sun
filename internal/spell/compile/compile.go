package compile

import (
	"github.com/alan-b-lima/sun/internal/spell"
	"github.com/alan-b-lima/sun/internal/spell/generation"
	"github.com/alan-b-lima/sun/internal/spell/lexer"
	"github.com/alan-b-lima/sun/internal/spell/parser"
	"github.com/alan-b-lima/sun/internal/spell/scanner"
	"github.com/alan-b-lima/sun/internal/spell/semantics"
)

func Compile(source []byte) (spell.Spell, error) {
	scanner := scanner.New(source)
	stream := lexer.New(scanner)

	tree, err := parser.Parse(stream)
	if err != nil {
		return spell.Spell{}, err
	}

	outcome, err := semantics.Analyze(tree)
	if err != nil {
		return spell.Spell{}, err
	}

	spell := generation.Generate(outcome)
	
	return spell, nil
}
