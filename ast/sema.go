package ast

import "fmt"

type analyzer struct {
	scopes  []binding
	paramID int
}

type binding struct {
	name string
	id   int
}

func (a *analyzer) nextID() int {
	a.paramID++
	return a.paramID
}

// peek is kept for convenience when making changes/degubbing, but currently unused
func (a analyzer) peek() binding {
	if len(a.scopes) > 0 {
		return a.scopes[len(a.scopes)-1]
	}
	return binding{}
}

func (a *analyzer) pop() {
	a.scopes = a.scopes[:len(a.scopes)-1]
}

func (a *analyzer) push(f LambdaE) {
	a.scopes = append(a.scopes, binding{name: f.ParamName, id: f.ParamId})
}

func (a *analyzer) lookVarUp(x VariableE) (id int, matched bool) {
	matched = false
	name := x.Name
	for i := len(a.scopes) - 1; i >= 0; i = i - 1 {
		if a.scopes[i].name == name {
			return a.scopes[i].id, true
		}
	}
	return 0, matched

}

func Resolve(e Expression) (Expression, error) {
	a := &analyzer{}
	return a.resolve(e)
}

func (a *analyzer) resolve(E Expression) (Expression, error) {
	switch e := E.(type) {
	case IntE:
		return e, nil
	case ApplicationE:
		resolvedArg, err := a.resolve(e.Argument)
		if err != nil {
			return nil, err
		}
		resolvedFun, err := a.resolve(e.Function)
		if err != nil {
			return nil, err
		}
		e.Function = resolvedFun
		e.Argument = resolvedArg
		return e, err
	case BuiltInFunc:
		return e, nil
	case LambdaE:
		id := a.nextID()
		e.ParamId = id
		a.push(e)
		var err error
		e.Body, err = a.resolve(e.Body)
		if err != nil {
			return nil, err
		}
		a.pop()
		return e, nil
	case VariableE:
		var ok bool
		e.ID, ok = a.lookVarUp(e)
		if ok == false {
			return nil, fmt.Errorf("unbound var %v", e)
		}
		return e, nil
	default:
		return nil, fmt.Errorf("ast contains expression of unkown type (%T) of unimplemented semantic analysis", E)
	}
}
