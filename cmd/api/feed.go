package main

import (
	"net/http"

	"example.com/Go_Land/internal/env/store"
)

// GetUserFeedHandler godoc
// @Summary      Get user feed
// @Description  Get the feed for a user
// @Tags         users
// @Produce      json
// @Success      200  {array}   store.Post
// @Router       /api/v1/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	feed, err := app.store.Posts.GetUserFeed(ctx, int64(42), fq)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if len(feed) == 0 {
		app.writeJSONErrorResponse(w, http.StatusNotFound, nil)
		return
	}

	err = writeJSON(w, http.StatusOK, feed)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
