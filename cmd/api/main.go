package main

import (
	"github.com/sjxiang/social/internal/auth"
	"github.com/sjxiang/social/internal/logger"
	"github.com/sjxiang/social/internal/mailer"
	"github.com/sjxiang/social/internal/ratelimiter"
	"github.com/sjxiang/social/internal/token"
	"github.com/sjxiang/social/internal/utils"
)


func main() {
		
	logger := logger.Must("bbs")
	defer logger.Sync()

	// dotenv
	cfg, err := loadConfig()
	if err != nil {
		logger.Fatal(err)
	}

	// MySQL
	db, err := utils.NewMySQL(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err!= nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("initializing database support")


	// Redis
	if cfg.redis.enabled {
		logger.Fatal("未配置缓存")
	}

	// Rate limiter
	fixedWindowLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame, 
		cfg.rateLimiter.TimeFrame,
	)
	
	// Mailer
	sender := mailer.NewQQmailSender(
		"no-reply", 
		cfg.mail.fromEmail, 
		cfg.mail.sendGrid.apiKey,
	)
	
	// Authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.jwt.secretKey, 
		cfg.auth.jwt.issuer,
	)

	// Token Maker
	tokenMaker := token.NewJWTMaker(
		cfg.auth.jwt.secretKey, 
		cfg.auth.jwt.issuer,
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