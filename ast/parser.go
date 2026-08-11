package ast

import (
	"errors"
	"strconv"

	"github.com/ApoloLoxias/gotranspiler/lex"
)

func Parse(tokens []lex.Token) Expression {
	parser := parser{tokens, nil, 0}
	parser.parse()
	return parser.out
}

/* ------------------------- */

type parser struct {
	in  []lex.Token
	out Expression
	at  int
}

func (p *parser) parse() {
	parsing := parseNumber
	for parsing != nil {
		p.out, parsing = parsing(p)
	}
	return
}

func (p *parser) next() error {
	p.at++
	if p.at >= len(p.in) {
		return errEOF
	}

	return nil
}

func (p *parser) peekAhead() lex.Token {
	if p.at >= len(p.in) {
		return lex.EOFtoken
	}
	return p.in[p.at+1]
}

func (p *parser) peekBehind() lex.Token {
	if p.at <= 0 {
		return lex.SOFtoken
	}
	return p.in[p.at-1]
}

var errEOF = errors.New("EOF")

/* ------------------------- */

type parsingFunction func(*parser) (Expression, parsingFunction)

func parseNumber(p *parser) (Expression, parsingFunction) {
	value, _ := strconv.Atoi(p.in[p.at].Value)
	expr := IntE{Value: value}

	err := p.next()
	if err == errEOF {
		return expr, nil
	}

	if p.in[p.at].IsInfix() {
		return expr, parseInfixOperation
	}

	return nil, nil //TODO
}

func parseInfixOperation(p *parser) (Expression, parsingFunction) {
	arg1 := p.out

	var op BuiltInFunc
	switch p.in[p.at].Kind {
	case lex.TokenCROSS:
		op = Sum
	case lex.TokenHYPHEN:
		op = Sub
	case lex.TokenASTERISK:
		op = Mul
	case lex.TokenFORWARD_SLASH:
		op = Div
	}

	expr := ApplicationE{
		Function: op,
		Argument: arg1,
	}

	err := p.next()
	if err == errEOF {
		return expr, nil
	}
	return expr, parseInfixFollowUp
}

func parseInfixFollowUp(p *parser) (Expression, parsingFunction) {
	arg2, pFunc := parseNumber(p)
	expr := ApplicationE{
		Function: p.out,
		Argument: arg2,
	}

	return expr, pFunc
}
