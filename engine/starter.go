package engine

import (
	"context"
	"time"

	"github.com/themue/ghpm/analyze"
	"github.com/themue/ghpm/github"
)

// ReadJobs reads a list of jobs from a configuration file.
func ReadJobs() Jobs {
	jobs := Jobs{
		{
			ID:            "ghpm",
			Owner:         "themue",
			Repo:          "ghpm",
			GitHubOptions: []github.Option{},
			Interval:      1 * time.Minute,
			EventsAnalyzers: []analyze.EventsAnalyzer{
				analyze.TypeCounter,
			},
		},
	}
	return jobs
}

// StartPoller poller starts the individual poller for the passed jobs.
func StartPoller(ctx context.Context, jobs Jobs) <-chan *Result {
	resultc := make(chan *Result)
	for _, job := range jobs {
		SpawnPoller(ctx, job, resultc)
	}
	return resultc
}
