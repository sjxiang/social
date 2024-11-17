package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)


// 配置
type Config struct {
	Web       web
	Auth      auth
	DB        database
	RateLimit rateLimit
	Mail      mail
}

type web struct {
	Addr   string
	Env    string
	ApiURL string
}

type auth struct {
	Basic basic
	JWT   jwt
}

type basic struct {
	Username  string
	Password  string
}

type jwt struct {
	SecretKey string
	Issuer    string
	Expiry    time.Duration
}

type database struct {
	MySQLHost          string 
	MySQLPort          string 
	MySQLUser          string 
	MySQLPassword      string 
	MySQLDatabase      string 
	MySQLMaxIdleConns  int
	MySQLMaxOpenConns  int
	MySQLMaxIdleTime   time.Duration
	// cache config
	RedisHost          string 
	RedisPort          string 
	RedisPassword      string 
	RedisDatabase      int    
	RedisEnabled       bool
}


type mail struct {
	FromEmail  string
	ApiKey     string  // 授权码
	Expiry     time.Duration
}

type rateLimit struct {
	RequestsPerTimeFrame int  // 每个时间段的请求量
	TimeFrame            time.Duration
	Enabled              bool
}


func New() (Config, error) {
	
	cfg := Config{}

	// 加载 workspace 同级目录下的 .env 文件
	if err := godotenv.Load(); err != nil {
		return cfg, err
	}

	cfg.Web = web{
		Addr:   env("WEB_ADDR", ":8080"),
		Env:    env("WEB_ENV", "Realese"),
		ApiURL: env("API_URL", ":8080"),
	}

	cfg.DB = database{
		MySQLHost:          env("MYSQL_HOST", "localhost"),
		MySQLPort:          env("MYSQL_PORT", "13306"),
		MySQLUser:          env("MYSQL_USER", "root"),
		MySQLPassword:      env("MYSQL_PASSWORD", "my-secret-pw"),
		MySQLDatabase:      env("MYSQL_DATABASE", "social"),
		MySQLMaxIdleConns:  30,
		MySQLMaxOpenConns:  30,
		MySQLMaxIdleTime:   time.Minute * 15,

		RedisHost:      env("REDIS_HOST", "localhost"),
		RedisPort:      env("REDIS_PORT", "16379"),
		RedisPassword:  env("REDIS_PASSWORD", ""),
		RedisDatabase:  0,
		RedisEnabled:   false,  
	}
	
	cfg.Auth = auth{
		Basic: basic{
			Username:  env("BASIC_AUTH_USERNAME", "admin"),
			Password:  env("BASIC_AUTH_PASSWORD", "123456"),
		},
		JWT: jwt{
			SecretKey: env("JWT_SECRET_KEY", "8xEMrWkBARcDDYQ"),
			Issuer:    env("JWT_ISSUER", "gua@vip.cn"),
			Expiry:    time.Hour * 24 * 7,
		},
	}

	cfg.Mail = mail{
		FromEmail:  env("MAILER_FROM_EMAIL", "gua@vip.cn"),
		ApiKey:     env("MAILER_API_KEY", "xxxooo"),
		Expiry:     time.Minute * 30,
	}

	cfg.RateLimit = rateLimit{
		RequestsPerTimeFrame: 20,
		TimeFrame:            time.Minute * 1,
		Enabled:              true,
	}

	return cfg, nil 
}

func (cfg *Config) FormattedMySQLAddr() string {
	
	// 例, "root:my-secret-pw@tcp(127.0.0.1:13306)/social?charset=utf8&parseTime=True&loc=Local"
	
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", 
			cfg.DB.MySQLUser, 
			cfg.DB.MySQLPassword, 
			cfg.DB.MySQLHost, 
			cfg.DB.MySQLPort, 
			cfg.DB.MySQLDatabase,
		)
}

func (cfg *Config) FormattedRedisAddr() string {
	// 例, "localhost:16379"
	return fmt.Sprintf("%s:%s", cfg.DB.RedisHost, cfg.DB.RedisPort)
}


func env(key, fallbackValue string) string {
	s, exists := os.LookupEnv(key)	
	
	if !exists {
		return fallbackValue
	} 
	return s
}