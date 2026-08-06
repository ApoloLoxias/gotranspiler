package main

import "slices"
import "fmt"

func Parse(e Expression, s []rune) (Expression, []rune) {
	return e, s
}

func NextToken(s []rune) []rune {
	nextToken := make([]rune, 0)

	for i := 0; i < len(s); i++ {
		if slices.Contains([]rune("0123456789"), s[i]) {
			nextToken = append(nextToken, s[i])
		} else {
			break
		}
	}

	return nextToken
}

func main() {
	s := []rune("123aw+7\n")
	fmt.Println(string(NextToken(s)))
	fmt.Println(string(NextTokenAlt(s)))
}

func NextTokenAlt(s []rune) []rune {
	nextToken := make([]rune, 0)

	for {
		nextToken = append(nextToken, s[0])
		s = s[1:]

		if slices.Contains([]rune(" \n()[]{}'\"`´~^\\|,.<>;:/?!@#$%¨&*=+"), s[0]) {
			return nextToken
		}
	}
}
