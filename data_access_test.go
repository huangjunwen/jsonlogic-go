package jsonlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpVar(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpVar(jl)
	TestCases{
		// http://jsonlogic.com/operations.html#var
		{Logic: `{"var":["a"]}`, Data: `{"a":1,"b":2}`, Result: float64(1)},
		{Logic: `{"var":"a"}`, Data: `{"a":1,"b":2}`, Result: float64(1)},
		{Logic: `{"var":["z",26]}`, Data: `{"a":1,"b":2}`, Result: float64(26)},
		{Logic: `{"var":"champ.name"}`, Data: `{"champ":{"name":"Fezzig","height":223},"challenger":{"name":"DreadPirateRoberts","height":183}}`, Result: "Fezzig"},
		{Logic: `{"var":1}`, Data: `["zero", "one", "two"]`, Result: "one"},
		{Logic: `{"var":""}`, Data: `"Dolly"`, Result: "Dolly"},
		// Default.
		{Logic: `{"var":["a",["def"]]}`, Data: `{"c":"d"}`, Result: []interface{}{"def"}},
		// Using logic in key/default.
		{Logic: `{"var":["a",{"var":"a_default"}]}`, Data: `{"c":"d","a_default":"aaa"}`, Result: "aaa"},
		{Logic: `{"var":{"var":"pointer"}}`, Data: `{"pointer":"x","x":1.1}`, Result: float64(1.1)},
		// null or "" returns whole data.
		{Logic: `{"var":null}`, Data: `{"a":"b"}`, Result: map[string]interface{}{"a": "b"}},
		{Logic: `{"var":""}`, Data: `{"a":"b"}`, Result: map[string]interface{}{"a": "b"}},
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
		// Err.
		{Logic: `{"var":[]}`, Data: `null`, Err: true},
		{Logic: `{"var":[[]]}`, Data: `null`, Err: true},
	}.Run(assert, jl)

}

func TestOpMissing(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpVar(jl)
	AddOpMissing(jl)
	TestCases{
		// http://jsonlogic.com/operations.html#missing
		{Logic: `{"missing":["a","b"]}`, Data: `{"a":"apple","c":"carrot"}`, Result: []interface{}{"b"}},
		{Logic: `{"missing":["a","b"]}`, Data: `{"a":"apple", "b":"banana"}`, Result: []interface{}{}},
		// null or "" are treated as  missing
		{Logic: `{"missing":"a"}`, Data: `{"a":null}`, Result: []interface{}{"a"}}, // Surprise!
		{Logic: `{"missing":"a"}`, Data: `{"a":""}`, Result: []interface{}{"a"}},   // Surprise again!
		{Logic: `{"missing":"a"}`, Data: `{"a":false}`, Result: []interface{}{}},
		{Logic: `{"missing":"a"}`, Data: `{"a":0}`, Result: []interface{}{}},
		{Logic: `{"missing":"a"}`, Data: `{"a":[]}`, Result: []interface{}{}},
		{Logic: `{"missing":"a"}`, Data: `{"a":{}}`, Result: []interface{}{}},
		// First param is an array.
		{Logic: `{"missing":[["a"]]}`, Data: `{"a":"1"}`, Result: []interface{}{}},
		{Logic: `{"missing":[["a"]]}`, Data: `{"b":"1"}`, Result: []interface{}{"a"}},
		// Using logic in param.
		{Logic: `{"missing":{"var":"pointer"}}`, Data: `{"a":"xxx","pointer":"a"}`, Result: []interface{}{}},
		{Logic: `{"missing":{"var":"pointer"}}`, Data: `{"a":"xxx","pointer":"b"}`, Result: []interface{}{"b"}},
	}.Run(assert, jl)
}

func TestOpMissingSome(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpVar(jl)
	AddOpMissingSome(jl)
	TestCases{
		// http://jsonlogic.com/operations.html#missing_some
		{Logic: `{"missing_some":[1,["a","b","c"]]}`, Data: `{"a":"apple"}`, Result: []interface{}{}},
		{Logic: `{"missing_some":[2,["a","b","c"]]}`, Data: `{"a":"apple"}`, Result: []interface{}{"b", "c"}},
		// Using logic in param.
		{Logic: `{"missing_some":[2,{"var":"pointer"}]}`, Data: `{"pointer":["x","y"]}`, Result: []interface{}{"x", "y"}},
	}.Run(assert, jl)
}
