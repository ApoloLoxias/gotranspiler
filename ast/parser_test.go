package ast

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ApoloLoxias/gotranspiler/lex"
)

func lexString(t *testing.T, source string) ([]lex.Token, error) {
	t.Helper()

	if len(source) == 0 {
		return nil, fmt.Errorf("no input to lex")
	}

	tokens := lex.Lex(source)
	if tokens == nil || len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens emitted")
	}

	return tokens, nil
}

func parseTokens(t *testing.T, tokens []lex.Token) (Expression, error) {
	t.Helper()

	if tokens == nil || len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens to parse")
	}

	ast, parsingErr := Parse(tokens)

	if ast == nil && parsingErr != nil {
		return nil, fmt.Errorf("couldn't create an ast, despite no parsingError")
	}

	return ast, parsingErr
}

func evalAst(t *testing.T, ast Expression) (Expression, error) {
	t.Helper()

	if ast == nil {
		return nil, fmt.Errorf("No ast to evaluate")
	}

	value := ast.Evaluate()

	if value == nil {
		return nil, fmt.Errorf("nil evaluation of %v", ast)
	}

	return ast, nil
}

func evalString(t *testing.T, source string) (
	tokens []lex.Token,
	ast Expression,
	value Expression,
	err error,
) {
	t.Helper()

	tokens, err = lexString(t, source)
	if err != nil {
		return tokens, nil, nil, err
	}

	ast, err = parseTokens(t, tokens)
	if err != nil {
		return tokens, ast, nil, err
	}

	value, err = evalAst(t, ast)
	return tokens, ast, value, err
}

func assertInteger(t *testing.T, GOT Expression, want int) {
	t.Helper()

	switch got := GOT.(type) {
	case IntE:
		if got.Value != want {
			t.Errorf("wrong value: %d, want: %d", got, want)
		}
	default:
		t.Errorf("want IntE(%d), got %t: %v", want, got, got)
	}
}

func assertBuiltInPartial(t *testing.T, GOT Expression, want string) {
	t.Helper()

	switch got := GOT.(type) {
	case BuiltInFunc:
		if got.Name != want {
			t.Errorf("want %s, got %s", got.Name, want)
		}
	default:
		t.Errorf("want BuiltInFunc(%s), got %v", want, got)
	}
}
