package ast

import (
	"errors"
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
	switch f := a.Function.(type) {
	case BuiltInFunc:
		result, _ := f.Func(a.Argument.Evaluate())
		return result
	case ApplicationE:
		fEval := f.Evaluate()
		argEval := a.Argument.Evaluate()
		appEval := ApplicationE{Function: fEval, Argument: argEval}
		return appEval.Evaluate()
	default:
		break
	}
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
	Func func(Expression) (Expression, error)
}

func (f BuiltInFunc) Evaluate() Expression { //TODO
	return f
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

var Sum = BuiltInFunc{Name: "Sum", Func: SumBuiltIn} //const
func SumBuiltIn(x Expression) (Expression, error) {
	switch X := x.(type) {
	case IntE:
		f := BuiltInFunc{
			Name: fmt.Sprintf("Sum%d", X.Value),
			Func: func(y Expression) (Expression, error) {
				switch Y := y.(type) {
				case IntE:
					return IntE{Value: X.Value + Y.Value}, nil
				default:
					return nil, errors.New("curried Sum typeError")
				}
			},
		}
		return f, nil
	default:
		return nil, errors.New("Sum typeError")
	}
}

var Sub = BuiltInFunc{Name: "Subtraction", Func: SubBuiltIn} //const
func SubBuiltIn(x Expression) (Expression, error) {
	switch X := x.(type) {
	case IntE:
		f := BuiltInFunc{
			Name: fmt.Sprintf("Sub%d", X.Value),
			Func: func(y Expression) (Expression, error) {
				switch Y := y.(type) {
				case IntE:
					return IntE{Value: X.Value - Y.Value}, nil
				default:
					return nil, errors.New("curried Sub typeError")
				}
			},
		}
		return f, nil
	default:
		return nil, errors.New("Sub typeError")
	}
}

var Mul = BuiltInFunc{Name: "Multiplication", Func: MulBuiltIn} //const
func MulBuiltIn(x Expression) (Expression, error) {
	switch X := x.(type) {
	case IntE:
		f := BuiltInFunc{
			Name: fmt.Sprintf("Mul%d", X.Value),
			Func: func(y Expression) (Expression, error) {
				switch Y := y.(type) {
				case IntE:
					return IntE{Value: X.Value * Y.Value}, nil
				default:
					return nil, errors.New("curried Mul typeError")
				}
			},
		}
		return f, nil
	default:
		return nil, errors.New("Mul typeError")
	}
}

var Div = BuiltInFunc{Name: "Division", Func: DivBuiltIn} // const
func DivBuiltIn(x Expression) (Expression, error) {
	switch X := x.(type) {
	case IntE:
		f := BuiltInFunc{
			Name: fmt.Sprintf("Div%d", X.Value),
			Func: func(y Expression) (Expression, error) {
				switch Y := y.(type) {
				case IntE:
					return IntE{Value: X.Value / Y.Value}, nil
				default:
					return nil, errors.New("curried Div typeError")
				}
			},
		}
		return f, nil
	default:
		return nil, errors.New("Div typeError")
	}
}
