package main

import (
	"github.com/sjxiang/social/internal/auth"
	"github.com/sjxiang/social/internal/logger"
	"github.com/sjxiang/social/internal/mailer"
	"github.com/sjxiang/social/internal/ratelimiter"
	"github.com/sjxiang/social/internal/streamer"
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
		cfg.db.dsn,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err!= nil {
		logger.Fatal(err)
	}

	defer db.Close()

	// Redis
	if cfg.useCaching {
		logger.Fatal("未配置缓存")
	}

	logger.Info("initializing database support")

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

	const numWorkers = 4
	videoQueue := make(chan streamer.VideoProcessingJob, numWorkers)
	defer close(videoQueue)

	wp := streamer.New(videoQueue, numWorkers)
	wp.Run()

	
	// app
	app := &application{
		config:        cfg,
		mailer:        sender,
		logger:        logger,
		tokenMaker:    tokenMaker,
		auth:          jwtAuthenticator,
		rateLimiter:   fixedWindowLimiter,
		videoQueue:    videoQueue,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}