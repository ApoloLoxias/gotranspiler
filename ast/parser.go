package ast

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ApoloLoxias/gotranspiler/lex"
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

/*
func (p *parser) parse(previousPriority int) Expression {

	// Parse First
	argToken := p.current()
	var arg Expression
	var fun Expression

	if argToken.IsPrefix() {
		priority := lex.PrefixPriority[argToken.Kind]
		if p.next() == errEOF { // since arg is nill at this arm, could probably delete this, and delete p.next() from elif and else arms and just if p.next() == errEOF{return arg} out of the conditional, assuming, of course that p.next() on the else branch does indeed execute for all of the Terminals //I can't actually do that
			return nil
		}

		prefixArg := p.parse(priority) //Note: p.next() is called inside the nested p.parse call
		arg = ApplicationE{Function: Neg, Argument: prefixArg}

		// guard clause
		if p.at >= len(p.in) { //guard clause is necessary because p.next() within a nested p.parse call can't guarantee the safety of p.at/p.current() without
			return arg
		}
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
		case lex.TokenOPEN_PARENTHESIS: //TODO
			if p.next() == errEOF {
				return nil
			}
			return p.parse(0)
		}

		if p.next() == errEOF { //executes on case lex.TokenNumber for now, but left out of switch because it will probably be a default step for Terminals (i.e. non-prefix, non-function, non-parenthesis, non-EOF tokens)
			return arg
		}
	}

	//parseFollowUp
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
			case lex.TokenCLOSE_PARENTHESIS: //TODO
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
*/

/*
func (p *parser) parse(previousPriority int) Expression {
	first := p.parseFirst()
	//intE{1}

	if p.at >= len(p.in) {
		return first
	}

	return p.parseFollowUp(first, previousPriority)
}
*/

func (p *parser) parseFirst() Expression {
	first := p.current()

	if first.IsPrefix() {
		priority := lex.PrefixPriority[first.Kind]
		if p.next() == errEOF {
			return nil
		}
		prefixArg := p.parse(priority)
		return ApplicationE{Function: Neg, Argument: prefixArg}
	}

	if first.IsFunction() {
		var fun Expression
		switch first.Kind {
		case lex.TokenCROSS_FUNC:
			fun = Sum
		case lex.TokenHYPHEN_FUNC:
			fun = Sub
		case lex.TokenASTERISK_FUNC:
			fun = Mul
		case lex.TokenFORWARD_SLASH_FUNC:
			fun = Div
		}
		if p.next() == errEOF {
			return fun
		}
		arg := p.parse(0) //returns nil on `)`
		if arg == nil {   //Lest `)` be treated as fun's argument
			//		p.next() // Get ')' out of the way. No need to check for EOF since we'll return anyway. // Adding this skip didn't work lol
			return fun
		}
		return ApplicationE{Function: fun, Argument: arg}
	}

	fmt.Println("first=", first, " and is of kind ", first.Kind, "We are in a switch, and if it is open_parenthesus, we'll go to the breaking change!")
	switch first.Kind {
	case lex.TokenNUMBER:
		value, _ := strconv.Atoi(first.Value)
		p.next()
		return IntE{value}
	case lex.TokenOPEN_PARENTHESIS:
		fmt.Println("First is '(' we we'll advance and check for EOF")
		if p.next() == errEOF {
			return nil
		}
		fmt.Println("Now at p.at=", p.at, " and p.current()=", p.current(), " will return toReturn:=p.parseFirst(), This is the breaking line!!!!!!!!")
		toReturn := p.parseFirst() // An expression that follows "(" is a standalone expression and its first token must be a First. This type of call was why I separated parseFirst() into its own subfunction in the first place, so I should not call p.parse(0) here!
		// p.parseFirst() breaks "( 1 ) + - 2 * ( 8 \n )", but p.parse() breaks ""(1+0-(42/1)+1-(0))""
		if p.at < len(p.in) && p.current().Kind == lex.TokenCLOSE_PARENTHESIS { // toReturn is p.parseFirst(), because what follows "(" is a first, so we skip the infix loop, but than we don't really handle ")" in pareFirst() as we do in parse(0), so we handle it here
			p.next()
		}
		fmt.Println("returning toReturn:=p.parseFirst()=", toReturn)
		return toReturn
	}

	return nil
}

/*
func (p *parser) parseFollowUp(first Expression, previousPriority int) Expression {
	//	var result Expression
	var current lex.Token
	for p.at < len(p.in) {
		current = p.current()
		if current.IsInfix() || current.IsParenthesis() {
			for p.at < len(p.in) {
				if current.Kind == lex.TokenCLOSE_PARENTHESIS {
					p.next()
					return first
				}
				var operation Expression
				switch current.Kind {
				case lex.TokenCROSS:
					operation = Sum
				case lex.TokenHYPHEN:
					operation = Sub
				case lex.TokenASTERISK:
					operation = Mul
				case lex.TokenFORWARD_SLASH:
					operation = Div
				}
				priority := lex.InfixPriority[current.Kind]
				if priority <= previousPriority {
					return first
				}
				if p.next() == errEOF {
					return ApplicationE{Function: operation, Argument: first}
				}
				result = ApplicationE{
					Function: ApplicationE{Function: operation, Argument: first},
					Argument: p.parse(priority),
				}
			}
			//return result
		} else {
			fun := first
			arg := p.parseFirst()
			return ApplicationE{Function: fun, Argument: arg}
		}
	}
	return first
}
*/

func (p *parser) parse(previousPriority int) Expression {
	first := p.parseFirst()
	if p.at >= len(p.in) {
		return first
	}

	var fun Expression
	//	Opening parenthesis mustn't trigger the loop, lest we update priority and arg2 := p.parse(priority), accepting `(` as the argument
	//	if p.current().IsInfix() || p.current().IsParenthesis() {
	if p.current().IsInfix() || p.current().Kind == lex.TokenCLOSE_PARENTHESIS {
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
				return first
			}

			priority := lex.InfixPriority[operator.Kind]
			if priority <= previousPriority {
				return first
			}

			if p.next() == errEOF {
				return ApplicationE{Function: operation, Argument: first}
			}

			arg2 := p.parse(priority)

			first = ApplicationE{
				Function: ApplicationE{Function: operation, Argument: first},
				Argument: arg2,
			}
		}

		return first
	} else { //If dealing with "Expression1 Expression2" (will likely be preceded by else if p.current().IsPostfix())
		// arg is Expression1, now we must apply Expression1(Expression2), by means of fun = (E1=arg); arg = E2; return fun(arg) (or assign fun(arg) to arg, which is the deafult return value
		fun = first
		first = p.parse(0)
		return ApplicationE{Function: fun, Argument: first} //maybe just assign it to arg and let arg be returned at the tail by default
	}

	return first //currently unreachable except for prefix followed by EOF, which won't likely be a valid expression anyway // seems its not reachable even in that case now, huh!?
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
