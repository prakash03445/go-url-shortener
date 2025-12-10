package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os" // Needed for os.Getenv in this simple handler setup

	"go-url-shortener/internal/model"
)


// URLService defines the methods the handler needs from the service layer.
type URLService interface {
	Shorten(longURL string) (*model.URL, error)
}

// --- Handler Struct ---

// ShortenerHandler holds the URL service dependency and necessary configuration.
type ShortenerHandler struct {
	Service URLService
	BaseURL string // e.g., "http://localhost:8080/"
}

// NewShortenerHandler creates a new handler instance, injecting the dependencies.
func NewShortenerHandler(s URLService) *ShortenerHandler {
	// For simplicity, we get the base URL from the environment. 
    // In a production system, this would be set via configuration.
    baseURL := os.Getenv("BASE_URL")
    if baseURL == "" {
        baseURL = "http://localhost:8080/" // Default for local testing
    }
    
	return &ShortenerHandler{
        Service: s,
        BaseURL: baseURL,
    }
}

// ShortenURLHandler handles POST /api/v1/shorten
func (h *ShortenerHandler) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var req model.ShortenRequest

	// 1. Decode Request Body (Read JSON)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Handler Error: Invalid request body: %v", err)
		http.Error(w, "Invalid request body format (expected JSON)", http.StatusBadRequest)
		return
	}
    
	// 2. Input Validation
    // NOTE: We are relying on the 'url' validation later. For now, a simple check.
    if req.LongURL == "" {
        http.Error(w, "Field 'long_url' is required.", http.StatusBadRequest)
        return
    }

	// 3. Call Service Layer to Handle Business Logic (Generate Code, Save to DB)
	urlModel, err := h.Service.Shorten(req.LongURL)
	if err != nil {
		log.Printf("Service Error: Failed to save URL: %v", err)
		// Return a generic error to the client
		http.Error(w, "Failed to process shortening request.", http.StatusInternalServerError)
		return
	}

	// 4. Prepare and Send Response
	response := model.ShortenResponse{
		ShortURL: h.BaseURL + urlModel.ShortCode,
		LongURL:  urlModel.LongURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // HTTP 201 Created
	if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Handler Error: Failed to encode response: %v", err)
        // Note: Cannot change status code after writing header, but log the error.
    }
}