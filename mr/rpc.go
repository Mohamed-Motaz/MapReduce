package mr

//
// RPC definitions.
//

import (
	"os"
	"strconv"
)

const (
	Map = 1
	Reduce = 2
	Done = 3
)

type KeyValue struct{
	Key   string
	Value string
}

type TaskType int;

//no fields needed for get task args from worker to coordinator
type GetTaskArgs struct{}

type GetTaskReply struct{
	//type of the task, (map, reduce, or done)
	TaskType TaskType

	//number of the task, for both map and reduce
	TaskNum int

	//Name of file, for map
	FileName string

	//Number of reduce tasks, for the map 
	NumReduceTasks int

	//Number of map tasks, for the reduce 
	NumMapTasks int
	
	//When true, the worker exits as there are no more tasks to be given
	Done bool
}

type FinishedTaskArgs struct{
	TaskType TaskType
	TaskNum int
}

//No reply needed from the coordinators to the worker when he finished
type FinishedTaskReply struct{}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the coordinator.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func coordinatorSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
