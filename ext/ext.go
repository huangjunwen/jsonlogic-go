package ext

import (
	"fmt"

	"github.com/huangjunwen/jsonlogic-go"
)

// AddOpRange adds "range" operation to the JSONLogic instance. The op accept 1 to 3 numeric params and
// generate a range of numbers, examples:
//   - {"range":null} -> []
//   - {"range":0} -> []
//   - {"range":2} -> [0,1]
//   - {"range":-2} -> [0,-1]
//   - {"range":[3,6]} -> [3,4,5]
//   - {"range":[6,3]} -> [6,5,4]
//   - {"range":[3,6,2]} -> [3,5]
//   - {"range":[6,3,-2]} -> [6,4]
func AddOpRange(jl *jsonlogic.JSONLogic) {
	jl.AddOperation("range", opRange)
}

func opRange(apply jsonlogic.Applier, params []interface{}, data interface{}) (result interface{}, err error) {
	if len(params) < 1 {
		return nil, fmt.Errorf("range: expect at least one param")
	}

	params, err = jsonlogic.ApplyParams(apply, params, data)
	if err != nil {
		return
	}
	var b, e, s interface{}
	switch len(params) {
	default:
		fallthrough
	case 3:
		b = params[0]
		e = params[1]
		s = params[2]
	case 2:
		b = params[0]
		e = params[1]
	case 1:
		b = float64(0)
		e = params[0]
	}

	bf64, err := jsonlogic.ToNumeric(b)
	if err != nil {
		return nil, fmt.Errorf("range: get begin error %s", err.Error())
	}
	ef64, err := jsonlogic.ToNumeric(e)
	if err != nil {
		return nil, fmt.Errorf("range: get end error %s", err.Error())
	}
	sf64, err := jsonlogic.ToNumeric(s)
	if err != nil {
		return nil, fmt.Errorf("range: get step error %s", err.Error())
	}

	var (
		begin = int(bf64)
		end   = int(ef64)
		step  = int(sf64)
	)
	if end == begin {
		return []interface{}{}, nil
	}
	if end > begin {
		if step == 0 {
			step = 1
		}
		if step < 0 {
			return nil, fmt.Errorf("range: end > begin but got negative step")
		}
		ret := []interface{}{}
		for i := begin; i < end; i += step {
			ret = append(ret, float64(i))
		}
		return ret, nil
	}

	if step == 0 {
		step = -1
	}
	if step > 0 {
		return nil, fmt.Errorf("range: end < begin but got postive step")
	}
	ret := []interface{}{}
	for i := begin; i > end; i += step {
		ret = append(ret, float64(i))
	}
	return ret, nil
}
