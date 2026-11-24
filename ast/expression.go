package ast

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Fundamental AST nodes
type Expression interface {
	//Evaluate reduces the ast to an equivalent, simplified one
	Evaluate(map[int]Expression) Expression

	String() string
	pretty(int, int) string //helper for Pretty
	Pretty(...int) string   //variadic as a trick for allowing a default value

	Type() TypeE //Expressions have types, which are expressions
	SetType(TypeE) (Expression, error)
}

type assignment struct{ id int }

// pretty-printing configuration
const defaultTabStop = 8

var TypeErr = errors.New("Invalid Type assignment")

/* --- Primitive Expressions ---*/
// A primitive is interpreted as a constant function that allways returns itself

// Integer literals
type IntE struct {
	Value int
}

func (i IntE) Evaluate(scope map[int]Expression) Expression {
	return i
}

func (i IntE) String() string {
	return fmt.Sprintf("IntE(%d)", i.Value)
}

func (i IntE) pretty(int, int) string {
	return i.String()
}

func (i IntE) Pretty(tabSpace ...int) string {
	tab := defaultTabStop
	if len(tabSpace) > 0 {
		tab = tabSpace[0]
	}
	return i.pretty(tab, 0)
}

func (i IntE) Type() TypeE {
	return TypeE{
		Simple: TypeINT,
	}
}

func (i IntE) SetType(t TypeE) (Expression, error) {
	if t.Simple != TypeINT {
		return i, TypeErr
	}
	return i, nil
}

// Function application
// includes arithmetic expressions, be they prefix, infix or postfix operations
type ApplicationE struct {
	Function Expression
	Argument Expression

	T TypeE
}

func (a ApplicationE) Evaluate(scope map[int]Expression) Expression { //TODO
	switch f := a.Function.(type) {
	case BuiltInFunc:
		result, _ := f.Func(a.Argument.Evaluate(scope))
		return result
	case ApplicationE:
		fEval := f.Evaluate(scope)
		argEval := a.Argument.Evaluate(scope)
		appEval := ApplicationE{Function: fEval, Argument: argEval}
		return appEval.Evaluate(scope)
	case IntE:
		return f
	case VariableE:
		return f //probably a type error when implementing types?A variable could be a legitimate type and,thus a function...
	case LambdaE:
		scope[f.ParamId] = a.Argument.Evaluate(scope)
		// Argument evaluated before expanding scope (x: {(x: {f x}) 2}) 1 =  f 2
		return f.Body.Evaluate(scope)
	default:
		break
	}
	return nil
}

func (a ApplicationE) String() string {
	return fmt.Sprintf("ApplicationE{Function: %s, Argument: %s, Type: %s}", a.Function, a.Argument, a.T)
}

func (a ApplicationE) pretty(tabSpace, lvl int) string {
	tab := strings.Repeat(" ", tabSpace)
	indent := strings.Repeat(tab, lvl)

	return fmt.Sprintf("ApplicationE{\n%s%sFunction: %s,\n%s%sArgument: %s,\n%s%sType: %s\n}", indent, tab, a.Function.pretty(tabSpace, lvl+1), indent, tab, a.Argument.pretty(tabSpace, lvl+1), indent, tab, a.T.pretty(tabSpace, lvl+1))
}

func (a ApplicationE) Pretty(tabsSpace ...int) string {
	tab := defaultTabStop
	if len(tabsSpace) > 0 {
		tab = tabsSpace[0]
	}
	return a.pretty(tab, 0)
}

func (a ApplicationE) Type() TypeE {
	return a.T
}

func (a ApplicationE) SetType(t TypeE) (Expression, error) {
	if (a.T.Simple == "" && a.T.Compound == nil) || a.T.Simple == TypeNO_TYPE {
		a.T = t
		return a, nil
	}
	if !a.T.EqualTo(t) {
		return a, TypeErr
	}
	return a, nil
}

// Native/standard/pre-defined functions
// Include arithmetic operators, be they infix, suffix or prefix
type BuiltInFunc struct {
	Name string //for human-friendly display
	Func func(Expression) (Expression, error)
}

func (f BuiltInFunc) Evaluate(scope map[int]Expression) Expression {
	return f
}

