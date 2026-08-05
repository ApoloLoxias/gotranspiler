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

var Sub BuiltinFunc = func(x Expression) (Expression, error) {
	ex, _ := EvalArithmetics(x)
	switch X := ex.(type) {
	case IntLiteral:
		var f BuiltinFunc = func(y Expression) (Expression, error) {
			switch Y := y.(type) {
			case IntLiteral:
				return IntLiteral{Value: X.Value - Y.Value}, nil
			default:
				return nil, errors.New("a")
			}
		}
		return f, nil
	default:
		return nil, errors.New("a")
	}
}

var Prod BuiltinFunc = func(x Expression) (Expression, error) {
	ex, _ := EvalArithmetics(x)
	switch X := ex.(type) {
	case IntLiteral:
		var f BuiltinFunc = func(y Expression) (Expression, error) {
			switch Y := y.(type) {
			case IntLiteral:
				return IntLiteral{Value: X.Value * Y.Value}, nil
			default:
				return nil, errors.New("a")
			}
		}
		return f, nil
	default:
		return nil, errors.New("A")
	}
}

var Div BuiltinFunc = func(x Expression) (Expression, error) {
	ex, _ := EvalArithmetics(x)
	switch X := ex.(type) {
	case IntLiteral:
		var f BuiltinFunc = func(y Expression) (Expression, error) {
			switch Y := y.(type) {
			case IntLiteral:
				return IntLiteral{Value: X.Value / Y.Value}, nil
			default:
				return nil, errors.New("")
			}
		}
		return f, nil
	default:
		return nil, errors.New("")
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
		return Node, nil

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
	X, _ := EvalArithmetics(x)

	switch F := f.(type) {
	case IntLiteral:
		return EvalAriInt(F)
	case Call:
		g, _ := EvalAriCall(F)
		return EvalAriCall(Call{Argument: X, Function: g})
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
	E3, _ := EvalArithmetics(simpleProd)
	E4, _ := EvalArithmetics(curProd)
	E5, _ := EvalArithmetics(mix)
	E6, _ := EvalArithmetics(mix2)
	fmt.Println("1 + 2 = ", E1)
	fmt.Println("1 + 2 + 3 = ", E2)
	fmt.Println("1*2 = ", E3)
	fmt.Println("1 * 2 * 3 = ", E4)
	fmt.Println("1*2 + 3 = ", E5)
	fmt.Println("1 + 2*3 = ", E6)
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

// 1*2

var simpleProd = Call{
	Argument: IntLiteral{2},
	Function: Call{
		Argument: IntLiteral{1},
		Function: Prod,
	},
}

// 1*2*3

var curProd = Call{
	Argument: IntLiteral{3},
	Function: Call{
		Argument: simpleProd,
		Function: Prod,
	},
}

// 1*2+3

var mix = Call{
	Argument: IntLiteral{3},
	Function: Call{
		Argument: simpleProd,
		Function: Sum,
	},
}

// 1+2*3

var mix2 = Call{
	Argument: Call{
		Argument: IntLiteral{3},
		Function: Call{
			Argument: IntLiteral{2},
			Function: Prod,
		},
	},
	Function: Call{
		Argument: IntLiteral{1},
		Function: Sum,
	},
}
