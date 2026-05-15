package dag

import (
	"errors"
)

func ExecutionOrderGenerator(graph map[string][]string) ([][]string,error) {
	var order [][]string
	inDegree := make(map[string]int)

	for k := range graph {
		for _,j:= range graph[k] {
			inDegree[j] += 1;
		}
	}
	var queue []string

	for i:= range graph {
		if inDegree[i]== 0 {
			queue = append(queue, i)
		}
	}

	total := 0

	for len(queue)>0 {
		total += len(queue)
		order = append(order, queue)
		NextQueue := []string{}
		for _, node := range queue {
			for _,i:= range graph[node] {
				inDegree[i]--
				if inDegree[i] == 0{
					NextQueue = append(NextQueue, i)

				}
			}
		}
		queue = NextQueue
	}
	if total != len(graph) {
		return nil, errors.New("Invalid workflow, it contain cycle")
	}
	return order,nil

}