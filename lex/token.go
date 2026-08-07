package lex

import "fmt"

type Token struct {
	Value string
	Kind  TokenKind
}

type TokenKind string // Enum

const (
	NUMBER TokenKind = "Number"

	CROSS         TokenKind = "PlusToken"
	HYPHEN        TokenKind = "HyphenToken"
	ASTERISK      TokenKind = "AsteriskToken"
	FORWARD_SLASH TokenKind = "ForwardSlashToken"
)

func (t Token) String() string {
	return fmt.Sprintf("Token{%s, '%s'}", t.Kind, t.Value)
}

func (t Token) IsOfKind(kinds ...TokenKind) bool {
	for _, kind := range kinds {
		if t.Kind == kind {
			return true
		}
	}
	return false
}
