package main

import (
	"log"
	"time"

	"github.com/themue/ghpm/analyze"
	"github.com/themue/ghpm/github"
)

// Result contains the result of a job.
type Result struct {
	ID        string
	Aggregate analyze.Aggregate
	Err       error
}

// Job contains the parameters for the pollers work.
type Job struct {
	ID              string
	Owner           string
	Repo            string
	GitHubOptions   []github.Option
	Frequency       time.Duration
	EventsAnalyzers []analyze.EventsAnalyzer
}

// Poller periodically polls events of a repository
// and analyzes it.
type Poller struct {
	job     *Job
	eventor *github.RepoEventor
	resultc chan *Result
	stopc   chan struct{}
}

// NewPoller creates a repository poller and starts it
// in the background.
func NewPoller(job *Job, resultc chan *Result) *Poller {
	p := &Poller{
		job:     job,
		eventor: github.NewRepoEventor(job.Owner, job.Repo, job.GitHubOptions...),
		resultc: resultc,
		stopc:   make(chan struct{}),
	}
	return p
}

// Stop terminates the backend of the poller.
func (p *Poller) Stop() {
	p.stopc <- struct{}{}
}

// backend polls the repository and analyzes it periodically
// in the background.
func (p *Poller) backend() {
	ticker := time.NewTicker(p.job.Frequency)
	defer ticker.Stop()
	for {
		select {
		case <-p.stopc:
			return
		case <-ticker.C:
			p.resultc <- p.analyze()
		}
	}
}

// analyze performs a poll and an analyzing.
func (p *Poller) analyze() *Result {
	r := &Result{
		ID: p.job.ID,
	}
	events, err := p.eventor.Get()
	if err != nil {
		log.Printf("error polling %q: %v", p.job.ID, err)
		r.Err = err
		return r
	}
	a, err := analyze.Events(events, p.job.EventsAnalyzers...)
	if err != nil {
		log.Printf("error analyzing %q: %v", p.job.ID, err)
		r.Err = err
		return r
	}
	r.Aggregate = a
	return r
}
