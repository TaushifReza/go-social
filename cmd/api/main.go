package main

import (
	"expvar"
	"runtime"
	"time"

	"github.com/TaushifReza/go-social/internal/auth"
	"github.com/TaushifReza/go-social/internal/cache"
	"github.com/TaushifReza/go-social/internal/db"
	"github.com/TaushifReza/go-social/internal/env"
	"github.com/TaushifReza/go-social/internal/mailer"
	"github.com/TaushifReza/go-social/internal/ratelimiter"
	"github.com/TaushifReza/go-social/internal/store"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

//	@title			Go Social API
//	@description	This is api docs for Go Social.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic("No .env file found (using environment variables)")
	}

	config := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://postgres:root786@172.21.64.1/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisConfig: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "172.21.64.1:6379"),
			pw:      env.GetString("REDIS_PASSWORD", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLE", false),
		},
		env:     env.GetString("ENV", "development"),
		version: env.GetString("VERSION", "0.0.1"),
		mail: mailConfig{
			exp: time.Hour * 3, // 3 days
			sendGrid: sendGridConfig{
				apiKey:    env.GetString("SEND_GRID_API_KEY", "apiKey"),
				fromEmail: env.GetString("SEND_GRID_FROM_EMAIL", "noreply@taushif.com"),
			},
			mailTrap: mailTrapConfig{
				apiKey:    env.GetString("MAIL_TRAP_API_KEY", "apiKey"),
				fromEmail: env.GetString("MAIL_TRAP_FROM_EMAIL", "noreply@taushifreza.com.np"),
			},
		},
		auth: authConfig{
			token: tokenConfig{
				jwtSecret: env.GetString("JWT_SECRET", "jwt-secret"),
				aud:       env.GetString("JWT_AUD", "go-social"),
				iss:       env.GetString("JWT_ISS", "go-social"),
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		config.db.addr,
		config.db.maxOpenConns,
		config.db.maxIdleConns,
		config.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal("database connection error: ", err)
	}

	defer db.Close()
	logger.Info("database connection established successfully.")

	// cache
	var rdb *redis.Client
	if config.redisConfig.enabled {
		rdb = cache.NewRedisClient(
			config.redisConfig.addr,
			config.redisConfig.pw,
			config.redisConfig.db,
		)
	}

	store := store.NewStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)

	mailtrap, err := mailer.NewMailTrapClient(config.mail.mailTrap.apiKey, config.mail.mailTrap.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		config.rateLimiter.RequestsPerTimeFrame,
		config.rateLimiter.TimeFrame,
	)

	authenticator := auth.NewJWTAuthenticator(
		config.auth.token.jwtSecret,
		config.auth.token.aud,
		config.auth.token.iss,
	)

	app := &application{
		config:        config,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailtrap,
		authenticator: authenticator,
		rateLimiter:   rateLimiter,
	}

	// Metrics collected
	expvar.NewString("version").Set(config.version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
