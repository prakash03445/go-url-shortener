// internal/router/api_router.go
package router

import (
	"github.com/go-chi/chi/v5"
	"go-url-shortener/internal/handler" 
)

// APIRouter defines all API endpoints and applies a /api/v1 prefix.
// It accepts the ShortenerHandler instance 'h'.
func APIRouter(r chi.Router, h *handler.ShortenerHandler) { 
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/shorten", h.ShortenURLHandler) 
	})
}