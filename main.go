package main

import (
	"streak/app"
	"streak/app/common"
)

func main() {
	common.InitLogger()

	common.InitConfig()

	common.InitInfluxdb()

	app.RunAgent()

}
