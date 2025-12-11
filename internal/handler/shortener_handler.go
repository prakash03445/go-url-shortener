package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"errors"

	"go-url-shortener/internal/model"
	"github.com/go-chi/chi/v5"
	"go-url-shortener/internal/service"
)


type URLService interface {
	Shorten(longURL string) (*model.URL, error)
	ResolveURL(shortCode string) (string, error)
}


type ShortenerHandler struct {
	Service URLService
	BaseURL string
}

func NewShortenerHandler(s URLService) *ShortenerHandler {
    baseURL := os.Getenv("BASE_URL")
    if baseURL == "" {
        baseURL = "http://localhost:8080/"
    }
    
	return &ShortenerHandler{
        Service: s,
        BaseURL: baseURL,
    }
}

func (h *ShortenerHandler) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	var req model.ShortenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Handler Error: Invalid request body: %v", err)
		http.Error(w, "Invalid request body format (expected JSON)", http.StatusBadRequest)
		return
	}
    
    if req.LongURL == "" {
        http.Error(w, "Field 'long_url' is required.", http.StatusBadRequest)
        return
    }

	urlModel, err := h.Service.Shorten(req.LongURL)
	if err != nil {
		log.Printf("Service Error: Failed to save URL: %v", err)
		http.Error(w, "Failed to process shortening request.", http.StatusInternalServerError)
		return
	}

	response := model.ShortenResponse{
		ShortURL: h.BaseURL + urlModel.ShortCode,
		LongURL:  urlModel.LongURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Handler Error: Failed to encode response: %v", err)
    }
}

func (h *ShortenerHandler) ResolveURLHandler(w http.ResponseWriter, r *http.Request) {

	shortCode := chi.URLParam(r, "short_code")
	if shortCode == "" {
		http.Error(w, "Short code is missing.", http.StatusBadRequest) 
		return
	}

	longURL, err := h.Service.ResolveURL(shortCode)

	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			http.Error(w, "Short URL not found.", http.StatusNotFound)
			return
		}
		
		log.Printf("Error resolving URL %s: %v", shortCode, err)
		http.Error(w, "Internal server error during lookup.", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound) 
}