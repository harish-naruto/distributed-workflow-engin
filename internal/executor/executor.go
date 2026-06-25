package executor

import (
	"time"

	"github.com/Harish-Naruto/Distributed-Workflow-Engin/gen/go/workflow"
)

func Executor(Task *workflow.TaskDetail) *workflow.TaskStatus{
	time.Sleep(3*time.Second)
	return &workflow.TaskStatus{
		Id: Task.Id,
		Status: "Completed",
	}
}