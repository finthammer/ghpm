package logic

import (
	"log"
	"net/http"

	"github.com/themue/ghpm/analyze"
	"github.com/themue/ghpm/api/infra"
	"github.com/themue/ghpm/engine"
)

// AccumulationsHandler handles the individual requests regarding the
// analyzing jobs.
type AccumulationsHandler struct {
	collector *engine.Collector
}

// NewAccumulationsHandler creates a new handler managing the accumulations
// of a job.
func NewAccumulationsHandler(collector *engine.Collector) http.Handler {
	return &AccumulationsHandler{
		collector: collector,
	}
}

// ServeHTTPGet implements infra.GetHandler.
func (jh *AccumulationsHandler) ServeHTTPGet(w http.ResponseWriter, r *http.Request) {
	jobID, _ := infra.PathAt(r.URL.Path, 1)
	accumulationID, ok := infra.PathAt(r.URL.Path, 3)
	if ok {
		// Got an accumulation value for a job.
		log.Printf("requested accumulation %q for job %q", accumulationID, jobID)
		value := jh.collector.GetAccumulation(jobID, accumulationID)
		if value == nil {
			http.Error(w, "accumulated value not found", http.StatusNotFound)
			return
		}
		infra.ReplyJSON(w, struct {
			JobID          string
			AccumulationID string
			Value          analyze.Value
		}{
			JobID:          jobID,
			AccumulationID: accumulationID,
			Value:          value,
		})
		return
	}
	// Requesting list of accumulation IDs.
	accumulationIDs := jh.collector.GetAccumulationIDs(jobID)
	log.Printf("requested accumulations of job %q", jobID)
	infra.ReplyJSON(w, accumulationIDs)
}

// ServeHTTP implements http.Handler.
func (jh *AccumulationsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "cannot handle request", http.StatusMethodNotAllowed)
}
