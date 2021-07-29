package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestScanner(t *testing.T) {
	for _, tt := range []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "given",
			input: "1 - (2 + 3)",
			expected: []Token{
				{Kind: KindNumber, Lexeme: "1", Literal: 1},
				{Kind: KindMinus, Lexeme: "-"},
				{Kind: KindLeftParen, Lexeme: "("},
				{Kind: KindNumber, Lexeme: "2", Literal: 2},
				{Kind: KindPlus, Lexeme: "+"},
				{Kind: KindNumber, Lexeme: "3", Literal: 3},
				{Kind: KindRightParen, Lexeme: ")"},
				{Kind: KindEOF},
			},
		},
		{
			name:  "tricky multiplication",
			input: "5*-20",
			expected: []Token{
				{Kind: KindNumber, Lexeme: "5", Literal: 5},
				{Kind: KindStar, Lexeme: "*"},
				{Kind: KindMinus, Lexeme: "-"},
				{Kind: KindNumber, Lexeme: "20", Literal: 20},
				{Kind: KindEOF},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewScanner(strings.NewReader(tt.input))
			if err != nil {
				t.Fatalf("failed to initialize scanner: %s", err)
			}

			actual, err := s.ScanTokens()
			if err != nil {
				t.Errorf("while scanning input: %s", err)
			}

			assertTokensEqual(t, tt.expected, actual)
		})
	}
}

func assertTokensEqual(t *testing.T, expected, actual []Token) {
	t.Helper()

	if diff := cmp.Diff(expected, actual, cmpopts.IgnoreFields(Token{}, "PosStart", "PosEnd")); diff != "" {
		t.Errorf("non-zero diff (-expected +actual):\n%s", diff)
	}
}
