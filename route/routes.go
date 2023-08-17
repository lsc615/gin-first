package route

import (
	"github.com/gin-gonic/gin"
	"github.com/shicli/gin-first/controller"
	_ "github.com/shicli/gin-first/docs"
	"github.com/shicli/gin-first/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware())
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, "go 1.19")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/api/auth/register", controller.Register)
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	return r
}
