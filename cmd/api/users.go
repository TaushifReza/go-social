package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/TaushifReza/go-social/internal/dto"
	"github.com/TaushifReza/go-social/internal/model"
	"github.com/TaushifReza/go-social/internal/store"
	"github.com/TaushifReza/go-social/internal/utils"
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

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	model.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{id} [get]
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

// GetUser godoc
//
//	@Summary		Register User
//	@Description	Register User
//	@Tags			users
//	@Accept			json
//	@Produce		json
//
// @Param           payload body dto.UserRegisterationDto
//
//	@Success		200	{object}	model.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/auth/users/ [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var dto dto.UserRegisterationDto
	if err := readJSON(w, r, &dto); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	if err := Validate.Struct(dto); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Validation failed", formatValidationErrors(err))
		return
	}

	hashPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Error while hashing password", err.Error())
		return
	}

	user := &model.User{
		UserName: dto.UserName,
		Email:    dto.Email,
		Password: hashPassword,
	}

	ctx := r.Context()

	_, hashedToken := utils.CreateToken()

	// store the user
	if err := app.store.Users.CreateAndInvite(ctx, user, hashedToken, app.config.mail.exp); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			writeJSONError(w, http.StatusInternalServerError, "Email already exists.", err)
		case store.ErrDuplicateUsername:
			writeJSONError(w, http.StatusBadRequest, "Username already exists.", err)
		default:
			writeJSONError(w, http.StatusBadRequest, "something went wrong. please try again later.", err)
		}
		return
	}

	if err := writeJSONSuccess(w, http.StatusCreated, "user registered", user); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later.", err)
		return
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if err := app.store.Users.Activate(r.Context(), token); err != nil {
		switch {
		case err.Error() == "invitation not found or invalid":
			writeJSONError(w, http.StatusNotFound, err.Error(), nil)
		case err.Error() == "invitation has expired":
			writeJSONError(w, http.StatusGone, err.Error(), nil)
		default:
			// Log the actual error for debugging, but hide it from the user
			app.logger.Errorw("activation failed", "error", err)
			writeJSONError(w, http.StatusInternalServerError, "internal server error", nil)
		}
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}
