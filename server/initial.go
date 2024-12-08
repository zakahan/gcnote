// -------------------------------------------------
// Package server
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package server

import (
	"context"
	"fmt"
	"gcnote/server/config"
	"gcnote/server/model"
	"github.com/allegro/bigcache/v3"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/schema"
	"io"
	"os"
	"strings"
	"time"
	//"github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitConfig() {
	configFileName := "./etc/config.yaml"
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		zap.S().Panicf("读取配置文件失败%v", err)
	}
	c := config.ServerCfg
	if err := v.Unmarshal(&c); err != nil {
		zap.S().Panicf("解析配置文件失败 %v", err)
	}

	config.ServerCfg = c
	//fmt.Printf("配置文件：%+v", c)
	fmt.Println("配置文件加载成功")
	fmt.Println("---------------")
}

func InitMysql() {
	c := config.ServerCfg.MysqlConf
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.DB)

	// 连接
	var err error
	db, err := gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		},
	) // end of gorm.open
	// 创建表
	config.DB = db
	config.DB.AutoMigrate(&model.User{})
	config.DB.AutoMigrate(&model.Index{})
	config.DB.AutoMigrate(&model.KBFile{})
	if err != nil {
		zap.S().Panicf("初始化数据库失败 err:%v", err)
		fmt.Println("数据库初始化失败")
	}
}

func InitRedis() {
	ctx := context.Background()
	c := config.ServerCfg
	redisClient := redis.NewUniversalClient(
		&redis.UniversalOptions{
			Addrs:    strings.Split(c.RedisConf.Host, ","),
			Password: c.RedisConf.Password,
		},
	)

	err := redisClient.Ping(ctx).Err()
	if err != nil {
		zap.S().Panicf("初始化Redis失败 err:%+v", err)
	}
	config.RedisClient = redisClient
}

func InitLocalCache() {
	interval := 8760 * 100 * time.Hour
	c := bigcache.DefaultConfig(interval)
	var err error
	ctx := context.Background()
	localCache, err := bigcache.New(ctx, c)
	if err != nil {
		zap.S().Panicf("初始化Redis失败 err:%+v", err)
	}
	config.LocalCache = localCache
}

func InitLogger() {
	encoder := getEncoder()
	loggerInfo := getLogWriterInfo()
	logLevel := zapcore.DebugLevel
	switch config.ServerCfg.LogConf.Level {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	}

	coreInfo := zapcore.NewCore(encoder, loggerInfo, logLevel)
	logger := zap.New(coreInfo)
	zap.ReplaceGlobals(logger)
}

func getEncoder() zapcore.Encoder {
	productionEncoderConfig := zap.NewProductionEncoderConfig()
	productionEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	productionEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(productionEncoderConfig)
}

func getLogWriterInfo() zapcore.WriteSyncer {
	logPath := config.ServerCfg.LogConf.Path + "/" + config.ServerCfg.Name + ".log"
	l := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    config.ServerCfg.LogConf.MaxSize,    //最大MB
		MaxBackups: config.ServerCfg.LogConf.MaxBackups, //最大备份
		MaxAge:     config.ServerCfg.LogConf.MaxAge,     //保留7天
		Compress:   true,
	}

	var ws io.Writer
	if config.ServerCfg.Mode == "release" {
		ws = io.MultiWriter(l)
	} else {
		//如果不是开发环境，那么会打印日志到日志文件和标准输出，也就是控制台
		ws = io.MultiWriter(l, os.Stdout)
	}

	return zapcore.AddSync(ws)
}
