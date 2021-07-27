# My implementation of MapReduce based on the labs of MIT 6.824 Spring 2021

6.824 is a core graduate subject with lectures, labs, an optional project, a mid-term exam, and a final exam. 6.824 is 12 units.

This project is a simple MapReduce that simply counts the occurence of all unique words in given text files and runs the computation in parallel, 
on as many workers as you'd like. Currently, the machines communicate via RPC over UNIX sockets. The next step is for them to communicate via 
TCP/IP and be on different machines, and therefore use a distributed file sharing system.

## Deploy
Install golang, and setup golang environment variables and directories. [Click here](https://golang.org/doc/install) to learn it.

```
cd $GOPATH
git clone https://github.com/Mohamed247/MapReduce.git
cd MIT-6.824
export GOPATH=$GOPATH:$(pwd)
```

## Run
To run the MapReduce, just type the following command after you cd into src/main

TestTextFiles folder should contain all the .txt files that you wish to run the MapReduce job on

```
go run -race 5 mrcoordinator.go TestTextFiles/pg-*.txt
```

This creates one coordinator (master). Please avoid creating more than one master at a time

To alter the number of reduce tasks, just change the number after "-race" in the above command to the number you'd like

To create a single worker, run the command
```
go run mrworker.go
```

## Output
The output files will be in the main project directory with the following format

```
mr-out-*
```

where * is the number of the reduce job, starting from 0


Remeber, you can run this on as many workers as you'd want! Have fun!
