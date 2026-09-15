package main

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ApoloLoxias/gotranspiler/ast"
	"github.com/ApoloLoxias/gotranspiler/lex"
)

func main() {
	printTitle("basic infix")
	test("1")
	test("1+2")
	test("12-23")
	test("1000*12345")
	test("4/2")
	test("5/2")

	printTitle("Partial infix")
	test("1+")
	test("23-")
	test("4*")
	test("0/")

	printTitle("Multiple infixes")
	test("1+3+2")
	test("1+3*2")
	test("1+2*3*4")
	test("1+2*3+4")
	test("1*2+3")
	test("1*2+3*4")
	test("1+2*3*4+5")
	test("1-2+4")

	printTitle("Prifinx '-'")
	test("-1")
	test("-1+2")
	test("-1*2")

	printTitle("Parenthesis")
	test("(1")
	test("(1)")
	test("(-1)")
	test("-(1)")
	test("(1+2)")
	test("(1+2)*3")
	test("1*(2+3)")
	test("(1+0-(42/1)+1-(0))")
	test("1*((2+3)/4)")

	printTitle("Whitespace")
	test(" -   1")
	test("2    ")
	test(" 3 ")
	test("1\n+ -7")
	test("( 1 ) + - 2 * ( 8 \n )")

	printTitle("Functionalized infixes")
	test("+$")
	test("*$")
	test("-$")
	test("/$")
	println()
	test("+$1")
	test("-$ 2")
	test("(*$)3")
	test("/$(4)")
	println()
	test("+$1 -2")
	test("(+$1)2")
	test("/$ 17(8)")

	printTitle("Ints as functions")
	test("(1)2")
	test("(1)(2)")
	test("1(2)")
	test("1 2")
	test("1 2 3")

	printTitle("Functions and Arithmetics")
	test("1 + (*$ -2 3) - -5")
	test("1 + -$ 2 3")
	test("+$ 1 - 2 3")
}

/*
func main() {
	test("1 2")
	test("(+$1)2")
}
*/

func test(s string) {
	fmt.Println(s)

	tokens := lex.Lex(s)
	fmt.Println(tokens)

	expr := ast.Parse(tokens)
	fmt.Println(expr.Pretty())
	fmt.Println(expr.Evaluate().Pretty())

	fmt.Println("------------------")
}

func printTitle(name string) {
	bufferLen := 0
	if utf8.RuneCountInString(name) < 0 {
		bufferLen = (24 - len(name)) / 2
	}
	buffer := strings.Repeat(" ", bufferLen)

	fmt.Print(
		"\n",
		"========================\n",
		buffer, name, buffer, "\n",
		"========================\n",
	)
}
