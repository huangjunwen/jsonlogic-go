package jsonlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLogic(t *testing.T) {
	assert := assert.New(t)

	for i, testCase := range []struct {
		Obj     interface{}
		IsLogic bool
	}{
		{Obj: nil, IsLogic: false},
		{Obj: true, IsLogic: false},
		{Obj: float64(1.22), IsLogic: false},
		{Obj: "x", IsLogic: false},
		{Obj: []interface{}{"x"}, IsLogic: false},
		{Obj: map[string]interface{}{}, IsLogic: false},
		{Obj: map[string]interface{}{"var": ""}, IsLogic: true},
		{Obj: map[string]interface{}{"x": "y", "a": 2}, IsLogic: false},
	} {
		assert.Equal(testCase.IsLogic, isLogic(testCase.Obj), "test case %d", i)
	}
}

func TestGetLogic(t *testing.T) {
	assert := assert.New(t)

	for i, testCase := range []struct {
		Obj    interface{}
		Op     string
		Params []interface{}
	}{
		{
			Obj: map[string]interface{}{
				"var": nil,
			},
			Op:     "var",
			Params: []interface{}{nil},
		},
		{
			Obj: map[string]interface{}{
				"var": []interface{}{},
			},
			Op:     "var",
			Params: []interface{}{},
		},
		{
			Obj: map[string]interface{}{
				"var": []interface{}{[]interface{}{}},
			},
			Op:     "var",
			Params: []interface{}{[]interface{}{}},
		},
	} {
		op, params := getLogic(testCase.Obj)
		assert.Equal(testCase.Op, op, "test case %d", i)
		assert.Equal(testCase.Params, params, "test case %d", i)
	}
}

func TestToBool(t *testing.T) {
	assert := assert.New(t)

	for i, testCase := range []struct {
		Obj    interface{}
		IsTrue bool
	}{
		{Obj: nil, IsTrue: false},
		{Obj: true, IsTrue: true},
		{Obj: false, IsTrue: false},
		{Obj: float64(0), IsTrue: false},
		{Obj: float64(1.1), IsTrue: true},
		{Obj: float64(-1.1), IsTrue: true},
		{Obj: "", IsTrue: false},
		{Obj: "x", IsTrue: true},
		{Obj: []interface{}{}, IsTrue: false},
		{Obj: []interface{}{"x"}, IsTrue: true},
		{Obj: map[string]interface{}{}, IsTrue: true},
		{Obj: map[string]interface{}{"var": "x"}, IsTrue: true},
	} {
		assert.Equal(testCase.IsTrue, toBool(testCase.Obj), "test case %d", i)
	}
}

func TestToNumeric(t *testing.T) {
	assert := assert.New(t)

	for i, testCase := range []struct {
		Obj interface{}
		N   float64
		Err bool
	}{
		{Obj: nil, N: 0},
		{Obj: true, N: 1},
		{Obj: false, N: 0},
		{Obj: float64(1.11), N: 1.11},
		{Obj: "1.11", N: 1.11},
		{Obj: "a", Err: true},
		{Obj: []interface{}{1}, Err: true},
		{Obj: map[string]interface{}{}, Err: true},
	} {
		n, err := toNumeric(testCase.Obj)
		if testCase.Err {
			assert.Error(err, "test case %d", i)
		} else {
			assert.NoError(err, "test case %d", i)
			assert.Equal(testCase.N, n, "test case %d", i)
		}
	}
}

func TestCompareValues(t *testing.T) {
	assert := assert.New(t)

	for i, testCase := range []struct {
		Symbol compSymbol
		Left   interface{}
		Right  interface{}
		Result bool
		Err    bool
	}{
		{Symbol: lt, Left: "1", Right: "a", Result: true},          // This is compared as strings.
		{Symbol: lt, Left: "1", Right: float64(1.1), Result: true}, // This is compared as numerics.
		{Symbol: lt, Left: float64(1), Right: "a", Err: true},
		{Symbol: lt, Left: "a", Right: float64(1), Err: true},
		{Symbol: lt, Left: []interface{}{}, Right: float64(1), Err: true},
	} {
		result, err := compareValues(testCase.Symbol, testCase.Left, testCase.Right)
		if testCase.Err {
			assert.Error(err, "test case %d", i)
		} else {
			assert.NoError(err, "test case %d", i)
			assert.Equal(testCase.Result, result, "test case %d", i)
		}
	}
}
