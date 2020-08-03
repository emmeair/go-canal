package main

import (
	"github.com/syndtr/goleveldb/leveldb"
	"testing"
)

func TestLocal(t *testing.T) {

	var err error
	localDb, err = leveldb.OpenFile("source/local/db", nil)

	if err != nil {

		t.Log(err.Error())
	}

}
