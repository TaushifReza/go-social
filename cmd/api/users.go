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

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followUser := getUserFromCtx(r)

	// TODO get user from auth middleware
	var userID int64 = 12

	ctx := r.Context()

	if err := app.store.Users.Follow(ctx, userID, followUser.ID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again", err)
		return
	}

	if err := writeJSONSuccess(w, http.StatusOK, "Follow user success", "Follow user success"); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later", err)
		return
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unFollowUser := getUserFromCtx(r)

	// TODO get user from auth middleware
	var userID int64 = 12

	ctx := r.Context()

	if err := app.store.Users.UnFollow(ctx, userID, unFollowUser.ID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again", err)
		return
	}

	if err := writeJSONSuccess(w, http.StatusOK, "UnFollow user success", "UnFollow user success"); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later", err)
		return
	}
}

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
