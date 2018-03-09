package fsm

import (
	"errors"
	"strconv"

	"github.com/kawamuray/jsonpath"
)

// Operator ...
type Operator interface {
	Evaluate(string) (bool, error)
}

// StringEqualsOperator ...
type StringEqualsOperator struct {
	comp string
}

// Evaluate ...
func (se StringEqualsOperator) Evaluate(data string) (bool, error) {
	return data == se.comp, nil
}

// NumberEqualsOperator ...
type NumberEqualsOperator struct {
	comp float64
}

// Evaluate ...
func (ne NumberEqualsOperator) Evaluate(data string) (bool, error) {
	f, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return false, err
	}
	return f == ne.comp, nil
}

// Choice ...
type Choice interface {
	Evaluate(*Execution) (bool, error)
}

// AndOperator ...
type AndOperator struct {
	Choices []BaseChoice
}

// Evaluate ...
func (op AndOperator) Evaluate(exec *Execution) (bool, error) {
	for _, choice := range op.Choices {
		if ok, err := choice.Evaluate(exec); !ok || err != nil {
			return false, err
		}
	}
	return true, nil
}

// OrOperator ...
type OrOperator struct {
	Choices []BaseChoice
}

// Evaluate ...
func (op OrOperator) Evaluate(exec *Execution) (bool, error) {
	result := 0
	for _, choice := range op.Choices {
		ok, err := choice.Evaluate(exec)
		if err != nil {
			return false, err
		}
		if ok {
			result++
		}
	}
	if result == 0 {
		return false, nil
	}
	return true, nil
}

// NotOperator ...
type NotOperator struct {
	Choice BaseChoice
}

// Evaluate ...
func (op NotOperator) Evaluate(exec *Execution) (bool, error) {
	result, err := op.Choice.Evaluate(exec)
	if err != nil {
		return false, err
	}
	return !result, nil
}

// BaseChoice ...
type BaseChoice struct {
	Variable string // string path to variable
	Op       Operator
}

// Evaluate ...
func (op BaseChoice) Evaluate(exec *Execution) (bool, error) {
	// TODO: save paths instead if Variable
	paths, err := jsonpath.ParsePaths(op.Variable)
	if err != nil {
		return false, err
	}
	eval, err := jsonpath.EvalPathsInBytes(exec.data, paths)
	if err != nil {
		return false, err
	}

	// TODO: check if there is only one result
	result, ok := eval.Next()
	if !ok { // no match
		return false, errors.New("json path - no match")
	}

	return op.Op.Evaluate(string(result.Value))
}
