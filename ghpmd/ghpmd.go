package main

import (
	"context"
	"log"
	"time"

	"github.com/themue/ghpm/api"
	"github.com/themue/ghpm/engine"
)

func main() {
	log.Printf("GitHub Process Monitor started ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	collector := engine.NewCollector(ctx)
	api.SpawnAPI(ctx, collector)
	engine.SpawnPollers(ctx, engine.ReadJobs(), collector)
	<-ctx.Done()
	log.Printf("GitHub Process Monitor done!")
}
