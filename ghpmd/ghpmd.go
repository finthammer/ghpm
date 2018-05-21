package main

import (
	"context"
	"log"
	"time"

	"github.com/themue/ghpm/engine"
)

func main() {
	log.Printf("GitHub Process Monitor started ...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	collector := engine.NewCollector(ctx)
	engine.SpawnPollers(ctx, engine.ReadJobs(), collector)
	<-ctx.Done()
	log.Printf("GitHub Process Monitor done!")
}
