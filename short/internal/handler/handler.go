package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/polishedfeedback/go-utils/short/internal/storage"
)

type Handler struct {
	storage storage.Storage
}

func NewHandler(s storage.Storage) *Handler {
	return &Handler{storage: s}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	switch r.Method {
	case http.MethodGet:
		longUrl, err := h.storage.Get(shortCode)
		if err != nil {
			http.Error(w, "Error getting the URL", http.StatusNotFound)
			return
		}
		h.storage.IncrementClicks(shortCode)
		http.Redirect(w, r, longUrl, http.StatusFound)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading the body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		err = h.storage.Save(shortCode, string(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("URL created: " + shortCode))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
