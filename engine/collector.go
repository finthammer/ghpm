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
	resultc     chan *Result
	accumulates map[string]analyze.Accumulation
}

// NewCollector starts the collecting goroutine.
func NewCollector(ctx context.Context) *Collector {
	c := &Collector{
		ctx:         ctx,
		resultc:     make(chan *Result),
		accumulates: make(map[string]analyze.Accumulation),
	}
	go c.backend()
	return c
}

// ResultC returns the channel where the users of
// the collector can write their results in.
func (c *Collector) ResultC() chan<- *Result {
	return c.resultc
}

// backend receives the individual results of the pollers
// and aggregates them.
func (c *Collector) backend() {
	defer close(c.resultc)
	for {
		select {
		case <-c.ctx.Done():
			return
		case r := <-c.resultc:
			if r.Err != nil {
				log.Printf("error in poll job %q: %v", r.Job.ID, r.Err)
				continue
			}
			acc, ok := c.accumulates[r.Job.ID]
			if !ok {
				acc = analyze.Accumulation{}
			}
			log.Printf("accumulating job %q ...", r.Job.ID)
			c.accumulates[r.Job.ID] = r.Job.Accumulate(acc, r.Accumulation)
		}
	}
}
