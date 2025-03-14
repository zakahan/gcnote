// -------------------------------------------------
// Package config
// Author: hanzhi
// Date: 2024/12/8
// -------------------------------------------------

package config

import (
	"github.com/allegro/bigcache/v3"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServerConfig struct {
	Name        string              `mapstructure:"name" json:"name"`                   // 服务名称
	Host        string              `mapstructure:"host" json:"host"`                   // 主机地址
	Port        int                 `mapstructure:"port" json:"port"`                   // 启动端口
	Mode        string              `mapstructure:"mode" json:"mode"`                   // 启动端口
	RedisConf   redisConfig         `mapstructure:"redis" json:"redis"`                 // Redis配置
	MysqlConf   mysqlConfig         `mapstructure:"mysql" json:"mysql"`                 // Mysql配置
	LogConf     logsConfig          `mapstructure:"logs" json:"logs"`                   // 日志配置
	ElasticConf elasticSearchConfig `mapstructure:"elasticsearch" json:"elasticsearch"` // es的配置
}

type redisConfig struct {
	Host     string `mapstructure:"host" json:"host"`         // Redis地址。集群用多个逗号分割
	Port     string `mapstructure:"port" json:"port"`         // Redis端口
	Password string `mapstructure:"password" json:"password"` // Redis密码
}

type mysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`         // Mysql地址
	Port     int    `mapstructure:"port" json:"port"`         // Mysql端口
	DB       string `mapstructure:"db" json:"db"`             // 数据库
	User     string `mapstructure:"user" json:"user"`         // Mysql用户
	Password string `mapstructure:"password" json:"password"` // Mysql密码
}

type logsConfig struct {
	Path       string `mapstructure:"path" json:"path"`               // 配置文件路径
	Level      string `mapstructure:"level" json:"level"`             // 日志级别 debug、info、warn、error
	MaxAge     int    `mapstructure:"max_age" json:"max_age"`         // 最大保存时间（单位天
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups"` //最大备份数
	MaxSize    int    `mapstructure:"max_size" json:"max_size"`       // 最大Size MB
	Compress   int    `mapstructure:"compress" json:"compress"`       // 是否压缩
}

type elasticSearchConfig struct {
	Address  string `mapstructure:"address" json:"address"`     // elasticsearch地址
	Username string `mapstructure:"user_name" json:"user_name"` // es账户
	Password string `mapstructure:"password" json:"password"`   // es密码
	CertPath string `mapstructure:"cert_path" json:"cert_path"` // 许可证路径
	UseCert  bool   `mapstructure:"use_cert" json:"use_cert"`   // 是否使用许可证
}

var ServerCfg ServerConfig
var DB *gorm.DB
var RedisClient redis.UniversalClient
var LocalCache *bigcache.BigCache
var ElasticClient *elasticsearch.Client
