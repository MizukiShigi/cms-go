package helper

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/MizukiShigi/cms-go/internal/domain/valueobject"
)

func RespondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encoding failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"type":"internal_error","message":"Failed to generate response"}`))
	}
}

func RespondWithError(w http.ResponseWriter, err error) {
	var domainErr *valueobject.MyError
	if errors.As(err, &domainErr) {
		RespondWithJSON(w, domainErr.StatusCode(), domainErr)
		return
	}

	myError := valueobject.InternalServerError
	RespondWithJSON(w, myError.StatusCode(), myError)
}
