package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: replace with authenticated user
	userID := int64(1)
	feed, err := app.store.Posts.GetUserFeed(ctx, userID)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
