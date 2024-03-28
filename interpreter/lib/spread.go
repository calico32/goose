package lib

type SpreadValue struct {
	Value Value
}

func (s *SpreadValue) Elements() []Value {
	switch v := s.Value.(type) {
	case *Array:
		return v.Elements
	default:
		return []Value{s.Value}
	}
}
