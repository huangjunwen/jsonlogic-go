package jsonlogic

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpVar(t *testing.T) {
	assert := assert.New(t)
	mustUnmarshal := func(src string, target interface{}) {
		if err := json.NewDecoder(strings.NewReader(src)).Decode(target); err != nil {
			panic(err)
		}
	}
	jl := NewEmptyJSONLogic()
	AddOpVar(jl)

	for i, testCase := range []struct {
		Logic  string
		Data   string
		Result interface{}
	}{
		// null or "" returns whole data.
		{Logic: `{"var":null}`, Data: `{"a":"b"}`, Result: map[string]interface{}{"a": "b"}},
		{Logic: `{"var":[null]}`, Data: `{"a":"b"}`, Result: map[string]interface{}{"a": "b"}},
		{Logic: `{"var":""}`, Data: `{"a":"b"}`, Result: map[string]interface{}{"a": "b"}},
		{Logic: `{"var":[""]}`, Data: `{"a":"b"}`, Result: map[string]interface{}{"a": "b"}},
		// String key and map.
		{Logic: `{"var":"a"}`, Data: `{"c":"d"}`, Result: nil},
		{Logic: `{"var":"a"}`, Data: `{"a":"b"}`, Result: "b"},
		{Logic: `{"var":"a.c"}`, Data: `{"a":"b"}`, Result: nil},
		{Logic: `{"var":"a.c"}`, Data: `{"a":{"c":"d"}}`, Result: "d"},
		// Numeric key and map.
		{Logic: `{"var":0.0}`, Data: `{"1":2}`, Result: nil},
		{Logic: `{"var":0.0}`, Data: `{"0":1}`, Result: float64(1)},
		{Logic: `{"var":0.1}`, Data: `{"0.1":1}`, Result: nil}, // Surprise! Since 0.1 -> "0.1" -> data["0"]["1"]
		{Logic: `{"var":0.1}`, Data: `{"0":{"1":3}}`, Result: float64(3)},
		{Logic: `{"var":-0.1}`, Data: `{"0":{"1":3}}`, Result: nil},
		{Logic: `{"var":-0.1}`, Data: `{"-0":{"1":3}}`, Result: float64(3)},
		// String key and array.
		{Logic: `{"var":"-1"}`, Data: `["a","b"]`, Result: nil},
		{Logic: `{"var":"0"}`, Data: `["a","b"]`, Result: "a"},
		{Logic: `{"var":"1"}`, Data: `["a","b"]`, Result: "b"},
		{Logic: `{"var":"2"}`, Data: `["a","b"]`, Result: nil},
		{Logic: `{"var":"a"}`, Data: `["a","b"]`, Result: nil},
		// Numeric key and array.
		{Logic: `{"var":-1.0}`, Data: `["a","b"]`, Result: nil},
		{Logic: `{"var":0.0}`, Data: `["a","b"]`, Result: "a"},
		{Logic: `{"var":1.0}`, Data: `["a","b"]`, Result: "b"},
		{Logic: `{"var":2.0}`, Data: `["a","b"]`, Result: nil},
		{Logic: `{"var":0.1}`, Data: `["a","b"]`, Result: nil},
		// Mix.
		{Logic: `{"var":"1.a"}`, Data: `["a",{"a":"b"}]`, Result: "b"},
		{Logic: `{"var":"a.0"}`, Data: `{"a":[8,9,10]}`, Result: float64(8)},
		{Logic: `{"var":{"var":"pointer"}}`, Data: `{"pointer":"x","x":1.1}`, Result: float64(1.1)},
		// Default.
		{Logic: `{"var":["a",["def"]]}`, Data: `{"c":"d"}`, Result: []interface{}{"def"}},
	} {
		var (
			logic, data, result interface{}
		)
		mustUnmarshal(testCase.Logic, &logic)
		mustUnmarshal(testCase.Data, &data)

		result, err := jl.Apply(logic, data)
		assert.NoError(err, "test case %d", i)
		assert.Equal(testCase.Result, result, "test case %d", i)
	}
}
