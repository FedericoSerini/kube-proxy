package api

import (
	"fmt"
	"kube-proxy/internal/pods/v1/models"
	"kube-proxy/pkg/kubernetes"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AddPodsRoutes(c *gin.RouterGroup) {
	podsRoutes := c.Group("/pods")
	podsRoutes.GET("/metadata", getAllPodsMetadata)
	podsRoutes.Static("/view", "./static")
}

func getAllPodsMetadata(c *gin.Context) {
	clientset, err := kubernetes.CreateClient()
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	pods, err := clientset.CoreV1().Pods("").List(c, v1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching pods: %v", err)})
		return
	}

	var podList []models.Pod
	for _, pod := range pods.Items {
		podList = append(podList, models.Pod{
			Namespace: pod.Namespace,
			Name:      pod.Name,
			Status:    string(pod.Status.Phase),
		})
	}

	c.JSON(http.StatusOK, podList)
}
