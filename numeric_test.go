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
	TestCases{
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
	}.Run(assert, jl)
}

func TestOpMinMax(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpMin(jl)
	AddOpMax(jl)
	TestCases{
		{Logic: `{"min":[]}`, Data: `null`, Result: nil},
		{Logic: `{"max":[]}`, Data: `null`, Result: nil},
		{Logic: `{"min":[1,"3","-1",2]}`, Data: `null`, Result: float64(-1)},
		{Logic: `{"max":[1,"3","-1",2]}`, Data: `null`, Result: float64(3)},
		{Logic: `{"min":["a"]}`, Data: `null`, Err: true},
		{Logic: `{"max":["a"]}`, Data: `null`, Err: true},
		{Logic: `{"min":[1,"-Inf"]}`, Data: `null`, Err: true},
		{Logic: `{"max":[1,"+Inf"]}`, Data: `null`, Err: true},
	}.Run(assert, jl)
}

func TestOpAdd(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpAdd(jl)
	TestCases{
		{Logic: `{"+":[]}`, Data: `null`, Result: float64(0)},
		{Logic: `{"+":["1"]}`, Data: `null`, Result: float64(1)},
		{Logic: `{"+":[1,"-2",33]}`, Data: `null`, Result: float64(32)},
		{Logic: `{"+":["a"]}`, Data: `null`, Err: true},
		{Logic: `{"+":["inf"]}`, Data: `null`, Err: true},
		{Logic: `{"+":[179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.000000,179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.000000]}`, Data: `null`, Err: true},
	}.Run(assert, jl)
}

func TestOpMul(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpMul(jl)
	TestCases{
		{Logic: `{"*":[]}`, Data: `null`, Err: true},
		{Logic: `{"*":["3"]}`, Data: `null`, Result: float64(3)},
		{Logic: `{"*":[2,"-2",2]}`, Data: `null`, Result: float64(-8)},
		{Logic: `{"*":["a"]}`, Data: `null`, Err: true},
		{Logic: `{"*":["inf"]}`, Data: `null`, Err: true},
		{Logic: `{"*":[179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.000000,179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.000000]}`, Data: `null`, Err: true},
	}.Run(assert, jl)
}

func TestOpMinus(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpMinus(jl)
	TestCases{
		{Logic: `{"-":[]}`, Data: `null`, Err: true},
		{Logic: `{"-":["-3.1"]}`, Data: `null`, Result: float64(3.1)},
		{Logic: `{"-":[4,2]}`, Data: `null`, Result: float64(2)},
		{Logic: `{"-":["a"]}`, Data: `null`, Err: true},
		{Logic: `{"-":["inf"]}`, Data: `null`, Err: true},
		{Logic: `{"-":[-179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.000000,179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.000000]}`, Data: `null`, Err: true},
	}.Run(assert, jl)
}

func TestOpDiv(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpDiv(jl)
	TestCases{
		{Logic: `{"/":[]}`, Data: `null`, Err: true},
		{Logic: `{"/":[1]}`, Data: `null`, Err: true},
		{Logic: `{"/":[4,2]}`, Data: `null`, Result: float64(2)},
		{Logic: `{"/":[1,"a"]}`, Data: `null`, Err: true},
		{Logic: `{"/":[1,0]}`, Data: `null`, Err: true},
	}.Run(assert, jl)
}

func TestOpMod(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpMod(jl)
	TestCases{
		{Logic: `{"%":[]}`, Data: `null`, Err: true},
		{Logic: `{"%":[1]}`, Data: `null`, Err: true},
		{Logic: `{"%":[101,2]}`, Data: `null`, Result: float64(1)},
		{Logic: `{"%":[1,"a"]}`, Data: `null`, Err: true},
		{Logic: `{"%":[1,0]}`, Data: `null`, Err: true},
	}.Run(assert, jl)
}
