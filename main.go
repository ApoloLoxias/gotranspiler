package main

import "fmt"
import "github.com/ApoloLoxias/gotranspiler/lex"
import "github.com/ApoloLoxias/gotranspiler/ast"

func main() {
	test("1+2")
	test("11-22")
	test("123*456")
	test("01/2000")
}

func test(s string) {
	fmt.Println(s)

	tokens := lex.Lex(s)
	fmt.Println(tokens)

	str := ast.Parse(tokens).String()
	fmt.Println(str)

	fmt.Println("------------------")
}
