package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	appcontext "github.com/MizukiShigi/cms-go/internal/domain/context"
)

func TestLoggingMiddleware_SetsRequestID(t *testing.T) {
	var capturedRequestID string
	var capturedMethod string
	var capturedURL string

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequestID = r.Context().Value(appcontext.RequestIDKey).(string)
		capturedMethod = r.Context().Value(appcontext.MethodKey).(string)
		capturedURL = r.Context().Value(appcontext.URLKey).(string)
		w.WriteHeader(http.StatusOK)
	})

	middleware := LoggingMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test-path", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, capturedRequestID)
	assert.Equal(t, "GET", capturedMethod)
	assert.Equal(t, "/test-path", capturedURL)

	assert.Len(t, capturedRequestID, 36)
	assert.Contains(t, capturedRequestID, "-")
}

func TestLoggingMiddleware_DifferentMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			var capturedMethod string

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedMethod = r.Context().Value(appcontext.MethodKey).(string)
				w.WriteHeader(http.StatusOK)
			})

			middleware := LoggingMiddleware(nextHandler)

			req := httptest.NewRequest(method, "/test", nil)
			rr := httptest.NewRecorder()

			middleware.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			assert.Equal(t, method, capturedMethod)
		})
	}
}

func TestLoggingMiddleware_DifferentURLs(t *testing.T) {
	urls := []string{
		"/",
		"/api/v1/posts",
		"/api/v1/posts/123",
		"/users/456/profile",
		"/search?q=test&limit=10",
	}

	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			var capturedURL string

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				capturedURL = r.Context().Value(appcontext.URLKey).(string)
				w.WriteHeader(http.StatusOK)
			})

			middleware := LoggingMiddleware(nextHandler)

			req := httptest.NewRequest("GET", url, nil)
			rr := httptest.NewRecorder()

			middleware.ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
			assert.Equal(t, url, capturedURL)
		})
	}
}

func TestLoggingMiddleware_UniqueRequestIDs(t *testing.T) {
	requestIDs := make(map[string]bool)
	numRequests := 100

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(appcontext.RequestIDKey).(string)
		if requestIDs[requestID] {
			t.Errorf("Duplicate request ID found: %s", requestID)
		}
		requestIDs[requestID] = true
		w.WriteHeader(http.StatusOK)
	})

	middleware := LoggingMiddleware(nextHandler)

	for i := 0; i < numRequests; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
	}

	assert.Len(t, requestIDs, numRequests)
}

func TestLoggingMiddleware_PreservesExistingContext(t *testing.T) {
	existingValue := "existing-value"
	existingKey := "existing-key"

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retrievedValue := r.Context().Value(existingKey)
		assert.Equal(t, existingValue, retrievedValue)

		requestID := r.Context().Value(appcontext.RequestIDKey)
		assert.NotEmpty(t, requestID)

		w.WriteHeader(http.StatusOK)
	})

	middleware := LoggingMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(context.WithValue(req.Context(), existingKey, existingValue))
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestLoggingMiddleware_HandlerError(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	})

	middleware := LoggingMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "Internal Server Error", rr.Body.String())
}

func TestLoggingMiddleware_HandlerPanic(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	middleware := LoggingMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	assert.Panics(t, func() {
		middleware.ServeHTTP(rr, req)
	})
}

func TestLoggingMiddleware_AllContextValuesSet(t *testing.T) {
	var allValuesSet bool

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(appcontext.RequestIDKey)
		method := r.Context().Value(appcontext.MethodKey)
		url := r.Context().Value(appcontext.URLKey)

		allValuesSet = requestID != nil && method != nil && url != nil
		w.WriteHeader(http.StatusOK)
	})

	middleware := LoggingMiddleware(nextHandler)

	req := httptest.NewRequest("POST", "/api/test", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, allValuesSet)
}