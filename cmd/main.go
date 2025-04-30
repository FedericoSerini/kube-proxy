package main

import (
	"fmt"
	kustomizationRoutes "kube-proxy/internal/kustomizations/v1/api"
	podsRoutes "kube-proxy/internal/pods/v1/api"
	"kube-proxy/internal/static"
	"kube-proxy/pkg/auth"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	getRoutes().Run(":8080")
}

func getRoutes() *gin.Engine {
	router := gin.Default()
	static.AddFrontEndRoute(router)

	v1 := router.Group("/v1")
	podsRoutes.AddPodsRoutes(v1)
	kustomizationRoutes.AddKustomizationRoutes(v1)
	initKeycloakOidc(v1)
	return router
}

func initKeycloakOidc(apiGroup *gin.RouterGroup) {
	issuer := os.Getenv("KEYCLOAK_ISSUER")
	clientID := os.Getenv("KEYCLOAK_CLIENT_ID")

	if issuer == "" || clientID == "" {
		panic("OIDC config missing")
	}

	if err := auth.InitKeycloakAuth(issuer, clientID); err != nil {
		panic(fmt.Sprintf("failed to initialize Keycloak authentication: %v", err))
	}

	apiGroup.Use(auth.AuthMiddleware())
}
