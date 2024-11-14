package config

import (
	"time"
	"fmt"

	"github.com/joho/godotenv"
)

// 配置
type Config struct {
	Web       web
	RateLimit rateLimit
	Mail      mail
	Auth      auth
	Database  database
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
	MySQLPort          int
	MySQLUser          string
	MySQLPassword      string
	MySQLDatabaseName  string
	MySQLMaxIdleConns  int
	MySQLMaxOpenConns  int
	MySQLMaxIdleTime   string
	RedisHost          string
	RedisPort          int
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
	
	cfg.Database = database{
		MySQLHost:          defaultEnvString("MYSQL_HOST", "localhost"),
		MySQLPort:          defaultEnvNumeric("MYSQL_PORT", 13306),
		MySQLUser:          defaultEnvString("MYSQL_USER", "root"),
		MySQLPassword:      defaultEnvString("MYSQL_PASSWORD", "my-secret-pw"),
		MySQLDatabaseName:  defaultEnvString("MYSQL_DATABASE_NAME", "social"),
		MySQLMaxIdleConns:  defaultEnvNumeric("MYSQL_MAX_IDLE_CONNS", 30),
		MySQLMaxOpenConns:  defaultEnvNumeric("MYSQL_MAX_OPEN_CONNS", 30),
		MySQLMaxIdleTime:   defaultEnvString("MYSQL_MAX_IDLE_TIME", "15m"),
		RedisHost:          defaultEnvString("REDIS_HOST", "localhost"),
		RedisPort:          defaultEnvNumeric("REDIS_PORT", 16379),
		RedisPassword:      defaultEnvString("REDIS_PASSWORD", ""),
		RedisDatabase:      defaultEnvNumeric("REDIS_DATABASE", 0),
		RedisEnabled:       defaultEnvBoolean("REDIS_ENABLED", false),
	}
	
	cfg.Auth = auth{
		Basic: basic{
			Username:  defaultEnvString("BASIC_AUTH_USERNAME", "admin"),
			Password:  defaultEnvString("BASIC_AUTH_PASSWORD", "123456"),
		},
		JWT: jwt{
			SecretKey: defaultEnvString("JWT_SECRET_KEY", "8xEMrWkBARcDDYQ"),
			Issuer:    defaultEnvString("JWT_ISSUER", "gua@vip.cn"),
			Expiry:    time.Hour * 24 * 7,
		},
	}

	cfg.Web = web{
		Addr:   defaultEnvString("WEB_ADDR", ":8080"),
		Env:    defaultEnvString("WEB_ENV", "Realese"),
		ApiURL: defaultEnvString("WEB_API_URL", "localhost:8080"),
	}

	cfg.Mail = mail{
		FromEmail:  defaultEnvString("MAIL_FROM_EMAIL", "EMAIL"),
		ApiKey:     defaultEnvString("MAIL_API_KEY", "API_KEY"),
		Expiry:     time.Minute * 30,
	}

	cfg.RateLimit = rateLimit{
		RequestsPerTimeFrame: defaultEnvNumeric("RATE_LIMITER_REQUESTS_PER_TIME_FRAME", 20),
		TimeFrame:            time.Minute,
		Enabled:              defaultEnvBoolean("RATE_LIMITER_ENABLED", true),
	}

	return cfg, nil 
}


func (cfg *Config) FormattedMySQLAddr() string {
	// "root:my-secret-pw@tcp(127.0.0.1:13306)/social?charset=utf8&parseTime=True&loc=Local"
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", 
			cfg.Database.MySQLUser, 
			cfg.Database.MySQLPassword, 
			cfg.Database.MySQLHost, 
			cfg.Database.MySQLPort, 
			cfg.Database.MySQLDatabaseName,)
}

func (cfg *Config) FormattedRedisAddr() string {
	// "localhost:16379"
	return fmt.Sprintf("%s:%d", cfg.Database.RedisHost, cfg.Database.RedisPort)
}

