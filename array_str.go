package jsonlogic

import (
	"fmt"
	"strings"
)

// AddOpMap adds "map" operation to the JSONLogic instance. Param restriction:
//   - At least two params: the first evaluated to an array and the second the logic.
func AddOpMap(jl *JSONLogic) {
	jl.AddOperation("map", opMap)
}

func opMap(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("map: expect at least two params")
	}

	scopedData, err := apply(params[0], data)
	if err != nil {
		return
	}
	scopedLogic := params[1]

	arr, ok := scopedData.([]interface{})
	if !ok {
		return []interface{}{}, nil
	}

	mappedArr := []interface{}{}
	for _, item := range arr {
		mappedItem, err := apply(scopedLogic, item)
		if err != nil {
			return nil, err
		}
		mappedArr = append(mappedArr, mappedItem)
	}
	return mappedArr, nil
}

// AddOpFilter adds "filter" operation to the JSONLogic instance. Param restriction:
//   - At least two params: the first evaluated to an array and the second the logic.
func AddOpFilter(jl *JSONLogic) {
	jl.AddOperation("filter", opFilter)
}

func opFilter(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("filter: expect at least two params")
	}
	scopedData, err := apply(params[0], data)
	if err != nil {
		return
	}
	scopedLogic := params[1]

	arr, ok := scopedData.([]interface{})
	if !ok {
		return []interface{}{}, nil
	}

	filteredArr := []interface{}{}
	for _, item := range arr {
		r, err := apply(scopedLogic, item)
		if err != nil {
			return nil, err
		}
		if ToBool(r) {
			filteredArr = append(filteredArr, item)
		}
	}
	return filteredArr, nil
}

// AddOpReduce adds "reduce" operation to the JSONLogic instance. Param restriction:
//   - At least three params: the first evaluated to an array, the second the logic and the third the initial value.
func AddOpReduce(jl *JSONLogic) {
	jl.AddOperation("reduce", opReduce)
}

func opReduce(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 3 {
		return nil, fmt.Errorf("filter: expect at least three params")
	}
	scopedData, err := apply(params[0], data)
	if err != nil {
		return
	}
	scopedLogic := params[1]
	initial := params[2]

	arr, ok := scopedData.([]interface{})
	if !ok {
		return initial, nil
	}

	for _, item := range arr {
		r, err := apply(scopedLogic, map[string]interface{}{
			"current":     item,
			"accumulator": initial,
		})
		if err != nil {
			return nil, err
		}
		initial = r
	}

	return initial, nil
}

// AddOpAll adds "all" operation to the JSONLogic instance. Param restriction:
//   - At least two params: the first evaluated to an array and the second the logic.
func AddOpAll(jl *JSONLogic) {
	jl.AddOperation("all", opAll)
}

func opAll(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("all: expect at least two params")
	}

	scopedData, err := apply(params[0], data)
	if err != nil {
		return
	}
	scopedLogic := params[1]

	arr, ok := scopedData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("all: expect array as the first param but got %T", scopedData)
	}

	if len(arr) == 0 {
		return false, nil
	}
	for _, item := range arr {
		r, err := apply(scopedLogic, item)
		if err != nil {
			return nil, err
		}
		if !ToBool(r) {
			return false, nil
		}
	}
	return true, nil
}

// AddOpNone adds "none" operation to the JSONLogic instance. Param restriction:
//   - At least two params: the first evaluated to an array and the second the logic.
func AddOpNone(jl *JSONLogic) {
	jl.AddOperation("none", opNone)
}

func opNone(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("none: expect at least two params")
	}

	scopedData, err := apply(params[0], data)
	if err != nil {
		return
	}
	scopedLogic := params[1]

	arr, ok := scopedData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("none: expect array as the first param but got %T", scopedData)
	}

	if len(arr) == 0 {
		return true, nil
	}
	for _, item := range arr {
		r, err := apply(scopedLogic, item)
		if err != nil {
			return nil, err
		}
		if ToBool(r) {
			return false, nil
		}
	}
	return true, nil
}

