package main

import (
	"net/http"

	"github.com/TaushifReza/go-social/internal/dto"
	"github.com/TaushifReza/go-social/internal/model"
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
