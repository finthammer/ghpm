package main

import (
	"fmt"

	"github.com/themue/ghpm/github"
)

func main() {
	fmt.Println("GitHub Process Monitor")
	e := github.NewRepoEventor("themue", "ghpm")
	events, err := e.Get()
	if err != nil {
		fmt.Printf("error %v\n", err)
	}
	fmt.Printf("events %v\n", events)
}
