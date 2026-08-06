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
}
