package analyze

// Accumulation is a generic key/value type.
type Accumulation map[string]interface{}

// Keys returns a list of the accumulation keys.
func (acc Accumulation) Keys() []string {
	var keys []string
	for key := range acc {
		keys = append(keys, key)
	}
	return keys
}

// Add adds the key value of addend to the same of the accumulator
// if both are ints, uints, or float64.
func (acc Accumulation) Add(key string, addend Accumulation) bool {
	switch a := addend[key].(type) {
	case int:
		if _, ok := acc[key]; !ok {
			acc[key] = a
			return true
		} else if i, ok := acc[key].(int); ok {
			acc[key] = i + a
			return true
		}
	case uint:
		if _, ok := acc[key]; !ok {
			acc[key] = a
			return true
		} else if ui, ok := acc[key].(uint); ok {
			acc[key] = ui + a
			return true
		}
	case float64:
		if _, ok := acc[key]; !ok {
			acc[key] = a
			return true
		} else if f, ok := acc[key].(float64); ok {
			acc[key] = f + a
			return true
		}
	}
	return false
}

// AddAll performs Add() for all keys of the addend.
func (acc Accumulation) AddAll(addend Accumulation) bool {
	for key := range addend {
		if !acc.Add(key, addend) {
			return false
		}
	}
	return true
}
