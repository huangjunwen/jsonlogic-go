package jsonlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpIf(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpIf(jl)
	runJSONLogicTestCases(assert, jl, []jsonLogicTestCase{
		// http://jsonlogic.com/operations.html#if
		{Logic: `{"if":[true,"yes","no"]}`, Data: `null`, Result: "yes"},
		{Logic: `{"if":[false,"yes","no"]}`, Data: `null`, Result: "no"},
		// Zero param.
		{Logic: `{"if":[]}`, Data: `null`, Result: nil},
		// One param.
		{Logic: `{"if":"xxx"}`, Data: `null`, Result: "xxx"},
		// Two params.
		{Logic: `{"if":[true,"1"]}`, Data: `null`, Result: "1"},
		{Logic: `{"if":[false,"1"]}`, Data: `null`, Result: nil},
	})

}

func TestOpStrictEqual(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpStrictEqual(jl)
	runJSONLogicTestCases(assert, jl, []jsonLogicTestCase{
		// http://jsonlogic.com/operations.html
		{Logic: `{"===":[1,1]}`, Data: `null`, Result: true},
		{Logic: `{"===":[1,"1"]}`, Data: `null`, Result: false},
		// Zero param.
		{Logic: `{"===":[]}`, Data: `null`, Result: true},
		// One param.
		{Logic: `{"===":[null]}`, Data: `null`, Result: false},
		// Two params, primitives.
		{Logic: `{"===":[null,null]}`, Data: `null`, Result: true},
		{Logic: `{"===":[false,false]}`, Data: `null`, Result: true},
		{Logic: `{"===":[3.0,3]}`, Data: `null`, Result: true},
		{Logic: `{"===":["",""]}`, Data: `null`, Result: true},
		{Logic: `{"===":["",3.0]}`, Data: `null`, Result: false},
		// Non-primitives.
		{Logic: `{"===":["",[]]}`, Data: `null`, Err: true},
	})

}

func TestOpStrictNotEqual(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpStrictNotEqual(jl)
	runJSONLogicTestCases(assert, jl, []jsonLogicTestCase{
		// http://jsonlogic.com/operations.html
		{Logic: `{"!==":[1,2]}`, Data: `null`, Result: true},
		{Logic: `{"!==":[1,"1"]}`, Data: `null`, Result: true},
		// Zero param.
		{Logic: `{"!==":[]}`, Data: `null`, Result: false},
		// One param.
		{Logic: `{"!==":[null]}`, Data: `null`, Result: true},
		// Two params, primitives.
		{Logic: `{"!==":[null,null]}`, Data: `null`, Result: false},
		{Logic: `{"!==":[false,false]}`, Data: `null`, Result: false},
		{Logic: `{"!==":[3.0,3]}`, Data: `null`, Result: false},
		{Logic: `{"!==":["",""]}`, Data: `null`, Result: false},
		{Logic: `{"!==":["",3.0]}`, Data: `null`, Result: true},
		// Non-primitives.
		{Logic: `{"!==":["",[]]}`, Data: `null`, Err: true},
	})

}

func TestOpNegative(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpNegative(jl)
	runJSONLogicTestCases(assert, jl, []jsonLogicTestCase{
		// http://jsonlogic.com/operations.html
		{Logic: `{"!":[true]}`, Data: `null`, Result: false},
		{Logic: `{"!":true}`, Data: `null`, Result: false},
		// Zero param.
		{Logic: `{"!":[]}`, Data: `null`, Result: true},
		// One param.
		{Logic: `{"!":null}`, Data: `null`, Result: true},
		{Logic: `{"!":false}`, Data: `null`, Result: true},
		{Logic: `{"!":true}`, Data: `null`, Result: false},
		{Logic: `{"!":0}`, Data: `null`, Result: true},
		{Logic: `{"!":-1}`, Data: `null`, Result: false},
		{Logic: `{"!":""}`, Data: `null`, Result: true},
		{Logic: `{"!":"x"}`, Data: `null`, Result: false},
		{Logic: `{"!":[[]]}`, Data: `null`, Result: true},
		{Logic: `{"!":[["x"]]}`, Data: `null`, Result: false},
		{Logic: `{"!":{}}`, Data: `null`, Result: false},
		{Logic: `{"!":{"1":2,"3":4}}`, Data: `null`, Result: false},
	})
}

func TestOpDoubleNegative(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpDoubleNegative(jl)
	runJSONLogicTestCases(assert, jl, []jsonLogicTestCase{
		// http://jsonlogic.com/operations.html
		{Logic: `{"!!":[true]}`, Data: `null`, Result: true},
		{Logic: `{"!!":true}`, Data: `null`, Result: true},
		// Zero param.
		{Logic: `{"!!":[]}`, Data: `null`, Result: false},
		// One param.
		{Logic: `{"!!":null}`, Data: `null`, Result: false},
		{Logic: `{"!!":false}`, Data: `null`, Result: false},
		{Logic: `{"!!":true}`, Data: `null`, Result: true},
		{Logic: `{"!!":0}`, Data: `null`, Result: false},
		{Logic: `{"!!":-1}`, Data: `null`, Result: true},
		{Logic: `{"!!":""}`, Data: `null`, Result: false},
		{Logic: `{"!!":"x"}`, Data: `null`, Result: true},
		{Logic: `{"!!":[[]]}`, Data: `null`, Result: false},
		{Logic: `{"!!":[["x"]]}`, Data: `null`, Result: true},
		{Logic: `{"!!":{}}`, Data: `null`, Result: true},
		{Logic: `{"!!":{"1":2,"3":4}}`, Data: `null`, Result: true},
	})
}

func TestOpAnd(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpAnd(jl)
	runJSONLogicTestCases(assert, jl, []jsonLogicTestCase{
		// http://jsonlogic.com/operations.html
		{Logic: `{"and":[true,true]}`, Data: `null`, Result: true},
		{Logic: `{"and":[true,false]}`, Data: `null`, Result: false},
		{Logic: `{"and":[true,"a",3]}`, Data: `null`, Result: float64(3)},
		{Logic: `{"and":[true,"",3]}`, Data: `null`, Result: ""},
		// Zero param.
		{Logic: `{"and":[]}`, Data: `null`, Err: true},
	})
}

func TestOpOr(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpOr(jl)
	runJSONLogicTestCases(assert, jl, []jsonLogicTestCase{
		// http://jsonlogic.com/operations.html
		{Logic: `{"or":[true,false]}`, Data: `null`, Result: true},
		{Logic: `{"or":[false,true]}`, Data: `null`, Result: true},
		{Logic: `{"or":[false,"a"]}`, Data: `null`, Result: "a"},
		{Logic: `{"or":[false,0,"a"]}`, Data: `null`, Result: "a"},
		// Zero param.
		{Logic: `{"or":[]}`, Data: `null`, Err: true},
	})
}
