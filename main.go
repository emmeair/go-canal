package main

import (
	"canal/app"
	_ "canal/config"
	"canal/db"
	_ "canal/db"
	_ "canal/server"
)

func main() {

	defer closeMain()

	go app.InitSyncFramework()

	select {}
}

func closeMain() {

	db.LocalClose()
}
