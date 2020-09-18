package jsonlogic

import (
	"encoding/json"
	"strings"

	"github.com/stretchr/testify/assert"
)

type jsonLogicTestCase struct {
	Logic  string
	Data   string
	Result interface{}
	Err    bool
}

func runJSONLogicTestCases(a *assert.Assertions, jl *JSONLogic, cases []jsonLogicTestCase) {
	mustUnmarshal := func(src string, target interface{}) {
		if err := json.NewDecoder(strings.NewReader(src)).Decode(target); err != nil {
			panic(err)
		}
	}

	for i, testCase := range cases {
		var (
			logic, data, result interface{}
		)
		mustUnmarshal(testCase.Logic, &logic)
		mustUnmarshal(testCase.Data, &data)

		result, err := jl.Apply(logic, data)
		if testCase.Err {
			a.Error(err, "test case %d", i)
		} else {
			a.NoError(err, "test case %d", i)
			a.Equal(testCase.Result, result, "test case %d", i)
		}
	}
}
