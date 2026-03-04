package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) getCommentByPostID(w http.ResponseWriter, r *http.Request) {
	postID, err := strconv.ParseInt(chi.URLParam(r, "postID"), 10, 64)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id", err)
		return
	}

	ctx := r.Context()

	comments, err := app.store.Comments.GetCommentByPostID(ctx, postID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later", err)
		return
	}

	if err := writeJSONSuccess(w, http.StatusOK, "comment retrived successfully", comments); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later", err)
		return
	}
}
