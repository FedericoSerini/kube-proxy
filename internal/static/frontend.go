package static

import (
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AddFrontEndRoute(router *gin.Engine) {
	router.GET("/", exposeFrontEnd)
}

func exposeFrontEnd(c *gin.Context) {
	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "template error")
		return
	}

	data := map[string]string{
		"KeycloakBase":     os.Getenv("KEYCLOAK_BASE_URL"),
		"KeycloakRealm":    os.Getenv("KEYCLOAK_REALM"),
		"KeycloakClientID": os.Getenv("KEYCLOAK_CLIENT_ID"),
	}

	c.Header("Content-Type", "text/html")
	tmpl.Execute(c.Writer, data)
}
