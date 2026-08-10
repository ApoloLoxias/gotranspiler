package lex

import (
	"errors"
	"fmt"
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
	l.lex()
	return l.out
}

func (l *lexer) lex() {
	var err error

	err = l.next()
	if err == errEOF {
		return
	}

	switch l.currentKind() {
	case runeDIGIT:
		l.lexDECIMAL()
	case runeSYMBOL:
		l.lexSymbol()
	default:
		fmt.Errorf("could not lex unknown rune: %q", l.current)
	}
}

/* -------------------------- */

func (l *lexer) next() error {
	if l.length <= l.at {
		l.current = utf8.RuneError
		return errEOF
	}

	l.at = l.at + l.width

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

func (l *lexer) produceToken(width int, kind TokenKind) {
	token := Token{Value: l.in[l.at : l.at+width], Kind: TokenNUMBER}
	l.out = append(l.out, token)
}

/* -------------------------- */

type lexingFunction func(*lexer) lexingFunction

func lexDecimal(l *lexer) lexingFunction {
	i := -1
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
