# Runtime Values

## Primitive golang values
- int
- float
- bool
- string

## Golang Slices
- []primitive
- []closure //for slices of functions
- [][]

## Environment
- Holds runtime values assigned to variables

## Closures
- Group expressions/AST nodes with bound arguments in a given runtime environment


# Expression Evaluation

## PrimitiveLiterals
- IntLiteral evaluates to an integer value
- FloatLiteral evaluates to a float value
- BoolLiteral evaluates to a bool value
- StringLiteral Evaluates to a string value

## FunctionLiteral
- Evaluates to a closure capturing the current environment

## SliceLiteral
- Evaluates to a slice value

## Variable
- Evaluates to its bound value

## Call
- Evaluates the function expression
- Evaluates the argument expression
- Applies the argument to the function
- Evaluates to the application

## Condition
- Evaluates condition
- Evaluates to Yes or No branch accordingly

## Block
- Creates a new lexical environment
- Evaluates bindings (in order)
- Evaluates to the value of the final expression in the enviornment
- Wants special eval ordering:
    - If Assess is a Block, evaluate Assess first
    - If Assess is not a Block, evaluate Parent first
