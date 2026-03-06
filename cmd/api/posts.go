package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/TaushifReza/go-social/internal/dto"
	"github.com/TaushifReza/go-social/internal/model"
	"github.com/go-chi/chi/v5"
)

type postKey string

const postCtx postKey = "post"

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
	post := getPostFromCtx(r)

	if err := writeJSONSuccess(w, http.StatusOK, "Post retrived successfully.", post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later", err)
		return
	}
}

func (app *application) postDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid post id", err)
		return
	}

	ctx := r.Context()
	err = app.store.Posts.DeletePostByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			writeJSONError(w, http.StatusNotFound, "post not found", fmt.Errorf("invalid post id"))
		default:
			writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later", err)
		}
		return
	}

	if err := writeJSONSuccess(w, http.StatusOK, "Post deleted", "Post deleted"); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later", err)
		return
	}
}

func (app *application) postUpdateHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var dto dto.UpdatePostDto
	if err := readJSON(w, r, &dto); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	if err := Validate.Struct(dto); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid request payload", formatValidationErrors(err))
		return
	}

	if dto.Content != nil {
		post.Content = *dto.Content
	}
	if dto.Title != nil {
		post.Title = *dto.Title
	}

	ctx := r.Context()

	if err := app.store.Posts.Update(ctx, post); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			writeJSONError(w, http.StatusNotFound, "post not found", "post not found")
		default:
			writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again", err)
		}
		return
	}

	if err := writeJSONSuccess(w, http.StatusOK, "Post updated", post); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again", err)
		return
	}
}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx = context.WithValue(ctx, postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *model.Posts {
	post, _ := r.Context().Value(postCtx).(*model.Posts)
	return post
}
