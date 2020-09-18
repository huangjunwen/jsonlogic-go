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
