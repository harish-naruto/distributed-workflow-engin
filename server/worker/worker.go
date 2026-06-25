package main

import (
	"log"
	"net"

	"github.com/Harish-Naruto/Distributed-Workflow-Engin/gen/go/health"
	"github.com/Harish-Naruto/Distributed-Workflow-Engin/gen/go/workflow"
	grpcservices "github.com/Harish-Naruto/Distributed-Workflow-Engin/internal/grpcServices"
	"google.golang.org/grpc"
)

type Server struct {
	health.UnimplementedHealthServiceServer
}

func main() {
	lis, err := net.Listen("tcp",":9000")
	if err != nil {
		log.Fatalf("fail to listen on port 9000: %s",err.Error())
	}
	
	grpcServer := grpc.NewServer()
	s := Server{} 
	ts := &grpcservices.TaskService{}

	health.RegisterHealthServiceServer(grpcServer,&s)
	workflow.RegisterTaskServiceServer(grpcServer,ts)	
	
	if err := grpcServer.Serve(lis); err!=nil {
		log.Fatalln("failed to start grpc server, err: ",err)
	}

}