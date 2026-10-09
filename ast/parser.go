package ast

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ApoloLoxias/gotranspiler/lex"
)

/* --------EXPORTS---------- */

func Parse(tokens []lex.Token) (Expression, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("parser error: no tokens were input")
	}

	parser := parser{tokens, 0}
	expr := parser.parse(0)

	if parser.at < len(parser.in) {
		return expr,
			fmt.Errorf(
				"parsing stoped at token %v, at index %d",
				parser.current(),
				parser.at,
			)
	}

	if expr == nil {
		return nil, fmt.Errorf("malformed/unparsable stream of tokens: %v", tokens)
	}

	return expr, nil
}

/* ----MAIN PARSING LOGIC--- */

type parser struct {
	in []lex.Token
	at int
}

// The First token of an expression is special because there is no token before it
// Returns nil if called at a token an expression is not allowed to start with
func (p *parser) parseFirst() Expression {
	first := p.current()

	if first.IsPrefix() {
		priority := lex.PrefixPriority[first.Kind]
		if p.next() == errEOF {
			return nil
		}
		prefixArg := p.parse(priority)
		//consider a break guard for nil prefixArg if we want a lone
		//prefix to be parsed as its unaplied function
		return ApplicationE{Function: Neg, Argument: prefixArg}
	}

	//A Application branch ignroed for now
	/*if first.IsFunction() {
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
	} */

	switch first.Kind {
	case lex.TokenNUMBER:
		value, _ := strconv.Atoi(first.Value)
		p.next()
		return IntE{value}
	case lex.TokenOPEN_PARENTHESIS:
		if p.next() == errEOF {
			return nil
		}
		// consider the implications of p.parse(0) vs p.pareFirst() when reimplementing function application
		toReturn := p.parse(0) // An expression that follows "(" is a standalone expression and its first token must be a First. This type of call was why I separated parseFirst() into its own subfunction in the first place, so I should not call p.parse(0) here!
		// p.parseFirst() breaks "( 1 ) + - 2 * ( 8 \n )", but p.parse() breaks ""(1+0-(42/1)+1-(0))""
		if toReturn == nil { //propagates failed grouping (i.e. what comes after `(` can't start an expression, to allow us to parse(0) rather than parseFirst()
			return nil
		}

		if p.at < len(p.in) && p.current().Kind == lex.TokenCLOSE_PARENTHESIS { // toReturn is p.parseFirst(), because what follows "(" is a first, so we skip the infix loop, but than we don't really handle ")" in pareFirst() as we do in parse(0), so we handle it here
			/*
			* Note: I was going to compare handling ")" after exiting parseFirst to handling it before exiting parseFirst(i.e. in the parseFirst function)
			* but this is already  the parseFirst function, and we are right bbefore return anyway
			* so this is handling it inside parseFirst before returning
			* so handling it inside parseFirst before returning or handling it outside parseFirst after returning from the call triggered by "(" is the same handling
			* so I don't need to do this comparison
			 */
			p.next()
			return toReturn
		}
		// was return toReturn
		// tucked return to Return on the block/conditional abocee
		// Parse's trailing token cehck won't catch unbalacned "("
		// The following return is accessed by ending the input before closing all "("
		return nil
	}

	// return for:
	//tokens that can't start expressions/aren't firsts
	//closing parenthesis that weren't opened
	//tokens we can't handle (yet)
	return nil // Falls through if First is not a Prefix, a Function, a Number, nor an TokenOPEN_PARENTHESIS -> In other words, if First is not a valid First token (i.e. an infix)
	// I guess such a fall-through could be used to assign different behaviour to non-first tokens in first position, i.e. let infixes be functionalized in first position even if they are not followed by $
}

// returns nil if it can't parse
// Parse returns an error when parse return nil
func (p *parser) parse(previousPriority int) Expression {
	//First token must be a prefix or a standalone expression (or ´(´)
	first := p.parseFirst()
	if first == nil { //propagate the parsing failure
		return nil
	}
	if p.at >= len(p.in) { //nothing after first. It is the whole expression
		return first
	}

	for p.at < len(p.in) {
		if p.current().Kind == lex.TokenCLOSE_PARENTHESIS {
			return first
		}

		if !(p.current().IsInfix() || p.current().IsSuffix()) {
			// functional application branch
			// triggered by opening parenthesis or any regular first
			// in other words, juxtaposed expressions with no operator
			// or whitespace as an operator
		}

		operation, ok := operationFromOperator[p.current().Kind]
		if ok == false {
			return nil
		}

		var priority int
		if p.current().IsInfix() {
			priority = lex.InfixPriority[p.current().Kind]
		} else {
			priority = lex.SuffixPriority[p.current().Kind]
		}
		associativity := lex.AssociativityOf(p.current().Kind)

		if checkPriority(priority, previousPriority, associativity) {
			return first
		}

		if p.current().IsSuffix() {
			first = ApplicationE{Function: operation, Argument: first}
			p.next()
			continue
		}

		//infix branch
		if p.next() == errEOF { //partial application of infix
			return ApplicationE{Function: operation, Argument: first}
		}

		arg := p.parse(priority) //Second arg of full infix applicaiton. First is the first arg
		if arg == nil {          // second arg fails to parse
			return nil
		}

		first = ApplicationE{
			Function: ApplicationE{Function: operation, Argument: first},
			Argument: arg,
		}

	}

	return first
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

// Note: peekAhead, peekBehind not currently used, but kept as they are usefull
// tp have arround in the codebase when making changes and debugging
func (p *parser) peekAhead() lex.Token {
	if p.at >= len(p.in)-1 {
		return lex.EOFtoken
	}
	return p.in[p.at+1]
}

// Note: see peekAhead
func (p *parser) peekBehind() lex.Token {
	if p.at <= 0 {
		return lex.SOFtoken
	}
	return p.in[p.at-1]
}

// Doens't range check p.at vs len(p.in) and lets the caller guarantee safety
// Implemented so to avoid redundant checks with p.next()
func (p *parser) current() lex.Token {
	return p.in[p.at]
}

var errEOF = errors.New("EOF")

//export to a config file later; explore as a means of operator overload
var operationFromOperator = map[lex.TokenKind]Expression{ //const
	lex.TokenCROSS:             Sum,
	lex.TokenHYPHEN:            Sub,
	lex.TokenASTERISK:          Mul,
	lex.TokenFORWARD_SLASH:     Div,
	lex.TokenASTERISK_ASTERISK: Pow,
	lex.TokenEXCLAMATION:       Fac,
}

// nextPriority < previousPriority for RightAssociativity
// nextPreiirty <= previousPriority for LeftAssociativity
func checkPriority(nextPriority, previousPriority int, associativy lex.Associativity) bool {
	if nextPriority < previousPriority {
		return true
	}
	if associativy == lex.LeftAssociativity && nextPriority == previousPriority {
		return true
	}
	return false
}
