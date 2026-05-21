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

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("unauthorized basic error",
		"error", err.Error(),
		"method", r.Method,
		"path", r.URL.Path,
	)

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted" charset="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}
