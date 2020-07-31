package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/siddontang/go-log/log"
	"github.com/spf13/viper"
	"os"
)

var SupportSchema []string
var MysqlInfo map[string]interface{}

func init() {

	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatalln("读取配置文件失败")
		os.Exit(0)
	}

	//监控的库名
	SupportSchema = viper.GetStringSlice("schema")

	//读取mysql相关数据
	MysqlInfo = viper.GetStringMap("mysqlInfo")

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {

		//库名进行热更新
		SupportSchema = viper.GetStringSlice("schema")
		log.Warnln("配置文件进行了更改，如果更改的是连接配置项，请重启服务")
	})
}
