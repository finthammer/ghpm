package api

import (
	"context"
	"net/http"

	"github.com/themue/ghpm/api/infra"
	"github.com/themue/ghpm/api/logic"
	"github.com/themue/ghpm/engine"
)

// SpawnAPI starts the HTTP server providing the API.
func SpawnAPI(ctx context.Context, collector *engine.Collector) {
	// Prepare handler and multiplexing.
	jobsHandler := infra.NewNestedHandler()
	jobsHandler.AppendHandler("jobs", infra.NewMetaMethodHandler(logic.NewJobsHandler(collector)))
	jobsHandler.AppendHandler("accumulations", infra.NewMetaMethodHandler(logic.NewAccumulationsHandler(collector)))

	mux := http.NewServeMux()
	mux.Handle(
		"/api/jobs/",
		http.StripPrefix(
			"/api/",
			jobsHandler,
		),
	)
	// Spawn the server.
	go func() {
		panic(http.ListenAndServe(":1337", mux))
	}()
}
