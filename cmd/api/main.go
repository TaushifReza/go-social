package main

import (
	"time"

	"github.com/TaushifReza/go-social/internal/auth"
	"github.com/TaushifReza/go-social/internal/db"
	"github.com/TaushifReza/go-social/internal/env"
	"github.com/TaushifReza/go-social/internal/mailer"
	"github.com/TaushifReza/go-social/internal/store"
	"github.com/joho/godotenv"
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

	store := store.NewStorage(db)

	mailtrap, err := mailer.NewMailTrapClient(config.mail.mailTrap.apiKey, config.mail.mailTrap.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	authenticator := auth.NewJWTAuthenticator(
		config.auth.token.jwtSecret,
		config.auth.token.aud,
		config.auth.token.iss,
	)

	app := &application{
		config:        config,
		store:         store,
		logger:        logger,
		mailer:        mailtrap,
		authenticator: authenticator,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
