package app

import (
	"canal/config"
	"canal/db"
	"canal/server"
	"encoding/json"
	"errors"
	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/canal"
	"github.com/siddontang/go-mysql/mysql"
	"os/exec"
	"time"
)

func InitSyncFramework() {

	var err error

	cfg := canal.NewDefaultConfig()
	cfg.Addr = config.MysqlInfo["addr"].(string)
	cfg.User = config.MysqlInfo["user"].(string)
	cfg.Password = config.MysqlInfo["password"].(string)
	cfg.HeartbeatPeriod = 200 * time.Millisecond
	cfg.ServerID = 1
	cfg.Dump.ExecutionPath = "mysqldump"

	_, err = exec.LookPath(cfg.Dump.ExecutionPath)

	if err != nil {

		log.Warnln("本地无法找到 " + cfg.Dump.ExecutionPath + " ,请检查环境变量")
		return
	}

	c, err := canal.NewCanal(cfg)

	if err != nil {

		log.Fatalln(err.Error())
		return
	}

	defer c.Close()

	c.SetEventHandler(&eventHandler{})

	var position mysql.Position
	position, err = c.GetMasterPos()

	if err != nil {

		log.Fatalln("无法读取log文件位置")
		return
	}

	//读取本地数据坐标，如果有则使用本地坐标数据,恢复上次暂停的进度
	data, err := db.LocalDb.Get([]byte("mysql_pos"), nil)

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

type eventHandler struct {
	canal.DummyEventHandler
}

func (h *eventHandler) OnRow(ev *canal.RowsEvent) error {

	next := false

	for _, supportable := range config.GetSupportSchema() {

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

		if server.Conn != nil {

			_, err := server.Conn.Write(s)

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

func (h *eventHandler) OnXID(p mysql.Position) error {

	strData, err := json.Marshal(p)

	if err != nil {

		return errors.New("数据库读取坐标转化JSON失败")
	}

	err = db.LocalDb.Put([]byte("mysql_pos"), strData, nil)

	if err != nil {

		return errors.New("记录数据库坐标失败")
	}

	return nil
}
