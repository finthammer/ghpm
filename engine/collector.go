package engine

import (
	"context"
	"log"

	"github.com/themue/ghpm/analyze"
)

// Accumulator combines old and new accumulated results.
type Accumulator func(accOld, accNew analyze.Accumulation) analyze.Accumulation

// Result contains the result of a job.
type Result struct {
	Job          *Job
	Accumulation analyze.Accumulation
	Err          error
}

// Collector provides a result channel for the pollers
// and collects their results.
type Collector struct {
	ctx         context.Context
	messagec    chan func()
	accumulates map[string]analyze.Accumulation
}

// NewCollector starts the collecting goroutine.
func NewCollector(ctx context.Context) *Collector {
	c := &Collector{
		ctx:         ctx,
		messagec:    make(chan func()),
		accumulates: make(map[string]analyze.Accumulation),
	}
	go c.backend()
	return c
}

// HandleResult handles a new result passed by any poller.
func (c *Collector) HandleResult(result *Result) {
	c.messagec <- func() {
		if result.Err != nil {
			log.Printf("error in poll job %q: %v", result.Job.ID, result.Err)
			return
		}
		acc, ok := c.accumulates[result.Job.ID]
		if !ok {
			acc = analyze.Accumulation{}
		}
		log.Printf("accumulating job %q ...", result.Job.ID)
		c.accumulates[result.Job.ID] = result.Job.Accumulate(acc, result.Accumulation)
	}
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
