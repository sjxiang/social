package config

import (
	"time"
	"fmt"
)

// 配置
type Config struct {
	Web         WebOptions
	Auth        AuthOptions
	Redis       RedisOptions
	MySQL       MySQLOptions
	Mail        MailOptions
	RateLimiter RateLimiterOptions
}

type WebOptions struct {
	Addr        string
	Env         string 
	ApiURL      string
}

type AuthOptions struct {
	Basic BasicOptions
	Token TokenOptions
}

type BasicOptions struct {
	Username  string
	Password  string
}

type TokenOptions struct {
	SecretKey string
	Issuer    string
	Exp       time.Duration
}

type MySQLOptions struct {
	Host           string
	Port           int
	Username       string
	Password       string
	Database       string
	MaxIdleConns   int
	MaxOpenConns   int   
	MaxIdleTime    string
}

type RedisOptions struct {
	Host      string 
	Port      int 
	Password  string 
	DB        int 
	Enabled   bool
}

type MailOptions struct {
	Password    string  // 授权码
	FromEmail   string
	Exp         time.Duration
}

type RateLimiterOptions struct {
	RequestsPerTimeFrame int  // 每个时间段的请求量
	TimeFrame            time.Duration
	Enabled              bool
}

func (cfg *Config) FormattedMySQLAddr() string {
	// "root:my-secret-pw@tcp(127.0.0.1:13306)/social?charset=utf8&parseTime=True&loc=Local"
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", 
			cfg.MySQL.Username, 
			cfg.MySQL.Password, 
			cfg.MySQL.Host, 
			cfg.MySQL.Port, 
			cfg.MySQL.Database,)
}

func (cfg *Config) FormattedRedisAddr() string {
	// "localhost:16379"
	return fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
}
