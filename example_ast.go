package main

//Simple Operations

// 1 + 2
var Add = Call{
	Argument: IntLiteral{Value: 2},
	Function: Call{
		Argument: IntLiteral{Value: 1},
		Function: Variable{Name: "+"},
	},
}

var AddEval1 = Call{
	Argument: IntLiteral{Value: 2},
	Function: Call{
		Argument: IntLiteral{Value: 1},
		Function: ADD,
	},
}
var AddEval2 = Call{
	Argument: IntLiteral{2},
	Function: Builtin("ADD1"),
}
var AddEval3 = IntLiteral{Value: 1 + 2}
var AddEvalDone = 3

// 1 + 2 + 3
var AddCurried = Call{
	Argument: IntLiteral{Value: 3},
	Function: Call{
		Argument: Add,
		Function: Variable{Name: "+"},
	},
}
var AddCEval1 = Call{
	Argument: IntLiteral{3},
	Function: Call{
		Argument: IntLiteral{3},
		Function: ADD,
	},
}
var AddCEval2 = Call{
	Argument: IntLiteral{3},
	Function: Builtin("ADD3"),
}
var AddCEval3 = IntLiteral{3 + 3}
var AddCEvalDone = 6

// 1 + 2 * 3
var Precedence = Call{ // 1 + (2*3)
	Argument: Call{ //2*3,
		Argument: IntLiteral{3}, //3
		Function: Call{ // "Times 2"
			Argument: IntLiteral{2},
			Function: Variable{Name: "*"},
		},
	},
	Function: Call{
		Argument: IntLiteral{1},
		Function: Variable{Name: "+"},
	},
}
var PrecEval1 = Call{ // 1 + (2*3)
	Argument: Call{ //2*3,
		Argument: IntLiteral{3}, //3
		Function: Call{ // "Times 2"
			Argument: IntLiteral{2},
			Function: MUL,
		},
	},
	Function: Call{
		Argument: IntLiteral{1},
		Function: ADD,
	},
}
var PrecEval2 = Call{ // 1 + (2*3)
	Argument: Call{ //2*3,
		Argument: IntLiteral{3}, //3
		Function: Builtin("MUL2"),
	},
	Function: Builtin("Add1"),
}
var PrecEval3 = Call{
	Argument: IntLiteral{2 * 3},
	Function: Builtin("Add1"),
}
var PrecEval4 = Call{
	Argument: IntLiteral{6},
	Function: Builtin("Add1"),
}
var PrecEval5 = IntLiteral{1 + 6}
var PrecEvalDone = 6

/* ------------------ */

// Variables, blocks, let/in bindings and lexical scoping

// x = 1
// x + 2

var Var = Block{
	Bind: IntLiteral{1},
	To:   Variable{"x"},
	Assess: Call{
		Argument: IntLiteral{2},
		Function: Call{
			Argument: Variable{"x"},
			Function: Variable{"+"},
		},
	},
}
var VarEval1 = Call{
	Argument: IntLiteral{2},
	Function: Call{
		Argument: IntLiteral{1},
		Function: Variable{"+"}, // This one evaluates after "x" because we treat initial environment as being a block itself, so we are implicilty withing a Block{Bind: Builtins, Asses: everything}[
	},
}
var VarEval2 = Call{
	Argument: IntLiteral{2},
	Function: Call{
		Argument: IntLiteral{1},
		Function: Add,
	},
}
var VarEval3 = Call{
	Argument: IntLiteral{2},
	Function: Builtin("Add1"),
}
var VarEval4 = IntLiteral{1 + 2}
var VarEvalDone = 3

// x = 1
// y = 2
// x + y

var uMultiVar = Block{
	Bind: IntLiteral{1},
	To:   Variable{"x"},
	Assess: Block{
		Bind: IntLiteral{2},
		To:   Variable{"y"},
		Assess: Call{
			Argument: Variable{"y"},
			Function: Call{
				Argument: Variable{"x"},
				Function: Variable{"+"},
			},
		},
	},
}
var uMultiVarEval1 = Block{
	Bind: IntLiteral{2},
	To:   Variable{"y"},
	Assess: Call{
		Argument: Variable{"y"},
		Function: Call{
			Argument: IntLiteral{1},
			Function: Variable{"+"},
		},
	},
}
var uMultiVarEval = Call{
	Argument: IntLiteral{2},
	Function: Call{
		Argument: IntLiteral{1},
		Function: Variable{"+"},
	},
}
var uMultivarEvaol3 = Call{
	Argument: IntLiteral{2},
	Function: Call{
		Argument: IntLiteral{1},
		Function: ADD,
	},
}
var uMultivarEvalDone = 1 + 2

