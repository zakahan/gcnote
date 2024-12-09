package main

import (
	"gcnote/server"
	"gcnote/server/config"
	"gcnote/server/router"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

func main() {
	server.InitConfig()
	server.InitLogger()
	server.InitMysql()
	server.InitRedis()
	server.InitLocalCache()
	route := router.InitRouter()
	app := &http.Server{
		Addr:           "0.0.0.0:" + strconv.Itoa(config.ServerCfg.Port),
		Handler:        route,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // http头部最大字节数1024
	}
	err := app.ListenAndServe()
	if err != nil {
		zap.S().Panicf("监听失败, err:%v", err.Error())
	}
}
