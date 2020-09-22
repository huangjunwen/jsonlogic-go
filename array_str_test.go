package jsonlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpMapFilterReduce(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpVar(jl)
	AddOpAdd(jl)
	AddOpMul(jl)
	AddOpMod(jl)
	AddOpMap(jl)
	AddOpFilter(jl)
	AddOpReduce(jl)
	TestCases{
		// http://jsonlogic.com/operations.html#map-reduce-and-filter
		{Logic: `{"map":[{"var":"integers"},{"*":[{"var":""},2]}]}`, Data: `{"integers":[1,2,3,4,5]}`, Result: []interface{}{float64(2), float64(4), float64(6), float64(8), float64(10)}},
		{Logic: `{"filter":[{"var":"integers"},{"%":[{"var":""},2]}]}`, Data: `{"integers":[1,2,3,4,5]}`, Result: []interface{}{float64(1), float64(3), float64(5)}},
		{Logic: `{"reduce":[{"var":"integers"},{"+":[{"var":"current"},{"var":"accumulator"}]},0]}`, Data: `{"integers":[1,2,3,4,5]}`, Result: float64(15)},
		// Boring cases.
		{Logic: `{"map":[[1,2]]}`, Data: `null`, Err: true},
		{Logic: `{"map":[1,{"var":""}]}`, Data: `null`, Result: []interface{}{}},
		{Logic: `{"filter":[[1,2]]}`, Data: `null`, Err: true},
		{Logic: `{"filter":[1,{"var":""}]}`, Data: `null`, Result: []interface{}{}},
		{Logic: `{"reduce":[[1,2],{"var":""}]}`, Data: `null`, Err: true},
	}.Run(assert, jl)

}

func TestOpAllNoneSome(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpVar(jl)
	AddOpGreaterThan(jl)
	AddOpStrictEqual(jl)
	AddOpAll(jl)
	AddOpNone(jl)
	AddOpSome(jl)
	TestCases{
		// http://jsonlogic.com/operations.html#map-reduce-and-filter
		{Logic: `{"all":[[1,2,3],{">":[{"var":""},0]}]}`, Data: `null`, Result: true},
		{Logic: `{"none":[[-3,-2,-1],{">":[{"var":""},0]}]}`, Data: `null`, Result: true},
		{Logic: `{"some":[[-1,0,1],{">":[{"var":""},0]}]}`, Data: `null`, Result: true},
		{Logic: `{"some":[{"var":"pies"},{"===":[{"var":"filling"},"apple"]}]}`, Data: `{"pies":[{"filling":"pumpkin","temp":110},{"filling":"rhubarb","temp":210},{"filling":"apple","temp":310}]}`, Result: true},
		//
		{Logic: `{"all":[[1,2,-3],{">":[{"var":""},0]}]}`, Data: `null`, Result: false},
		{Logic: `{"all":[[],{">":[{"var":""},0]}]}`, Data: `null`, Result: false},
		{Logic: `{"none":[[-3,2,-1],{">":[{"var":""},0]}]}`, Data: `null`, Result: false},
		{Logic: `{"none":[[],{">":[{"var":""},0]}]}`, Data: `null`, Result: true},
		{Logic: `{"some":[[-1,0,-2],{">":[{"var":""},0]}]}`, Data: `null`, Result: false},
		{Logic: `{"some":[[],{">":[{"var":""},0]}]}`, Data: `null`, Result: false},
	}.Run(assert, jl)
}

func TestOpMerge(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpMissing(jl)
	AddOpIf(jl)
	AddOpVar(jl)
	AddOpMerge(jl)
	TestCases{
		// http://jsonlogic.com/operations.html#merge
		{Logic: `{"merge":[[1,2],[3,4]]}`, Data: `null`, Result: []interface{}{float64(1), float64(2), float64(3), float64(4)}},
		{Logic: `{"missing":{"merge":["vin",{"if":[{"var":"financing"},["apr","term"],[]]}]}}`, Data: `{"financing":true}`, Result: []interface{}{"vin", "apr", "term"}},
		{Logic: `{"missing":{"merge":["vin",{"if":[{"var":"financing"},["apr","term"],[]]}]}}`, Data: `{"financing":false}`, Result: []interface{}{"vin"}},
		//
		{Logic: `{"merge":[]}`, Data: `null`, Result: []interface{}{}},
	}.Run(assert, jl)
}

func TestOpIn(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpIn(jl)
	TestCases{
		// http://jsonlogic.com/operations.html#in
		{Logic: `{"in":["Ringo",["John","Paul","George","Ringo"]]}`, Data: `null`, Result: true},
		// http://jsonlogic.com/operations.html#in-1
		{Logic: `{"in":["Spring","Springfield"]}`, Data: `null`, Result: true},
		// 'in' uses strict equal in array mode.
		{Logic: `{"in":[1,["1"]]}`, Data: `null`, Result: false},
		// 'in' convert to string in string mode.
		{Logic: `{"in":[1,"1"]}`, Data: `null`, Result: true},
		{Logic: `{"in":[null,"xnull"]}`, Data: `null`, Result: true},
	}.Run(assert, jl)
}

func TestOpCat(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpVar(jl)
	AddOpCat(jl)
	TestCases{
		// http://jsonlogic.com/operations.html#cat
		{Logic: `{"cat":["I love"," pie"]}`, Data: `null`, Result: "I love pie"},
		{Logic: `{"cat":["I love ",{"var":"filling"}," pie"]}`, Data: `{"filling":"apple","temp":110}`, Result: "I love apple pie"},
		//
		{Logic: `{"cat":[]}`, Data: `null`, Result: ""},
		{Logic: `{"cat":[1,true,3.11,null]}`, Data: `null`, Result: "1true3.11null"},
	}.Run(assert, jl)
}

func TestOpSubstr(t *testing.T) {
	assert := assert.New(t)
	jl := NewEmpty()
	AddOpSubstr(jl)
	TestCases{
		// http://jsonlogic.com/operations.html#substr
		{Logic: `{"substr":["jsonlogic",4]}`, Data: `null`, Result: "logic"},
		{Logic: `{"substr":["jsonlogic",-5]}`, Data: `null`, Result: "logic"},
		{Logic: `{"substr":["jsonlogic",1,3]}`, Data: `null`, Result: "son"},
		{Logic: `{"substr":["jsonlogic",4,-2]}`, Data: `null`, Result: "log"},
		// Multi bytes chars.
		{Logic: `{"substr":["大家好",1]}`, Data: `null`, Result: "家好"},
	}.Run(assert, jl)
}
