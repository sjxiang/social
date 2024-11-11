package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/sjxiang/social/internal/env"
)

func LoadConf() (Config, error) {
	
	// 从当前目录下的 .env 文件中加载环境变量
	if err := godotenv.Load(); err != nil {
		return Config{}, nil 
	}

	// 从环境变量中, 加载配置
	cfg := Config{
		Web: WebOptions{
			Addr:   env.GetString("WEB_ADDR", ":8080"), 
			Env:    env.GetString("WEB_ENV", "Realese"),
			ApiURL: env.GetString("WEB_API_URL", "localhost:8080"),
		},
		Auth: AuthOptions{
			Basic: BasicOptions{
				Username:     env.GetString("AUTH_BASIC_USERNAME", "admin"),
				Password:     env.GetString("AUTH_BASIC_PASSWORD", "123456"),
			},
			Token: TokenOptions{
				SecretKey: env.GetString("AUTH_TOKEN_SECRET_KEY", "8xEMrWkBARcDDYQ"),
				Exp:    time.Hour * 24 * 3, // 3 days
			},
		},
		MySQL: MySQLOptions{
			Username:     env.GetString("MYSQL_USERNAME", "root"),
			Password:     env.GetString("MYSQL_PASSWORD", "my-secret-pw"),
			Host:         env.GetString("MYSQL_HOST", "127.0.0.1"),
			Port:         env.GetInteger("MYSQL_PORT", 13306),
			Database:     env.GetString("MYSQL_DATABASE_NAME", "social"),
			MaxOpenConns: env.GetInteger("MYSQL_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInteger("MYSQL_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("MYSQL_MAX_IDLE_TIME", "15m"),
		},
		Redis: RedisOptions{
			Host:     env.GetString("REDIS_HOST", "IP_ADDRESS"),
			Port:     env.GetInteger("REDIS_PORT", 6379),
			Password: env.GetString("REDIS_PASSWORD", ""),
			DB:       env.GetInteger("REDIS_DB", 0),
			Enabled:  env.GetBool("REDIS_ENABLED", false),
		},
		Mail: MailOptions{
			Password:  env.GetString("MAIL_PASSWORD", ""),
			FromEmail: env.GetString("MAIL_FROM_EMAIL", "gua@vip.cn"),
			Exp:       time.Hour * 24 * 3,
		},
		RateLimiter: RateLimiterOptions{
			RequestsPerTimeFrame: env.GetInteger("RATE_LIMITER_REQUESTS_PER_TIME_FRAME", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
		},

	}
	
	return cfg, nil 
}	
