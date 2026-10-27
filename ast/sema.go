package ast

import "fmt"

type analyzer struct {
	ast     Expression
	scopes  []map[string]int
	paramID int
}

func (a *analyzer) nextID() int {
	a.paramID++
	return a.paramID
}

func Resolve(e Expression) (Expression, error) {
	return resolve(e, nil, 0)
}

func resolve(e Expression, s []map[string]int, i int) (Expression, error) {
	if s == nil {
		s = make([]map[string]int, 0)
	}
	a := analyzer{
		ast:     e,
		scopes:  s,
		paramID: i,
	}

	err := a.resolve()

	return a.ast, err
}

func (a *analyzer) resolve() error {
	switch e := a.ast.(type) {
	case IntE:
		return nil
	case ApplicationE:
		resolvedArg, err := resolve(e.Argument, a.scopes, a.paramID)
		if err != nil {
			return err
		}
		resolvedFun, err := resolve(e.Function, a.scopes, a.paramID)
		if err != nil {
			return err
		}
		a.ast = ApplicationE{Function: resolvedFun, Argument: resolvedArg}
		return nil
	case BuiltInFunc:
		return nil
	case LambdaE:
		id := a.nextID()
		e.ParamId = id
		a.ast = e
		a.scopes = append(a.scopes, map[string]int{e.ParamName: id})
		var err error
		e.Body, err = resolve(e.Body, a.scopes, a.paramID)
		if err != nil {
			return err
		}
		a.ast = e
		a.scopes = a.scopes[:len(a.scopes)-1]
		return nil
	case VariableE:
		for _, m := range a.scopes {
			id, ok := m[e.Name]
			if ok {
				e.ID = id
				a.ast = e
				break
			}
		}
		if e.ID == 0 {
			return fmt.Errorf("unbound variable %v", e.Pretty())
		} else {
			return nil
		}
	default:
		return fmt.Errorf("ast contains expression of unkown type of unimplemented semantic analysis")
	}
}
