package router

import (
	"github.com/go-chi/chi/v5"
	"go-url-shortener/internal/handler" 
)

func APIRouter(r chi.Router, h *handler.ShortenerHandler) { 
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/shorten", h.ShortenURLHandler) 
	})
}