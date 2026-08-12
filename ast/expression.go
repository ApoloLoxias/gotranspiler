package ast

import (
	"fmt"
	"strings"
)

type Expression interface {
	Evaluate() Expression
	String() string
	pretty(int, int) string
	Pretty(...int) string
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

func (i IntE) pretty(int, int) string {
	return i.String()
}

func (i IntE) Pretty(tabSpace ...int) string {
	tab := 8
	if len(tabSpace) > 0 {
		tab = tabSpace[0]
	}
	return i.pretty(tab, 0)
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

func (a ApplicationE) pretty(tabSpace, lvl int) string {
	tab := strings.Repeat(" ", tabSpace)
	indent := strings.Repeat(tab, lvl)

	return fmt.Sprintf("ApplicationE{\n%s%sFunction: %s,\n%s%sArgument: %s,\n%s}", indent, tab, a.Function.pretty(tabSpace, lvl+1), indent, tab, a.Argument.pretty(tabSpace, lvl+1), indent)
}

func (a ApplicationE) Pretty(tabsSpace ...int) string {
	tab := 8
	if len(tabsSpace) > 0 {
		tab = tabsSpace[0]
	}
	return a.pretty(tab, 0)
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

func (f BuiltInFunc) pretty(int, int) string {
	return f.String()
}

func (f BuiltInFunc) Pretty(tabSpace ...int) string {
	tab := 8
	if len(tabSpace) > 0 {
		tab = tabSpace[0]
	}
	return f.pretty(tab, 0)
}

var Sum = BuiltInFunc{Name: "Sum"} //const

var Sub = BuiltInFunc{Name: "Subtraction"} //const

var Mul = BuiltInFunc{Name: "Multiplication"} //const

var Div = BuiltInFunc{Name: "Division"} // const
