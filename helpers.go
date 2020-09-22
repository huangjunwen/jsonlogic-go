package jsonlogic

import (
	"fmt"
	"math"
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
	panic(fmt.Errorf("no operator in logic"))
}

// IsPrimitive returns true if obj is json primitive (null/bool/float64/string).
func IsPrimitive(obj interface{}) bool {
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
		panic(fmt.Errorf("IsPrimitive not support type %T", obj))
	}
}

// ToBool returns the truthy of a json object.
// ref:
//   - http://jsonlogic.com/truthy.html
//   - json-logic-js/logic.js::truthy
func ToBool(obj interface{}) bool {
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
		panic(fmt.Errorf("ToBool got non-json type %T", obj))
	}
}

// ToNumeric converts json primitive to numeric. It should be the same as JavaScript's Number(), except:
//   - an error is returned if obj is not a json primitive.
//   - an error is returned if obj is string but not well-formed.
//   - the number is NaN or +Inf/-Inf.
func ToNumeric(obj interface{}) (f float64, err error) {
	defer func() {
		if err == nil {
			if math.IsNaN(f) {
				f = 0
				err = fmt.Errorf("ToNumeric got NaN")
			} else if math.IsInf(f, 0) {
				f = 0
				err = fmt.Errorf("ToNumeric got +Inf/-Inf")
			}
		}
	}()

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
	case []interface{}, map[string]interface{}:
		return 0, fmt.Errorf("ToNumeric not support type %T", obj)
	default:
		panic(fmt.Errorf("ToNumeric got non-json type %T", obj))
	}
}

// ToString converts json primitive to string. It should be the same as JavaScript's String(), except:
//   - an error is returned if obj is not a json primitive.
//   - obj is number NaN or +Inf/-Inf.
func ToString(obj interface{}) (string, error) {
	switch o := obj.(type) {
	case nil:
		return "null", nil
	case bool:
		if o {
			return "true", nil
		}
		return "false", nil
	case float64:
		if math.IsNaN(o) {
			return "", fmt.Errorf("ToString got NaN")
		}
		if math.IsInf(o, 0) {
			return "", fmt.Errorf("ToString got +Inf/-Inf")
		}
		return strconv.FormatFloat(o, 'f', -1, 64), nil
	case string:
		return o, nil
	case []interface{}, map[string]interface{}:
		return "", fmt.Errorf("ToString not support type %T", obj)
	default:
		panic(fmt.Errorf("ToString got non-json type %T", obj))
	}
}

// CompSymbol represents compare operator.
type CompSymbol string

const (
	LT CompSymbol = "<"
	LE CompSymbol = "<="
	GT CompSymbol = ">"
	GE CompSymbol = ">="
	EQ CompSymbol = "==="
	NE CompSymbol = "!=="
)

// CompareValues compares json primitives. It should be the same as JavaScript's "<"/"<="/">"/">="/"==="/"!==", except:
//   - an error is returned if any value is not a json primitive.
//   - any error retuend by ToNumeric.
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
func CompareValues(symbol CompSymbol, left, right interface{}) (bool, error) {
	if !IsPrimitive(left) || !IsPrimitive(right) {
		return false, fmt.Errorf("only primitive values can be compared")
	}
	switch symbol {
	case EQ:
		return left == right, nil
	case NE:
		return left != right, nil
	}

	leftStr, leftIsStr := left.(string)
	rightStr, rightIsStr := right.(string)
	if leftIsStr && rightIsStr {
		switch symbol {
		case LT:
			return leftStr < rightStr, nil
		case LE:
			return leftStr <= rightStr, nil
		case GT:
			return leftStr > rightStr, nil
		case GE:
			return leftStr >= rightStr, nil
		default:
			panic(fmt.Errorf("Impossible branch"))
		}
	}

	leftNum, err := ToNumeric(left)
	if err != nil {
		return false, err
	}

	rightNum, err := ToNumeric(right)
	if err != nil {
		return false, err
	}
	switch symbol {
	case LT:
		return leftNum < rightNum, nil
	case LE:
		return leftNum <= rightNum, nil
	case GT:
		return leftNum > rightNum, nil
	case GE:
		return leftNum >= rightNum, nil
	default:
		panic(fmt.Errorf("Impossible branch"))
	}

}

// ApplyParams apply data to an array of params. Useful in operation implementation.
func ApplyParams(apply Applier, params []interface{}, data interface{}) ([]interface{}, error) {
	r, err := apply(params, data)
	if err != nil {
		return nil, err
	}
	return r.([]interface{}), nil
}
