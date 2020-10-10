package jsonlogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInherit(t *testing.T) {
	assert := assert.New(t)

	parent := NewEmpty()
	parent.AddOperation("xxx", func(apply Applier, params []interface{}, data interface{}) (interface{}, error) {
		return "parent", nil
	})

	child := NewInherit(parent)
	child.AddOperation("xxx", func(apply Applier, params []interface{}, data interface{}) (interface{}, error) {
		return "child", nil
	})

	notFoundChild := NewInherit(parent)
	notFoundChild.AddOperation("xxx", nil)

	logic := map[string]interface{}{
		"xxx": nil,
	}

	{
		res, err := child.Apply(logic, nil)
		assert.NoError(err)
		assert.Equal("child", res)
	}

	{
		_, err := notFoundChild.Apply(logic, nil)
		assert.Error(err)
	}

	{
		res, err := parent.Apply(logic, nil)
		assert.NoError(err)
		assert.Equal("parent", res)
	}

}
