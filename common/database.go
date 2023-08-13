package common

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	//"github.com/jinzhu/gorm"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() (DB *gorm.DB) {
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, database)
	println(dsn)
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Errorf("连接数据库失败：%s", err))
	}
	return DB
	//// 创建数据对象
	//probe := Probe{
	//	Name: "c",
	//	Age:  3,
	//}
	//
	//// 插入数据
	//result := db.Create(&probe)
	//if result.Error != nil {
	//	log.Fatalf("无法插入数据：%v", result.Error)
	//}

	//// 监控配置文件变化
	//viper.WatchConfig()
	//// 注意！！！配置文件发生变化后要同步到全局变量Conf
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	fmt.Println("配置文件被修改...")
	//	if err := viper.Unmarshal(&config.yml); err != nil {
	//		panic(fmt.Errorf("unmarshal conf failed, err:%s \n", err))
	//	}
	//})
}

func GetDB() *gorm.DB {
	return DB
}
