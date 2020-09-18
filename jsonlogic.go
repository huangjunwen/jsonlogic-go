package jsonlogic

import (
	"fmt"
)

type JSONLogic struct {
	ops map[string]Operation
}

type Applier func(logic, data interface{}) (res interface{}, err error)

type Operation func(apply Applier, params []interface{}, data interface{}) (interface{}, error)

func NewJSONLogic() *JSONLogic {
	ret := NewEmptyJSONLogic()
	AddOpVar(ret)
	return ret
}

func NewEmptyJSONLogic() *JSONLogic {
	return &JSONLogic{
		ops: make(map[string]Operation),
	}
}

// Apply data to logic and returns a result. Both logic/data must be one of 'encoding/json' supported types:
//   - nil
//   - bool
//   - float64
//   - string
//   - []interface{} with items of supported types
//   - map[string]interface{} with values of supported types
func (jl *JSONLogic) Apply(logic, data interface{}) (res interface{}, err error) {
	// An array of rules.
	if arr, ok := logic.([]interface{}); ok {
		ret := []interface{}{}
		for _, item := range arr {
			res, err := jl.Apply(item, data)
			if err != nil {
				return nil, err
			}
			ret = append(ret, res)
		}
		return ret, nil
	}

	// Primitive.
	if !isLogic(logic) {
		return logic, nil
	}

	if data == nil {
		data = map[string]interface{}{}
	}

	op, params := getLogic(logic)

	defer func() {
		if e := recover(); e != nil {
			var ok bool
			err, ok = e.(error)
			if !ok {
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	opFn := jl.ops[op]
	if opFn == nil {
		return nil, fmt.Errorf("Apply: operator %q not found", op)
	}

	return opFn(jl.Apply, params, data)
}

func (jl *JSONLogic) AddOperation(name string, op Operation) {
	jl.ops[name] = op
}
