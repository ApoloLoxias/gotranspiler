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
	Let:  []Variable{{"x"}},
	Bind: []Expression{IntLiteral{1}},
	In: Call{
		Argument: IntLiteral{2},
		Function: Call{
			Argument: Variable{"x"},
			Function: Variable{"+"},
		},
	},
}
var VarEval1 = Block{
	Let:  []Variable{{"x"}},
	Bind: []Expression{IntLiteral{1}},
	In: Call{
		Argument: IntLiteral{2},
		Function: Call{
			Argument: Variable{"x"},
			Function: ADD,
		},
	},
}
var VarEval2 = Block{
	Let:  []Variable{{"x"}},
	Bind: []Expression{IntLiteral{1}},
	In: Call{
		Argument: IntLiteral{2},
		Function: Builtin("Add'x'"),
	},
}
var VarEval3 = Block{
	Let:  []Variable{{"x"}},
	Bind: []Expression{IntLiteral{1}},
	In: Call{
		Argument: IntLiteral{2},
		Function: Builtin("Add'x'"),
	},
}
var VarEval4 = Block{
	Let:  []Variable{{"x"}},
	Bind: []Expression{IntLiteral{1}},
	In:   IntLiteral{'x' + 2},
}
var VarEval5 = IntLiteral{1 + 2}
var VarEvalDone = 3

var VarEvalAlt1 = Call{
	Argument: IntLiteral{2},
	Function: Call{
		Argument: IntLiteral{1},
		Function: ADD,
	},
}
var VarEvalAlt2 = Call{
	Argument: IntLiteral{2},
	Function: Builtin("ADD1"),
}
var VarEvalAlt3 = IntLiteral{1 + 2}
var VarEvalAltDone = 3
