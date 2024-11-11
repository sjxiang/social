package main

import (

	"github.com/sjxiang/social/internal/config"
	"github.com/sjxiang/social/internal/logger"
	"github.com/sjxiang/social/internal/token"
	"github.com/sjxiang/social/internal/auth"
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
	
	// Mailer
	
	// Authenticator
	authenticator := auth.NewJWTAuthenticator(cfg.Auth.Token.SecretKey, "JUEJIN")

	// Token Maker
	tokenMaker := token.NewJWTMaker(cfg.Auth.Token.SecretKey, "JUEJIN")

	app := &application{
		config:        cfg,
		logger:        logger,
		tokenMaker:    tokenMaker,
		authenticator: authenticator,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}