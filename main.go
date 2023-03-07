package main

import (
	"streak/app"
	"streak/app/common"
)

func main() {
	common.InitConfig()

	common.InitLogger()

	common.InitInfluxdb()

	app.RunAgent()

}
