package jsonlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsLogic(t *testing.T) {
	assert := assert.New(t)

	for i, testCase := range []struct {
		Obj     interface{}
		IsLogic bool
	}{
		{Obj: nil, IsLogic: false},
		{Obj: true, IsLogic: false},
		{Obj: float64(1.22), IsLogic: false},
		{Obj: "x", IsLogic: false},
		{Obj: []interface{}{"x"}, IsLogic: false},
		{Obj: map[string]interface{}{}, IsLogic: false},
		{Obj: map[string]interface{}{"var": ""}, IsLogic: true},
	} {
		assert.Equal(testCase.IsLogic, isLogic(testCase.Obj), "test case %d", i)
	}
}

func TestIsTrue(t *testing.T) {
	assert := assert.New(t)

	for i, testCase := range []struct {
		Obj    interface{}
		IsTrue bool
	}{
		{Obj: nil, IsTrue: false},
		{Obj: true, IsTrue: true},
		{Obj: false, IsTrue: false},
		{Obj: float64(0), IsTrue: false},
		{Obj: float64(1.1), IsTrue: true},
		{Obj: float64(-1.1), IsTrue: true},
		{Obj: "", IsTrue: false},
		{Obj: "x", IsTrue: true},
		{Obj: []interface{}{}, IsTrue: false},
		{Obj: []interface{}{"x"}, IsTrue: true},
	} {
		assert.Equal(testCase.IsTrue, isTrue(testCase.Obj), "test case %d", i)
	}
}
