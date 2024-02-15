package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	initCli()
	fmt.Println("docker cli OK")
	r := gin.Default()

	r.POST("/docker/start", func(c *gin.Context) {
		var requestBody RequestCreateContainer

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := createContainer(requestBody)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Docker container started successfully"})
	})

	r.GET("/docker/runnings", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{"runningContainers": getAllDockersRunning()})
	})

	// Run the server
	port := 8080
	fmt.Printf("Server is running on :%d\n", port)
	r.Run(fmt.Sprintf(":%d", port))
}
