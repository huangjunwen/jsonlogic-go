package jsonlogic

import (
	"fmt"
)

// AddOpLessThan adds "<" operation to the JSONLogic instance. Param restriction:
//   - At least two params.
//   - Must be evaluated to json primitives.
//   - If comparing numerics, then params must be able to convert to numeric. (See toNumeric)
func AddOpLessThan(jl *JSONLogic) {
	jl.AddOperation(string(lt), opCompare(lt))
}

// AddOpLessEqual adds "<=" operation to the JSONLogic instance. Param restriction: the same as "<".
func AddOpLessEqual(jl *JSONLogic) {
	jl.AddOperation(string(le), opCompare(le))
}

// AddOpGreaterThan adds ">" operation to the JSONLogic instance. Param restriction: the same as "<".
func AddOpGreaterThan(jl *JSONLogic) {
	jl.AddOperation(string(gt), opCompare(gt))
}

// AddOpGreaterEqual adds ">=" operation to the JSONLogic instance. Param restriction: the same as "<".
func AddOpGreaterEqual(jl *JSONLogic) {
	jl.AddOperation(string(ge), opCompare(ge))
}

// ref:
//   - https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Less_than
func opCompare(symbol compSymbol) Operation {
	return func(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
		if len(params) < 2 {
			return nil, fmt.Errorf("%s: expect at least two params", symbol)
		}
		params, err = applyParams(apply, params, data)
		if err != nil {
			return
		}

		r0, err := compareValues(symbol, params[0], params[1])
		if err != nil {
			return nil, fmt.Errorf("%s: %s", symbol, err.Error())
		}

		var r1 = true
		if len(params) > 2 {
			r1, err = compareValues(symbol, params[1], params[2])
			if err != nil {
				return nil, fmt.Errorf("%s: %s", symbol, err.Error())
			}
		}

		return r0 && r1, nil
	}
}

// AddOpMin adds "min" operation to the JSONLogic instance. Param restriction:
//   - Must be evaluated to json primitives that can convert to numeric.
func AddOpMin(jl *JSONLogic) {
	jl.AddOperation("min", opMin)
}

func opMin(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	for _, param := range params {
		r, err := apply(param, data)
		if err != nil {
			return nil, err
		}

		n, err := toNumeric(r)
		if err != nil {
			return nil, err
		}

		if res == nil || res.(float64) > n {
			res = n
		}
	}
	return
}

// AddOpMax adds "max" operation to the JSONLogic instance. Param restriction: the same as "and".
func AddOpMax(jl *JSONLogic) {
	jl.AddOperation("max", opMax)
}

func opMax(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	for _, param := range params {
		r, err := apply(param, data)
		if err != nil {
			return nil, err
		}

		n, err := toNumeric(r)
		if err != nil {
			return nil, err
		}

		if res == nil || res.(float64) < n {
			res = n
		}
	}
	return
}

// AddOpAdd adds "+" operation to the JSONLogic instance. Param restriction:
//   - Must be evaluated to json primitives that can convert to numeric.
func AddOpAdd(jl *JSONLogic) {
	jl.AddOperation("+", opAdd)
}

func opAdd(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	sum := float64(0)
	for _, param := range params {
		r, err := apply(param, data)
		if err != nil {
			return nil, err
		}

		n, err := toNumeric(r)
		if err != nil {
			return nil, err
		}
		sum += n
	}
	return sum, nil
}

// AddOpMul adds "*" operation to the JSONLogic instance. Param restriction:
//   - At least one param.
//   - Must be evaluated to json primitives that can convert to numeric.
func AddOpMul(jl *JSONLogic) {
	jl.AddOperation("*", opMul)
}

func opMul(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 1 {
		return nil, fmt.Errorf("*: expect at least one param")
	}
	prod := float64(1)
	for _, param := range params {
		r, err := apply(param, data)
		if err != nil {
			return nil, err
		}

		n, err := toNumeric(r)
		if err != nil {
			return nil, err
		}
		prod *= n
	}
	return prod, nil
}
