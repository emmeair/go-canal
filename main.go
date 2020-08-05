package main

import (
	"canal/app"
	_ "canal/config"
	"canal/db"
	_ "canal/db"
	_ "canal/server"
	"os"
	"time"
)

func main() {

	defer closeMain()

	go InitSyncFramework(0)

	select {}
}

func InitSyncFramework(timeout time.Duration) {
	if timeout != 0 {
		go func() {

			for range time.After(timeout) {

				os.Exit(0)
			}
		}()
	}

	app.InitSyncFramework()
}

func closeMain() {

	db.LocalClose()
}
