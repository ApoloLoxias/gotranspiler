package main

import "fmt"
import "unicode"

func Parse(e Expression, s []rune) (Expression, []rune) {
	return e, s
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
	NumToken   TokenKind = "numerical token"
	OtherToken           = "other token"
)

type Token struct {
	Token []rune
	Kind  TokenKind
}

type TokenThing struct {
	Tokens []Token
	Runes  []rune
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

	s := "12 +2*3 1"
	tt := TokenThing{Runes: []rune(s)}
	tt.Tokenize()

	fmt.Println(s)
	for _, token := range tt.Tokens {
		fmt.Printf("%s, %s\n", string(token.Token), token.Kind)
	}
}
