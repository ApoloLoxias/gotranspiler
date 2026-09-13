package main

import "fmt"
import "github.com/ApoloLoxias/gotranspiler/lex"
import "github.com/ApoloLoxias/gotranspiler/ast"

/*
func main() {
	test("1+2")
	test("11-22")
	test("123*456")
	test("01/2000")
	test("1+3+2")
	test("1+3*2")
	test("1+2*3*4")
	test("1+2*3+4")
	test("1*2+3")
	test("1*2+3*4")
	test("1+2*3*4+5")
	test("1-2+4")
	test("1+")

	fmt.Println("\n==============\n")

	test("1*2")
	test("1-2")
	test("1/2")

	fmt.Println("\n==============\n")

	test("(1")
	test("(1)")
	test("(1+2)")
	test("(1+2)*3")
	test("1*(2+3)")
	test("(1+0-(42/1)+1-(0))")
	test("1*((2+3)/4)")

	fmt.Println("\n==============\n")

	test("-1")
	test("-1+2")
	test("-1*2")
	test("1*(-2)+3")
	test("-(1+2)")

	fmt.Println("\n==============\n")

	test("1   + \n(    - 3 )")

	fmt.Println("\n==============\n")

	test("+$")
}
*/

func main() {
	test("+$")

	fmt.Println("\n============\n")

	tk := []lex.Token{
		{Value: "+$", Kind: lex.TokenCROSS_FUNC},
		{Value: "-", Kind: lex.TokenHYPHEN},
		{Value: "1", Kind: lex.TokenNUMBER},
	}
	fmt.Println(tk)

	exp := ast.Parse(tk)
	fmt.Println(exp.Pretty())

	fmt.Println("\n============\n")

	tk = []lex.Token{
		{Value: "+$", Kind: lex.TokenCROSS_FUNC},
		{Value: "1", Kind: lex.TokenNUMBER},
	}
	fmt.Println(tk)

	exp = ast.Parse(tk)
	fmt.Println(exp.Pretty())

	fmt.Println("\n============\n")

	tk = []lex.Token{
		{Value: "+$", Kind: lex.TokenCROSS_FUNC},
		{Value: "1", Kind: lex.TokenNUMBER},
		{Value: "1", Kind: lex.TokenNUMBER},
	}
	fmt.Println(tk)

	exp = ast.Parse(tk)
	fmt.Println(exp.Pretty())

	fmt.Println("\n============\n")

	test("1")
	test("1+1")
}

func test(s string) {
	fmt.Println(s)

	tokens := lex.Lex(s)
	fmt.Println(tokens)

	expr := ast.Parse(tokens)
	fmt.Println(expr.Pretty())
	fmt.Println(expr.Evaluate().Pretty())

	fmt.Println("------------------")
}
