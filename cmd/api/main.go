package main

import (
	"musicAPI/server"
)

func main() {
	app := server.NewApp()
	app.Run()
}
