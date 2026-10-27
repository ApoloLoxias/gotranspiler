# Second Spec

## Types

### Primitives

- intE -> will later transpile to int64
- float64 [defered]
- bool [defered]
- string [defered]

### Composites

- Via algebraic data types (Sums and Products) [defered]
- Finite domain functions [defered]
- Composite types expose a "fundamental function" that represents them in the surface syntax, but isn't necessarily the memory/runtime/post-transpilation  representation
  e.g. a slice [1,2] could be f :: int -> Sum(int, bool); f = i: if (i == 0) (1) (if (i==1) (2) ) (false)

## Expressions

- Only expressions. No statements.
- All expressions are evaluatable to a single value of the formerly declared types
- A file contains a single expression whose value is exported to a variable in the transpiled go code
- A primitive type T is just a function of type (any -> T) whose elements t are constant functions that return themselves (t = x: {t}) 

## Variables and  Scoping

- Lexical scoping
- Let bindings [achieved via  immediatly invoked functions and syntax sugar]
- No mutations

## Types
- A type is a function T :: any -> bool, that returns true for elements that it contains and false for elements it doesnt contain (effectively the set of elements of that type)
- Therefore 'type type' evaluates to true
- Also, 'type' is a type, but not a primitive type

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

- for now, only for integers: +-/*! and **
- appending '$' to the operator token exposes its function (i.e. ordinary "prefix"/juxtaposition application semantics)
- allow partial application of infixes
  e.g. '1+' = 'x: 1+x'

## "Cannonical" function syntax
- 'a b' application by juxtaposition
- 'parameter: {expression}' to define the function
- Multi-argument/currying like 'x:{y: {expr}}'

## Syntactical sugar for functions
- 'param: body', 'param {body}' desugar into 'param: {body}' (but perhaps only for easy to parse cases)
- 'param1: param2: {body}' and 'param1: param2: body' desugar into 'param1: { param2: {body} }'
- perhaps also allow things like 'param1, param2: {body}'; 'param1, param2 {body}' and 'param1, param2: body' for the curried form, if they aren't confusing with chosen product and sum syntax

