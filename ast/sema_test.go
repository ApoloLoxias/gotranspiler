package ast

import (
	"testing"
)

type manualTests struct {
	From []string
	In   []Expression
	Err  []bool
	Want []Expression
	Got  []Expression
}

type manualTest struct {
	From string
	In   Expression
	Err  bool
	Want Expression
	Got  Expression
}

var manualCases = []manualTest{
	{
		From: "x: {0}",
		In: LambdaE{
			ParamName: "x",
			ParamId:   0,
			Body:      IntE{0},
		},
		Err: false,
		Want: LambdaE{
			ParamName: "x",
			ParamId:   1,
			Body:      IntE{0},
		},
	},
	{
		From: "y: {y}",
		In: LambdaE{
			ParamName: "y",
			ParamId:   0,
			Body:      VariableE{"y", 0},
		},
		Err: false,
		Want: LambdaE{
			ParamName: "y",
			ParamId:   1,
			Body:      VariableE{"y", 1},
		},
	},
	{
		From: "x: {y}",
		In: LambdaE{
			ParamName: "x",
			ParamId:   0,
			Body:      VariableE{"y", 0},
		},
		Err: true,
	},
	{
		From: "x: {x + y:{y}}",
		In: LambdaE{
			ParamName: "x",
			ParamId:   0,
			Body: ApplicationE{
				Function: ApplicationE{
					Function: Sum,
					Argument: VariableE{"x", 0},
				},
				Argument: LambdaE{
					ParamName: "y",
					ParamId:   0,
					Body:      VariableE{"y", 0},
				},
			},
		},
		Err: false,
		Want: LambdaE{
			ParamName: "x",
			ParamId:   1,
			Body: ApplicationE{
				Function: ApplicationE{
					Function: Sum,
					Argument: VariableE{"x", 1},
				},
				Argument: LambdaE{
					ParamName: "y",
					ParamId:   2,
					Body:      VariableE{"y", 2},
				},
			},
		},
	},
	{
		From: "x: {x + x:{x}}",
		In: LambdaE{
			ParamName: "x",
			ParamId:   0,
			Body: ApplicationE{
				Function: ApplicationE{
					Function: Sum,
					Argument: VariableE{"x", 0},
				},
				Argument: LambdaE{
					ParamName: "x",
					ParamId:   0,
					Body:      VariableE{"x", 0},
				},
			},
		},
		Err: false,
		Want: LambdaE{
			ParamName: "x",
			ParamId:   1,
			Body: ApplicationE{
				Function: ApplicationE{
					Function: Sum,
					Argument: VariableE{"x", 1},
				},
				Argument: LambdaE{
					ParamName: "x",
					ParamId:   2,
					Body:      VariableE{"x", 2},
				},
			},
		},
	},
	{
		From: "x: {y: {x} x: {x}}",
		In: LambdaE{
			ParamName: "x",
			ParamId:   0,
			Body: ApplicationE{
				Function: LambdaE{
					ParamName: "y",
					ParamId:   0,
					Body:      VariableE{"x", 0},
				},
				Argument: LambdaE{
					ParamName: "x",
					ParamId:   0,
					Body:      VariableE{"x", 0},
				},
			},
		},
		Err: false,
		Want: LambdaE{ //Note: on applications, argument is indexed before function, so x: {x} gets index 2 and y: {x} gets index 3
			ParamName: "x",
			ParamId:   1,
			Body: ApplicationE{
				Function: LambdaE{
					ParamName: "y",
					ParamId:   3,
					Body:      VariableE{"x", 1},
				},
				Argument: LambdaE{
					ParamName: "x",
					ParamId:   2,
					Body:      VariableE{"x", 2},
				},
			},
		},
	},
	{
		From: "x: {x + x:{x + x:{x}}}",
		In: LambdaE{
			ParamName: "x",
			ParamId:   0,
			Body: ApplicationE{
				Function: ApplicationE{Sum, VariableE{"x", 0}},
				Argument: LambdaE{
					ParamName: "x",
					ParamId:   0,
					Body: ApplicationE{
						Function: ApplicationE{Sum, VariableE{"x", 0}},
						Argument: LambdaE{
							ParamName: "x",
							ParamId:   0,
							Body:      VariableE{"x", 0},
						},
					},
				},
			},
		},
		Err: false,
		Want: LambdaE{
			ParamName: "x",
			ParamId:   1,
			Body: ApplicationE{
				Function: ApplicationE{Sum, VariableE{"x", 1}},
				Argument: LambdaE{
					ParamName: "x",
					ParamId:   2,
					Body: ApplicationE{
						Function: ApplicationE{Sum, VariableE{"x", 2}},
						Argument: LambdaE{
							ParamName: "x",
							ParamId:   3,
							Body:      VariableE{"x", 3},
						},
					},
				},
			},
		},
	},
}

