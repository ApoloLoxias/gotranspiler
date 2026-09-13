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
	var fun Expression

	if argToken.IsPrefix() {
		priority := lex.PrefixPriority[argToken.Kind]
		if p.next() == errEOF { // since arg is nill at this arm, could probably delete this, and delete p.next() from elif and else arms and just if p.next() == errEOF{return arg} out of the conditional, assuming, of course that p.next() on the else branch does indeed execute for all of the Terminals
			return nil
		}

		prefixArg := p.parse(priority)
		arg = ApplicationE{Function: Neg, Argument: prefixArg}
	} else if argToken.IsFunction() {
		switch argToken.Kind {
		case lex.TokenCROSS_FUNC:
			fun = Sum
		}
		if p.next() == errEOF {
			return fun
		}
		arg = p.parse(0)
		arg = ApplicationE{Function: fun, Argument: arg}
		if p.next() == errEOF {
			return arg
		}
	} else {
		switch argToken.Kind {
		case lex.TokenNUMBER:
			value, _ := strconv.Atoi(argToken.Value)
			arg = IntE{Value: value}
		case lex.TokenOPEN_PARENTHESIS:
			if p.next() == errEOF {
				return nil
			}
			return p.parse(0)
		}

		if p.next() == errEOF { //executes on case lex.TokenNumber for now, but left out of switch because it will probably be a default step for Terminals (i.e. non-prefix, non-function, non-parenthesis, non-EOF tokens)
			return arg
		}
	}

	if p.current().IsInfix() || p.current().IsParenthesis() {
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
			case lex.TokenCLOSE_PARENTHESIS:
				p.next()
				return arg
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
	} else { //If dealing with "Expression1 Expression2" (will likely be preceded by else if p.current().IsPostfix())
		// arg is Expression1, now we must apply Expression1(Expression2), by means of fun = (E1=arg); arg = E2; return fun(arg) (or assign fun(arg) to arg, which is the deafult return value
		fun = arg
		arg = p.parse(0)
		return ApplicationE{Function: fun, Argument: arg} //maybe just assign it to arg and let arg be returned at the tail by default
	}

	return arg //currently unreachable except for prefix followed by EOF, which won't likely be a valid expression anyway // seems its not reachable even in that case now, huh!?
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
