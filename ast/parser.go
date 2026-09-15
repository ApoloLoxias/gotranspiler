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
	fmt.Println("Iinside p.parseFirst() for p.at=", p.at, "and p.current()=", p.current())
	first := p.current()
	fmt.Println("first := p,current()=", p.current())

	fmt.Println("Checking if first=", first, "is prefix")
	if first.IsPrefix() {
		fmt.Println("It is!")
		priority := lex.PrefixPriority[first.Kind]
		fmt.Println("priority := PrefixPriority of first = ", priority)
		fmt.Println("will advance and check for EOF")
		if p.next() == errEOF {
			return nil
		}
		fmt.Println("Now p.at=", p.at, " and p.current()=", p.current())
		fmt.Println("Will call prefixArg := p.parse(priority), for pirority=", priority)
		prefixArg := p.parse(priority)
		fmt.Println("exited prefixArg := p.parse(priority) = ", prefixArg)
		fmt.Println("will return ", ApplicationE{Function: Neg, Argument: prefixArg})
		return ApplicationE{Function: Neg, Argument: prefixArg}
	}

	if first.IsFunction() {
		fmt.Println("First = ", first, " is a function!")
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
		fmt.Println("fun = ", fun)
		fmt.Println("will advance and check for EOF")
		if p.next() == errEOF {
			return fun
		}
		fmt.Print("Now at p.at=", p.at, " p.curren()t=", p.current())
		fmt.Println("Will call arg := p.parse(0)!!!! This is the non-breaking line!")
		arg := p.parse(0) //returns nil on `)`
		fmt.Println("Now exiting arg := p.parse(0)=", arg, " If arg is nill, it we should be at a parenthesis. Now, p,at=", p.at, " p.current()=", p.current())
		if arg == nil { //Lest `)` be treated as fun's argument
			//		p.next() // Get ')' out of the way. No need to check for EOF since we'll return anyway. // Adding this skip didn't work lol
			fmt.Println("Will return fun=", fun, " because arg=nil")
			return fun
		}
		fmt.Println("Will return ApplicationE{Function: fun, Argument: arg}=", ApplicationE{Function: fun, Argument: arg}, " because arg is not nil")
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
		return p.parse(0) // An expression that follows "(" is a standalone expression and its first token must be a First. This type of call was why I separated parseFirst() into its own subfunction in the first place, so I should not call p.parse(0) here!
		// p.parseFirst() breaks "( 1 ) + - 2 * ( 8 \n )", but p.parse() breaks ""(1+0-(42/1)+1-(0))""
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
	fmt.Println("Calling parse, p.at=", p.at) //TODO
	first := p.parseFirst()
	fmt.Println("finished executing parseFirst and assigned its return value to first=", first) //TODO
	fmt.Println("length check")
	if p.at >= len(p.in) {
		return first
	}

	var fun Expression
	//	Opening parenthesis mustn't trigger the loop, lest we update priority and arg2 := p.parse(priority), accepting `(` as the argument
	//	if p.current().IsInfix() || p.current().IsParenthesis() {
	fmt.Println("p.current().IsInfix() || p.current().Kind == lex.TokenCLOSE_PARENTHESIS=", p.current().IsInfix() || p.current().Kind == lex.TokenCLOSE_PARENTHESIS) //TODO
	if p.current().IsInfix() || p.current().Kind == lex.TokenCLOSE_PARENTHESIS {
		for p.at < len(p.in) {
			operator := p.current()
			fmt.Println("Inside loop with operator = p.current() = ", operator)

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

			fmt.Println("Operation =", operation)
			priority := lex.InfixPriority[operator.Kind]
			fmt.Println("priority=", priority)
			if priority <= previousPriority {
				fmt.Println("return first=", first, "since priority <=previousPriority")
				return first
			}
			fmt.Println("not return first=", first, "since priority !<=previousPriority")

			fmt.Println("advancing and checking EOF")
			if p.next() == errEOF {
				return ApplicationE{Function: operation, Argument: first}
			}
			fmt.Println("p.at=", p.at, "p.current()=", p.current())

			fmt.Println("will call arg2:=p.parse(pirority), with priority=", priority)
			arg2 := p.parse(priority)
			fmt.Println("returned from arg2:=p.parse(priority) call with arg2=", arg2)

			first = ApplicationE{
				Function: ApplicationE{Function: operation, Argument: first},
				Argument: arg2,
			}
			fmt.Println("first=", first)
		}

		fmt.Println("Will return first=", first)
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
