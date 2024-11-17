package main

import (
	"github.com/sjxiang/social/internal/auth"
	"github.com/sjxiang/social/internal/config"
	"github.com/sjxiang/social/internal/logger"
	"github.com/sjxiang/social/internal/mailer"
	"github.com/sjxiang/social/internal/ratelimiter"
	"github.com/sjxiang/social/internal/token"
	"github.com/sjxiang/social/internal/utils"
)


const version = "1.1.0"

func main() {
	
	
	logger := logger.Must("社区")
	defer logger.Sync()

	// dotenv
	cfg, err := config.New()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("配置清单", cfg)

	// MySQL
	db, err := utils.NewMySQL(
		cfg.FormattedMySQLAddr(),
		cfg.DB.MySQLMaxOpenConns,
		cfg.DB.MySQLMaxIdleConns,
		cfg.DB.MySQLMaxIdleTime,
	)
	if err!= nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("initializing database support")

	// Redis

	// Rate limiter
	fixedWindowLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.RateLimit.RequestsPerTimeFrame, 
		cfg.RateLimit.TimeFrame,
	)
	
	// Mailer
	sender := mailer.NewQQmailSender("no-reply", cfg.Mail.FromEmail, cfg.Mail.ApiKey)
	
	// Authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.Auth.JWT.SecretKey, 
		cfg.Auth.JWT.Issuer,
	)

	// Token Maker
	tokenMaker := token.NewJWTMaker(
		cfg.Auth.JWT.SecretKey, 
		cfg.Auth.JWT.Issuer,
	)

	app := &application{
		config:        cfg,
		mailer:        sender,
		logger:        logger,
		tokenMaker:    tokenMaker,
		auth:          jwtAuthenticator,
		rateLimiter:   fixedWindowLimiter,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}