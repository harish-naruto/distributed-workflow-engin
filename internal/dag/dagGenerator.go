package dag

import (
	"errors"

	"github.com/Harish-Naruto/Distributed-Workflow-Engin/internal/models"
)

/*
DagGenerator generate a graph for workflow in map[string][]string
*/
func DagGenerator(workflow models.Workflow, graph map[string][]string) error {
	// add all nodes to map
	for _,i := range workflow.Tasks {
		graph[i.Id] = []string{}
	}
	// add edges to nodes
	for _,i := range workflow.Tasks {
		for _,j:= range i.DependOn{
			if _,ok := graph[j];ok{
				graph[i.Id] = append(graph[i.Id], j)
			}else{
				return errors.New("Invalid dependency for a task")
			}
		}
	}
	return nil

}