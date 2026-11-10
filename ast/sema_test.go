package ast

import (
	"errors"
	"fmt"
	"github.com/ApoloLoxias/gotranspiler/lex"
	"reflect"
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
				Function: ApplicationE{
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
				Argument: LambdaE{},
			},
		},
		Err: false,
		Want: LambdaE{
			ParamName: "x",
			ParamId:   1,
			Body: ApplicationE{
				Function: ApplicationE{
					Function: LambdaE{
						ParamName: "y",
						ParamId:   2,
						Body:      VariableE{"x", 1},
					},
					Argument: LambdaE{
						ParamName: "x",
						ParamId:   3,
						Body:      VariableE{"x", 3},
					},
				},
				Argument: LambdaE{},
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
