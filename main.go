package main

import "fmt"
import "github.com/ApoloLoxias/gotranspiler/lex"
import "github.com/ApoloLoxias/gotranspiler/ast"

func main() {
	tokens := lex.Lex("1")
	str := ast.Parse(tokens).String()
	fmt.Println(str)
}
