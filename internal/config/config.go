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
	MySQL mysql
	Redis redis
}

type mysql struct {
	Host          string 
	Port          string 
	User          string 
	Password      string 
	Database      string 
	MaxIdleConns  int
	MaxOpenConns  int
	MaxIdleTime   time.Duration
}

type redis struct {
	Host       string 
	Port       string 
	Password   string 
	Database   int    
	Enabled    bool
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
		MySQL: mysql{
			Host:          env("MYSQL_HOST", "localhost"),
			Port:          env("MYSQL_PORT", "13306"),
			User:          env("MYSQL_USER", "root"),
			Password:      env("MYSQL_PASSWORD", "my-secret-pw"),
			Database:      env("MYSQL_DATABASE", "social"),
			MaxIdleConns:  30,
			MaxOpenConns:  30,
			MaxIdleTime:   time.Minute * 15,
		},
		Redis: redis{
			Host:      env("REDIS_HOST", "localhost"),
			Port:      env("REDIS_PORT", "16379"),
			Password:  env("REDIS_PASSWORD", ""),
			Database:  0,
			Enabled:   false, 
		}, 
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
			cfg.DB.MySQL.User, 
			cfg.DB.MySQL.Password, 
			cfg.DB.MySQL.Host, 
			cfg.DB.MySQL.Port, 
			cfg.DB.MySQL.Database,
		)
}

func (cfg *Config) FormattedRedisAddr() string {

	// 例, "localhost:16379"
	return fmt.Sprintf("%s:%s", cfg.DB.Redis.Host, cfg.DB.Redis.Port)
}


func env(key, fallbackValue string) string {
	s, exists := os.LookupEnv(key)	
	
	if !exists {
		return fallbackValue
	} 
	return s
}