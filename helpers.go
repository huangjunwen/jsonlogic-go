package jsonlogic

import (
	"fmt"
	"strconv"
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

// getLogic gets operator and params from logic object.
// Must check isLogic before calling this.
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

// isTrue returns the truthy of an json object.
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

// isPrimitive returns true if obj is json primitive.
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

// toNumeric do a subset job of Number() in js:
//   - obj must be json primitive
//   - if obj is a string containing non numeric format, returns an error instead of NaN (like js)
func toNumeric(obj interface{}) (float64, error) {
	switch o := obj.(type) {
	case nil:
		return 0, nil
	case bool:
		if o {
			return 1, nil
		}
		return 0, nil
	case float64:
		return o, nil
	case string:
		return strconv.ParseFloat(o, 64)
	default:
		return 0, fmt.Errorf("toNumeric not support %T", obj)
	}
}

type compSymbol string

// compare symbol can be "<"/"<="/">"/">="
const (
	lt compSymbol = "<"
	le compSymbol = "<="
	gt compSymbol = ">"
	ge compSymbol = ">="
)

// compareValues only accepts json primitives
//
// ref:
//   - https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Less_than
//     > First, objects are converted to primitives using Symbol.ToPrimitive with the hint parameter be 'number'.
//     > If both values are strings, they are compared as strings, based on the values of the Unicode code points they contain.
//     > Otherwise JavaScript attempts to convert non-numeric types to numeric values:
//     > Boolean values true and false are converted to 1 and 0 respectively.
//     >   null is converted to 0.
//     >   undefined is converted to NaN.
//     >   Strings are converted based on the values they contain, and are converted as NaN if they do not contain numeric values.
//     > If either value is NaN, the operator returns false.
//     > Otherwise the values are compared as numeric values.
func compareValues(symbol compSymbol, left, right interface{}) (bool, error) {
	if !isPrimitive(left) || !isPrimitive(right) {
		return false, fmt.Errorf("only primitive values can be compared")
	}
	leftStr, leftIsStr := left.(string)
	rightStr, rightIsStr := right.(string)
	if leftIsStr && rightIsStr {
		switch symbol {
		case lt:
			return leftStr < rightStr, nil
		case le:
			return leftStr <= rightStr, nil
		case gt:
			return leftStr > rightStr, nil
		case ge:
			return leftStr >= rightStr, nil
		default:
			panic(fmt.Errorf("Impossible branch"))
		}
	}

	leftNum, err := toNumeric(left)
	if err != nil {
		return false, err
	}

	rightNum, err := toNumeric(right)
	if err != nil {
		return false, err
	}
	switch symbol {
	case lt:
		return leftNum < rightNum, nil
	case le:
		return leftNum <= rightNum, nil
	case gt:
		return leftNum > rightNum, nil
	case ge:
		return leftNum >= rightNum, nil
	default:
		panic(fmt.Errorf("Impossible branch"))
	}

}
