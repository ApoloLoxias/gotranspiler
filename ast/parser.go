package ast

import (
	"errors"
	"github.com/ApoloLoxias/gotranspiler/lex"
	"strconv"
)

func Parse(tokens []lex.Token) Expression {
	parser := parser{tokens, nil, 0}
	return parser.parse(0)
}

/* ------------------------- */

type parser struct {
	in  []lex.Token
	out Expression
	at  int
}

func (p *parser) parse(previousPriority int) Expression {
	/*if p.next() == errEOF {
		return nil
	}*/

	argToken := p.current()
	var arg Expression
	switch argToken.Kind {
	case lex.TokenNUMBER:
		value, _ := strconv.Atoi(argToken.Value)
		arg = IntE{Value: value}
	}

	if p.next() == errEOF {
		return arg
	}

	for p.at < len(p.in) {
		operator := p.current()

		var operation Expression
		switch operator.Kind {
		case lex.TokenCROSS:
			operation = Sum
		case lex.TokenHYPHEN:
			operation = Sub
		case lex.TokenASTERISK:
			operation = Mul
		case lex.TokenFORWARD_SLASH:
			operation = Div
		}

		priority := lex.InfixPriority[operator.Kind]
		if priority <= previousPriority {
			return arg
		}

		if p.next() == errEOF {
			return ApplicationE{Function: operation, Argument: arg}
		}

		arg2 := p.parse(priority)

		arg = ApplicationE{
			Function: ApplicationE{Function: operation, Argument: arg},
			Argument: arg2,
		}
	}

	return arg
}

func (p *parser) next() error {
	p.at++
	if p.at >= len(p.in) {
		return errEOF
	}

	return nil
}

func (p *parser) peekAhead() lex.Token {
	if p.at >= len(p.in)-1 {
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

func (p *parser) current() lex.Token {
	return p.in[p.at]
}

var errEOF = errors.New("EOF")

/* ------------------------- */
