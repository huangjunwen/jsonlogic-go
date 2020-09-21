package jsonlogic

import (
	"fmt"
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
		if toBool(r) {
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
		if !toBool(r) {
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
		if toBool(r) {
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
		if toBool(r) {
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
	params, err = applyParams(apply, params, data)
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
//   - At least two params: the first to check and the second evaluated to an array.
//   - All items must be evaluated to json primitives.
func AddOpIn(jl *JSONLogic) {
	jl.AddOperation("in", opIn)
}

func opIn(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("in: expect at least two params")
	}
	params, err = applyParams(apply, params, data)
	if err != nil {
		return
	}

	arr, ok := params[1].([]interface{})
	if !ok {
		return nil, fmt.Errorf("in: expect array as the second param but got %T", params[1])
	}

	param0 := params[0]
	if !isPrimitive(param0) {
		return nil, fmt.Errorf("in: expect json primitive")
	}

	for _, item := range arr {
		if !isPrimitive(item) {
			return nil, fmt.Errorf("in: expect json primitives")
		}
		if param0 == item {
			return true, nil
		}
	}
	return false, nil

}
