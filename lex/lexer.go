package lex

import (
	"errors"
	"fmt"
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
	case runeNUMBER:
		l.lexNumber()
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
