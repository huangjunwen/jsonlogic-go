package jsonlogic

import (
	"fmt"
)

// AddOpLessThan adds "<" operation to the JSONLogic instance.
func AddOpLessThan(jl *JSONLogic) {
	jl.AddOperation(string(lt), opCompare(lt))
}

// AddOpLessEqual adds "<=" operation to the JSONLogic instance.
func AddOpLessEqual(jl *JSONLogic) {
	jl.AddOperation(string(le), opCompare(le))
}

// AddOpGreaterThan adds ">" operation to the JSONLogic instance.
func AddOpGreaterThan(jl *JSONLogic) {
	jl.AddOperation(string(gt), opCompare(gt))
}

// AddOpGreaterEqual adds ">=" operation to the JSONLogic instance.
func AddOpGreaterEqual(jl *JSONLogic) {
	jl.AddOperation(string(ge), opCompare(ge))
}

// opCompare returns false if len(params) < 2, since compare undefined with any returns false.
//
// ref:
//   - https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Less_than
func opCompare(symbol compSymbol) Operation {
	return func(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
		if len(params) < 2 {
			return false, nil
		}
		var (
			params0, params1, params2 interface{}
			between                   bool
		)
		switch len(params) {
		default:
			fallthrough
		case 3:
			between = true
			params2, err = apply(params[2], data)
			if err != nil {
				return
			}
			fallthrough
		case 2:
			params1, err = apply(params[1], data)
			if err != nil {
				return
			}
			params0, err = apply(params[0], data)
			if err != nil {
				return
			}
		}

		r0, err := compareValues(symbol, params0, params1)
		if err != nil {
			return nil, fmt.Errorf("%s: %s", symbol, err.Error())
		}

		var r1 = true
		if between {
			switch symbol {
			case lt, le:
				r1, err = compareValues(symbol, params1, params2)
				if err != nil {
					return nil, fmt.Errorf("%s: %s", symbol, err.Error())
				}
			}
		}

		return r0 && r1, nil
	}
}
