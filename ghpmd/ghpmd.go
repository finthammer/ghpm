package main

import (
	"fmt"
	"time"

	"github.com/themue/ghpm/analyze"
	"github.com/themue/ghpm/github"
)

func main() {
	fmt.Println("GitHub Process Monitor")
	e := github.NewRepoEventor("themue", "ghpm")
	events, err := e.Get()
	if err != nil {
		fmt.Printf("error %v\n", err)
	}
	a, err := analyze.Events(events, analyze.TypeCounter)
	if err != nil {
		fmt.Printf("error %v\n", err)
	}
	fmt.Printf("1st events %v\n", a)
	time.Sleep(time.Second)
	events, err = e.Get()
	if err != nil {
		fmt.Printf("error %v\n", err)
	}
	a, err = analyze.Events(events, analyze.TypeCounter)
	if err != nil {
		fmt.Printf("error %v\n", err)
	}
	fmt.Printf("2nd events %v\n", a)
}
