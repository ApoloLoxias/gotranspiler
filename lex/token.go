package lex

import "fmt"

/* ---- TOKEN DEFINITION AND KIND ENUMERATION --- */

type Token struct {
	Value string
	Kind  TokenKind
}

func (t Token) String() string {
	return fmt.Sprintf("%s('%s')", t.Kind, t.Value)
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

	TokenEXCLAMATION TokenKind = "ExclamationPointToken"

	TokenASTERISK_ASTERISK TokenKind = "DoubleAsteriskToken"

	TokenCROSS_FUNC         TokenKind = "PlusFuncToken"
	TokenHYPHEN_FUNC        TokenKind = "SubFuncToken"
	TokenASTERISK_FUNC      TokenKind = "MulFuncToken"
	TokenFORWARD_SLASH_FUNC TokenKind = "DivFuncToken"

	TokenOPEN_PARENTHESIS  TokenKind = "OpenParenthesis"
	TokenCLOSE_PARENTHESIS TokenKind = "CloseParenthesis"

	TokenUNDEFINED_SYMBOL TokenKind = "UndefinedSymbolToken"
)

var EOFtoken = Token{"EOF", TokenEOF} //const
var SOFtoken = Token{"SOF", TokenSOF} //const

/* --- TOKEN CLASSIFICATION --- */

// A single terminal is a valid expression
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

var ParenthesisTokens = []TokenKind{ //const
	TokenOPEN_PARENTHESIS,
	TokenCLOSE_PARENTHESIS,
}

var InfixTokens = []TokenKind{ //const
	TokenCROSS,
	TokenHYPHEN,
	TokenASTERISK,
	TokenFORWARD_SLASH,
	TokenASTERISK_ASTERISK,
}

var InfixPriority = map[TokenKind]int{ //const
	TokenCROSS:         1,
	TokenHYPHEN:        1,
	TokenASTERISK:      2,
	TokenFORWARD_SLASH: 2,

	TokenASTERISK_ASTERISK: 4,
}

// const ApplicationPriority int = 4 // when reintroducing application, consider asterisk_asterisk: 4 and application: 5

var PrefixTokens = []TokenKind{ //const
	TokenHYPHEN,
}

var PrefixPriority = map[TokenKind]int{ //const
	TokenHYPHEN: 3,
}

var SuffixTokens = []TokenKind{ //const
	TokenEXCLAMATION,
}

var SuffixPriority = map[TokenKind]int{ //const
	TokenEXCLAMATION: 5, //consider 6 if application gets bumped to 5
}

// Firsts can start an expression
func defineFirstTokens() []TokenKind {
	var FirstTokens2 = make([]TokenKind, 0, len(PrefixTokens)+len(TerminalTokens))
	var FirstTokens1 = append(FirstTokens2, PrefixTokens...)
	var FirstTokens0 = append(FirstTokens1, TerminalTokens...)
	return FirstTokens0
}

var FirstTokens []TokenKind = defineFirstTokens() //const

// Classification checkers
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
	return t.IsOfKind(FunctionTokens...)
}

func (t Token) IsParenthesis() bool {
	return t.IsOfKind(ParenthesisTokens...)
}

func (t Token) IsInfix() bool {
	return t.IsOfKind(InfixTokens...)
}

func (t Token) IsPrefix() bool {
	return t.IsOfKind(PrefixTokens...)
}

func (t Token) IsSuffix() bool {
	return t.IsOfKind(SuffixTokens...)
}

func (t Token) IsFirst() bool {
	return t.IsOfKind(FirstTokens...)
}

/* --- ASSOCIATIVITY --- */

type Associativity string

const (
	LeftAssociativity  Associativity = "left"
	RightAssociativity Associativity = "right"
)

var rightAssociatives = map[TokenKind]Associativity{ //left associativity by default
	TokenASTERISK_ASTERISK: RightAssociativity,
}

func AssociativityOf(kind TokenKind) Associativity {
	_, ok := rightAssociatives[kind]
	if ok {
		return RightAssociativity
	}
	return LeftAssociativity
}
