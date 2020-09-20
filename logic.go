package jsonlogic

import (
	"fmt"
)

// AddOpIf adds "if"/"?:" operation to the JSONLogic instance.
func AddOpIf(jl *JSONLogic) {
	jl.AddOperation("if", opIf)
	jl.AddOperation("?:", opIf)
}

func opIf(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	var i int
	for i = 0; i < len(params)-1; i += 2 {
		r, err := apply(params[i], data)
		if err != nil {
			return nil, err
		}
		if toBool(r) {
			return apply(params[i+1], data)
		}
	}

	if len(params) == i+1 {
		return apply(params[i], data)
	}

	return nil, nil
}

// AddOpStrictEqual adds "===" operation to the JSONLogic instance. Param restriction:
//   - At least two params.
//   - Params must be evaluated to json primitives.
func AddOpStrictEqual(jl *JSONLogic) {
	jl.AddOperation("===", opStrictEqual)
}

func opStrictEqual(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("===/!==: expect at least two params")
	}
	params, err = applyParams(apply, params, data)
	if err != nil {
		return
	}

	if !isPrimitive(params[0]) || !isPrimitive(params[1]) {
		return nil, fmt.Errorf("===/!==: params must be json primitives")
	}
	return params[0] == params[1], nil
}

// AddOpStrictNotEqual adds "!==" operation to the JSONLogic instance. Param restriction: the same as "===".
func AddOpStrictNotEqual(jl *JSONLogic) {
	jl.AddOperation("!==", opStrictNotEqual)
}

func opStrictNotEqual(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	r, err := opStrictEqual(apply, params, data)
	if err != nil {
		return
	}
	return !r.(bool), nil
}

// AddOpNegative adds "!" operation to the JSONLogic instance. Param restriction:
//   - At least one param.
func AddOpNegative(jl *JSONLogic) {
	jl.AddOperation("!", opNegative)
}

func opNegative(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 1 {
		return nil, fmt.Errorf("!/!!: expect at least one param")
	}
	res, err = apply(params[0], data)
	if err != nil {
		return
	}
	return !toBool(res), nil
}

// AddOpDoubleNegative adds "!!" operation to the JSONLogic instance. Param Restriction: the same as "!".
func AddOpDoubleNegative(jl *JSONLogic) {
	jl.AddOperation("!!", opDoubleNegative)
}

func opDoubleNegative(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	r, err := opNegative(apply, params, data)
	if err != nil {
		return
	}
	return !r.(bool), nil
}

// AddOpAnd adds "and" operation to the JSONLogic instance. Param restriction:
//   - At least one param.
func AddOpAnd(jl *JSONLogic) {
	jl.AddOperation("and", opAnd)
}

func opAnd(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 1 {
		return nil, fmt.Errorf("and: expect at least one params")
	}
	for _, param := range params {
		res, err = apply(param, data)
		if err != nil {
			return
		}
		if !toBool(res) {
			return res, nil
		}
	}
	return
}

// AddOpOr adds "or" operation to the JSONLogic instance. Param restriction:
//   - At least one param.
func AddOpOr(jl *JSONLogic) {
	jl.AddOperation("or", opOr)
}

func opOr(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 1 {
		return nil, fmt.Errorf("or: expect at least one params")
	}
	for _, param := range params {
		res, err = apply(param, data)
		if err != nil {
			return
		}
		if toBool(res) {
			return res, nil
		}
	}
	return
}
