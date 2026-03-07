package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/TaushifReza/go-social/internal/dto"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx postKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := writeJSONSuccess(w, http.StatusOK, "User retrived successfully", user); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later", err)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid id", err)
			return
		}
		ctx := r.Context()

		user, err := app.store.Users.GetUserbyID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, sql.ErrNoRows):
				writeJSONError(w, http.StatusNotFound, "user not found", "user not found")
			default:
				writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later", err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromCtx(r *http.Request) *dto.UserResponseDto {
	user, _ := r.Context().Value(userCtx).(*dto.UserResponseDto)
	return user
}
