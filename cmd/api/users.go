package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id", "invalid id")
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

	if err := writeJSONSuccess(w, http.StatusOK, "User retrived successfully", user); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "something went wrong. please try again later", err)
		return
	}
}
