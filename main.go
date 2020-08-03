package main

import (
	"encoding/json"
	"errors"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/spf13/viper"
	"net"
	"time"
)

var conn net.Conn

func main() {

	defer closeMain()

	//连接服务器
	go LinkServer()

	cfg := canal.NewDefaultConfig()
	cfg.Addr = MysqlInfo["addr"].(string)
	cfg.User = MysqlInfo["user"].(string)
	cfg.Password = MysqlInfo["password"].(string)
	cfg.HeartbeatPeriod = 200 * time.Millisecond
	cfg.ServerID = 1

	c, err := canal.NewCanal(cfg)

	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	defer c.Close()

	// Register a handler to handle RowsEvent
	c.SetEventHandler(&MyEventHandler{})

	var position mysql.Position
	position, err = c.GetMasterPos()

	if err != nil {

		log.Fatalln("无法读取log文件位置")
		return
	}

	//读取本地数据坐标，如果有则使用本地坐标数据,恢复上次暂停的进度
	data, err := localDb.Get([]byte("mysql_pos"), nil)

	if err == nil {

		var mysqlPos mysql.Position

		err = json.Unmarshal(data, &mysqlPos)

		if err == nil {

			position = mysqlPos
		}
	}

	// Start canal
	_ = c.RunFrom(position)
}

//主动连接服务器
func LinkServer() {
	serverInfo := viper.GetStringMap("server")

	for {

		var err error
		conn, err = net.DialTimeout(serverInfo["network"].(string), serverInfo["addr"].(string), 3*time.Second)
		if err != nil {

			log.Warnln("连接服务器失败,正在重试...")
		} else {

			log.Infoln("连接推送服务器成功!!!")
			doRead(conn)
			if conn != nil {

				_ = conn.Close()
			}
		}
		conn = nil
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

type MyEventHandler struct {
	canal.DummyEventHandler
}

func (h *MyEventHandler) OnRow(ev *canal.RowsEvent) error {

	next := false

	for _, supportable := range SupportSchema {

		if ev.Table.Schema == supportable {

			next = true
			break
		}
	}

	if !next {
		//如果配置未找到需要读取的表，则直接跳过
		return nil
	}

	sendRow := make(map[string]interface{})

	sendRow["SchemaName"] = ev.Table.Schema //库名
	sendRow["TableName"] = ev.Table.Name    //表名
	sendRow["Action"] = ev.Action           //行为

	//此处是参考 https://github.com/gitstliu/MysqlToAll 里面的获取字段和值的方法
	ColumnData := make(map[string]interface{})

	for columnIndex, currColumn := range ev.Table.Columns {

		//字段名，字段对应的值
		ColumnData[currColumn.Name] = ev.Rows[len(ev.Rows)-1][columnIndex]

	}

	sendRow["ColumnData"] = ColumnData

	go func(sendRow map[string]interface{}) {

		s, err := json.Marshal(sendRow)
		s = append(s, []byte("\n")...)
		if err != nil {

			return
		}

		if conn != nil {

			_, err := conn.Write(s)

			if err != nil {

				log.Infoln("向服务器推送数据失败")
				return
			}
			log.Infoln("数据库数据已推送成功")
		}

		return

	}(sendRow)

	return nil
}

func (h *MyEventHandler) OnXID(p mysql.Position) error {

	strData, err := json.Marshal(p)

	if err != nil {

		return errors.New("数据库读取坐标转化JSON失败")
	}

	err = localDb.Put([]byte("mysql_pos"), strData, nil)

	if err != nil {

		return errors.New("记录数据库坐标失败")
	}

	return nil
}

func closeMain() {

	LocalClose()
}
