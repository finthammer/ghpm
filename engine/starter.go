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
			Interval:      10 * time.Second,
			EventsAnalyzers: []analyze.EventsAnalyzer{
				analyze.TypeCounter,
			},
			Accumulate: func(accOld, accNew analyze.Accumulation) analyze.Accumulation {
				log.Printf("adding %v ...", accNew.Keys())
				if !accOld.AddAll(accNew) {
					log.Printf("cannot accumulate correctly")
					return analyze.Accumulation{}
				}
				return accOld
			},
		}, {
			ID:            "tideland-goaudit",
			Owner:         "tideland",
			Repo:          "goaudit",
			GitHubOptions: []github.Option{},
			Interval:      10 * time.Second,
			EventsAnalyzers: []analyze.EventsAnalyzer{
				analyze.CreateActorFilter("themue"),
				analyze.TypeCounter,
			},
			Accumulate: func(accOld, accNew analyze.Accumulation) analyze.Accumulation {
				log.Printf("adding %v ...", accNew.Keys())
				if !accOld.AddAll(accNew) {
					log.Printf("cannot accumulate correctly")
					return analyze.Accumulation{}
				}
				return accOld
			},
		}, {
			ID:            "tideland-gocells",
			Owner:         "tideland",
			Repo:          "gocells",
			GitHubOptions: []github.Option{},
			Interval:      10 * time.Second,
			EventsAnalyzers: []analyze.EventsAnalyzer{
				analyze.CreateActorFilter("themue"),
				analyze.TypeCounter,
			},
			Accumulate: func(accOld, accNew analyze.Accumulation) analyze.Accumulation {
				log.Printf("adding %v ...", accNew.Keys())
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

// SpawnPollers starts the individual pollers for the passed jobs.
func SpawnPollers(ctx context.Context, jobs Jobs, collector *Collector) {
	for _, job := range jobs {
		SpawnPoller(ctx, job, collector)
	}
}
