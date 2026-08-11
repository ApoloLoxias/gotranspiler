package main

import "fmt"
import "github.com/ApoloLoxias/gotranspiler/lex"

func main() {
	fmt.Println(lex.Lex("1+2-3*4/5"))
}
