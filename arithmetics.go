package main

import "errors"
import "fmt"

type BuiltinFunc func(Expression) (Expression, error)

func (BuiltinFunc) ExpressionMarker() {}

var Sum BuiltinFunc = func(x Expression) (Expression, error) {
	ex, _ := EvalArithmetics(x)
	switch X := ex.(type) {
	case IntLiteral:
		var f BuiltinFunc = func(y Expression) (Expression, error) {
			switch Y := y.(type) {
			case IntLiteral:
				return IntLiteral{Value: X.Value + Y.Value}, nil
			default:
				return nil, errors.New("a")
			}
		}
		return f, nil
	default:
		return nil, errors.New("a")
	}
}

func EvalArithmetics(Node Expression) (Expression, error) {
	switch node := Node.(type) {
	case IntLiteral:
		return EvalAriInt(node)
	case Call:
		return EvalAriCall(node)
	case BuiltinFunc:
		//TODO

	default:
		return nil, errors.New("not arithmetics or not implemented yet")
	}

	return nil, errors.New("this shouldn't happen")
}

func EvalAriInt(x IntLiteral) (IntLiteral, error) {
	return x, nil
}

func EvalAriCall(node Call) (Expression, error) {
	f, x := node.Function, node.Argument

	switch F := f.(type) {
	case IntLiteral:
		return EvalAriInt(F)
	case Call:
		g, _ := EvalAriCall(F)
		return EvalAriCall(Call{Argument: x, Function: g})
	case BuiltinFunc:
		return F(x)

	default:
		return nil, errors.New("not ari or not implemented")
	}
}

/*
1+2
*/

var simpleAdition = Call{
	Argument: IntLiteral{2},
	Function: Call{
		Argument: IntLiteral{1},
		Function: Sum,
	},
}

func main() {
	E1, _ := EvalArithmetics(simpleAdition)
	E2, _ := EvalArithmetics(tripleAddition)
	fmt.Println("1 + 2 = ", E1)
	fmt.Println("1 + 2 + 3 = ", E2)
}

/*
1+2+3
*/

var tripleAddition = Call{
	Argument: IntLiteral{3},
	Function: Call{
		Function: Sum,
		Argument: simpleAdition,
	},
}
