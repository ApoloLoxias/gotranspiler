package lex

import (
	//	"errors"
	//	"strings"
	"testing"
	// "unicode/utf8"
)

func extractTokenKinds(tokens []Token) []TokenKind {
	result := make([]TokenKind, 0, len(tokens))
	for _, token := range tokens {
		result = append(result, token.Kind)
	}
	return result
}

func assertTokenKinds(t *testing.T, source string, want []TokenKind) {
	t.Helper()
	got := extractTokenKinds(Lex(source))

	if len(got) != len(want) {
		t.Errorf("Lex(%v) gets %d tokens: `%v`, want %d tokens: `%v`", source, len(got), got, len(want), want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Errorf(
				"Lex(%v) produces %s at index %d, want %s",
				source,
				got[i],
				i,
				want[i],
			)
		}
	}

}

func TestKindsForLexTokens(t *testing.T) {
	type test struct {
		source string
		want   []TokenKind
	}
	type block struct {
		name  string
		tests []test
	}

	blocks := []block{
		{
			"tokens",
			[]test{
				{"1", []TokenKind{TokenNUMBER}},
				{"+", []TokenKind{TokenCROSS}},
				{"-", []TokenKind{TokenHYPHEN}},
				{"*", []TokenKind{TokenASTERISK}},
				{"/", []TokenKind{TokenFORWARD_SLASH}},
				{"!", []TokenKind{TokenEXCLAMATION}},
				{"**", []TokenKind{TokenASTERISK_ASTERISK}},
				{"()", []TokenKind{
					TokenOPEN_PARENTHESIS,
					TokenCLOSE_PARENTHESIS,
				}},
			},
		},
		{
			"infixes",
			[]test{
				{"0+1", []TokenKind{
					TokenNUMBER,
					TokenCROSS,
					TokenNUMBER,
				}},
				{"01+234-5/6*7**8", []TokenKind{
					TokenNUMBER,
					TokenCROSS,
					TokenNUMBER,
					TokenHYPHEN,
					TokenNUMBER,
					TokenFORWARD_SLASH,
					TokenNUMBER,
					TokenASTERISK,
					TokenNUMBER,
					TokenASTERISK_ASTERISK,
					TokenNUMBER,
				}},
			},
		},
	}

	for _, block := range blocks {
		for _, test := range block.tests {
			t.Run(
				test.source,
				func(t *testing.T) {
					assertTokenKinds(t, test.source, test.want)
				},
			)
		}
	}
}
