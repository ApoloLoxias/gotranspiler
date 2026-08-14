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
	lexing := lexNumOrParen(l)
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

	if r == runeOPEN_PARENTHESIS || r == runeCLOSE_PARENTHESIS {
		return runePARENTHESIS
	}

	return runeUNKNOWN
}

type runeKind string

const (
	runeUNKNOWN runeKind = "unkown rune"

	runeDIGIT  runeKind = "numeric rune"                      // 0123456789
	runeSYMBOL runeKind = "Arithmetic operation synmbol rune" // +-/*

	runePARENTHESIS runeKind = "parenthesis rune"
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

	runeOPEN_PARENTHESIS  = rune("("[0])
	runeCLOSE_PARENTHESIS = rune(")"[0])
)

func (l *lexer) produceToken(start int, width int, kind TokenKind) {
	token := Token{Value: l.in[start : start+width], Kind: kind}
	l.out = append(l.out, token)
}

/* -------------------------- */

type lexingFunction func(*lexer) lexingFunction

func lexDecimal(l *lexer) lexingFunction {
	startAt := l.at
	totalWidth := 0
	var err error = nil

	for l.currentKind() == runeDIGIT {
		totalWidth += l.width
		err = l.next()

		if err == errEOF {
			l.produceToken(startAt, totalWidth, TokenNUMBER)
			return nil
		}
	}

	if totalWidth != 0 {
		l.produceToken(startAt, totalWidth, TokenNUMBER)
	}
	return lexSymbol
}

func lexSymbol(l *lexer) lexingFunction {
	switch l.current {
	case rune(runeCROSS):
		l.produceToken(l.at, l.width, TokenCROSS)
	case rune(runeHYPHEN):
		l.produceToken(l.at, l.width, TokenHYPHEN)
	case rune(runeASTERISK):
		l.produceToken(l.at, l.width, TokenASTERISK)
	case rune(runeFORWARD_SLASH):
		l.produceToken(l.at, l.width, TokenFORWARD_SLASH)
	case runeCLOSE_PARENTHESIS:
		l.produceToken(l.at, l.width, TokenCLOSE_PARENTHESIS)
	case runeOPEN_PARENTHESIS: // This one should probably never happen on valid code
		l.produceToken(l.at, l.width, TokenOPEN_PARENTHESIS)
	}

	err := l.next()

	if err == errEOF {
		return nil
	}
	return lexNumOrParen
}

func lexNumOrParen(l *lexer) lexingFunction {
	if l.current == runeOPEN_PARENTHESIS {
		l.produceToken(l.at, l.width, TokenOPEN_PARENTHESIS)
		err := l.next()
		if err == errEOF {
			return nil
		}
		return lexNumOrParen
	}
	if l.current == runeCLOSE_PARENTHESIS {
		l.produceToken(l.at, l.width, TokenCLOSE_PARENTHESIS)
		err := l.next()
		if err == errEOF {
			return nil
		}
		return lexNumOrParen
	}

	return lexDecimal
}
