package lex

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

type lexer struct {
	in      string
	out     []Token
	at      int
	current rune
	width   int
	length  int
}

func Lex(s string) []Token {
	l := lexer{s, []Token{}, 0, utf8.RuneError, 0, len(s)}
	err := l.next()

	if err == errEOF {
		return []Token{}
	}

	l.lex()
	return l.out
}

func (l *lexer) lex() {
	lexing := lexDecimal(l)
	for lexing != nil {
		lexing = lexing(l)
	}
}

/* -------------------------- */

func (l *lexer) next() error {
	l.at = l.at + l.width

	if l.length <= l.at {
		l.current = utf8.RuneError
		return errEOF
	}

	l.current, l.width = utf8.DecodeRuneInString(l.in[l.at:])

	return nil
}

var errEOF = errors.New("end of file")

func (l *lexer) currentKind() runeKind {
	r := l.current

	if unicode.IsDigit(r) {
		return runeDIGIT
	}
	if isOperator(r) {
		return runeSYMBOL
	}

	return runeUNKNOWN
}

type runeKind string

const (
	runeUNKNOWN runeKind = "unkown rune"

	runeDIGIT  runeKind = "numeric rune"                      // 0123456789
	runeSYMBOL runeKind = "Arithmetic operation synmbol rune" // +-/*
)

func isOperator(r rune) bool {
	for _, Rune := range operatorCharacters {
		if r == Rune {
			return true
		}
	}
	return false
}

var operatorCharacters = []rune("+-*/") // wish it were const

var (
	runeCROSS         = rune("+"[0])
	runeHYPHEN        = rune("-"[0])
	runeASTERISK      = rune("*"[0])
	runeFORWARD_SLASH = rune("/"[0])
)

func (l *lexer) produceToken(width int, kind TokenKind) {
	token := Token{Value: l.in[l.at : l.at+width], Kind: kind}
	l.out = append(l.out, token)
}

/* -------------------------- */

type lexingFunction func(*lexer) lexingFunction

func lexDecimal(l *lexer) lexingFunction {
	i := 0
	var err error = nil

	for l.currentKind() == runeDIGIT {
		i += l.width
		err = l.next()

		if err == errEOF {
			l.produceToken(i, TokenNUMBER)
			return nil
		}
	}

	l.produceToken(i, TokenNUMBER)
	return lexSymbol
}

func lexSymbol(l *lexer) lexingFunction {
	switch l.current {
	case rune(runeCROSS):
		l.produceToken(l.width, TokenCROSS)
	case rune(runeHYPHEN):
		l.produceToken(l.width, TokenHYPHEN)
	case rune(runeASTERISK):
		l.produceToken(l.width, TokenASTERISK)
	case rune(runeFORWARD_SLASH):
		l.produceToken(l.width, TokenFORWARD_SLASH)
	}

	err := l.next()

	if err == errEOF {
		return nil
	}
	return lexDecimal
}
