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

	return nil, nil //TODO
}
