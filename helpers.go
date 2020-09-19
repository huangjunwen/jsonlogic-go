package jsonlogic

import (
	"fmt"
)

// isLogic returns true when obj is a map[string]interface{} with length 1.
// ref:
//   - json-logic-js/logic.js::is_logic
func isLogic(obj interface{}) bool {
	if obj == nil {
		return false
	}

	m, ok := obj.(map[string]interface{})
	if !ok {
		return false
	}

	return len(m) == 1
}

// must check isLogic before calling this.
func getLogic(obj interface{}) (op string, params []interface{}) {
	var ok bool
	for key, value := range obj.(map[string]interface{}) {
		op = key
		params, ok = value.([]interface{})
		// {"var": "x"} - > {"var": ["x"]}
		if !ok {
			params = []interface{}{value}
		}
		return
	}
	panic(fmt.Errorf("getLogic: no operator in logic"))
}

// ref:
//   - http://jsonlogic.com/truthy.html
func isTrue(obj interface{}) bool {
	switch o := obj.(type) {
	case nil:
		return false
	case bool:
		return o
	case float64:
		return o != 0
	case string:
		return o != ""
	case []interface{}:
		return len(o) != 0
	case map[string]interface{}:
		// Always true
		return true
	default:
		panic(fmt.Errorf("isTrue: %T's truthy is not supported", obj))
	}
}

func isPrimitive(obj interface{}) bool {
	switch obj.(type) {
	case nil:
		return true
	case bool:
		return true
	case float64:
		return true
	case string:
		return true
	case []interface{}:
		return false
	case map[string]interface{}:
		return false
	default:
		panic(fmt.Errorf("isPrimitive: %T is not supported", obj))
	}
}
