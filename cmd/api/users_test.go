package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/zavista/social-api/internal/store"
	"github.com/zavista/social-api/internal/store/cache"
)

func TestGetUserHandler(t *testing.T) {
	withRedis := config{
		redisCfg: redisConfig{enabled: true},
	}

	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		app := newTestApplication(t, withRedis)
		mux := app.mount()

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("on a cache miss it fetches from the store and populates the cache", func(t *testing.T) {
		app := newTestApplication(t, withRedis)
		mux := app.mount()

		testToken, err := app.authenticator.GenerateToken(1)
		if err != nil {
			t.Fatal(err)
		}

		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)
		mockCacheStore.On("Get", int64(1)).Return(nil, nil)
		mockCacheStore.On("Set", mock.Anything).Return(nil)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)
		// getUser is called twice: once by authTokenMiddleware (for the
		// authenticated caller, claims.UserID) and once by getUserHandler
		// (for the requested profile, the {userID} path param). Both
		// resolve to user 1 here, so each hits the cache.
		mockCacheStore.AssertNumberOfCalls(t, "Get", 2)
		mockCacheStore.AssertNumberOfCalls(t, "Set", 2)
	})

	t.Run("on a cache hit it returns early without writing to the cache", func(t *testing.T) {
		app := newTestApplication(t, withRedis)
		mux := app.mount()

		testToken, err := app.authenticator.GenerateToken(1)
		if err != nil {
			t.Fatal(err)
		}

		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)
		cached := &store.UserWithRole{User: store.User{ID: 1}}
		mockCacheStore.On("Get", int64(1)).Return(cached, nil)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)
		// Both lookups (auth middleware + handler) resolve to user 1 and
		// hit the cache, so Set is never needed.
		mockCacheStore.AssertNumberOfCalls(t, "Get", 2)
		mockCacheStore.AssertNotCalled(t, "Set")
	})

	t.Run("should not hit the cache when redis is disabled", func(t *testing.T) {
		noRedis := config{
			redisCfg: redisConfig{enabled: false},
		}

		app := newTestApplication(t, noRedis)
		mux := app.mount()

		testToken, err := app.authenticator.GenerateToken(1)
		if err != nil {
			t.Fatal(err)
		}

		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)
		mockCacheStore.AssertNotCalled(t, "Get")
		mockCacheStore.AssertNotCalled(t, "Set")
	})
}
