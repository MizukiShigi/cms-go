package helper

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

func TestRespondWithJSON(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		data       interface{}
		wantStatus int
		wantBody   string
	}{
		{
			name:       "Simple string response",
			statusCode: http.StatusOK,
			data:       "test message",
			wantStatus: http.StatusOK,
			wantBody:   `"test message"`,
		},
		{
			name:       "JSON object response",
			statusCode: http.StatusCreated,
			data:       map[string]string{"message": "success"},
			wantStatus: http.StatusCreated,
			wantBody:   `{"message":"success"}`,
		},
		{
			name:       "Nil data",
			statusCode: http.StatusNoContent,
			data:       nil,
			wantStatus: http.StatusNoContent,
			wantBody:   "null",
		},
		{
			name:       "Error status with data",
			statusCode: http.StatusBadRequest,
			data:       map[string]string{"error": "invalid input"},
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"invalid input"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			
			RespondWithJSON(w, tt.statusCode, tt.data)
			
			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("RespondWithJSON() status = %v, want %v", w.Code, tt.wantStatus)
			}
			
			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("RespondWithJSON() content-type = %v, want application/json", contentType)
			}
			
			// Check body content (removing newlines for comparison)
			body := strings.TrimSpace(w.Body.String())
			if body != tt.wantBody {
				t.Errorf("RespondWithJSON() body = %v, want %v", body, tt.wantBody)
			}
		})
	}
}

func TestRespondWithJSON_InvalidJSON(t *testing.T) {
	w := httptest.NewRecorder()
	
	// Create data that cannot be marshaled to JSON
	invalidData := make(chan int)
	
	RespondWithJSON(w, http.StatusOK, invalidData)
	
	// Should return internal server error when JSON encoding fails
	if w.Code != http.StatusInternalServerError {
		t.Errorf("RespondWithJSON() with invalid data status = %v, want %v", w.Code, http.StatusInternalServerError)
	}
	
	body := strings.TrimSpace(w.Body.String())
	expectedBody := `{"type":"internal_error","message":"Failed to generate response"}`
	if body != expectedBody {
		t.Errorf("RespondWithJSON() with invalid data body = %v, want %v", body, expectedBody)
	}
}

func TestRespondWithError_MyError(t *testing.T) {
	tests := []struct {
		name       string
		err        *valueobject.MyError
		wantStatus int
	}{
		{
			name:       "Invalid error",
			err:        valueobject.NewMyError(valueobject.InvalidCode, "Invalid input"),
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Not found error",
			err:        valueobject.NewMyError(valueobject.NotFoundCode, "Resource not found"),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "Unauthorized error",
			err:        valueobject.NewMyError(valueobject.UnauthorizedCode, "Unauthorized"),
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Forbidden error",
			err:        valueobject.NewMyError(valueobject.ForbiddenCode, "Forbidden"),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "Conflict error",
			err:        valueobject.NewMyError(valueobject.ConflictCode, "Conflict"),
			wantStatus: http.StatusConflict,
		},
		{
			name:       "Internal server error",
			err:        valueobject.NewMyError(valueobject.InternalServerErrorCode, "Internal error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			
			RespondWithError(w, tt.err)
			
			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("RespondWithError() status = %v, want %v", w.Code, tt.wantStatus)
			}
			
			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("RespondWithError() content-type = %v, want application/json", contentType)
			}
			
			// Check that response contains error information
			body := w.Body.String()
			if !strings.Contains(body, string(tt.err.Code)) {
				t.Errorf("RespondWithError() body should contain error code %v", tt.err.Code)
			}
			if !strings.Contains(body, tt.err.Message) {
				t.Errorf("RespondWithError() body should contain error message %v", tt.err.Message)
			}
		})
	}
}

func TestRespondWithError_StandardError(t *testing.T) {
	w := httptest.NewRecorder()
	
	// Test with standard Go error
	err := errors.New("standard error")
	
	RespondWithError(w, err)
	
	// Should default to internal server error
	if w.Code != http.StatusInternalServerError {
		t.Errorf("RespondWithError() with standard error status = %v, want %v", w.Code, http.StatusInternalServerError)
	}
	
	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("RespondWithError() content-type = %v, want application/json", contentType)
	}
	
	// Check that response contains internal server error
	body := w.Body.String()
	if !strings.Contains(body, string(valueobject.InternalServerErrorCode)) {
		t.Error("RespondWithError() should return internal server error for standard errors")
	}
}

func TestRespondWithError_NilError(t *testing.T) {
	w := httptest.NewRecorder()
	
	RespondWithError(w, nil)
	
	// Should default to internal server error
	if w.Code != http.StatusInternalServerError {
		t.Errorf("RespondWithError() with nil error status = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}