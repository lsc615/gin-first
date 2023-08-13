package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shicli/gin-first/route"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// @title			gin first
// @version		1.0
// @description	This is gin
// @contact.name	shicli
//
// @BasePath		controller
func main() {
	InitConfig()
	//common.InitDB()
	r := gin.Default()
	r = route.CollectRoute(r)
	port := viper.GetString("server.port")
	if err := r.Run(":" + port); err != nil {
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

	viper.SetConfigFile(ConfigPath)

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("无法读取配置文件：%s", err))
	}
}
