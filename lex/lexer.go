package lex

import (
	"fmt"
	"unicode/utf8"
)

type lexer struct {
	in      string
	out     []Token
	at      int
	current rune
	width   int
}

func Lex(s string) []Token {
	l := lexer{s, []Token{}, 0, utf8.RuneError, 0}
	l.lex()
	return l.out
}

func (l *lexer) lex() {
	var err error

	l.current, l.width, err = l.next()
	if err == errEOF {
		return nil
	}

	switch l.currentKind() {
	case runeNUMBER:
		l.lexNumber()
	case runeSYMBOL:
		l.lexSymbol()
	default:
		fmt.Errorf("could not lex unknown rune: %q", l.current)
	}
}

/* -------------------------- */

func (l *lexer) next() rune, int, error{
	if len(l.in) <= l.at{
		return utf8.RuneError, 0, errEOF
	}
}
