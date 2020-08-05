package server

import (
	"canal/config"
	"github.com/siddontang/go-log/log"
	"github.com/spf13/viper"
	"net"
	"time"
)

var Conn net.Conn

func init() {

	//主动连接服务器
	switch config.ServerInfo["network"].(string) {

	case "tcp":

		go LinkTcpServer()

		break

	default:

		log.Println("推送服务器配置错误")

		break
	}

}

func LinkTcpServer() {
	serverInfo := viper.GetStringMap("server")

	for {

		var err error
		Conn, err = net.DialTimeout(serverInfo["network"].(string), serverInfo["addr"].(string), 3*time.Second)
		if err != nil {

			log.Warnln("连接服务器失败,正在重试...")
		} else {

			log.Infoln("连接推送服务器成功!!!")
			doRead(Conn)
			if Conn != nil {

				_ = Conn.Close()
			}
		}
		Conn = nil
		time.Sleep(3 * time.Second)

	}
}

func doRead(conn net.Conn) {
	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])

		if err != nil {

			log.Warnln("推送服务断线啦,进行重连")
			break
		}
		if n > 0 {

			log.Println(string(buf[0:n]))
		}
	}

}
