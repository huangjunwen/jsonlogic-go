package jsonlogic

import (
	"encoding/json"
	"strings"

	"github.com/stretchr/testify/assert"
)

// TestCase is a single test case.
type TestCase struct {
	Logic  string
	Data   string
	Result interface{}
	Err    bool
}

// TestCases is a set of test cases.
type TestCases []TestCase

// Run a single test case
func (tc TestCase) Run(a *assert.Assertions, jl *JSONLogic) {
	logic := tc.mustUnmarshal(tc.Logic)
	data := tc.mustUnmarshal(tc.Data)
	result, err := jl.Apply(logic, data)
	if tc.Err {
		a.Error(err, "test case logic=%s data=%s", tc.Logic, tc.Data)
	} else {
		a.NoError(err, "test case logic=%s data=%s", tc.Logic, tc.Data)
		a.Equal(tc.Result, result, "test case logic=%s data=%s", tc.Logic, tc.Data)
	}
}

func (tc TestCase) mustUnmarshal(src string) interface{} {
	var res interface{}
	if err := json.NewDecoder(strings.NewReader(src)).Decode(&res); err != nil {
		panic(err)
	}
	return res
}

// Run a set of test cases.
func (tcs TestCases) Run(a *assert.Assertions, jl *JSONLogic) {
	for _, tc := range tcs {
		tc.Run(a, jl)
	}
}
