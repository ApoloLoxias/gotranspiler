package ast

import (
	"fmt"
)

type Expression interface {
	Evaluate() Expression
	String() string
}

/* --------------------------- */

type IntE struct {
	Value int
}

func (i IntE) Evaluate() Expression {
	return i
}

func (i IntE) String() string {
	return fmt.Sprintf("IntE(%d)", i.Value)
}

//

type ApplicationE struct {
	Function Expression
	Argument Expression
}

func (a ApplicationE) Evaluate() Expression { //TODO
	return nil
}

func (a ApplicationE) String() string {
	return fmt.Sprintf("ApplicationE{Function: %s, Argument: %s}", a.Function, a.Argument)
}

//

type BuiltInFunc struct {
	Name string
}

func (f BuiltInFunc) Evaluate() Expression { //TODO
	return nil
}

func (f BuiltInFunc) String() string {
	return fmt.Sprintf("BuiltInFunc(%s)", f.Name)
}

var Sum = BuiltInFunc{Name: "Sum"} //const

var Sub = BuiltInFunc{Name: "Subtraction"} //const

var Mul = BuiltInFunc{Name: "Multiplication"} //const

var Div = BuiltInFunc{Name: "Division"} // const
