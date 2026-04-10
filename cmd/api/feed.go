package main

import (
	"net/http"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(5))

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again", err)
		return
	}

	if err := writeJSONSuccess(w, http.StatusOK, "Post updated", feed); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again", err)
		return
	}
}
