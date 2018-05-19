package engine

import (
	"context"
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
	Interval        time.Duration
	EventsAnalyzers []analyze.EventsAnalyzer
}

// Jobs contains a number of jobs.
type Jobs []*Job

// Poller periodically polls events of a repository
// and analyzes it.
type Poller struct {
	ctx     context.Context
	job     *Job
	eventor *github.RepoEventor
	resultc chan *Result
}

// SpawnPoller creates a repository poller and starts it
// in the background.
func SpawnPoller(ctx context.Context, job *Job, resultc chan *Result) {
	p := &Poller{
		ctx:     ctx,
		job:     job,
		eventor: github.NewRepoEventor(job.Owner, job.Repo, job.GitHubOptions...),
		resultc: resultc,
	}
	go p.backend()
}

// backend polls the repository and analyzes it periodically
// in the background.
func (p *Poller) backend() {
	ticker := time.NewTicker(p.job.Interval)
	defer ticker.Stop()
	defer close(p.resultc)
	for {
		select {
		case <-p.ctx.Done():
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
