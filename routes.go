package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.GET("/version", func(c *gin.Context) {
		c.String(http.StatusOK, "go 1.19")
	})
	return r
}
