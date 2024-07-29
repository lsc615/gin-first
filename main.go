package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shicli/gin-first/common"
	"github.com/shicli/gin-first/controller"
	"github.com/shicli/gin-first/route"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log"
	"time"
)

// @title			gin first
// @version		1.0
// @description	This is gin
// @contact.name	shicli
func main() {
	InitConfig()
	common.InitDB()

	cleanup, err := controller.InitTrace()
	log.Println(cleanup)
	if err != nil {
		log.Printf("Failed to initialize tracer: %v", err)
		return
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		// 清理 OTLP exporter
		if err := cleanup(ctx); err != nil {
			log.Printf("Failed to shut down OTLP exporter: %v", err)
		}
	}()

	r := gin.Default()
	r.Use(otelgin.Middleware(controller.ServiceName))
	r = route.CollectRoute(r)

	port := viper.GetString("server.port")
	if err := r.Run("127.0.0.1:" + port); err != nil {
		panic(err)
	}
}

func InitConfig() {
	var ConfigPath string
	pflag.StringVarP(&ConfigPath, "", "c", "", "配置文件路径")
	pflag.Parse()

	if ConfigPath == "" {
		fmt.Println("请提供配置文件路径")
		return
	}

	//设置文件
	viper.SetConfigFile(ConfigPath)

	//读取文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("无法读取配置文件：%s", err))
	}
}
