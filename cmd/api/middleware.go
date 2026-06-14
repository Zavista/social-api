package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/zavista/social-api/internal/store"
)

func (app *application) rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.rateLimiter.Enabled {
			if allow, retryAfter := app.rateLimiter.Allow(r.RemoteAddr); !allow {
				app.rateLimitExceededResponse(w, r, retryAfter.String())
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) basicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// read auth header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		// parse to get token
		// Authorization: Basic <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Basic" {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorized header is malformed"))
			return
		}

		// decode token
		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			app.unauthorizedBasicErrorResponse(w, r, err)
			return
		}

		username := app.config.auth.basic.user
		password := app.config.auth.basic.pass

		creds := strings.SplitN(string(decoded), ":", 2)
		if len(creds) != 2 || creds[0] != username || creds[1] != password {
			app.unauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials"))
			return
		}

		app.logger.Info("successfully logged in basic auth")
		next.ServeHTTP(w, r)
	})
}

func (app *application) authTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorized header is malformed"))
			return
		}

		token := parts[1]

		claims, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.getUser(ctx, claims.UserID)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtxKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) checkPostOwnership(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromCtx(r)
		post := getPostFromCtx(r)

		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		// role precedence check
		allowed, err := app.checkRolePrecedence(r.Context(), user, requiredRole)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		if !allowed {
			app.forbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.UserWithRole, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}

func (app *application) getUser(ctx context.Context, userID int64) (*store.UserWithRole, error) {
	if !app.config.redisCfg.enabled {
		return app.store.Users.GetByIDWithRole(ctx, userID)
	}

	user, err := app.cacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	// cache miss
	if user == nil {
		user, err = app.store.Users.GetByIDWithRole(ctx, userID)
		if err != nil {
			return nil, err
		}

		err = app.cacheStorage.Users.Set(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}
