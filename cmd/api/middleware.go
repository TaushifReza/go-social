package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/TaushifReza/go-social/internal/auth"
	"github.com/TaushifReza/go-social/internal/dto"
)

// Define a private type for context keys to prevent third-party collisions
type contextKey string

const (
	UserCtxKey  contextKey = "user_id"
	EmailCtxKey contextKey = "email"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeUnauthorized(w, "authorization header is missing", "authorization header is missing")
				return
			}
			// parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				writeUnauthorized(w, "authorization header is malformed", "authorization header is malformed")
				return
			}
			// decode it
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				writeUnauthorized(w, "authorization header is malformed", "authorization header is malformed")
				return
			}
			creds := strings.SplitN(string(decoded), ":", 2)
			// check credentials
			if len(creds) != 2 || creds[0] != "admin" || creds[1] != "admin" {
				writeUnauthorized(w, "invalid credentials", "invalid credentials")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) AuthTokenMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeUnauthorized(w, "authorization header is missing", "authorization header is missing")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				writeUnauthorized(w, "authorization header is missing", "authorization header is malformed")
				return
			}

			token := parts[1]
			claims, err := app.authenticator.VerifyToken(token, auth.AccessTokenType)
			if err != nil {
				writeUnauthorized(w, "unauthorized", "unauthorized")
				return
			}

			ctx := r.Context()

			// Use type-safe context keys
			ctx = context.WithValue(ctx, UserCtxKey, claims.ID)
			ctx = context.WithValue(ctx, EmailCtxKey, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext extracts the user ID and returns an "ok" boolean if it exists
func (app *application) GetUserIDFromContext(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(UserCtxKey).(int64)
	return id, ok
}

// GetEmailFromContext extracts the email and returns an "ok" boolean if it exists
func (app *application) GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailCtxKey).(string)
	return email, ok
}

func (app *application) checkPostOwnership(role string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromCtx(r)
		post := getPostFromCtx(r)

		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		allowed, err := app.checkRolePrecedence(r.Context(), user, role)
		if err != nil {
			writeJSONError(w, http.StatusForbidden, "something went wrong", err)
			return
		}

		if !allowed {
			writeForbidden(w, "Access denied", "Access denied")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *dto.UserResponseDto, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}

func (app *application) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.rateLimiter.Enabled {
			if allow, retryAfter := app.rateLimiter.Allow(r.RemoteAddr); !allow {
				writeJSONError(w, 429, "Limit exceed", retryAfter.String())
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