func assertManualCase(t *testing.T, test *manualTest) {
	t.Helper()

	if test.In == nil {
		tokens, err := lexString(t, test.From)
		if err != nil {
			t.Errorf("failed to lex '%s'", test.In)
		}
		ast, err := parseTokens(t, tokens)
		if err != nil {
			t.Errorf("failed to parse %s, lexed from %s", tokens, test.In)
		}
		test.In = ast
	}

	got, err := Resolve(test.In)
	if err != nil && test.Err == false {
		t.Errorf("failed to resolve %s, with errror %s", test.From, err)
	}
	if err != nil && test.Err == true {
		return
	}
	if err == nil && test.Err == true {
		t.Errorf("didn't failt to resolve %s", test.From)
	}
	test.Got = got

	if test.Got.String() != test.Want.String() {
		t.Errorf("%s, evaluates to:\n%s\n, want:\n%s", test.From, test.Got, test.Want)
	}
}

func TestVarIDs(t *testing.T) {
	for _, test := range manualCases {
		t.Run(
			test.From,
			func(t *testing.T) { assertManualCase(t, &test) },
		)
	}
}

// Properties to test:
// No lambdas have the same ID
// Varialbes with same ID/Lambdas with same ParamId have same Name/ParamName
// No free variable, whose ID isn't tied to lambda
// All variables are inside lambds with their exact ID

func assertNoLambdasHaveTheSameId(t *testing.T, test *manualTest) {
	t.Helper()

	if test.In == nil {
		tokens, err := lexString(t, test.From)
		if err != nil {
			t.Errorf("failed to lex '%s'", test.In)
		}
		ast, err := parseTokens(t, tokens)
		if err != nil {
			t.Errorf("failed to parse %s, lexed from %s", tokens, test.In)
		}
		test.In = ast
	}

	got, err := Resolve(test.In)
	if err != nil && test.Err == false {
		t.Errorf("failed to resolve %s, with errror %s", test.From, err)
	}
	if err != nil && test.Err == true {
		return
	}
	if err == nil && test.Err == true {
		t.Errorf("didn't failt to resolve %s", test.From)
	}
	test.Got = got

	lambdaAndVarIdsHelper(t, test.Got, []intString{})
}

type intString struct {
	int
	string
}

func lambdaAndVarIdsHelper(t *testing.T, expr Expression, oldStack []intString) []intString {
	t.Helper()

	switch EXPR := expr.(type) {
	case LambdaE:
		stack := append(oldStack, intString{EXPR.ParamId, EXPR.ParamName})
		lambdaAndVarIdsHelper(t, EXPR.Body, stack)
	case VariableE:
		for _, pair := range oldStack {
			if EXPR.ID == pair.int {
				if EXPR.Name == pair.string {
					return oldStack
				}
				t.Errorf("non unique var/lambda ID for var %s", EXPR.Pretty())
			}
		}
		t.Errorf("%s, in scope with variable IDs %v", EXPR, oldStack)
	case ApplicationE:
		lambdaAndVarIdsHelper(t, EXPR.Function, oldStack)
		lambdaAndVarIdsHelper(t, EXPR.Argument, oldStack)
	}

	return []intString{}
}

func TestLambdasAssignUniqueIDsToTheirParameters(t *testing.T) {
	for _, test := range manualCases {
		t.Run(
			test.From,
			func(t *testing.T) { assertNoLambdasHaveTheSameId(t, &test) },
		)
	}
}
