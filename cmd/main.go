package main

import (
	"github.com/spf13/viper"
	"myproject/Apigateway/infrastructure/mlog"
)

func main(){
	mlog.ConfigLocalFilesystemLogger(viper.GetString("log.path"),viper.GetString("log.level"))
	rout
}
