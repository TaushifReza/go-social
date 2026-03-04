package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/TaushifReza/go-social/internal/dto"
	"github.com/TaushifReza/go-social/internal/model"
	"github.com/go-chi/chi/v5"
)

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var dto dto.CreatePostDto
	if err := readJSON(w, r, &dto); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
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
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	ctx := r.Context()

	post, err := app.store.Posts.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			writeJSONError(w, http.StatusNotFound, "post not found")
		default:
			writeJSONError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
}
