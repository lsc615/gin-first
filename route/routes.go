package route

import (
	"github.com/gin-gonic/gin"
	"github.com/shicli/gin-first/controller"
	_ "github.com/shicli/gin-first/docs"
	"github.com/shicli/gin-first/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"html/template"
	"net/http"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware())
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, "go 1.20")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.SetFuncMap(template.FuncMap{
		"safe": func(str string) template.HTML {
			return template.HTML(str)
		},
	})
	r.LoadHTMLFiles("templates/index.html")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", "<a href='https://liwenzhou.com'>李文周的博客</a>")
	})
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	// 中间件AuthMiddleware实现用户的认证，保护路由
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info)

	return r
}
