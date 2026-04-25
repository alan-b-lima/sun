package lexer_test

import (
	"testing"

	. "github.com/alan-b-lima/sun/internal/spell/lexer"
	"github.com/alan-b-lima/sun/internal/spell/scanner"
)

func TestLexer(t *testing.T) {
	type TokenEx struct {
		Token   Token
		Content string
	}

	type Test struct {
		Name string
		In   string
		Want []TokenEx
	}

	tests := []Test{
		{
			Name: "many semicolons",
			In:   `;;;;`,
			Want: []TokenEx{
				{Token: Semicolon},
				{Token: EOF},
			},
		},
		{
			Name: "semicolon insertion in empty block",
			In:   `{}`,
			Want: []TokenEx{
				{Token: LBrace},
				{Token: RBrace},
				{Token: Semicolon},
				{Token: EOF},
			},
		},
		{
			Name: "atoms",
			In:   `~air~ ~dirt~ ~water~`,
			Want: []TokenEx{
				{Token: Atom, Content: "air"},
				{Token: Atom, Content: "dirt"},
				{Token: Atom, Content: "water"},
				{Token: Semicolon},
				{Token: EOF},
			},
		},
		{
			Name: "spell declaration",
			In:   `spell name { ~air~ 0[~dirt~] release[0] face name }`,
			Want: []TokenEx{
				{Token: Spell},
				{Token: Identifier, Content: "name"},
				{Token: LBrace},
				{Token: Atom, Content: "air"},
				{Token: Identifier, Content: "0"},
				{Token: LBrack},
				{Token: Atom, Content: "dirt"},
				{Token: RBrack},
				{Token: Release},
				{Token: LBrack},
				{Token: Identifier, Content: "0"},
				{Token: RBrack},
				{Token: Face},
				{Token: Identifier, Content: "name"},
				{Token: Semicolon},
				{Token: RBrace},
				{Token: Semicolon},
				{Token: EOF},
			},
		},
	}

	for _, test := range tests {
		source := scanner.NewSource([]byte(test.In))
		stream := NewStream(source)

		for j, want := range test.Want {
			stream.Next()
			got := stream.Token()

			if want.Token != got {
				t.Errorf(
					"%s @ token %d\nexpected %v\ngot %v\n",
					test.Name, j, want.Token, got,
				)
				continue
			}

			if want.Content != "" {
				if content := stream.Content(); want.Content != content {
					t.Errorf(
						"%s @ token %d\nexpected\n%v(%s)\ngot %v(%s)\n",
						test.Name, j, want.Token, want.Content, got, content,
					)
				}
			}
		}
	}
}
