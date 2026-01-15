package handler

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/polishedfeedback/paste/internal/storage"
)

type Handler struct {
	storage storage.Storage
}

func NewHandler(s storage.Storage) *Handler {
	return &Handler{storage: s}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimPrefix(r.URL.Path, "/")
	if url == "" {
		http.Error(w, "URL path required", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		content, err := h.storage.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(content))
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error parsing body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		expiryStr := r.URL.Query().Get("expires")
		var expiresIn time.Duration
		if expiryStr != "" {
			expiresIn, err = time.ParseDuration(expiryStr)
			if err != nil {
				http.Error(w, "Invalid expiry format: Use 1h, 30m, 24h, etc.", http.StatusBadRequest)
				return
			}
		}
		err = h.storage.Save(url, string(body), expiresIn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Paste created: " + url))
	case http.MethodDelete:
		err := h.storage.Delete(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Paste deleted: " + url))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
