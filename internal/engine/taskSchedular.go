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
	//Task sender go Routine
	go func() {
		for msg := range sendChannel {
			sendLock.Lock()
			if err := stream.Send(msg); err !=nil {
				//handle fail send 
				taskWait.Done()
				log.Println("Send Error: ",err.Error())
			}
			log.Println("task send: ",msg.Id)
			sendLock.Unlock()
		}
	}()

	// send zero degree task to send channel
	for task := range graph {
		if inDegree[task] == 0 {
			taskWait.Add(1)
			sendChannel <- &workflow.TaskDetail{
				Id: task,
			}
		}
	}

	// receive finished task
	go func() {
		defer close(recivChannel)
		for {
			data,err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Print("Receive Error: ",err.Error())
				return
			}
			recivChannel <- data
		}
	}()

	go func() {
		for msg := range recivChannel {
			if msg.Status == "Completed" {
				//change this to a method implementation
				for _,task := range graph[msg.Id] {
					inDegree[task]--
					if inDegree[task] == 0 {
						taskWait.Add(1)
						sendChannel <- &workflow.TaskDetail{
							Id: task,
						}
					}
				}
			}else{
				// Handle retry logic here 
				// temp retry
				taskWait.Add(1)
				sendChannel <- &workflow.TaskDetail{
					Id: msg.GetId(),
				}
			}
			log.Println("Task Completed : ",msg.Id)
			taskWait.Done()
		}		
	}()
	taskWait.Wait()
	stream.CloseSend()
}

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
func (wfi *WorkflowInfo) UpdateInDegree(taskWait *sync.WaitGroup)  {
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