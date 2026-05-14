package main

import (
	"log"
	"github.com/Harish-Naruto/Distributed-Workflow-Engin/internal/api"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "server running",
		})
	})

	router.POST("/send", func(ctx *gin.Context) {
		file, err := ctx.FormFile("workflow")
		if err != nil {
			log.Println("error{ while uploading file:  ", err)
			ctx.JSON(400, gin.H{
				"error": "file not found",
			})
			return
		}

		if err := api.WorkflowUploadHandler(file);err != nil {
			ctx.JSON(400,gin.H{
				"error":err.Error(),
			})
			return
		}
		ctx.JSON(200, gin.H{
			"file": file,
		})
	})

	router.Run()
}
