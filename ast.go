package main

// AST nodes

type Expression interface {
	ExpressionMarker()
}

// Primitive literals

type IntLiteral struct {
	Value int64
}

func (IntLiteral) ExpressionMarker() {}

type FloatLiteral struct {
	Value float64
}

func (FloatLiteral) ExpressionMarker() {}

type BoolLiteral struct {
	Value bool
}

func (BoolLiteral) ExpressionMarker() {}

type StringLiteral struct {
	Value string
}

func (StringLiteral) ExpressionMarker() {}

// Compound literals

type SliceLiteral struct {
	Value []Expression
}

func (SliceLiteral) ExpressionMarker() {}

type FunctionLiteral struct {
	Argument string
	Image    Expression
}

func (FunctionLiteral) ExpressionMarker() {}

// Runtime values TODO

type Value interface {
	ValueMarker()
}

// Others

type Variable struct {
	Name string
}

func (Variable) ExpressionMarker() {}

type Block struct {
	Bind   Expression
	To     Variable
	Assess Expression
}

func (Block) ExpressionMarker() {}

type Call struct {
	Argument Expression
	Function Expression
}

func (Call) ExpressionMarker() {}

type Conditional struct {
	Check Expression
	Yes   Expression
	No    Expression
}

func (Conditional) ExpressionMarker() {}

//Builtins TODO

type Builtin string

func (Builtin) ExpressionMarker() {}

const (
	ADD  Builtin = "+"
	SUB  Builtin = "-"
	MUL  Builtin = "*"
	DIV  Builtin = "/"
	MOD  Builtin = "%"
	EQ   Builtin = "=="
	LT   Builtin = "<"
	LOE  Builtin = "<="
	GT   Builtin = ">"
	GOE  Builtin = ">="
	DIFF Builtin = "!="
	NOT  Builtin = "!"
	OR   Builtin = "||"
	AND  Builtin = "&&"
)
