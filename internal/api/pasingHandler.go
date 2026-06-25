package api

import (
	"errors"
	"io"
	"log"
	"mime/multipart"

	"github.com/Harish-Naruto/Distributed-Workflow-Engin/internal/dag"
	"github.com/Harish-Naruto/Distributed-Workflow-Engin/internal/models"
	"github.com/Harish-Naruto/Distributed-Workflow-Engin/internal/validator"
	"gopkg.in/yaml.v3"
)

func WorkflowUploadHandler(file *multipart.FileHeader) error {
	
	// Extraccting file content
	content,err := ExtractContent(file)
	if err!= nil {
		return err
	}
	//store workflow in-memory
	workflow,WorkflowErr := ParseYAML(content)
	if WorkflowErr != nil {
		log.Println("Parsing error: ",WorkflowErr)
		return errors.New("Error while parsing file")
	}

	// generate garph and validate the graph
	graph := make(map[string][]string)
	if err := dag.DagGenerator(workflow,graph); err!= nil {
		return errors.New(err.Error())
	}

	order,err := dag.ExecutionOrderGenerator(graph)
	if err != nil {
		return err
	}
	log.Print("order of execution for this graph is : ", order)
	return nil
}


/*
Validate YAML multipart file and create a byte slice for the file
content, if any error then return error with nil slice
*/
func ExtractContent(file *multipart.FileHeader) ([]byte,error) {
	if err := validator.ValidateYML(file);err!=nil{
		return nil,err
	}
	
	yamlOpen,_ := file.Open()
	defer yamlOpen.Close()

	content,err := io.ReadAll(yamlOpen)

	if err!=nil{
		return nil,err
	}
	return content,nil
}
/*
Parse yaml file into Workflow struct
*/	
func ParseYAML(Yamlfile []byte) (models.Workflow,error) {
	var temp models.Workflow
	if err := yaml.Unmarshal(Yamlfile,&temp); err!=nil {
		return temp,err
	}
	return temp,nil	
}
