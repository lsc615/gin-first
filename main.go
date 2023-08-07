package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"http/common"
)

func main() {
	InitConfig()
	common.InitDB()
	r := gin.Default()
	r = CollectRoute(r)
	port := viper.GetString("server.port")
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
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
