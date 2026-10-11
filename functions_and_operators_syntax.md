# Some notes

* 'f','g', ... will be often used as a shorthand for functions
* '+' will be often used as a shorthand for any infix operator
* '-' will be often used as a shorthand for any prefix operator
* '!' will be often used as a shorthand for any suffix operator
* 'n', 'm', ..., as well as 'x', 'y', ... will be often used as
  a shorthand for numerical values, primitives and operands of
  operators for clarity, though they aren't really distinct from
  functions at a language level

# Expressions

* Expressions are unary functions and n-arity is achieved via curried forms
* A primitive is a constant function that returns itself
* Functional application is done via juxtaposition: f arg = f(arg)
* Functional applicaiton is left-associative: f g h = (f g) h

# Operators

* Operators include symbols like +-*/!, and can be prefixes, infixes or suffixes

* They have special syntatictal treatment, but will eventually be intrepreted as a functional application

* Ordinary functional behaviour is exposed by appending and operator with $.
  In the future, consider preppending infixes, so we can differentite thigns like
  unary minus and subtraction in ordinary functional notation

* Infixes are binary curried functions and may be partially applied
  1 + 2 = Sum(1, 2) and 1 + = f such that f(x) = 1+x
  A "trailing infix" is, therefore, its own syntatictal category
  and distinct from a suffix

* Infixes are left-associative: 1 - 2 - 3 = -4 != 2

* Though prefix operators seem to share the funciton application juxtaposition form
  they still receive special syntactical treatment, not sharing all of the applicaiton
  rules

* Prefixes are right associative: - - n = - (- n)

* Suffixes are left associative ! ! = (n !) !

* Precedence of math operations is as expected:
    ! > application > ** (exponentiation) > *, / > +, -
    Note that unary minus is a prefix and follows the rules of prefixes,
    so its not just a precedence thing and the '-' on the former queu is
    binary

# Prefixes

* Though prefix operators seem to share the funciton application juxtaposition form
  they still receive special syntactical treatment, not sharing all of the applicaiton
  rules

* Prefixes operators get their argument by parsing everything
  up to (and excluding) an infix operator or (and including) a chain of suffix operators


# Desired behaviour and examples

* f g h = (f g) h
* - - n = - (- n)
* n ! ! = (n !) !
* n + m + o = (n + m) + o
* - n ! ! = - ((n !) !)
* - f x + y = - (f x) + y
* - x + = (- x) + = +$ (- x) // partial application of trailing infix
* - f x = - (f x)
* f - x = f (- x) but x - y is an infix?
