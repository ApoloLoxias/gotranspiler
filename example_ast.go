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

// f:x -> x
// f(1)
var NamedFunc = Block{
	Bind: Call{
		Argument: ARG{ID: 1}, //ARG us an AST element for literal functions, it holds an unique ID, perhaps a hash or some other unique string so we can model ARG as a Variable{Name: ID}
		Function: IntLiteral{ARG{ID: 1}},
	},
	To: Variable{"f"},
	Assess: Call{
		Argument: IntLiteral{1},
		Function: Variable{"f"},
	},
}
var NamedFuncEval1 = Call{
	Argument: IntLiteral{1},
	Function: Call{
		Argument: ARG{ID: 1},
		Function: IntLiteral{ARG{ID: 1}},
	},
}

// shadowing
// x = 1
// x = x+1
// x + 2 = 2 + 2 =4

var Shadow = Block{
	Bind: IntLiteral{1},
	To:   Variable{"x"},
	Assess: Block{
		Bind: Call{
			Argument: IntLiteral{1},
			Function: Call{
				Argument: Variable{"x"},
				Function: Variable{"+"},
			},
		},
		To: Variable{"x"},
		Assess: Call{
			Argument: IntLiteral{2},
			Function: Call{
				Argument: Variable{"x"},
				Function: Variable{"+"},
			},
		},
	},
}
var ShadowEval1 = Block{ //This eval is wronky/wrong
	Bind: IntLiteral{1},
	To:   Variable{"x"},
	Assess: Call{
		Argument: intLiteral{2},
		Function: Call{
			Argument: Call{
				Argument: intLiteral{1},
				Function: Call{
					Argument: Variable{"x"},
					Function: Variable{"+"},
				},
			},
			Function: Variable{"+"},
		},
	},
}
var ShadowEval2 = Call{
	Argument: IntLiteral{2},
	Function: Call{
		Argument: Call{
			Argument: IntLiteral{1},
			Function: Call{
				Argument: IntlLiteral{1},
				Function: Variable{"+"},
			},
		},
		Function: Variable{"+"},
	},
}
var ShadowEval3 = Call{
	Argument: IntLiteral{2},
	Function: Call{
		Argument: Call{
			Argument: intLiteral{1},
			Function: Call{
				Argument: IntLitral{1},
				Function: ADD,
			},
		},
		Function: ADD,
	},
}
var ShadowEval4 = Call{
	Argument: IntLiteral{2},
	Function: Call{
		Argument: Call{
			Argument: IntLiteral{1},
			Function: ADD(1),
		},
		Function: ADD,
	},
}
var ShadowEval5 = Call{
	Argument: IntLiteral{2},
	Function: Call{
		Argument: (ADD(ADD(1))) (1)
		Function ADD,
	},
}
var ShadowEval6 = Call{
	Argument: IntLiteral{2},
	Function: ADD( (ADD (ADD(1)) )(1) )
}
var ShadowEval 7 = ADD( (ADD (ADD(1)) )(1) )(2) = (ADD( (ADD(1)) )(1) + 2 =( (ADD(1)) ) + 1 + 2 = 1 + 1 + 2 = 4

// A Call node whose Function child is a Call whose Argument child is an ARG resolves as a block that binds the parents argument value to the childs ARG variable
var NamedFuncEval2 = Block{
	Bind: IntLiteral{1},
	To:   ARG{ID: 1},
	Assess: Call{
		Argument: ARG{ID: 1},
		Function: IntLiteral{ARG{ID: 1}},
	},
}
var NamedFuncEval3 = Call{
	Argument: IntLiteral{1},
	Function: IntLiteral{IntLiteral{1}},
}

// f = (x) -> x+2
// f(1)
var Named = Block{
	Bind: Block{ //x -> (Add(x))(2)
		Bind: ARG{ID: 1},
		To:   Variable{"x"},
		Assess: Call{
			Argument: IntLiteral{2},
			Function: Call{
				Argument: Variable{"x"},
				Function: Variable{"+"},
			},
		},
	},
	To: Variable{"f"},
	Assess: Call{
		Argument: IntLiteral{1},
		Function: Variable{"f"},
	},
}
var NamedEval1 = Call{
	Argument: IntLiteral{1},
	Function: Block{
		Bind: ARG{ID: 1},
		To:   Variable{"x"},
		Assess: Call{
			Argument: IntLiteral{2},
			Function: Call{
				Argument: Variable{"x"},
				Function: Variable{"+"},
			},
		},
	},
}
var NamedEval2 = Call{
	Argument: IntLiteral{1},
	Function: Call{
		Argument: IntLiteral{2},
		Function: Call{
			Argument: ARG{ID: 1},
			Function: Variable{"+"},
		},
	},
}
var NamedEval3 = Call{
	Argument: IntLiteral{1},
	Function: Block{
		Bind: IntLiteral{2},
		To:   ARG{ID: 1},
		Assess: Call{
			Argument: ARG{ID: 1},
			Function: Variable{"+"},
		},
	},
}
var NamedEval4 = Call{
	Argument: IntLiteral{1},
	Function: Call{
		Argument: IntLiteral{2},
		Function: Variable{"+"},
	},
}
var NamedEval5 = Call{
	Argument: IntLiteral{1},
	Function: Call{
		Argument: IntLiteral{2},
		Function: ADD,
	},
}
var NamedEvalDone = 1 + 2
