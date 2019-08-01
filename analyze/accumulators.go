package analyze

import (
	"encoding/json"
	"log"
)

// Accumulator combines old and new accumulated results.
type Accumulator func(accOld, accNew Accumulation) Accumulation

// AccumulateKeys adds the individual keys.
func AccumulateKeys(accOld, accNew Accumulation) Accumulation {
	if accNew.IsEmpty() {
		return accOld
	}
	if !accOld.AddAll(accNew) {
		log.Printf("cannot accumulate correctly")
		return Accumulation{}
	}
	b, err := json.Marshal(accOld)
	if err != nil {
		panic(err)
	}
	log.Printf("acc is now %v", string(b))
	return accOld
}
