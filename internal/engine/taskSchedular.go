package engine

import (
	"context"
	"io"
	"log"
	"sync"

	"github.com/Harish-Naruto/Distributed-Workflow-Engin/gen/go/workflow"
	"google.golang.org/grpc"
)

type WorkflowInfo struct {
	Stream grpc.BidiStreamingClient[workflow.TaskDetail, workflow.TaskStatus]
	Workflow map[string][]string
	SendChannel chan *workflow.TaskDetail
	RecivChannel chan *workflow.TaskStatus
	InDegree map[string]int	

}

// TaskSchedular schedules workflow tasks based on their dependencies.
// It sends ready tasks to the executor and schedules dependent tasks
// as their prerequisites complete.
func TaskSchedular(conn workflow.TaskServiceClient,graph map[string][]string){
	stream, err := conn.AssignTask(context.Background())
	if err != nil {
		log.Print("error at taskSchedular: ",err.Error())
		return
	}
	
	var sendLock sync.Mutex
	var taskWait sync.WaitGroup

	sendChannel := make(chan *workflow.TaskDetail)
	defer close(sendChannel)
	recivChannel := make(chan *workflow.TaskStatus)

	// Preprocess this indegree and store it in db
	inDegree := make(map[string]int)
	for dependency := range graph {
		for _,task := range graph[dependency] {
			inDegree[task]++;
		}
	}

	workflowInfo := WorkflowInfo{
		Stream: stream,
		Workflow: graph,
		SendChannel: sendChannel,
		RecivChannel: recivChannel,
		InDegree: inDegree,
	}
	
	go workflowInfo.WriteStream(&sendLock,&taskWait)

	go workflowInfo.ReadStream(&taskWait)

	// send first zero degree task to send channel
	for task := range graph {
		if inDegree[task] == 0 {
			taskWait.Add(1)
			sendChannel <- &workflow.TaskDetail{
				Id: task,
			}
		}
	}

	go workflowInfo.UpdateDependency(&taskWait)
	
	taskWait.Wait()
	stream.CloseSend()
}

/*
 WriteStream reads tasks from SendChannel and sends them to the
 executor through the gRPC stream.
*/
func (wfi *WorkflowInfo) WriteStream(sendLock *sync.Mutex, taskWait *sync.WaitGroup)  {
	for task := range wfi.SendChannel {
		sendLock.Lock()
		if err := wfi.Stream.Send(task); err != nil {
			log.Println("Send Error: ",err.Error())
			taskWait.Done()
		}	
		sendLock.Unlock()
	}
}

/*
ReadStream receives task execution status from the gRPC stream
and forwards it to RecivChannel.
*/
func (wfi *WorkflowInfo) ReadStream(taskWait *sync.WaitGroup)  {
	defer close(wfi.RecivChannel)
	for {
		status, err := wfi.Stream.Recv()
		if err == io.EOF{
			break
		}
		if err != nil {
			log.Print("Receive Error: ",err.Error())
			return
		}
		wfi.RecivChannel <- status
	}	
}
/*
UpdateDependency updates task dependencies after execution.
When a task's dependencies are resolved, it is sent for execution.
Failed tasks are rescheduled.
*/
func (wfi *WorkflowInfo) UpdateDependency(taskWait *sync.WaitGroup)  {
	for task := range wfi.RecivChannel {
		if task.Status == "Completed" {
			for _,node := range wfi.Workflow[task.Id] {
				wfi.InDegree[node]--
				if wfi.InDegree[node] == 0 {
					taskWait.Add(1)
					wfi.SendChannel <- &workflow.TaskDetail{
						Id: node,
					}
				}

			}
		}else {
			//handle retry logic here 
			taskWait.Add(1)
			wfi.SendChannel <- &workflow.TaskDetail{
				Id: task.Id,
			}
		}
		taskWait.Done()
	}
	
}