func (f BuiltInFunc) String() string {
	return fmt.Sprintf("BuiltInFunc(%s)", f.Name)
}

func (f BuiltInFunc) pretty(int, int) string {
	return f.String()
}

func (f BuiltInFunc) Pretty(tabSpace ...int) string {
	tab := defaultTabStop
	if len(tabSpace) > 0 {
		tab = tabSpace[0]
	}
	return f.pretty(tab, 0)
}

func (f BuiltInFunc) Type() TypeE {
	return TypeE{
		Simple:   TypeCOMPOUND,
		Compound: []TypeE{TypeE{Simple: TypeINT}, TypeE{Simple: TypeINT}},
	}
}

func (f BuiltInFunc) SetType(t TypeE) (Expression, error) {
	if t.EqualTo(TypeE{
		Simple: TypeCOMPOUND,
		Compound: []TypeE{TypeE{
			Simple: TypeINT,
		}, TypeE{
			Simple: TypeINT,
		},
		}}) {
		return f, nil
	}
	return f, TypeErr
}

/* --- BUILTIN OPERATORS --- */
// Binary and infix operators represent curried fucntions
// Those functions may be partially applied

var Sum = BuiltInFunc{Name: "Sum", Func: SumBuiltIn} //const
func SumBuiltIn(x Expression) (Expression, error) {
	switch X := x.(type) {
	case IntE:
		f := BuiltInFunc{
			Name: fmt.Sprintf("Sum%d", X.Value),
			Func: func(y Expression) (Expression, error) {
				switch Y := y.(type) {
				case IntE:
					return IntE{Value: X.Value + Y.Value}, nil
				default:
					return nil, errors.New("curried Sum typeError")
				}
			},
		}
		return f, nil
	default:
		return nil, errors.New("Sum typeError")
	}
}

var Sub = BuiltInFunc{Name: "Subtraction", Func: SubBuiltIn} //const
func SubBuiltIn(x Expression) (Expression, error) {
	switch X := x.(type) {
	case IntE:
		f := BuiltInFunc{
			Name: fmt.Sprintf("Sub%d", X.Value),
			Func: func(y Expression) (Expression, error) {
				switch Y := y.(type) {
				case IntE:
					return IntE{Value: X.Value - Y.Value}, nil
				default:
					return nil, errors.New("curried Sub typeError")
				}
			},
		}
		return f, nil
	default:
		return nil, errors.New("Sub typeError")
	}
}

var Mul = BuiltInFunc{Name: "Multiplication", Func: MulBuiltIn} //const
func MulBuiltIn(x Expression) (Expression, error) {
	switch X := x.(type) {
	case IntE:
		f := BuiltInFunc{
			Name: fmt.Sprintf("Mul%d", X.Value),
			Func: func(y Expression) (Expression, error) {
				switch Y := y.(type) {
				case IntE:
					return IntE{Value: X.Value * Y.Value}, nil
				default:
					return nil, errors.New("curried Mul typeError")
				}
			},
		}
		return f, nil
	default:
		return nil, errors.New("Mul typeError")
	}
}

var Div = BuiltInFunc{Name: "Division", Func: DivBuiltIn} // const
func DivBuiltIn(x Expression) (Expression, error) {
	switch X := x.(type) {
	case IntE:
		f := BuiltInFunc{
			Name: fmt.Sprintf("Div%d", X.Value),
			Func: func(y Expression) (Expression, error) {
				switch Y := y.(type) {
				case IntE:
					return IntE{Value: X.Value / Y.Value}, nil
				default:
					return nil, errors.New("curried Div typeError")
				}
			},
		}
		return f, nil
	default:
		return nil, errors.New("Div typeError")
	}
}

var Neg = BuiltInFunc{Name: "Unary negative", Func: NegBuiltIn} //const
func NegBuiltIn(x Expression) (Expression, error) {
	switch X := x.(type) {
	case IntE:
		return IntE{-X.Value}, nil
	default:
		return nil, errors.New("Neg typeError")
	}
}

