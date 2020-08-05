package main

import (
	"testing"
	"time"
)

func Test_sync(t *testing.T) {

	InitSyncFramework(10 * time.Second)
}
