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
		{Logic: `{">":[]}`, Data: `null`, Result: false},
		{Logic: `{">":[1]}`, Data: `null`, Result: false},
		{Logic: `{"<":[]}`, Data: `null`, Result: false},
		{Logic: `{"<":[1]}`, Data: `null`, Result: false},
		{Logic: `{">=":[]}`, Data: `null`, Result: false},
		{Logic: `{">=":[1]}`, Data: `null`, Result: false},
		{Logic: `{"<=":[]}`, Data: `null`, Result: false},
		{Logic: `{"<=":[1]}`, Data: `null`, Result: false},
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
		{Logic: `{">":["b",0]}`, Data: `null`, Err: true},
		{Logic: `{">":["1",0]}`, Data: `null`, Result: true},
		{Logic: `{">":["1",2]}`, Data: `null`, Result: false},
		{Logic: `{">=":["0",0]}`, Data: `null`, Result: true},
		{Logic: `{">=":["0",1]}`, Data: `null`, Result: false},
		{Logic: `{"<":["1",0]}`, Data: `null`, Result: false},
		{Logic: `{"<":["1",2]}`, Data: `null`, Result: true},
		{Logic: `{"<=":["0",0]}`, Data: `null`, Result: true},
		{Logic: `{"<=":["0",1]}`, Data: `null`, Result: true},
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
