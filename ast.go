package main

// AST nodes
type Expression interface {
	ExpressionMarker()
}

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

// Runtime values TODO
type Value interface {
	ValueMarker()
}
