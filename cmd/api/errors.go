package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("internal server error",
		"error", err.Error(),
		"method", r.Method,
		"path", r.URL.Path,
	)
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("bad request error",
		"error", err.Error(),
		"method", r.Method,
		"path", r.URL.Path,
	)
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundReponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("not found error",
		"error", err.Error(),
		"method", r.Method,
		"path", r.URL.Path,
	)
	writeJSONError(w, http.StatusNotFound, "not found")
}
