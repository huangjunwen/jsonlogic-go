package ext

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/huangjunwen/jsonlogic-go"
)

func TestOpRange(t *testing.T) {
	assert := assert.New(t)
	jl := jsonlogic.NewEmpty()
	AddOpRange(jl)
	jsonlogic.TestCases{
		{Logic: `{"range":null}`, Data: `null`, Result: []interface{}{}},
		{Logic: `{"range":0}`, Data: `null`, Result: []interface{}{}},
		{Logic: `{"range":2}`, Data: `null`, Result: []interface{}{float64(0), float64(1)}},
		{Logic: `{"range":-2}`, Data: `null`, Result: []interface{}{float64(0), float64(-1)}},
		{Logic: `{"range":[3,6]}`, Data: `null`, Result: []interface{}{float64(3), float64(4), float64(5)}},
		{Logic: `{"range":[6,3]}`, Data: `null`, Result: []interface{}{float64(6), float64(5), float64(4)}},
		{Logic: `{"range":[3,6,2]}`, Data: `null`, Result: []interface{}{float64(3), float64(5)}},
		{Logic: `{"range":[6,3,-2]}`, Data: `null`, Result: []interface{}{float64(6), float64(4)}},
		{Logic: `{"range":[3,6,-1]}`, Data: `null`, Err: true},
		{Logic: `{"range":[6,3,2]}`, Data: `null`, Err: true},
	}.Run(assert, jl)
}
