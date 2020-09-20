package jsonlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpCompare(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpLessThan(jl)
	AddOpLessEqual(jl)
	AddOpGreaterThan(jl)
	AddOpGreaterEqual(jl)
	runJSONLogicTestCases(assert, jl, []jsonLogicTestCase{
		// Zero/One param.
		{Logic: `{">":[]}`, Data: `null`, Err: true},
		{Logic: `{">":[1]}`, Data: `null`, Err: true},
		{Logic: `{"<":[]}`, Data: `null`, Err: true},
		{Logic: `{"<":[1]}`, Data: `null`, Err: true},
		{Logic: `{">=":[]}`, Data: `null`, Err: true},
		{Logic: `{">=":[1]}`, Data: `null`, Err: true},
		{Logic: `{"<=":[]}`, Data: `null`, Err: true},
		{Logic: `{"<=":[1]}`, Data: `null`, Err: true},
		// Two params, numeric compare.
		{Logic: `{">":[1.1,1]}`, Data: `null`, Result: true},
		{Logic: `{">":[1,1]}`, Data: `null`, Result: false},
		{Logic: `{">":[1,1.1]}`, Data: `null`, Result: false},
		{Logic: `{">=":[1.1,1]}`, Data: `null`, Result: true},
		{Logic: `{">=":[1,1]}`, Data: `null`, Result: true},
		{Logic: `{">=":[1,1.1]}`, Data: `null`, Result: false},
		{Logic: `{"<":[1.1,1]}`, Data: `null`, Result: false},
		{Logic: `{"<":[1,1]}`, Data: `null`, Result: false},
		{Logic: `{"<":[1,1.1]}`, Data: `null`, Result: true},
		{Logic: `{"<=":[1.1,1]}`, Data: `null`, Result: false},
		{Logic: `{"<=":[1,1]}`, Data: `null`, Result: true},
		{Logic: `{"<=":[1,1.1]}`, Data: `null`, Result: true},
		// Two params, string compare.
		{Logic: `{">":["b","a"]}`, Data: `null`, Result: true},
		{Logic: `{">":["b","b"]}`, Data: `null`, Result: false},
		{Logic: `{">":["a","b"]}`, Data: `null`, Result: false},
		{Logic: `{">=":["b","a"]}`, Data: `null`, Result: true},
		{Logic: `{">=":["b","b"]}`, Data: `null`, Result: true},
		{Logic: `{">=":["a","b"]}`, Data: `null`, Result: false},
		{Logic: `{"<":["b","a"]}`, Data: `null`, Result: false},
		{Logic: `{"<":["b","b"]}`, Data: `null`, Result: false},
		{Logic: `{"<":["a","b"]}`, Data: `null`, Result: true},
		{Logic: `{"<=":["b","a"]}`, Data: `null`, Result: false},
		{Logic: `{"<=":["b","b"]}`, Data: `null`, Result: true},
		{Logic: `{"<=":["a","b"]}`, Data: `null`, Result: true},
		// Two params, mix compare (as numeric).
		{Logic: `{">":["1",0]}`, Data: `null`, Result: true},
		{Logic: `{">":["1",2]}`, Data: `null`, Result: false},
		{Logic: `{">=":["0",0]}`, Data: `null`, Result: true},
		{Logic: `{">=":["0",1]}`, Data: `null`, Result: false},
		{Logic: `{"<":["1",0]}`, Data: `null`, Result: false},
		{Logic: `{"<":["1",2]}`, Data: `null`, Result: true},
		{Logic: `{"<=":["0",0]}`, Data: `null`, Result: true},
		{Logic: `{"<=":["0",1]}`, Data: `null`, Result: true},
		// Ill-form numeric.
		{Logic: `{"<":["b",0]}`, Data: `null`, Err: true},
		{Logic: `{"<=":["b",0]}`, Data: `null`, Err: true},
		{Logic: `{">":["b",0]}`, Data: `null`, Err: true},
		{Logic: `{">=":["b",0]}`, Data: `null`, Err: true},
		// Non-primitive.
		{Logic: `{"<":[1,[]]}`, Data: `null`, Err: true},
		{Logic: `{"<=":[1,[]]}`, Data: `null`, Err: true},
		{Logic: `{">":[1,[]]}`, Data: `null`, Err: true},
		{Logic: `{">=":[1,[]]}`, Data: `null`, Err: true},
		// Three params (between).
		{Logic: `{"<":["1",10,"100"]}`, Data: `null`, Result: true},
		{Logic: `{"<":["1",1,"100"]}`, Data: `null`, Result: false},
		{Logic: `{"<":["1",100,"100"]}`, Data: `null`, Result: false},
		{Logic: `{"<":["1",-1,"100"]}`, Data: `null`, Result: false},
		{Logic: `{"<":["1",101,"100"]}`, Data: `null`, Result: false},
		{Logic: `{"<=":["1",10,"100"]}`, Data: `null`, Result: true},
		{Logic: `{"<=":["1",1,"100"]}`, Data: `null`, Result: true},
		{Logic: `{"<=":["1",100,"100"]}`, Data: `null`, Result: true},
		{Logic: `{"<=":["1",-1,"100"]}`, Data: `null`, Result: false},
		{Logic: `{"<=":["1",101,"100"]}`, Data: `null`, Result: false},
	})
}

func TestOpMinMax(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpMin(jl)
	AddOpMax(jl)
	runJSONLogicTestCases(assert, jl, []jsonLogicTestCase{
		{Logic: `{"min":[]}`, Data: `null`, Result: nil},
		{Logic: `{"max":[]}`, Data: `null`, Result: nil},
		{Logic: `{"min":[1,"3","-1",2]}`, Data: `null`, Result: float64(-1)},
		{Logic: `{"max":[1,"3","-1",2]}`, Data: `null`, Result: float64(3)},
		{Logic: `{"min":["a"]}`, Data: `null`, Err: true},
		{Logic: `{"max":["a"]}`, Data: `null`, Err: true},
	})
}

func TestOpAdd(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpAdd(jl)
	runJSONLogicTestCases(assert, jl, []jsonLogicTestCase{
		{Logic: `{"+":[]}`, Data: `null`, Result: float64(0)},
		{Logic: `{"+":["1"]}`, Data: `null`, Result: float64(1)},
		{Logic: `{"+":[1,"-2",33]}`, Data: `null`, Result: float64(32)},
		{Logic: `{"+":["a"]}`, Data: `null`, Err: true},
	})
}

func TestOpMul(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpMul(jl)
	runJSONLogicTestCases(assert, jl, []jsonLogicTestCase{
		{Logic: `{"*":[]}`, Data: `null`, Err: true},
		{Logic: `{"*":["3"]}`, Data: `null`, Result: float64(3)},
		{Logic: `{"*":[2,"-2",2]}`, Data: `null`, Result: float64(-8)},
		{Logic: `{"*":["a"]}`, Data: `null`, Err: true},
	})
}
