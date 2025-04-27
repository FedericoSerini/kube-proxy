package main

import (
	podsRoutes "kube-proxy/internal/pods/v1/api"

	"github.com/gin-gonic/gin"
)

func main() {
	getRoutes().Run(":8080")
}

func getRoutes() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	podsRoutes.AddPodsRoutes(v1)
	return router
}
