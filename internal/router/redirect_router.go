package router

import (
	"github.com/go-chi/chi/v5"
	"go-url-shortener/internal/handler" 
)

func RedirectRouter(r chi.Router, h *handler.ShortenerHandler) {

	r.Get("/{short_code}", h.ResolveURLHandler) 
}