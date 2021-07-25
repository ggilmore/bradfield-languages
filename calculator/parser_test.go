package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParser(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    []Token
		expected Node
	}{
		{
			name: "(2 + 3)",
			input: []Token{
				{
					Kind:   KindLeftParen,
					Lexeme: "(",
				},
				{
					Kind:    KindNumber,
					Lexeme:  "2",
					Literal: 2,
				},
				{
					Kind:   KindPlus,
					Lexeme: "+",
				},
				{
					Kind:    KindNumber,
					Lexeme:  "3",
					Literal: 3,
				},
				{
					Kind:   KindRightParen,
					Lexeme: ")",
				},
			},
			expected: BinOp{
				Op:    Plus,
				Left:  Literal{2},
				Right: Literal{3},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.input)

			actual, err := p.Parse()
			if err != nil {
				t.Errorf("unexpected error while parsing: %s", err)
			}

			assertParseTreeEqual(t, actual, tt.expected)
		})
	}
}

func assertParseTreeEqual(t *testing.T, expected, actual Node) {
	t.Helper()

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("non-zero diff (-expected +actual):\n%s", diff)
	}
}
