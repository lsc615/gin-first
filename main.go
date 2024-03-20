package main

import (
	"fmt"
	"github.com/shicli/gin-first/common"

	"github.com/gin-gonic/gin"
	"github.com/shicli/gin-first/route"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// @title			gin first
// @version		1.0
// @description	This is gin
// @contact.name	shicli
func main() {
	InitConfig()
	common.InitDB()
	r := gin.Default()
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
