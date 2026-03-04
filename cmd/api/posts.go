package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TaushifReza/go-social/internal/dto"
	"github.com/TaushifReza/go-social/internal/model"
	"github.com/go-chi/chi/v5"
)

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var dto dto.CreatePostDto
	if err := readJSON(w, r, &dto); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if err := Validate.Struct(dto); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Validation failed", formatValidationErrors(err))
		return
	}

	userId := 1

	post := &model.Posts{
		Title:   dto.Title,
		Content: dto.Content,
		Tags:    dto.Tags,
		// TODO: change after auth
		UserID: int64(userId),
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later", err)
		return
	}

	if err := writeJSONSuccess(w, http.StatusCreated, "Post created successfully.", post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later", err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id", err)
		return
	}
	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			writeJSONError(w, http.StatusNotFound, "post not found", fmt.Errorf("invalid post id"))
		default:
			writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later", err)
		}
		return
	}

	if err := writeJSONSuccess(w, http.StatusOK, "Post retrived successfully.", post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later", err)
		return
	}
}
