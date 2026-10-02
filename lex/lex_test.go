package lex

import (
	//	"errors"
	"strings"
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

func extractTokenValues(tokens []Token) []string {
	result := make([]string, 0, len(tokens))
	for _, token := range tokens {
		result = append(result, token.Value)
	}
	return result
}

func assertTokenValues(t *testing.T, source string, want []string) {
	t.Helper()
	got := extractTokenValues(Lex(source))

	if len(got) != len(want) {
		t.Errorf("Lex(%s) gets %d tokens: `%v`, want %d tijebs: `%v`", source, len(got), got, len(want), want)
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
		{
			"prefix",
			[]test{
				{
					"-1",
					[]TokenKind{
						TokenHYPHEN,
						TokenNUMBER,
					},
				},
				{
					"- - 2",
					[]TokenKind{
						TokenHYPHEN,
						TokenHYPHEN,
						TokenNUMBER,
					},
				},
				{
					"- 1 + - 7 - 3 * - 4",
					[]TokenKind{
						TokenHYPHEN,
						TokenNUMBER,
						TokenCROSS,
						TokenHYPHEN,
						TokenNUMBER,
						TokenHYPHEN,
						TokenNUMBER,
						TokenASTERISK,
						TokenHYPHEN,
						TokenNUMBER,
					},
				},
			},
		},
		{
			"postfixes",
			[]test{
				{"1!", []TokenKind{TokenNUMBER, TokenEXCLAMATION}},
				{"0! !", []TokenKind{TokenNUMBER, TokenEXCLAMATION, TokenEXCLAMATION}},
				{
					"- 1 ! + 4 ! !",
					[]TokenKind{
						TokenHYPHEN,
						TokenNUMBER,
						TokenEXCLAMATION,
						TokenCROSS,
						TokenNUMBER,
						TokenEXCLAMATION,
						TokenEXCLAMATION,
					},
				},
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

func TestValuesForIntegers(t *testing.T) {
	tests := []struct {
		source string
		values []string
	}{
		{"1", []string{"1"}},
		{"3456789876542456789098765", []string{"3456789876542456789098765"}},
	}

	for _, test := range tests {
		testValues := func(t *testing.T) {
			assertTokenValues(t, test.source, test.values)
		}
		t.Run(test.source, testValues)
	}
}

func TestWhitespace(t *testing.T) {
	source := "\t\n\v\r\f 0\t\n\v\r\f 1\t\n\v\r\f "
	want := []TokenKind{TokenNUMBER, TokenNUMBER}
	assertTokenKinds(t, source, want)

}

func TestTokenValuesReconstructSourceWithNoWhitespace(t *testing.T) {
	tests := []struct {
		name  string
		chars string
	}{
		{"digits", "0 1 2 3 4 5 6 7 8 9"},
		{"operators", "+ - / * ! **"},
		{"parenthesis", "() (())"},
		{"random implemented stuff", "2345 **+-/  1 2 + - / / -12++++$"},
	}

	for _, test := range tests {
		testValues := extractTokenValues(Lex(test.chars))
		builder := strings.Builder{}
		for _, x := range testValues {
			builder.WriteString(x)
		}
		reconstructedString := builder.String()
		testNoSpaceString := strings.ReplaceAll(test.chars, " ", "")
		if reconstructedString != testNoSpaceString {
			t.Errorf("Trying to reconstitute %s from the values of its tokens, gets us %s instead of %s", test.chars, reconstructedString, testNoSpaceString)
		}
	}
}
