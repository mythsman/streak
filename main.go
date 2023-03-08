package main

import (
	"streak/app"
	"streak/app/cache"
	"streak/app/common"
)

func main() {
	common.InitConfig()

	common.InitLogger()

	cache.InitCache()

	common.InitInfluxdb()

	app.RunAgent()
}

