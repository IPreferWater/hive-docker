package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	err := initDockerEngine()

	if err != nil {
	panic(err)
	}

	r := gin.Default()
	r.POST("/docker/start", func(c *gin.Context) {
		var requestBody RequestCreateContainer

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := dEngine.createContainer(requestBody)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Docker container started successfully"})
	})

	r.GET("/docker/runnings", func(c *gin.Context) {

		allDocker, err := dEngine.getAllDockersRunning()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"runningContainers": allDocker})
	})

	port := 8080
	r.Run(fmt.Sprintf("localhost:%d", port))
}
