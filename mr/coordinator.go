package mr

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
	"time"
)


type Coordinator struct {

	mu sync.Mutex
	cond *sync.Cond

	NumReduceTasks int
	MapFiles []string


	MapTasksFinished []bool
	MapTasksIssued []time.Time

	ReduceTasksFinished []bool
	ReduceTasksIssued []time.Time

	isDone bool
}

//RPC handlers for the worker to call.

func (c *Coordinator) HandleGetTasks(args *GetTaskArgs, reply *GetTaskReply) error{
	c.mu.Lock()
	defer c.mu.Unlock()

	reply.NumMapTasks = len(c.MapFiles)
	reply.NumReduceTasks = c.NumReduceTasks

	//Issue all map tasks
	for {
		mapPhaseDone := true;
		for m, done := range c.MapTasksFinished{
			//only check tasks that havent been finished yet
			if !done{
				//all mapTasks are by default zeroed, or more than 10 seconds have passed since I heard back from them
				if c.MapTasksIssued[m].IsZero() ||
					time.Since(c.MapTasksIssued[m]).Seconds() > 10 {
					reply.FileName = c.MapFiles[m];
					reply.TaskType = Map;
					reply.TaskNum = m;
					c.MapTasksIssued[m] = time.Now()
					return nil; //succesful task given
				}else{   //there is a woker who is currently working on a map task
					mapPhaseDone = false
				}
			}
		}
		if !mapPhaseDone{
			//wait
			c.cond.Wait()
		}else{
			break
		}
	}//All map tasks are done

	//Issure all reduce tasks
	//fmt.Println("About to start reduce tasks")
	for {
		reducePhaseDone := true;
		for m, done := range c.ReduceTasksFinished{
			//only check tasks that havent been finished yet
			if !done{
				//all mapTasks are by default zeroed, or more than 10 seconds have passed since I heard back from them
				if c.ReduceTasksIssued[m].IsZero() ||
					time.Since(c.ReduceTasksIssued[m]).Seconds() > 10 {
					reply.TaskType = Reduce;
					reply.TaskNum = m;
					c.ReduceTasksIssued[m] = time.Now()
					return nil;
				}else{   //there is a woker who is currently working on a map task
					reducePhaseDone = false
				}
			}
		}
		if !reducePhaseDone{
			//wait
			c.cond.Wait()
		}else{
			break
		}
	}//All reduce tasks are done
	reply.TaskType = Done;
	c.isDone = true;

	return nil;
}

func (c *Coordinator) HandleFinishTasks(args *FinishedTaskArgs, reply *FinishedTaskReply) error{
	c.mu.Lock()
	defer c.mu.Unlock()

	switch args.TaskType{
	case Map:
		c.MapTasksFinished[args.TaskNum] = true
	case Reduce:
		c.ReduceTasksFinished[args.TaskNum] = true
	default:
		log.Fatalf("bad task type %v", args.TaskType);
	}
	//wake up task handler and notify him that the task is finished
	c.cond.Broadcast()
	return nil;
}

//
// start a thread that listens for RPCs from worker.go
//
func (c *Coordinator) server() {
	rpc.Register(c)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := coordinatorSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrcoordinator.go calls Done() periodically to find out
// if the entire job has finished.
//
func (c *Coordinator) Done() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.isDone{
		//remove all intermediate files
		c.removeAllIntermediateFiles()
	}
	return c.isDone
}

//
// create a Coordinator.
// main/mrcoordinator.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeCoordinator(files []string, nReduce int) *Coordinator {
	c := Coordinator{}

	c.cond = sync.NewCond(&c.mu)
	//fmt.Println("These are the files", files)
	c.NumReduceTasks = nReduce;
	c.MapFiles = files;
	
	c.MapTasksFinished = make([]bool, len(files))
	c.MapTasksIssued = make([]time.Time, len(files))
	
	c.ReduceTasksFinished = make([]bool, nReduce)
	c.ReduceTasksIssued = make([]time.Time, nReduce)

	//wake up task handler thread every second
	go func(){
		for {
			c.mu.Lock()
			c.cond.Broadcast()
			c.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	c.server()
	return &c
}

func (c *Coordinator) removeAllIntermediateFiles(){
	for m := 0; m < len(c.MapFiles); m++{
		for r := 0; r < c.NumReduceTasks; r++{
			fileName := fmt.Sprintf("mr-%d-%d", m, r);
			err := os.Remove(fileName)
			if err != nil{
				log.Fatal(err)
			}
		}
	}
}