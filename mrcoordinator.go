package main

//
// start the coordinator process, which is implemented
// in ../mr/coordinator.go
//
// go run mrcoordinator.go pg*.txt
//

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mohamed247/mapReduce/mr"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: mrcoordinator inputfiles...\n")
		os.Exit(1)
	}
	//fmt.Println("Coordinator is running.....")
	reduceTasksNum, err := strconv.Atoi(os.Args[1]);
	if err != nil {
        log.Fatal("Reduce tasks to run needs to be an integer")
    }
	m := mr.MakeCoordinator(os.Args[2:], reduceTasksNum)
	for m.Done() == false {
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)
}
