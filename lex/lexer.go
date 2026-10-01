package lex

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

/* --- ENTRY POINT --- */

func Lex(s string) []Token {
	l := lexer{s, []Token{}, 0, utf8.RuneError, 0, len(s)}
	err := l.next()

	if err == errEOF {
		return []Token{}
	}

	l.lex()
	return l.out
}

/* --- LEXER STATE AND CONTROL FLOW--- */

type lexer struct {
	//input and output
	in  string
	out []Token //in the future, channel (lex/parse simulteneously)

	//cursor
	at      int
	current rune

	//cached values for multiple access
	width  int
	length int
}

func (l *lexer) lex() {
	lexing := lexNumOrParen(l)
	for lexing != nil {
		lexing = lexing(l)
		for l.currentKind() == runeWHITESPACE {
			if l.next() == errEOF {
				return
			}
		}
	}
}

var errEOF = errors.New("end of file")

func (l *lexer) next() error {
	l.at = l.at + l.width

	if l.length <= l.at {
		l.current = utf8.RuneError
		return errEOF
	}

	nextRune, nextWidth := utf8.DecodeRuneInString(l.in[l.at:])
	l.current, l.width = nextRune, nextWidth

	return nil
}

func (l *lexer) produceToken(start int, width int, kind TokenKind) {
	token := Token{Value: l.in[start : start+width], Kind: kind}
	l.out = append(l.out, token)
}

/* --- RUNE CLASSIFICATION --- */

type runeKind string

const (
	runeUNKNOWN runeKind = "unkown rune"

	runeDIGIT  runeKind = "numeric rune"                      // 0123456789
	runeSYMBOL runeKind = "Arithmetic operation synmbol rune" // +-/*$!

	runePARENTHESIS runeKind = "parenthesis rune"

	runeWHITESPACE runeKind = "whitespace rune"
)

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

	if unicode.IsSpace(r) {
		return runeWHITESPACE
	}

	return runeUNKNOWN
}
func isOperator(r rune) bool {
	for _, Rune := range operatorCharacters {
		if r == Rune {
			return true
		}
	}
	return false
}

var operatorCharacters = []rune("+-*/$!") // wish it were const

var (
	runeCROSS         = rune("+"[0])
	runeHYPHEN        = rune("-"[0])
	runeASTERISK      = rune("*"[0])
	runeFORWARD_SLASH = rune("/"[0])

	runeEXCLAMATION = rune("!"[0])

	runeDOLLAR = rune("$"[0])

	runeOPEN_PARENTHESIS  = rune("("[0])
	runeCLOSE_PARENTHESIS = rune(")"[0])
)

/* --- LEXING FUNCTIONS --- */

type lexingFunction func(*lexer) lexingFunction

func lexDecimal(l *lexer) lexingFunction {
	startAt := l.at
	totalWidth := 0

	for l.currentKind() == runeDIGIT {
		totalWidth += l.width

		if l.next() == errEOF {
			l.produceToken(startAt, totalWidth, TokenNUMBER)
			return nil
		}
	}

	if totalWidth != 0 {
		l.produceToken(startAt, totalWidth, TokenNUMBER)
	}
	return lexSymbol
}

// Currently matches all consecutive symbols into a single rune.
// May want to change after definition of syntax
func lexSymbol(l *lexer) lexingFunction {
	var chars []rune
	var err error
	for l.currentKind() == runeSYMBOL {
		chars = append(chars, l.current)
		err = l.next()
		if err == errEOF {
			break
		}
	}

	if len(chars) == 0 {
		return lexNumOrParen
	}

	var token Token
	token.Value = string(chars)

	switch token.Value {
	case "+":
		token.Kind = TokenCROSS
	case "-":
		token.Kind = TokenHYPHEN
	case "*":
		token.Kind = TokenASTERISK
	case "/":
		token.Kind = TokenFORWARD_SLASH
	case "**":
		token.Kind = TokenASTERISK_ASTERISK

	case "!":
		token.Kind = TokenEXCLAMATION

	case "+$":
		token.Kind = TokenCROSS_FUNC
	case "-$":
		token.Kind = TokenHYPHEN_FUNC
	case "*$":
		token.Kind = TokenASTERISK_FUNC
	case "/$":
		token.Kind = TokenFORWARD_SLASH_FUNC

	default:
		token.Kind = TokenUNDEFINED_SYMBOL
	}

	l.out = append(l.out, token)

	var nextFunc lexingFunction
	if err != errEOF {
		nextFunc = lexNumOrParen
	}

	return nextFunc
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

	return lexDecimal //i.e. l.currentKind() != PARENTHESIS
}
