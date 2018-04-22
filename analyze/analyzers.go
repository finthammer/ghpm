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

// Counter simply counts the number of events.
func Counter(es github.Events, a Aggregate) (Aggregate, error) {
	a["total"] = len(es)
	return a, nil
}

// TypeCounter counts the different event types in the passed events.
func TypeCounter(es github.Events, a Aggregate) (Aggregate, error) {
	for _, e := range es {
		var c int
		t := "type(" + e.Type + ")"
		c, _ = a[t].(int)
		c++
		a[t] = c
	}
	return a, nil
}

// CreateActorFilter creates an events analyzer for actors based
// on a passed login.
func CreateActorFilter(login string) EventsAnalyzer {
	return func(es github.Events, a Aggregate) (Aggregate, error) {
		for _, e := range es {
			if e.Actor.Login != login {
				continue
			}
			var c int
			al := "actor(" + e.Actor.Login + ")"
			c, _ = a[al].(int)
			c++
			a[al] = c
		}
		return a, nil
	}
}
