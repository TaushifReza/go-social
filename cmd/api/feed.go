package main

import (
	"fmt"
	"net/http"

	"github.com/TaushifReza/go-social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	fmt.Println("Sort: ", fq.Sort)
	if err != nil {
		fmt.Println("Parse ERROR: ", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid pagination or sorting", err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		fmt.Println("Parse ERROR: ", err)
		writeJSONError(w, http.StatusBadRequest, "Invalid pagination or sorting", err)
		return
	}

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(5), fq)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again", err)
		return
	}

	if err := writeJSONSuccess(w, http.StatusOK, "Post feed", feed); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again", err)
		return
	}
}
