package main

//
// start a worker process, which is implemented
// in ../mr/worker.go. typically there will be
// multiple worker processes, talking to one coordinator.
//
// go run mrworker.go wc.so

import (
	functionality "github.com/mohamed247/mapReduce/Functionality"
	"github.com/mohamed247/mapReduce/mr"
)

func main() {
	mr.Worker(functionality.Map, functionality.Reduce)
}