// Declaring/scoping funcitons

// f = (x) -> x+2
// f(1)

var NamedFunc = Block{
	Bind: FunctionLiteral{
		Argument: Variable{"x"},
		Image: IntLiteral{Variable{"x"}+IntLiteral{2}},
	}, 
	To: Variable{"f"},
	Assess: Call{
		Argument: IntLiteral{1},
		Function: Variable{"f"},
	},
}
var NamedFuncEval1 = Call{
		Argument: IntLiteral{1},
		Function: FunctionLiteral{
			Argument: Variable{"x"},
			Image: IntLiteral{Variable{"x"}+IntLiteral{2}},
		},
	}
var NamedFuncEval2 = Call{
	Block{
		Bind: IntLiteral{1},
		To: Variable{"x"},
		Assess: IntLiteral{Variable{"x"}+IntLiteral{2}},
	},
}
var NamedFuncEvalDone = IntLiteral{1+2} // If a call has a FunctioNliteral as its Function, it evaluates to a block with Bind = Call.Argument, To = Call.Function.Argument and Assess = Call.Function.Image


// Composing functions

// f = (x) -> x + 1
// g = (x) --> x - 2
// (f(g))(4)

var NestedFun = Block{ // Using simplified, non-cannonical, FunctionLiteral notation for simplicity
	Bind: FunctionLiteral{"x" -> "x+1"},
	To: Variable{"f"},
	Assess: Block{
		Bind: FunctionLiteral{"x" -> "x-2"},
		To: Variable{"g"},
		Assess Call{
			Argument: IntLiteral{4},
			Function: Call{
				Argument: Variable("g"),
				Function: Varaible{"f"},
			},
		},
	},
}
var NestedFunEval1 = Block{
	Bind: FunctionLiteral{"x" -> "x+1"},
	To: Variable{"f"},
	Assess: Call{
		Argument: IntLiteral{4},
		Function: Call{
			Argument: FunctionLiteral{"x" -> "x-2"},
			Function: Varaible{"f"},
		},
	},
}
var NestedFunEval2 = Call{
	Argument:IntLiteral{4},
	Function: Call{
		Argument: FunctionLiteral{"x" -> "x-2"},
		Function: FunctionLiteral{"x" -> "x+1"},//=FunctionLiteral{Argument:Variable{"x"}, Image: IntLiteral{Variable{"x"}+IntLiteral{1^}}}
	}
}
var NestedFunEval3 = Call{
	Argument: IntLiteral{4},
	Function: Block{
		Bind: FunctionLiteral{"x" -> "x-2"},
		To: Variable{"x"},
		Function: IntLiteral{Variable{"x"}+IntLiteral{1}}
	},
}
var NestedFunEval4 = Call{
	Argument: IntLiteral{4},
	Function: IntLiteral{FunctionLiteral{"x" -> "x-2"}+IntLiteral{1}},
}
var NestedFunEval4 = Call{ //expanding inline...
	Argument: IntLiteral{4},
	Function: IntLiteral{FunctionLiteral{Argument: Variable{"x"}, Image: IntLiteral{Variable{"x"}-IntLiteral{2}}+IntLiteral{1}}},
}
// A call whose function is an expression with functionLiterals with a variable not declared in its environment passess its argument to that variable in a block
var NestedFunEval5 = Block{
	Bind: IntLiteral{4},
	To: Variable{x},
	Assess:  IntLiteral{FunctionLiteral{Argument: Variable{"x"}, Image: IntLiteral{Variable{"x"}-IntLiteral{2}}+IntLiteral{1}}},
}
var NestedFunEval6 = IntLiteral{{FunctionLiteral{Argument: IntLiteral{4}, Image: IntLiteral{IntLiteral{4}-IntLiteral{2}}+IntLiteral{1}}},
var NestedFunEval = IntLiteral{IntLiteral{IntLiteral{4}-IntLiteral{2}}+IntLiteral{1} = 4-2+1 = 3
