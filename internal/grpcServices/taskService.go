package grpcservices

import (
	"io"
	"log"
	"sync"

	"github.com/Harish-Naruto/Distributed-Workflow-Engin/gen/go/workflow"
	"github.com/Harish-Naruto/Distributed-Workflow-Engin/internal/executor"
	"google.golang.org/grpc"
)

type TaskService struct {
	workflow.UnimplementedTaskServiceServer
}

/*
AssignTask receive task from Grpc stream.
Assign task to executor engine for execution
*/
func (Ts *TaskService) AssignTask(stream grpc.BidiStreamingServer[workflow.TaskDetail,workflow.TaskStatus]) error {
	var wg sync.WaitGroup
	tasks := make(chan *workflow.TaskDetail)
	go func(){
		defer close(tasks)
		for {	
			data,err := stream.Recv()
			if err == io.EOF{
				break;
			}
			if err != nil {
				log.Println("Error at Assign task",err.Error())
				return 
			}
			tasks <- data
			log.Println("Recieved Task: ",data)
		}
	}()

	var SendMux sync.Mutex
	for task := range tasks{
		wg.Add(1)
		go func(t *workflow.TaskDetail) {
			defer wg.Done()
			ts := executor.Executor(t)
			log.Println("task executed: ",ts.Id)
			SendMux.Lock()
			defer SendMux.Unlock()
			stream.Send(ts)
		}(task)
	}
	wg.Wait()
	return nil

}