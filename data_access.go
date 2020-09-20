package jsonlogic

import (
	"fmt"
	"strconv"
	"strings"
)

// AddOpVar adds "var" operation to the JSONLogic instance. Param restriction:
//   - At least one param (the key).
//   - Keys must be evaluated to json primitives.
func AddOpVar(jl *JSONLogic) {
	jl.AddOperation("var", opVar)
}

func opVar(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 1 {
		return nil, fmt.Errorf("var: expect at least one param")
	}
	params, err = applyParams(apply, params, data)
	if err != nil {
		return
	}

	param0 := params[0]
	var param1 interface{}
	if len(params) >= 2 {
		param1 = params[1]
	}

	var key string
	switch k := param0.(type) {
	case nil:
		// Returns whole data if key is null
		return data, nil
	case bool:
		if k {
			key = "true"
		} else {
			key = "false"
		}
	case float64:
		key = strconv.FormatFloat(k, 'f', -1, 64)
	case string:
		// Returns whole data if key is empty string
		if k == "" {
			return data, nil
		}
		key = k
	default:
		return nil, fmt.Errorf("var: key must be json primitive but got %T", param0)
	}

	res = data
	// NOTE: key is not empty here
	for _, part := range strings.Split(key, ".") {
		switch r := res.(type) {
		case []interface{}:
			i, err := strconv.Atoi(part)
			if err == nil && i >= 0 && i < len(r) {
				res = r[i]
				continue
			}

		case map[string]interface{}:
			v, ok := r[part]
			if ok {
				res = v
				continue
			}
		}
		return param1, nil
	}
	return res, nil

}

// AddOpMissing adds "missing" operation to the JSONLogic instance.
// NOTE: null/"" is considered as missing, for example:
//   logic: {"missing":"a"}
//   data: {"a":null}
//   result will be ["a"]
// ref:
//   - json-logic-js/logic.js::"missing"
func AddOpMissing(jl *JSONLogic) {
	jl.AddOperation("missing", opMissing)
}

func opMissing(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) == 0 {
		return []interface{}{}, nil
	}
	params, err = applyParams(apply, params, data)
	if err != nil {
		return
	}

	keys, ok := params[0].([]interface{})
	if !ok {
		keys = params
	}

	missing := []interface{}{}
	for _, key := range keys {
		res, err := opVar(apply, []interface{}{key}, data)
		if err != nil {
			return nil, err
		}
		if res == nil || res == "" {
			missing = append(missing, key)
		}
	}
	return missing, nil

}

// AddOpMissingSome adds "missing_some" operation to the JSONLogic instance. Param restriction:
//   - At least 2 params.
//   - The first must be evaluated to a numeric and the second evaluated to an array.
func AddOpMissingSome(jl *JSONLogic) {
	jl.AddOperation("missing_some", opMissingSome)
}

func opMissingSome(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) != 2 {
		return nil, fmt.Errorf("missing_some: expect 2 params")
	}
	params, err = applyParams(apply, params, data)
	if err != nil {
		return
	}

	needed, ok := params[0].(float64)
	if !ok {
		return nil, fmt.Errorf("missing_some: expect number for param 0 but got %T", params[0])
	}
	keys, ok := params[1].([]interface{})
	if !ok {
		return nil, fmt.Errorf("missing_some: expect array for param 1 but got %T", params[1])
	}

	missing, err := opMissing(apply, keys, data)
	if err != nil {
		return nil, err
	}

	missingArr := missing.([]interface{})
	if len(keys)-len(missingArr) >= int(needed) {
		return []interface{}{}, nil
	} else {
		return missingArr, nil
	}

}
