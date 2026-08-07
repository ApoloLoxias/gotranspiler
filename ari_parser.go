package main

import "fmt"
import "unicode"
import "strconv"

func Parse(e Expression, tt []Token) (Expression, []Token, error) {
	j := 0
	for i, t := range tt {
		if t.Kind == SpaceToken {
			tt = append(tt[:i-j], tt[i-j+1:]...)
			j = j + 1
		}
	}

	if tt[0].Kind != NumToken {
		return nil, nil, fmt.Errorf("expression doesnt start with a number")
	}

	var arg1 Token
	var arg2 Token
	var fun Token

	arg1 = tt[0]
	tt = tt[1:]

	if tt[0].Kind != BinaryInfixToken {
		return nil, nil, fmt.Errorf("num followed by non operation")
	}

	fun = tt[0]
	tt = tt[1:]

	if tt[0].Kind != NumToken {
		return nil, nil, fmt.Errorf("operation followed by non num")
	}

	arg2 = tt[0]

	arg1V, _ := strconv.Atoi(string(arg1.Token))
	arg2V, _ := strconv.Atoi(string(arg2.Token))
	argument1 := IntLiteral{Value: int64(arg1V)}
	argument2 := IntLiteral{Value: int64(arg2V)}

	var function BuiltinFunc
	switch string(fun.Token) {
	case "+":
		function = Sum
	case "-":
		function = Sub
	case "*":
		function = Prod
	case "/":
		function = Div
	}

	expr := Call{
		Argument: argument2,
		Function: Call{
			Argument: argument1,
			Function: function,
		},
	}

	return expr, nil, nil
}

func IsPlusMinusSlashOrAsterisk(s rune) bool {
	for _, r := range []rune("+-/*") {
		if s == r {
			return true
		}
	}
	return false
}

type TokenKind string

const (
	NumToken         TokenKind = "numerical token"
	SpaceToken                 = "whitespace token"
	OtherToken                 = "other token"
	BinaryInfixToken           = "Binary Infix"
)

type Token struct {
	Token []rune
	Kind  TokenKind
}

func (t Token) Equals(u Token) bool {
	if t.Kind == u.Kind && string(t.Token) == string(u.Token) {
		return false
	}

	return false
}

type TokenThing struct {
	Tokens []Token
	Runes  []rune
}

func (t Token) Priority() (int, error) {
	if t.Kind == BinaryInfixToken {
		switch string(t.Token) {
		case "+", "-":
			return 0, nil
		case "*", "/":
			return 1, nil
		default:
			return 0, fmt.Errorf("malformed token")
		}
	}
	return 0, fmt.Errorf("token is not an operator")
}

func (t *TokenThing) Step() error {
	nextToken, err := NextToken(t.Runes)
	if err != nil {
		return err
	}
	if len(nextToken) == 0 {
		return fmt.Errorf("done")
	}

	nt := Token{Token: nextToken}
	nt.Kind = KindOfToken(nt.Token)

	t.Tokens = append(t.Tokens, nt)
	t.Runes = t.Runes[len(nextToken):]

	return nil
}

func KindOfToken(s []rune) TokenKind {
	if len(s) == 0 {
		return TokenKind("no token")
	}

	if unicode.IsDigit(s[0]) {
		return NumToken
	}
	if unicode.IsSpace(s[0]) {
		return SpaceToken
	}

	if len(s) == 1 {
		switch string(s[0]) {
		case "+", "-", "/", "*":
			return BinaryInfixToken
		}
	}

	return OtherToken
}

func (t *TokenThing) Tokenize() error {
	var err error

	for err == nil {
		err = t.Step()
		if err == fmt.Errorf("done") {
			return nil
		}
	}

	return err
}

func NextToken(s []rune) ([]rune, error) {
	nextToken := make([]rune, 0)
	var parsing rune

	for {
		if len(s) == 0 {
			return nextToken, nil
		}
		nextToken = append(nextToken, s[0])
		s = s[1:]

		parsing = nextToken[len(nextToken)-1]
		if unicode.IsDigit(parsing) {
			for {
				if len(s) == 0 || !unicode.IsDigit(s[0]) {
					return nextToken, nil
				}
				nextToken = append(nextToken, s[0])
				s = s[1:]
				parsing = nextToken[len(nextToken)-1]
			}
		}

		if IsPlusMinusSlashOrAsterisk(parsing) {
			return nextToken, nil
		}

		if unicode.IsSpace(parsing) {
			return nextToken, nil
		}

		return nil, fmt.Errorf("can't parse %q", parsing)
	}
}

func main() {
	twelve, _ := NextToken([]rune("12+2"))
	onetwo, _ := NextToken([]rune("12 2 3"))
	plus, _ := NextToken([]rune("+-"))
	_, err := NextToken([]rune("a"))
	fmt.Println(string(twelve), string(onetwo), string(plus), err)
	parse1, _, _ := Parse(nil, []Token{Token{Token: []rune("1"), Kind: NumToken}, Token{Token: []rune("+"), Kind: BinaryInfixToken}, Token{Token: []rune("12"), Kind: NumToken}})

	s := "12 +2*3 1"
	tt := TokenThing{Runes: []rune(s)}
	tt.Tokenize()

	fmt.Println(s)
	for _, token := range tt.Tokens {
		fmt.Printf("%s, %s\n", string(token.Token), token.Kind)
	}

	fmt.Println(EvalArithmetics(parse1))
}
