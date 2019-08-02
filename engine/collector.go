package engine

import (
	"context"
	"log"

	"github.com/themue/ghpm/analyze"
)

// Result contains the result of a job.
type Result struct {
	Job          *Job
	Accumulation analyze.Accumulation
	Err          error
}

// Collector provides a result channel for the pollers
// and collects their results.
type Collector struct {
	ctx           context.Context
	messagec      chan func()
	jobs          map[string]*Job
	accumulations analyze.Accumulations
}

// NewCollector starts the collecting goroutine.
func NewCollector(ctx context.Context) *Collector {
	c := &Collector{
		ctx:           ctx,
		messagec:      make(chan func(), 16),
		jobs:          make(map[string]*Job),
		accumulations: make(analyze.Accumulations),
	}
	go c.backend()
	return c
}

// HandleResult handles a new result passed by any poller.
func (c *Collector) HandleResult(result *Result) {
	// No need to wait here.
	c.messagec <- func() {
		if result.Err != nil {
			log.Printf("error in poll job %q: %v", result.Job.ID, result.Err)
			return
		}
		job, ok := c.jobs[result.Job.ID]
		if !ok {
			c.jobs[result.Job.ID] = result.Job
			job = result.Job
		}
		acc, ok := c.accumulations[job.ID]
		if !ok {
			acc = analyze.Accumulation{}
		}
		c.accumulations[job.ID] = job.Accumulate(acc, result.Accumulation)
	}
}

// GetJobs returns a list of all job IDs.
func (c *Collector) GetJobIDs() []string {
	var jobIDs []string
	c.do(func() {
		for jobID := range c.jobs {
			jobIDs = append(jobIDs, jobID)
		}
	})
	return jobIDs
}

// GetJob returns one job by id.
func (c *Collector) GetJob(id string) *Job {
	var job *Job
	c.do(func() {
		job = c.jobs[id]
	})
	return job
}

// GetAccumulationIDs returns a list of accumulation IDs of one job.
func (c *Collector) GetAccumulationIDs(jobID string) []string {
	var ids []string
	c.do(func() {
		acc, ok := c.accumulations[jobID]
		if !ok {
			return
		}
		for id := range acc {
			ids = append(ids, id)
		}
	})
	return ids
}

// GetAccumulation returns one accumulated value of one job.
func (c *Collector) GetAccumulation(jobID, id string) analyze.Value {
	var value analyze.Value
	c.do(func() {
		acc, ok := c.accumulations[jobID]
		if !ok {
			return
		}
		value = acc[id]
	})
	return value
}

// do performs a function and waits until it's done.
func (c *Collector) do(f func()) {
	donec := make(chan struct{})
	c.messagec <- func() {
		defer close(donec)
		f()
	}
	<-donec
}

// backend receives the individual results of the pollers
// and aggregates them.
func (c *Collector) backend() {
	defer close(c.messagec)
	for {
		select {
		case <-c.ctx.Done():
			return
		case method := <-c.messagec:
			method()
		}
	}
}
