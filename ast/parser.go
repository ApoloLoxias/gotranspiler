package ast

import (
	"errors"
	"strconv"

	"github.com/ApoloLoxias/gotranspiler/lex"
)

/* --------EXPORTS---------- */

func Parse(tokens []lex.Token) Expression {
	parser := parser{tokens, nil, 0}
	return parser.parse(0)
}

/* ----MAIN PARSING LOGIC--- */

type parser struct {
	in  []lex.Token
	out Expression
	at  int
}

// The First token of an expression is special because there is no token before it
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
		priority := lex.ApplicationPriority
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
		arg := p.parse(priority) //returns nil on `)`
		if arg == nil {          //Lest `)` be treated as fun's argument
			//		p.next() // Get ')' out of the way. No need to check for EOF since we'll return anyway. // Adding this skip didn't work lol
			return fun
		}
		return ApplicationE{Function: fun, Argument: arg}
	}

	switch first.Kind {
	case lex.TokenNUMBER:
		value, _ := strconv.Atoi(first.Value)
		p.next()
		return IntE{value}
	case lex.TokenOPEN_PARENTHESIS:
		if p.next() == errEOF {
			return nil
		}
		toReturn := p.parseFirst() // An expression that follows "(" is a standalone expression and its first token must be a First. This type of call was why I separated parseFirst() into its own subfunction in the first place, so I should not call p.parse(0) here!
		// p.parseFirst() breaks "( 1 ) + - 2 * ( 8 \n )", but p.parse() breaks ""(1+0-(42/1)+1-(0))""
		if p.at < len(p.in) && p.current().Kind == lex.TokenCLOSE_PARENTHESIS { // toReturn is p.parseFirst(), because what follows "(" is a first, so we skip the infix loop, but than we don't really handle ")" in pareFirst() as we do in parse(0), so we handle it here
			/*
			* Note: I was going to compare handling ")" after exiting parseFirst to handling it before exiting parseFirst(i.e. in the parseFirst function)
			* but this is already  the parseFirst function, and we are right bbefore return anyway
			* so this is handling it inside parseFirst before returning
			* so handling it inside parseFirst before returning or handling it outside parseFirst after returning from the call triggered by "(" is the same handling
			* so I don't need to do this comparison
			 */
			p.next()
		}
		return toReturn
	}

	return nil // Falls through if First is not a Prefix, a Function, a Number, nor an TokenOPEN_PARENTHESIS -> In other words, if First is not a valid First token (i.e. an infix)
	// I guess such a fall-through could be used to assign different behaviour to non-first tokens in first position, i.e. let infixes be functionalized in first position even if they are not followed by $
}

func (p *parser) parse(previousPriority int) Expression {
	//First token must be a prefix or a standalone expression (or ´(´)
	first := p.parseFirst()
	if p.at >= len(p.in) {
		return first
	}

	// Now we will parse the second term, i.e. the infix/sufix
	var fun Expression
	//	Opening parenthesis mustn't trigger the loop, lest we update priority and arg2 := p.parse(priority), accepting `(` as the argument
	if p.current().IsInfix() || p.current().Kind == lex.TokenCLOSE_PARENTHESIS { //start infix loop. When implementing sufixes, may want to have an aditional arm for the if-else (i.e. if-else if-else)
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
			if priority <= previousPriority { // <= for left-associativity; < for right-associativity
				return first
			}

			if p.next() == errEOF { // No second operand, so sufix instead of infix: "first opeator". When implementing actual sufixes will probably not be inside this arm of the outer if-else, and leave this just for partial application of incomplete infix expressions/tokens that are ordinarily infixes in suffix position
				return ApplicationE{Function: operation, Argument: first}
			}

			arg := p.parse(priority) // The second operand: "first operator arg"

			first = ApplicationE{
				Function: ApplicationE{Function: operation, Argument: first},
				Argument: arg,
			}
		}

		return first
	} else { //If dealing with "Expression1 Expression2" (will likely be preceded by else if p.current().IsPostfix())
		// arg is Expression1, now we must apply Expression1(Expression2), by means of fun = (E1=first); arg = E2; return fun(arg)

		for p.at < len(p.in) {
			fun = first
			// operator := whitespace

			priority := lex.ApplicationPriority
			if priority <= previousPriority {
				return fun
			}

			if p.next() == errEOF {
				return fun
			} else if p.current().Kind == lex.TokenCLOSE_PARENTHESIS {
				return fun
			}

			arg := p.parse(priority)
			return ApplicationE{Function: fun, Argument: arg} //maybe just assign it to a variable that is returned at the tail by default
		}
	}

	return first //currently unreachable except for prefix followed by EOF, which won't likely be a valid expression anyway // compiler says it is not reachable at all,huh!? //Currently using a if-else (posibly an if-else if-else when dealing with suffixes) which returns on both if and else. Makes sense that this is unreachable// SO the thing about assigning a var on else arm and letting it be returned  at the tail would make sense when implementing suffixes if suffixes and applications have semantically simillar return values? // I hope that doesn't happen
}

/* --------HELPERS---------- */

// Note: length check can't guarantee safety of the outer function on recursive calls
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
