package analyze

import (
	"github.com/themue/ghpm/github"
)

// Events performs the analyzers on the passed events.
func Events(es github.Events, eas ...EventsAnalyzer) (Accumulation, error) {
	var acc = Accumulation{}
	var err error
	for _, analyze := range eas {
		acc, err = analyze(es, acc)
		if err != nil {
			return nil, err
		}
	}
	return acc, nil
}
