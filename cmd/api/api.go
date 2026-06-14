package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TaushifReza/go-social/docs"
	"github.com/TaushifReza/go-social/internal/auth"
	"github.com/TaushifReza/go-social/internal/cache"
	"github.com/TaushifReza/go-social/internal/mailer"
	"github.com/TaushifReza/go-social/internal/ratelimiter"
	"github.com/TaushifReza/go-social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	rateLimiter   ratelimiter.Limiter
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	version     string
	mail        mailConfig
	auth        authConfig
	redisConfig redisConfig
	rateLimiter ratelimiter.Config
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type mailConfig struct {
	sendGrid sendGridConfig
	mailTrap mailTrapConfig
	exp      time.Duration
}

type sendGridConfig struct {
	apiKey    string
	fromEmail string
}

type mailTrapConfig struct {
	apiKey    string
	fromEmail string
}

type authConfig struct {
	token tokenConfig
}

type tokenConfig struct {
	jwtSecret string
	aud       string
	iss       string
}

type redisConfig struct {
	addr    string
	pw      string
	db      int
	enabled bool
}

func (app *application) mount() http.Handler {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	if app.config.rateLimiter.Enabled {
		r.Use(app.RateLimiterMiddleware)
	}

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.With(app.BasicAuthMiddleware()).Get("/health", app.healthCheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware())
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.checkPostOwnership("admin", app.postDeleteHandler))
				r.Patch("/", app.checkPostOwnership("moderator", app.postUpdateHandler))

				r.Route("/comments", func(r chi.Router) {
					r.Get("/", app.getCommentByPostIDHandler)

				})
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}/", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware())
				r.Use(app.userContextMiddleware)

				r.Get("/", app.getUserHandler)
				r.Put("/follow/", app.followUserHandler)
				r.Put("/unfollow/", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed/", app.getUserFeedHandler)
			})

		})

		// Public route
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register/", app.registerUserHandler)
			r.Post("/login/", app.loginHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = app.config.version
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/v1"

	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", s.String())

		shutdown <- server.Shutdown(ctx)
	}()

	app.logger.Info("server has started at http://localhost", app.config.addr)

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("server has been stop")

	return nil
}
