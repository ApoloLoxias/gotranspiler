package ast

import "fmt"

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
	return fmt.Sprintf("IntE(%i)", i.Value)
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
