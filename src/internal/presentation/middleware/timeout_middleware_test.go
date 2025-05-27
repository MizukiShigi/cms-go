package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeoutMiddleware_FastHandler(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	middleware := TimeoutMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "success", rr.Body.String())
}

func TestTimeoutMiddleware_SlowHandlerWithinTimeout(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("slow but ok"))
	})

	middleware := TimeoutMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	start := time.Now()
	middleware.ServeHTTP(rr, req)
	duration := time.Since(start)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "slow but ok", rr.Body.String())
	assert.True(t, duration < 10*time.Second)
	assert.True(t, duration >= 100*time.Millisecond)
}

func TestTimeoutMiddleware_ContextCancellation(t *testing.T) {
	var contextCancelled bool

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			contextCancelled = true
			return
		case <-time.After(15 * time.Second):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("should not reach here"))
		}
	})

	middleware := TimeoutMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	start := time.Now()
	middleware.ServeHTTP(rr, req)
	duration := time.Since(start)

	assert.True(t, contextCancelled)
	assert.True(t, duration >= 10*time.Second)
	assert.True(t, duration < 11*time.Second)
}

func TestTimeoutMiddleware_ContextValues(t *testing.T) {
	var hasDeadline bool
	var deadlineTime time.Time

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deadline, ok := r.Context().Deadline()
		hasDeadline = ok
		deadlineTime = deadline
		w.WriteHeader(http.StatusOK)
	})

	middleware := TimeoutMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	start := time.Now()
	middleware.ServeHTTP(rr, req)

	assert.True(t, hasDeadline)
	assert.True(t, deadlineTime.After(start))
	assert.True(t, deadlineTime.Before(start.Add(11*time.Second)))
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestTimeoutMiddleware_PreservesExistingContext(t *testing.T) {
	existingValue := "test-value"
	existingKey := "test-key"

	var retrievedValue interface{}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		retrievedValue = r.Context().Value(existingKey)
		w.WriteHeader(http.StatusOK)
	})

	middleware := TimeoutMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(context.WithValue(req.Context(), existingKey, existingValue))
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	assert.Equal(t, existingValue, retrievedValue)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestTimeoutMiddleware_MultipleRequests(t *testing.T) {
	handlerCallCount := 0

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCallCount++
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	middleware := TimeoutMiddleware(nextHandler)

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		rr := httptest.NewRecorder()

		middleware.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "success", rr.Body.String())
	}

	assert.Equal(t, 5, handlerCallCount)
}

func TestTimeoutMiddleware_HandlerError(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	})

	middleware := TimeoutMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "bad request", rr.Body.String())
}

func TestTimeoutMiddleware_HandlerPanic(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	middleware := TimeoutMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "/test", nil)
	rr := httptest.NewRecorder()

	assert.Panics(t, func() {
		middleware.ServeHTTP(rr, req)
	})
}

func TestTimeoutMiddleware_ContextDoneBeforeTimeout(t *testing.T) {
	var contextDone bool

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
			contextDone = true
			return
		case <-time.After(100 * time.Millisecond):
			w.WriteHeader(http.StatusOK)
		}
	})

	middleware := TimeoutMiddleware(nextHandler)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	middleware.ServeHTTP(rr, req)

	assert.True(t, contextDone)
}