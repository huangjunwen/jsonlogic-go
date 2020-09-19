package jsonlogic

import (
	"fmt"
	"strconv"
	"strings"
)

// AddOpVar adds "var" operation to the JSONLogic instance.
func AddOpVar(jl *JSONLogic) {
	jl.AddOperation("var", opVar)
}

func opVar(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {

	var keyObj, defObj interface{}

	switch len(params) {
	default:
		fallthrough

	case 2:
		defObj, err = apply(params[1], data)
		if err != nil {
			return
		}
		fallthrough

	case 1:
		keyObj, err = apply(params[0], data)
		if err != nil {
			return
		}

	case 0:
	}

	var key string
	switch k := keyObj.(type) {
	case nil:
		// Returns whole data if key is null
		return data, nil

	case float64:
		key = strconv.FormatFloat(k, 'f', -1, 64)

	case string:
		// Returns whole data if key is empty string
		if k == "" {
			return data, nil
		}
		key = k

	default:
		// XXX: This is different with jsonlogic-js
		return nil, fmt.Errorf("var: param should be null/number/string but got %T", keyObj)
	}

	res = data
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
		return defObj, nil
	}
	return res, nil

}

// AddOpMissing adds "missing" operation to the JSONLogic instance.
func AddOpMissing(jl *JSONLogic) {
	jl.AddOperation("missing", opMissing)
}

func opMissing(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {

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

// AddOpMissingSome adds "missing_some" operation to the JSONLogic instance.
func AddOpMissingSome(jl *JSONLogic) {
	jl.AddOperation("missing_some", opMissingSome)
}

func opMissingSome(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {

	if len(params) != 2 {
		return nil, fmt.Errorf("missing_some: expect 2 params")
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
