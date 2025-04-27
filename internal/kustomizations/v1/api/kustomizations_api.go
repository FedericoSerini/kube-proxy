package api

import (
	"fmt"
	"kube-proxy/internal/kustomizations/v1/models"
	"kube-proxy/pkg/kubernetes"
	"log"
	"net/http"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/gin-gonic/gin"
)

func AddKustomizationRoutes(c *gin.RouterGroup) {
	kustomizationRoutes := c.Group("/kustomizations")
	kustomizationRoutes.GET("/metadata", getAllKustomizationsMetadata)
}

func getAllKustomizationsMetadata(c *gin.Context) {
	_, dynamicClient, err := kubernetes.CreateClientCrd()
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %v", err)
	}

	resource := dynamicClient.Resource(
		schema.GroupVersionResource{
			Group:    "kustomize.toolkit.fluxcd.io",
			Version:  "v1beta1",
			Resource: "kustomizations",
		},
	)

	kustomizationList, err := resource.List(c, v1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching Kustomizations: %v", err)})
		return
	}

	var kustomizationListResp []models.Kustomization
	for _, item := range kustomizationList.Items {
		statusField := item.Object["status"]
		status := "Unknown"
		lastAppliedRevision := ""
		lastAttemptedRevision := ""
		message := ""

		if statusMap, ok := statusField.(map[string]interface{}); ok {
			// Get lastAppliedRevision if available
			if v, found := statusMap["lastAppliedRevision"].(string); found {
				lastAppliedRevision = v
			}
			// Get lastAttemptedRevision if available
			if v, found := statusMap["lastAttemptedRevision"].(string); found {
				lastAttemptedRevision = v
			}
			// Check conditions
			if conditions, found, _ := unstructured.NestedSlice(statusMap, "conditions"); found {
				for _, cond := range conditions {
					if condMap, ok := cond.(map[string]interface{}); ok {
						if condType, ok := condMap["type"].(string); ok && condType == "Ready" {
							if condStatus, ok := condMap["status"].(string); ok {
								status = condStatus
							}
							if condMessage, ok := condMap["message"].(string); ok {
								message = condMessage
							}
						}
					}
				}
			}
		}

		kustomizationListResp = append(kustomizationListResp, models.Kustomization{
			Namespace:             item.GetNamespace(),
			Name:                  item.GetName(),
			Status:                status,
			LastAppliedRevision:   lastAppliedRevision,
			LastAttemptedRevision: lastAttemptedRevision,
			Message:               message,
		})
	}
	c.JSON(http.StatusOK, kustomizationListResp)
}
