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
	var (
		keyObj interface{}
		defObj interface{}
	)

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
		return nil, fmt.Errorf("Op var: param should be null/number/string but got %T", keyObj)
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
