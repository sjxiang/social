package main

import (
	"time"
	"os"

	"github.com/joho/godotenv"

	"github.com/sjxiang/social/internal/ratelimiter"
)


type config struct {
	addr        string
	env         string
	apiURL      string
	version     string
	useCaching  bool 
	mail        mailConfig
	auth        authConfig
	redis       redisConfig
	db          dbConfig	
	rateLimiter ratelimiter.Config
}

type authConfig struct {
	basic basicConfig
	jwt   jwtConfig
}

type basicConfig struct {
	user string
	pass string
}

type jwtConfig struct {
	secretKey string
	issuer    string
	expiry    time.Duration
}

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

type redisConfig struct {
	addr     string
	pw       string
	db       int
}

type mailConfig struct {
	sendGrid  sendGridConfig  // send 网格 
	fromEmail string
	exp       time.Duration
}

type sendGridConfig struct {
	apiKey string  // 授权码
}


func loadConfig() (config, error) {
	
	// tips, 加载 workspace 同级目录下的 .env 文件
	if err := godotenv.Load(); err != nil {
		return config{}, nil
	}

	cfg := config{
		addr:       "localhost:8080",
		env:        "dev",
		apiURL:     "localhost:8080",
		version:    "1.1.0",
		useCaching: false,
	}

	cfg.auth = authConfig{
		basic: basicConfig{
			user: env("BASIC_USERNAME", ""),
			pass: env("BASIC_PASSWORD", ""),
		},
		jwt: jwtConfig{
			secretKey: env("JWT_SECRET_KEY", ""),
			issuer:    env("JWT_ISSUER", ""),
			expiry:    time.Hour * 7 * 24,
		},
	}

	cfg.db = dbConfig{
		dsn:         env("MYSQL_DSN", ""),
		maxOpenConns: 30,
		maxIdleConns: 30,
		maxIdleTime:  time.Minute * 15,
	}

	cfg.redis = redisConfig{
		addr:     "localhost:6379",
		pw:       "",
		db:       0,
	}

	cfg.mail = mailConfig{
		sendGrid: sendGridConfig{
			apiKey: env("SEND_GRID_API_KEY", "xxxxxx"),
		},
		fromEmail: env("FROM_EMAIL", "gua@vip.cn"),
		exp:       time.Minute * 15,
	}

	return cfg, nil 
} 



func env(key, fallbackValue string) string {
	s, exists := os.LookupEnv(key)	
	
	if !exists {
		return fallbackValue
	} 
	return s
}

