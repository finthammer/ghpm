package engine

import (
	"context"
	"errors"
	"log"
	"sort"

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
	ctx           context.Context
	messagec      chan func()
	accumulations analyze.Accumulations
}

// NewCollector starts the collecting goroutine.
func NewCollector(ctx context.Context) *Collector {
	c := &Collector{
		ctx:           ctx,
		messagec:      make(chan func()),
		accumulations: make(analyze.Accumulations),
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
		acc, ok := c.accumulations[result.Job.ID]
		if !ok {
			acc = analyze.Accumulation{}
		}
		log.Printf("accumulating job %q ...", result.Job.ID)
		c.accumulations[result.Job.ID] = result.Job.Accumulate(acc, result.Accumulation)
	}
}

// GetIndex returns all collected job IDs.
func (c *Collector) GetIndex() []string {
	var index []string
	c.messagec <- func() {
		for id := range c.accumulations {
			index = append(index, id)
		}
	}
	sort.Strings(index)
	return index
}

// GetAccumulation returns one accumulation by ID.
func (c *Collector) GetAccumulation(id string) (analyze.Accumulation, error) {
	var accumulation analyze.Accumulation
	var err error
	c.messagec <- func() {
		found := c.accumulations[id]
		if found == nil {
			err = errors.New("not found")
			return
		}
		accumulation := analyze.Accumulation{}
		for key, value := range found {
			accumulation[key] = value.Copy()
		}
	}
	return accumulation, err
}

// GetAccumulations returns all accumulations.
func (c *Collector) GetAccumulations() analyze.Accumulations {
	var accumulations analyze.Accumulations
	c.messagec <- func() {
		accumulations = c.accumulations.Copy()
	}
	return accumulations
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
