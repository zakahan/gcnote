// -------------------------------------------------
// Package server
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package server

import (
	"go.uber.org/zap"
	"testing"
)

func TestInitFunctions(t *testing.T) {
	InitConfig()
	InitLogger()
	InitMysql()
	InitRedis()
	InitLocalCache()

}

func TestInitLogger(t *testing.T) {
	InitConfig()
	// 初始化日志
	InitLogger()
	// 测试日志记录
	zap.S().Debug("zz 这是一个调试消息")
	zap.S().Info("zz 这是一个信息消息")
	zap.S().Warn("zz 这是一个警告消息")
	zap.S().Error("zz 这是一个错误消息")

}
