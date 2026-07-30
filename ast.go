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

type sliceLiteral struct {
	Value []Expression
}

func (sliceLiteral) ExpressionMarker() {}

type functionLiteral struct {
	Argument string
	Image    Expression
}

func (functionLiteral) ExpressionMarker() {}

// Runtime values TODO

type Value interface {
	ValueMarker()
}
