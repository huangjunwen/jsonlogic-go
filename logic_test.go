package jsonlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpIf(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmptyJSONLogic()
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
	jl := NewEmptyJSONLogic()
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
	})

}

func TestOpStrictNotEqual(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmptyJSONLogic()
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
	})

}
