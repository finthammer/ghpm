package analyze

import (
	"fmt"

	"github.com/themue/ghpm/github"
)

// EventsAnalyzer describes a function analyzing a number of events. Results
// are passed between analyzers by the aggregate. Analyzers can use data
// inside the passed aggregate too.
type EventsAnalyzer func(es github.Events, acc Accumulation) (Accumulation, error)

// MarshalJSON implements json.Marshaler.
func (ea EventsAnalyzer) MarshalJSON() ([]byte, error) {
	return []byte("\"EventsAnalyzer\""), nil
}

// Counter simply counts the number of events.
func Counter(es github.Events, acc Accumulation) (Accumulation, error) {
	acc["totalCounter"] = IntValue(len(es))
	return acc, nil
}

// TypeCounter counts the different event types in the passed events.
func TypeCounter(es github.Events, acc Accumulation) (Accumulation, error) {
	for _, e := range es {
		var c IntValue
		tc := "typeCounter(" + e.Type + ")"
		c, _ = acc[tc].(IntValue)
		c++
		acc[tc] = c
	}
	return acc, nil
}

// PayloadsCollector collects the payloads per actor.
func PayloadsCollector(es github.Events, acc Accumulation) (Accumulation, error) {
	for _, e := range es {
		for key, payload := range e.Payload {
			var pc StringsValue
			ap := "actorPayload(" + key + "@" + e.Actor.Login + ")"
			pc, _ = acc[ap].(StringsValue)
			pc = append(pc, fmt.Sprintf("%q", payload))
			acc[ap] = pc
		}
	}
	return acc, nil
}

// ActorCounter counts the different actors.
func ActorCounter(es github.Events, acc Accumulation) (Accumulation, error) {
	for _, e := range es {
		var c IntValue
		ac := "actorCounter(" + e.Actor.Login + ")"
		c, _ = acc[ac].(IntValue)
		c++
		acc[ac] = c
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
			var c IntValue
			af := "actorFilter(" + e.Actor.Login + ")"
			c, _ = acc[af].(IntValue)
			c++
			acc[af] = c
		}
		return acc, nil
	}
}
