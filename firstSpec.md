# First AST spec

## Types

### Primitives

- int64
- float64
- bool
- string

### Composites

- slice (including of functions)
- function

## Expressions

- Only expressions. No statements.
- All expressions are evaluatable to a single value of the formerly declared types
- A file contains a single expression whose value is exported to a variable in the transpiled go code

## Variables and  Scoping

- Lexical scoping
- Let bindings
- No mutations

## Functions

- First class citizens
- Have exactly one input/argument (no zero args functions!)
- Return exactly one value (which, naturally, can be a function)
- Are anonymous (can be bound to variables though)
- Can be recursive

## Control Flow

- Conditionals
- Function calls and recursion
- No loops
- No mutation


## Operations/Builtin functions

- (int, float): +,-,*,/,%, <, <=, >, >=
- (int, float, bool, string): ==, !=
- (bool): !, &&, ||

## AST nodes/expression types

- Literal
	- int
	- float
	- bool
	- string
- Function
- SliceLiteral
	- int
	- float
	- bool
	- string
	- function
	- sliceLiteral
- Variable
	- Literal
	- Function
	- SliceLiteral
- Call
- Conditional
- Block (inccludes let/in bindigns)
