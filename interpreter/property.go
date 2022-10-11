package interpreter

import "fmt"

func getProperty(x *GooseValue, prop *GooseValue) (*GooseValue, error) {
	switch x.Type {
	case GooseTypeArray, GooseTypeString:
		if x.Value == nil {
			return nil, fmt.Errorf("cannot index nil %s", x.Type)
		}

		if err := expectType(prop, GooseTypeNumeric); err != nil {
			return nil, err
		}
		idx := toInt64(prop.Value)

		var values []*GooseValue

		if x.Type == GooseTypeArray {
			values = x.Value.([]*GooseValue)
		} else {
			values = make([]*GooseValue, len(x.Value.(string)))
			for i, c := range x.Value.(string) {
				values[i] = wrap(string(c))
			}
		}

		if idx >= int64(len(values)) {
			return nil, fmt.Errorf("index %d out of bounds for array of length %d", idx, len(values))
		}

		if idx < 0 {
			idx = int64(len(values)) + idx
			if idx < 0 {
				idx = 0
			}
			if idx >= int64(len(values)) {
				return nil, fmt.Errorf("index %d out of bounds for array of length %d", toInt64(prop.Value), len(values))
			}
		}

		return values[idx], nil

	case GooseTypeComposite:
		switch prop.Type {
		case GooseTypeString, GooseTypeInt, GooseTypeFloat:
			// valid key
		default:
			return nil, fmt.Errorf("cannot index composite with type %s", prop.Type)
		}

		if x.Value == nil {
			return nil, fmt.Errorf("cannot index nil composite")
		}

		composite := x.Value.(GooseComposite)

		if val, ok := composite[prop.Value]; ok {
			return val, nil
		}

		return null, nil
	}

	return nil, fmt.Errorf("illegal property access on type %s", x.Type)
}

func setProperty(x *GooseValue, prop *GooseValue, val *GooseValue) error {
	if x.Value == nil {
		return fmt.Errorf("cannot index nil %s", x.Type)
	}

	switch x.Type {
	case GooseTypeArray:

		if err := expectType(prop, GooseTypeNumeric); err != nil {
			return err
		}
		idx := toInt64(prop.Value)

		var values []*GooseValue

		if x.Type == GooseTypeArray {
			values = x.Value.([]*GooseValue)
		} else {
			values = make([]*GooseValue, len(x.Value.(string)))
			for i, c := range x.Value.(string) {
				values[i] = wrap(string(c))
			}
		}

		if idx >= int64(len(values)) {
			return fmt.Errorf("index %d out of bounds for array of length %d", idx, len(values))
		}

		if idx < 0 {
			idx = int64(len(values)) + idx
			if idx < 0 {
				idx = 0
			}
			if idx >= int64(len(values)) {
				return fmt.Errorf("index %d out of bounds for array of length %d", toInt64(prop.Value), len(values))
			}
		}

		values[idx] = val

		return nil

	case GooseTypeComposite:
		switch prop.Type {
		case GooseTypeString, GooseTypeInt, GooseTypeFloat:
			// valid key
		default:
			return fmt.Errorf("cannot index composite with type %s", prop.Type)
		}

		if x.Value == nil {
			return fmt.Errorf("cannot index nil composite")
		}

		composite := x.Value.(GooseComposite)

		composite[prop.Value] = val

		return nil
	}

	return fmt.Errorf("illegal property assignment on type %s", x.Type)
}