var Pow = BuiltInFunc{Name: "Integer potentiation", Func: PowBuiltIn} //const
func PowBuiltIn(x Expression) (Expression, error) {
	switch X := x.(type) {
	case IntE:
		f := BuiltInFunc{
			Name: fmt.Sprintf("Pow%d", X.Value),
			Func: func(y Expression) (Expression, error) {
				switch Y := y.(type) {
				case IntE:
					return IntE{Value: intPow(X.Value, Y.Value)}, nil
				default:
					return nil, errors.New("curried Pow typeErro")
				}
			},
		}
		return f, nil
	default:
		return nil, errors.New("Pow typeError")
	}
}
func intPow(b, e int) int { //consider making this safe and eficient later
	B := float64(b)
	E := float64(e)
	BE := math.Pow(B, E)
	return int(BE)
}

var Fac = BuiltInFunc{Name: "Integer factorial", Func: FacBuiltIn} //const
func FacBuiltIn(x Expression) (Expression, error) {
	switch X := x.(type) {
	case IntE:
		if X.Value < 0 {
			return nil, errors.New("Factorial of negative number error")
		}
		f := func(y int) int {
			z := y
			for {
				if y == 0 {
					return 1
				}
				if y == 1 {
					return z
				}
				y = y - 1
				z = z * y
			}
		}
		return IntE{f(X.Value)}, nil
	default:
		return nil, errors.New("Fac typeError")
	}

}

/* --- For Declaring Functions and Binding Variables --- */

type LambdaE struct {
	ParamName string //metadata about the source code
	ParamId   int
	Body      Expression
	T         TypeE
}

func (l LambdaE) Evaluate(scope map[int]Expression) Expression {
	return LambdaE{
		ParamName: l.ParamName,
		ParamId:   l.ParamId,
		Body:      l.Body.Evaluate(scope),
	}
}

func (l LambdaE) String() string {
	return fmt.Sprintf("LambdaE{ParamName: %s, ParamId: %d, Body: %v, Type: %s}", l.ParamName, l.ParamId, l.Body, l.T)
}

func (l LambdaE) pretty(tabSpace, lvl int) string {
	tab := strings.Repeat(" ", tabSpace)
	indent := strings.Repeat(tab, lvl)

	return fmt.Sprintf("LambdaE{\n%s%sParamName %s,\n%s%sParamId %d,\n%s%sBody %v,\n%s%sType: %s\n}", indent, tab, l.ParamName, indent, tab, l.ParamId, indent, tab, l.Body.pretty(tabSpace, lvl+1), indent, tab, l.T.pretty(tabSpace, lvl+1))
}

func (l LambdaE) Pretty(tabSpace ...int) string {
	tab := defaultTabStop
	if len(tabSpace) > 0 {
		tab = tabSpace[0]
	}
	return l.pretty(tab, 0)
}

func (l LambdaE) Type() TypeE {
	return l.T
}

func (l LambdaE) SetType(t TypeE) (Expression, error) {
	if l.T.Simple == TypeNO_TYPE || (l.T.Simple == "" && l.T.Compound == nil) {
		l.T = t
		return l, nil
	}
	if l.T.EqualTo(t) {
		return l, nil
	}
	return l, TypeErr
}

/* --- */

type VariableE struct {
	Name string
	ID   int
	T    TypeE
}

func (v VariableE) Evaluate(scope map[int]Expression) Expression {
	value, ok := scope[v.ID]
	if ok {
		return value
	} else {
		return v
	}
}

func (v VariableE) String() string {
	return fmt.Sprintf("VariableE{Name: %s; ID: %d; Type: %s}", v.Name, v.ID, v.T)
}

func (v VariableE) pretty(tabSpace, lvl int) string {
	return fmt.Sprintf("VariableE{Name: %s; ID: %d; Type: %s}", v.Name, v.ID, v.T)
}

func (v VariableE) Pretty(tabSpace ...int) string {
	tab := defaultTabStop
	if len(tabSpace) > 0 {
		tab = tabSpace[0]
	}
	return v.pretty(tab, 0)
}

func (v VariableE) Type() TypeE {
	return v.T
}

func (v VariableE) SetType(t TypeE) (Expression, error) {
	if v.T.Simple == TypeNO_TYPE {
		v.T = t
		return v, nil
	}
	if v.T.Simple == "" && v.T.Compound == nil {
		v.T = t
		return v, nil
	}
	if v.T.EqualTo(t) {
		return v, nil
	}
	return v, TypeErr
}

/* --- */

