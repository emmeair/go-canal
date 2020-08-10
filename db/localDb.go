package db

import (
	"github.com/siddontang/go-log/log"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
)

var LocalDb *leveldb.DB

func init() {
	var err error
	LocalDb, err = leveldb.OpenFile("source/local/db", nil)

	if err != nil {
		log.Errorln(err.Error())
		os.Exit(0)
	}

}

func LocalClose() {

	_ = LocalDb.Close()
}
