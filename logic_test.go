package jsonlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpIf(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpIf(jl)
	TestCases{
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
	}.Run(assert, jl)

}

func TestOpStrictEqual(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpStrictEqual(jl)
	TestCases{
		// http://jsonlogic.com/operations.html
		{Logic: `{"===":[1,1]}`, Data: `null`, Result: true},
		{Logic: `{"===":[1,"1"]}`, Data: `null`, Result: false},
		// Zero/One param.
		{Logic: `{"===":[]}`, Data: `null`, Err: true},
		{Logic: `{"===":[null]}`, Data: `null`, Err: true},
		// Two params, primitives.
		{Logic: `{"===":[null,null]}`, Data: `null`, Result: true},
		{Logic: `{"===":[false,false]}`, Data: `null`, Result: true},
		{Logic: `{"===":[3.0,3]}`, Data: `null`, Result: true},
		{Logic: `{"===":["",""]}`, Data: `null`, Result: true},
		{Logic: `{"===":["",3.0]}`, Data: `null`, Result: false},
		// Non-primitives.
		{Logic: `{"===":["",[]]}`, Data: `null`, Err: true},
	}.Run(assert, jl)

}

func TestOpStrictNotEqual(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpStrictNotEqual(jl)
	TestCases{
		// http://jsonlogic.com/operations.html
		{Logic: `{"!==":[1,2]}`, Data: `null`, Result: true},
		{Logic: `{"!==":[1,"1"]}`, Data: `null`, Result: true},
		// Zero/One param.
		{Logic: `{"!==":[]}`, Data: `null`, Err: true},
		{Logic: `{"!==":[null]}`, Data: `null`, Err: true},
		// Two params, primitives.
		{Logic: `{"!==":[null,null]}`, Data: `null`, Result: false},
		{Logic: `{"!==":[false,false]}`, Data: `null`, Result: false},
		{Logic: `{"!==":[3.0,3]}`, Data: `null`, Result: false},
		{Logic: `{"!==":["",""]}`, Data: `null`, Result: false},
		{Logic: `{"!==":["",3.0]}`, Data: `null`, Result: true},
		// Non-primitives.
		{Logic: `{"!==":["",[]]}`, Data: `null`, Err: true},
	}.Run(assert, jl)

}

func TestOpNegative(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpNegative(jl)
	TestCases{
		// http://jsonlogic.com/operations.html
		{Logic: `{"!":[true]}`, Data: `null`, Result: false},
		{Logic: `{"!":true}`, Data: `null`, Result: false},
		// Zero param.
		{Logic: `{"!":[]}`, Data: `null`, Err: true},
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
	}.Run(assert, jl)
}

func TestOpDoubleNegative(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpDoubleNegative(jl)
	TestCases{
		// http://jsonlogic.com/operations.html
		{Logic: `{"!!":[true]}`, Data: `null`, Result: true},
		{Logic: `{"!!":true}`, Data: `null`, Result: true},
		// Zero param.
		{Logic: `{"!!":[]}`, Data: `null`, Err: true},
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
	}.Run(assert, jl)
}

func TestOpAnd(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpAnd(jl)
	TestCases{
		// http://jsonlogic.com/operations.html
		{Logic: `{"and":[true,true]}`, Data: `null`, Result: true},
		{Logic: `{"and":[true,false]}`, Data: `null`, Result: false},
		{Logic: `{"and":[true,"a",3]}`, Data: `null`, Result: float64(3)},
		{Logic: `{"and":[true,"",3]}`, Data: `null`, Result: ""},
		// Zero param.
		{Logic: `{"and":[]}`, Data: `null`, Err: true},
	}.Run(assert, jl)
}

func TestOpOr(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpOr(jl)
	TestCases{
		// http://jsonlogic.com/operations.html
		{Logic: `{"or":[true,false]}`, Data: `null`, Result: true},
		{Logic: `{"or":[false,true]}`, Data: `null`, Result: true},
		{Logic: `{"or":[false,"a"]}`, Data: `null`, Result: "a"},
		{Logic: `{"or":[false,0,"a"]}`, Data: `null`, Result: "a"},
		// Zero param.
		{Logic: `{"or":[]}`, Data: `null`, Err: true},
	}.Run(assert, jl)
}
