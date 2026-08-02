package main

import "errors"
import "fmt"

type BuiltinFunc func(Expression) (Expression, error)

func (BuiltinFunc) ExpressionMarker() {}

var Sum BuiltinFunc = func(x Expression) (Expression, error) {
	switch X := x.(type) {
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
		return EvalAriCall(F)
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
	fmt.Print(EvalArithmetics(simpleAdition))
}