// AddOpSome adds "some" operation to the JSONLogic instance. Param restriction:
//   - At least two params: the first evaluated to an array and the second the logic.
func AddOpSome(jl *JSONLogic) {
	jl.AddOperation("some", opSome)
}

func opSome(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("none: expect at least two params")
	}

	scopedData, err := apply(params[0], data)
	if err != nil {
		return
	}
	scopedLogic := params[1]

	arr, ok := scopedData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("some: expect array as the first param but got %T", scopedData)
	}

	if len(arr) == 0 {
		return false, nil
	}
	for _, item := range arr {
		r, err := apply(scopedLogic, item)
		if err != nil {
			return nil, err
		}
		if ToBool(r) {
			return true, nil
		}
	}
	return false, nil
}

// AddOpMerge adds "merge" operation to the JSONLogic instance.
func AddOpMerge(jl *JSONLogic) {
	jl.AddOperation("merge", opMerge)
}

func opMerge(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	params, err = ApplyParams(apply, params, data)
	if err != nil {
		return
	}

	ret := []interface{}{}
	for _, param := range params {
		if arr, ok := param.([]interface{}); ok {
			ret = append(ret, arr...)
		} else {
			ret = append(ret, param)
		}
	}
	return ret, nil
}

// AddOpIn adds "in" operation to the JSONLogic instance. Params restriction:
//   - At least two params: the first to check and the second evaluated to an array or string.
//   - All items must be evaluated to json primitives.
func AddOpIn(jl *JSONLogic) {
	jl.AddOperation("in", opIn)
}

func opIn(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("in: expect at least two params")
	}
	params, err = ApplyParams(apply, params, data)
	if err != nil {
		return
	}

	param0 := params[0]
	if !IsPrimitive(param0) {
		return nil, fmt.Errorf("in: expect json primitive as first param but got %T", param0)
	}

	switch param1 := params[1].(type) {
	case []interface{}:
		for _, item := range param1 {
			if !IsPrimitive(item) {
				return nil, fmt.Errorf("in: expect json primitives in array but got %T", item)
			}
			if param0 == item {
				return true, nil
			}
		}
		return false, nil
	case string:
		s, err := ToString(param0)
		if err != nil {
			return nil, err
		}
		return strings.Contains(param1, s), nil
	default:
		return nil, fmt.Errorf("in: expect array/string as the second param but got %T", params[1])
	}

}

// AddOpCat adds "cat" operation to the JSONLogic instance. Params restriction:
//   - All items must be evaluated to json primitives that can converted to string.
func AddOpCat(jl *JSONLogic) {
	jl.AddOperation("cat", opCat)
}

func opCat(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	params, err = ApplyParams(apply, params, data)
	if err != nil {
		return
	}
	parts := []string{}
	for _, param := range params {
		s, err := ToString(param)
		if err != nil {
			return nil, err
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, ""), nil
}

// AddOpSubstr adds "substr" operation to the JSONLogic instance. Params restriction:
//   - Two to three params: the first evaluated to string, and the second/third to numbers.
func AddOpSubstr(jl *JSONLogic) {
	jl.AddOperation("substr", opSubstr)
}

func opSubstr(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("substr: expect at least two params")
	}
	params, err = ApplyParams(apply, params, data)
	if err != nil {
		return
	}

	s, err := ToString(params[0])
	if err != nil {
		return nil, err
	}
	r := []rune(s)

	var start int
	{
		param1, err := ToNumeric(params[1])
		if err != nil {
			return nil, err
		}
		start = int(param1)
		if start < 0 {
			start += len(r)
			if start < 0 {
				start = 0
			}
		}
	}

	var (
		end    int
		hasEnd bool
	)
	if len(params) > 2 {
		param2, err := ToNumeric(params[2])
		if err != nil {
			return nil, err
		}
		end = int(param2)
		hasEnd = true
	}

	if !hasEnd {
		return string(r[start:]), nil
	}

	if end >= 0 {
		return string(r[start : start+end]), nil
	}

	end += len(r)
	if end < 0 {
		end = 0
	}
	return string(r[start:end]), nil

}
