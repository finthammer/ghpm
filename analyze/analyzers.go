package analyze

import (
	"github.com/themue/ghpm/github"
)

// EventsAnalyzer describes a function analyzing a number of events. Results
// are passed between analyzers by the aggregate. Analyzers can use data
// inside the passed aggregate too.
type EventsAnalyzer func(es github.Events, acc Accumulation) (Accumulation, error)

// Counter simply counts the number of events.
func Counter(es github.Events, acc Accumulation) (Accumulation, error) {
	acc["total"] = len(es)
	return acc, nil
}

// TypeCounter counts the different event types in the passed events.
func TypeCounter(es github.Events, acc Accumulation) (Accumulation, error) {
	for _, e := range es {
		var c int
		t := "type(" + e.Type + ")"
		c, _ = acc[t].(int)
		c++
		acc[t] = c
	}
	return acc, nil
}

// CreateActorFilter creates an events analyzer for actors based
// on a passed login.
func CreateActorFilter(login string) EventsAnalyzer {
	return func(es github.Events, acc Accumulation) (Accumulation, error) {
		for _, e := range es {
			if e.Actor.Login != login {
				continue
			}
			var c int
			al := "actor(" + e.Actor.Login + ")"
			c, _ = acc[al].(int)
			c++
			acc[al] = c
		}
		return acc, nil
	}
}
