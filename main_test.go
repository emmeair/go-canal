package main

import (
	"github.com/spf13/viper"
	"net"
	"testing"
	"time"
)

func TestLinkServer(t *testing.T) {

	serverInfo := viper.GetStringMap("server")

	var err error
	conn, err = net.DialTimeout(serverInfo["network"].(string), serverInfo["addr"].(string), 3*time.Second)
	if err != nil {

		t.Log("连接服务器失败,正 在重试...")
	} else {

		t.Log("连接推送服务器成功!!!")
		if conn != nil {

			_ = conn.Close()
		}
	}
	conn = nil
}
