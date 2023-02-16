package main

import (
	"streak/app"
)

func main() {
	app.InitLogger()

	app.InitConfig()

	app.InitInfluxdb()

	app.RunAgent()

}
