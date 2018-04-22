package analyze

import (
	"github.com/themue/ghpm/github"
)

// Aggregate is a generic key/value type.
type Aggregate map[string]interface{}

// EventsAnalyzer describes a function analyzing a number of events. Results
// are passed between analyzers by the aggregate. Analyzers can use data
// inside the passed aggregate too.
type EventsAnalyzer func(es github.Events, a Aggregate) (Aggregate, error)

// TypeCounter counts the different event types in the passed events.
func TypeCounter(es github.Events, a Aggregate) (Aggregate, error) {
	for _, e := range es {
		var c int
		c, _ = a[e.Type].(int)
		c++
		a[e.Type] = c
	}
	return a, nil
}
