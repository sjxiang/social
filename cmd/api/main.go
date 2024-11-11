package main

import (
	"github.com/sjxiang/social/internal/config"
	"github.com/sjxiang/social/internal/utils"
	"github.com/sjxiang/social/internal/logger"
)


const version = "1.1.0"

func main() {
	
	logger := logger.Must("社区")
	defer logger.Sync()

	// dotenv
	cfg, err := config.Load()
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
	
	
	app := &application{
		config:        cfg,
		logger:        logger,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}