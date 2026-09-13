package lex

import "fmt"

type Token struct {
	Value string
	Kind  TokenKind
}

type TokenKind string // Enum

const (
	TokenEOF TokenKind = "EndOfFileToken"
	TokenSOF TokenKind = "StarfOfFileToken"

	TokenNUMBER TokenKind = "NumericalToken"

	TokenCROSS         TokenKind = "PlusToken"
	TokenHYPHEN        TokenKind = "HyphenToken"
	TokenASTERISK      TokenKind = "AsteriskToken"
	TokenFORWARD_SLASH TokenKind = "ForwardSlashToken"

	TokenCROSS_FUNC         TokenKind = "PlusFuncToken"
	TokenHYPHEN_FUNC        TokenKind = "SubFuncToken"
	TokenASTERISK_FUNC      TokenKind = "MulFuncToken"
	TokenFORWARD_SLASH_FUNC TokenKind = "DivFuncToken"

	TokenOPEN_PARENTHESIS  TokenKind = "OpenParenthesis"
	TokenCLOSE_PARENTHESIS TokenKind = "CloseParenthesis"
)

var EOFtoken = Token{"EOF", TokenEOF} //const
var SOFtoken = Token{"SOF", TokenSOF} //const

var TerminalTokens = []TokenKind{ //const
	TokenNUMBER,
	TokenCROSS_FUNC,
	TokenHYPHEN_FUNC,
	TokenASTERISK_FUNC,
	TokenFORWARD_SLASH_FUNC,
}

var FunctionTokens = []TokenKind{ //const
	TokenCROSS_FUNC,
	TokenHYPHEN_FUNC,
	TokenASTERISK_FUNC,
	TokenFORWARD_SLASH_FUNC,
}

var InfixTokens = []TokenKind{ //const
	TokenCROSS,
	TokenHYPHEN,
	TokenASTERISK,
	TokenFORWARD_SLASH,
}

var InfixPriority = map[TokenKind]int{
	TokenCROSS:         1,
	TokenHYPHEN:        1,
	TokenASTERISK:      2,
	TokenFORWARD_SLASH: 2,
}

var PrefixTokens = []TokenKind{
	TokenHYPHEN,
	// TokenCROSS_FUNC,
}

var PrefixPriority = map[TokenKind]int{
	TokenHYPHEN: 3,
	// TokenCROSS_FUNC: 3,
}

func (t Token) String() string {
	return fmt.Sprintf("%s('%s')", t.Kind, t.Value)
}

func (t Token) IsOfKind(kinds ...TokenKind) bool {
	for _, kind := range kinds {
		if t.Kind == kind {
			return true
		}
	}
	return false
}

func (t Token) IsTerminal() bool {
	return t.IsOfKind(TerminalTokens...)
}

func (t Token) IsFunction() bool {
	return t.IsOfKind(TerminalTokens...)
}

func (t Token) IsInfix() bool {
	return t.IsOfKind(InfixTokens...)
}

func (t Token) IsPrefix() bool {
	return t.IsOfKind(PrefixTokens...)
}
