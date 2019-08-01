package analyze

// Accumulation is a generic key/value type.
type Accumulation map[string]Value

// Keys returns a list of the accumulation keys.
func (acc Accumulation) Keys() []string {
	var keys []string
	for key := range acc {
		keys = append(keys, key)
	}
	return keys
}

// Value returns the value for the given key.
func (acc Accumulation) Value(key string) Value {
	return acc[key]
}

// IsEmpty returns true if the accumulator contains no values.
func (acc Accumulation) IsEmpty() bool {
	return len(acc) == 0
}

// Add adds the key value of addend to the same of the accumulator
// if both are ints, uints, or float64.
func (acc Accumulation) Add(key string, addend Accumulation) bool {
	switch a := addend[key].(type) {
	case IntValue:
		if _, ok := acc[key]; !ok {
			acc[key] = a
			return true
		} else if i, ok := acc[key].(IntValue); ok {
			acc[key] = i + a
			return true
		}
	case UIntValue:
		if _, ok := acc[key]; !ok {
			acc[key] = a
			return true
		} else if ui, ok := acc[key].(UIntValue); ok {
			acc[key] = ui + a
			return true
		}
	case Float64Value:
		if _, ok := acc[key]; !ok {
			acc[key] = a
			return true
		} else if f, ok := acc[key].(Float64Value); ok {
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

// Accumulations is a map of accumulations by job ID.
type Accumulations map[string]Accumulation

// Copy creates a copy of the accumulations.
func (accs Accumulations) Copy() Accumulations {
	newAccs := Accumulations{}
	for id, acc := range accs {
		newAcc := Accumulation{}
		for key, value := range acc {
			newAcc[key] = value.Copy()
		}
		newAccs[id] = newAcc
	}
	return newAccs
}
