package logic

import (
	"log"
	"net/http"

	"github.com/themue/ghpm/api/infra"
	"github.com/themue/ghpm/engine"
)

// JobsHandler handles the individual requests regarding the
// analyzing jobs.
type JobsHandler struct {
	collector *engine.Collector
}

// NewJobsHandler creates a new handler managing the jobs.
func NewJobsHandler(collector *engine.Collector) http.Handler {
	return &JobsHandler{
		collector: collector,
	}
}

// ServeHTTPGet implements infra.GetHandler.
func (jh *JobsHandler) ServeHTTPGet(w http.ResponseWriter, r *http.Request) {
	jobID, ok := infra.PathAt(r.URL.Path, 1)
	if ok {
		// Got a job ID.
		log.Printf("requested job %q", jobID)
		job := jh.collector.GetJob(jobID)
		if job == nil {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}
		infra.ReplyJSON(w, job)
		return
	}
	// Requesting list of job IDs.
	jobIDs := jh.collector.GetJobIDs()
	log.Printf("requested %d job IDs", len(jobIDs))
	infra.ReplyJSON(w, jobIDs)
}

// ServeHTTP implements http.Handler.
func (jh *JobsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "cannot handle request", http.StatusMethodNotAllowed)
}
