package main

import (
	"log"

	"github.com/Harish-Naruto/Distributed-Workflow-Engin/gen/go/workflow"
	"github.com/Harish-Naruto/Distributed-Workflow-Engin/internal/engine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)
func main() {
	conn, err := grpc.NewClient(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("error while connection to server: ",err.Error())
	}
	defer conn.Close()

	taskConnecetion := workflow.NewTaskServiceClient(conn)
	
	graph := map[string][]string{
		"A": {"B", "C"},
		"B": {"D", "E"},
		"C": {"F"},
		"D": {"G"},
		"E": {"G"},
		"F": {"H"},
		"G": {},
		"H": {},
	}
	engine.TaskSchedular(taskConnecetion,graph)


}