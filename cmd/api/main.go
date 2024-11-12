package main

import (
	"github.com/sjxiang/social/internal/auth"
	"github.com/sjxiang/social/internal/config"
	"github.com/sjxiang/social/internal/logger"
	"github.com/sjxiang/social/internal/mail"
	"github.com/sjxiang/social/internal/ratelimiter"
	"github.com/sjxiang/social/internal/token"
	"github.com/sjxiang/social/internal/utils"
)


const version = "1.1.0"

func main() {
	
	
	logger := logger.Must("社区")
	defer logger.Sync()

	// dotenv
	cfg, err := config.LoadConf()
	if err != nil {
		logger.Fatal(err)
	}

	// MySQL
	db, err := utils.NewMySQL(
		cfg.FormattedMySQLAddr(),
		cfg.MySQL.MaxOpenConns,
		cfg.MySQL.MaxIdleConns,
		cfg.MySQL.MaxIdleTime,
	)
	if err!= nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("initializing database support")

	// Redis

	// Rate limiter
	fixedWindowLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.RateLimiter.RequestsPerTimeFrame, 
		cfg.RateLimiter.TimeFrame,
	)
	
	// Mailer
	sender := mail.NewQQmailSender("no-reply", cfg.Mail.FromEmail, cfg.Mail.Password)
	
	// Authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.Auth.Token.SecretKey, 
		cfg.Auth.Token.Issuer,
	)

	// Token Maker
	tokenMaker := token.NewJWTMaker(
		cfg.Auth.Token.SecretKey, 
		cfg.Auth.Token.Issuer,
	)

	app := &application{
		config:        cfg,
		mailer:        sender,
		logger:        logger,
		tokenMaker:    tokenMaker,
		authenticator: jwtAuthenticator,
		rateLimiter:   fixedWindowLimiter,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}