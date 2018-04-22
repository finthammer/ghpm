package analyze

import (
	"github.com/themue/ghpm/github"
)

// Events performs the analyzers on the passed events.
func Events(es github.Events, eas ...EventsAnalyzer) (Aggregate, error) {
	var a Aggregate = Aggregate{}
	var err error
	for _, analyze := range eas {
		a, err = analyze(es, a)
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}
