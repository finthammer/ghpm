package engine

import (
	"context"
	"log"
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
			Accumulate: func(accOld, accNew analyze.Accumulation) analyze.Accumulation {
				if !accOld.AddAll(accNew) {
					log.Printf("cannot accumulate correctly")
					return analyze.Accumulation{}
				}
				return accOld
			},
		}, {
			ID:            "status-themue",
			Owner:         "status-im",
			Repo:          "status-go",
			GitHubOptions: []github.Option{},
			Interval:      1 * time.Minute,
			EventsAnalyzers: []analyze.EventsAnalyzer{
				analyze.CreateActorFilter("themue"),
			},
			Accumulate: func(accOld, accNew analyze.Accumulation) analyze.Accumulation {
				if !accOld.AddAll(accNew) {
					log.Printf("cannot accumulate correctly")
					return analyze.Accumulation{}
				}
				return accOld
			},
		},
	}
	return jobs
}

// StartPollers starts the individual pollers for the passed jobs.
func StartPollers(ctx context.Context, jobs Jobs, collector *Collector) {
	for _, job := range jobs {
		SpawnPoller(ctx, job, collector)
	}
}
