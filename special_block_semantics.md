# On interesting block properties

## Block Struct
- Bind: []Binding = []{Variable, Expression}

## Evaluation Order
- If Assess is a block, evaluate assess first
- Else, evaluate the parent block/make the lexical substitutions before evaluating Asess

## Programs and runtimes
- A program is implicitly a block with the builtins bound to the initial environment:
    ```
    Block{
        Bind: []Binding{{Variable{"+"}, ADD}, {Variable{"-"}, SUB}, ... }
        Assess: Expression described by the program or file
    }
    ```

## Oportunities
- Naturally resolves lexical binding with variable shadowing
- Metaprogramming mechanism for operator/builtin overloading and macros (if we allow/implement that)
- Mechanism for resolving functional i/o problem in a monadic fashion
- Mecahanism for ffi/golang execution escape hatch 
