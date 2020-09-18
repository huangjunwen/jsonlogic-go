package jsonlogic

// AddOpIf adds "if"/"?:" operation to the JSONLogic instance.
func AddOpIf(jl *JSONLogic) {
	jl.AddOperation("if", opIf)
	jl.AddOperation("?:", opIf)
}

func opIf(apply Applier, params []interface{}, data interface{}) (res interface{}, err error) {

	var i int
	for i = 0; i < len(params)-1; i += 2 {
		r, err := apply(params[i], data)
		if err != nil {
			return nil, err
		}
		if isTrue(r) {
			return apply(params[i+1], data)
		}
	}

	if len(params) == i+1 {
		return apply(params[i], data)
	}

	return nil, nil
}
