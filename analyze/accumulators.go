package analyze

import (
	"log"
)

// Accumulator combines old and new accumulated results.
type Accumulator func(accOld, accNew Accumulation) Accumulation

// MarshalJSON implements json.Marshaler.
func (a Accumulator) MarshalJSON() ([]byte, error) {
	return []byte("\"Accumulator\""), nil
}

// AccumulateKeys adds the individual keys.
func AccumulateKeys(accOld, accNew Accumulation) Accumulation {
	if accNew.IsEmpty() {
		return accOld
	}
	if !accOld.AddAll(accNew) {
		log.Printf("cannot accumulate correctly")
		return Accumulation{}
	}
	return accOld
}
