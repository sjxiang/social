package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/sjxiang/social/internal/auth"
	"github.com/sjxiang/social/internal/config"
	"github.com/sjxiang/social/internal/data"
	"github.com/sjxiang/social/internal/mail"
	"github.com/sjxiang/social/internal/ratelimiter"
	"github.com/sjxiang/social/internal/token"
)


type application struct {
	logger        *zap.SugaredLogger
	config        config.Config
	store         data.MySQLStorage
	rateLimiter   ratelimiter.Limiter
	mailer        mail.EmailSender
	auth          auth.Authenticator
	tokenMaker    token.Maker
}


func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// 中间件, 拦截器
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{app.config.Web.Addr},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	
	// 限流
	if app.config.RateLimiter.Enabled {
		r.Use(app.RateLimiterMiddleware)
	}

	// 设置超时请求为 60s
	r.Use(middleware.Timeout(60 * time.Second))

	// 路由
	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
	})
	
	return r
}

func (app *application) run(mux http.Handler) error {
	
	// 创建一个 http server 实例
	srv := &http.Server{
		Addr:         app.config.Web.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		// 监听信号
		quit := make(chan os.Signal, 1)

		// 监听指定的信号 (kill, ctrl + c)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		
		// 阻塞在此, 直至从 quit 中读取到上述两种信号
		s := <-quit

		// 5s 超时, 冗余处理未完成的请求
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		// 关闭 server
		shutdown <- srv.Shutdown(ctx)
	}()

	app.logger.Infow("server has started", "addr", app.config.Web.Addr, "env", app.config.Web.Env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has stopped", "addr", app.config.Web.Addr, "env", app.config.Web.Env)

	return nil
}
