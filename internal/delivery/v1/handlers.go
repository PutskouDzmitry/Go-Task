package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func GetRoutes(engine *gin.Engine) *gin.RouterGroup {

	// Create public API routes.
	public := engine.Group("/api/v1")

	// Serve Swagger UI.
	public.GET("/swagger", redirect)
	public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return public
}

// redirect the user to the API documentation page using a 301 redirect
func redirect(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/api/v1/swagger/index.html")
}