type AppliedLambdaE struct {
	LambdaE
	Argument Expression
}

func (a AppliedLambdaE) Evaluate(scope map[int]Expression) Expression {
	desugared := ApplicationE{
		Function: LambdaE{
			ParamName: a.ParamName,
			ParamId:   a.ParamId,
			Body:      a.Body,
			T:         a.T,
		},
		Argument: a.Argument,
	}

	return desugared.Evaluate(scope)
}

func (a AppliedLambdaE) String() string {
	return fmt.Sprintf("AppLambdaE{ParamName: '%s', ParamId: %d, Body: %s, Argument: %s, Type: %s}", a.ParamName, a.ParamId, a.Body, a.Argument, a.T)
}

func (a AppliedLambdaE) pretty(tabSpace, lvl int) string {
	tab := strings.Repeat(" ", tabSpace)
	indent := strings.Repeat(tab, lvl)

	return fmt.Sprintf(
		"AppLambdaE{\n%s%sParamName: %s,\n%s%sParamId: %d,\n%s%sBody: %s,\n%s%sArgument: %s,\n%s$%sType: %s\n}",
		indent, tab, a.ParamName,
		indent, tab, a.ParamId,
		indent, tab, a.Body.pretty(tabSpace, lvl+1),
		indent, tab, a.Argument.pretty(tabSpace, lvl+1),
		indent, tab, a.T.pretty(tabSpace, lvl+1),
	)
}

func (a AppliedLambdaE) Pretty(tabSpace ...int) string {
	tab := defaultTabStop
	if len(tabSpace) > 0 {
		tab = tabSpace[0]
	}
	return a.pretty(tab, 0)
}

func (a AppliedLambdaE) Type() TypeE {
	return a.T
}

func (a AppliedLambdaE) SetType(t TypeE) (Expression, error) {
	if a.T.Simple == TypeNO_TYPE || (a.T.Simple == "" && a.T.Compound == nil) {
		a.T = t
		return a, nil
	}
	if a.T.EqualTo(t) {
		return a, nil
	}
	return a, TypeErr
}

/* --- */
type TypeE struct {
	Compound []TypeE
	Simple   SimpleType
}

type SimpleType string

const (
	TypeNO_TYPE  SimpleType = "untyped"
	TypeCOMPOUND SimpleType = "notASimpleType"
	TypeTYPE     SimpleType = "typeType"
	TypeINT      SimpleType = "intgerType"
)

func (t TypeE) Evaluate(scope map[int]Expression) Expression {
	return t
}

func (t TypeE) String() string {
	if t.Simple == TypeCOMPOUND {
		return fmt.Sprintf("TypeE{%s}", t.Compound)
	}
	return fmt.Sprintf("TypeE{%s}", t.Simple)
}

func (t TypeE) pretty(tabSpace, lvl int) string {
	//tab := strings.Repeat(" ", tabSpace)
	//indent := strings.Repeat(tab, lvl)

	if t.Compound == nil {
		return fmt.Sprintf("%s", t.Simple)
	}
	return fmt.Sprintf("%s", emprettyCompoundType(t.Compound))
}
func emprettyCompoundType(tt []TypeE) string {
	var result string
	for i, t := range tt {
		if t.Compound == nil {
			result += string(t.Simple)
			if i < len(tt) {
				result += "->"
			}
		} else {
			result += "(" + emprettyCompoundType(t.Compound) + ")"
			if i < len(tt) {
				result += "->"
			}
		}
	}
	return result
}

func (t TypeE) Pretty(tabsSpace ...int) string {
	tab := defaultTabStop
	if len(tabsSpace) > 0 {
		tab = tabsSpace[0]
	}
	return t.pretty(tab, 0)
}

func (t TypeE) Type() TypeE {
	return TypeE{
		Simple: TypeTYPE,
	}
}

func (t TypeE) EqualTo(T TypeE) bool {
	if t.Simple != T.Simple {
		return false
	}
	for i, ti := range t.Compound {
		if !ti.EqualTo(T.Compound[i]) {
			return false
		}
	}
	return true
}

func (t TypeE) SetType(T TypeE) (Expression, error) {
	if T.Simple == TypeTYPE {
		return t, nil
	}
	return t, TypeErr
}
