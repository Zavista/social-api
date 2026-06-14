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

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("unauthorized error",
		"error", err.Error(),
		"method", r.Method,
		"path", r.URL.Path,
	)
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warn("forbidden", "method", r.Method, "path", r.URL.Path)

	writeJSONError(w, http.StatusForbidden, "forbidden")
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warn("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
