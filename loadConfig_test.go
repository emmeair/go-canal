package main

import (
	"github.com/spf13/viper"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()

	if err != nil {
		t.Fatal("读取配置文件失败")
	}
}
