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
		if isTrue(r) {
			return apply(params[i+1], data)
		}
	}

	if len(params) == i+1 {
		return apply(params[i], data)
	}

	return nil, nil
}

// AddOpStrictEqual adds "===" operation to the JSONLogic instance.
func AddOpStrictEqual(jl *JSONLogic) {
	jl.AddOperation("===", opStrictEqual)
}

func opStrictEqual(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {

	var leftObj, rightObj interface{}

	switch len(params) {
	default:
		fallthrough

	case 2:
		leftObj, err = apply(params[0], data)
		if err != nil {
			return
		}
		rightObj, err = apply(params[1], data)
		if err != nil {
			return
		}

		if isPrimitive(leftObj) && isPrimitive(rightObj) {
			return leftObj == rightObj, nil
		}

		// XXX: This is different with jsonlogic-js
		return nil, fmt.Errorf("===/!==: params should be primitives")

	case 1:
		return false, nil

	case 0:
		return true, nil
	}

}

// AddOpStrictNotEqual adds "!==" operation to the JSONLogic instance.
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

// AddOpNegative adds "!" operation to the JSONLogic instance.
func AddOpNegative(jl *JSONLogic) {
	jl.AddOperation("!", opNegative)
}

func opNegative(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	var param interface{}
	if len(params) > 0 {
		param = params[0]
	}
	res, err = apply(param, data)
	if err != nil {
		return
	}
	return !isTrue(res), nil
}

// AddOpDoubleNegative adds "!!" operation to the JSONLogic instance.
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

// AddOpAnd adds "and" operation to the JSONLogic instance.
func AddOpAnd(jl *JSONLogic) {
	jl.AddOperation("and", opAnd)
}

func opAnd(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) == 0 {
		return nil, fmt.Errorf("and: expect at least one params")
	}
	for _, param := range params {
		res, err = apply(param, data)
		if err != nil {
			return
		}
		if !isTrue(res) {
			return res, nil
		}
	}
	return
}

// AddOpOr adds "or" operation to the JSONLogic instance.
func AddOpOr(jl *JSONLogic) {
	jl.AddOperation("or", opOr)
}

func opOr(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) == 0 {
		return nil, fmt.Errorf("or: expect at least one params")
	}
	for _, param := range params {
		res, err = apply(param, data)
		if err != nil {
			return
		}
		if isTrue(res) {
			return res, nil
		}
	}
	return
}
