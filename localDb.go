package main

import (
	"github.com/siddontang/go-log/log"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
)

var localDb *leveldb.DB

func init() {

	var err error
	localDb, err = leveldb.OpenFile("source/local/db", nil)

	if err != nil {

		log.Errorln(err.Error())
		os.Exit(0)
	}

}

func LocalClose() {

	_ = localDb.Close()
}
