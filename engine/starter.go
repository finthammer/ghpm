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
			Interval:      10 * time.Second,
			EventsAnalyzers: []analyze.EventsAnalyzer{
				analyze.TypeCounter,
			},
			Accumulate: analyze.AccumulateKeys,
		}, {
			ID:            "tideland-go",
			Owner:         "tideland",
			Repo:          "go",
			GitHubOptions: []github.Option{},
			Interval:      10 * time.Second,
			EventsAnalyzers: []analyze.EventsAnalyzer{
				analyze.CreateActorFilter("themue"),
				analyze.TypeCounter,
			},
			Accumulate: analyze.AccumulateKeys,
		}, {
			ID:            "kubernetes",
			Owner:         "kubernetes",
			Repo:          "kubernetes",
			GitHubOptions: []github.Option{},
			Interval:      10 * time.Second,
			EventsAnalyzers: []analyze.EventsAnalyzer{
				analyze.TypeCounter,
			},
			Accumulate: analyze.AccumulateKeys,
		}, {
			ID:            "kubeone",
			Owner:         "kubermatic",
			Repo:          "kubeone",
			GitHubOptions: []github.Option{},
			Interval:      10 * time.Second,
			EventsAnalyzers: []analyze.EventsAnalyzer{
				analyze.TypeCounter,
				analyze.ActorCounter,
			},
			Accumulate: analyze.AccumulateKeys,
		},
	}
	return jobs
}

// SpawnPollers starts the individual pollers for the passed jobs.
func SpawnPollers(ctx context.Context, jobs Jobs, collector *Collector) {
	for _, job := range jobs {
		SpawnPoller(ctx, job, collector)
	}
}
