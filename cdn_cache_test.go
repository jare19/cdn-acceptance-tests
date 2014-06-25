package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// Should cache first response for an unspecified period of time if when it
// doesn't specify it's own cache headers. Subsequent requests should return
// a cached response.
func TestCacheFirstResponse(t *testing.T) {
	ResetBackends(backendsByPriority)

	testRequestsCachedIndefinite(t, nil)
}

// Should cache responses for the period defined in a `Expires: n` response
// header.
func TestCacheExpires(t *testing.T) {
	ResetBackends(backendsByPriority)

	const cacheDuration = time.Duration(5 * time.Second)

	handler := func(w http.ResponseWriter) {
		headerValue := time.Now().UTC().Add(cacheDuration).Format(http.TimeFormat)
		w.Header().Set("Expires", headerValue)
	}

	testRequestsCachedDuration(t, handler, cacheDuration)
}

// Should cache responses for the period defined in a `Cache-Control:
// max-age=n` response header.
func TestCacheCacheControlMaxAge(t *testing.T) {
	ResetBackends(backendsByPriority)

	const cacheDuration = time.Duration(5 * time.Second)
	headerValue := fmt.Sprintf("max-age=%.0f", cacheDuration.Seconds())

	handler := func(w http.ResponseWriter) {
		w.Header().Set("Cache-Control", headerValue)
	}

	testRequestsCachedDuration(t, handler, cacheDuration)
}

// Should cache responses for the period defined in a `Cache-Control:
// max-age=n` response header when a `Expires: n*2` header is also present.
func TestCacheExpiresAndMaxAge(t *testing.T) {
	ResetBackends(backendsByPriority)

	const cacheDuration = time.Duration(5 * time.Second)
	const expiresDuration = cacheDuration * 2

	maxAgeValue := fmt.Sprintf("max-age=%.0f", cacheDuration.Seconds())

	handler := func(w http.ResponseWriter) {
		expiresValue := time.Now().UTC().Add(expiresDuration).Format(http.TimeFormat)

		w.Header().Set("Expires", expiresValue)
		w.Header().Set("Cache-Control", maxAgeValue)
	}

	testRequestsCachedDuration(t, handler, cacheDuration)
}

// Should cache responses with a `Cache-Control: no-cache` header. Varnish
// doesn't respect this by default.
func TestCacheCacheControlNoCache(t *testing.T) {
	ResetBackends(backendsByPriority)

	handler := func(w http.ResponseWriter) {
		w.Header().Set("Cache-Control", "no-cache")
	}

	testRequestsCachedIndefinite(t, handler)
}

// Should cache responses with a status code of 404. It's a common
// misconception that 404 responses shouldn't be cached; they should because
// they can be expensive to generate.
func TestCache404Response(t *testing.T) {
	ResetBackends(backendsByPriority)

	handler := func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusNotFound)
	}

	testRequestsCachedIndefinite(t, handler)
}
