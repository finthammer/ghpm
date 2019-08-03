package analyze

// Value defines methods an accumulated value has to provide.
type Value interface {
	// Copy helps that output can use values while they
	// are changed by the analyzer.
	Copy() Value
}

// IntValue implements Value for integers.
type IntValue int

// Copy implements Value.
func (v IntValue) Copy() Value {
	return v
}

// StringValue implements Value for single strings.
type StringValue string

// Copy implements Value.
func (v StringValue) Copy() Value {
	return v
}

// StringsValue implements Value for a number of strings.
type StringsValue []string

// Copy implements Value.
func (v StringsValue) Copy() Value {
	return v
}

// UIntValue implements Value for unsigned integers.
type UIntValue uint

// Copy implements Value.
func (v UIntValue) Copy() Value {
	return v
}

// Float64Value implements Value for float64.
type Float64Value float64

// Copy implements Value.
func (v Float64Value) Copy() Value {
	return v
}
