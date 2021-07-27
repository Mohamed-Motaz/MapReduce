package main

//
// start a worker process, which is implemented
// in ../mr/worker.go. typically there will be
// multiple worker processes, talking to one coordinator.
//
// go run mrworker.go wc.so

import (
	"github.com/mohamed247/mapReduce/mr"
)

func main() {
	mr.Worker(Map, Reduce)
}
