package main

import (
	"fmt"
	"net/http"

	"github.com/polishedfeedback/paste/internal/handler"
	"github.com/polishedfeedback/paste/internal/storage"
)

func main() {
	store := storage.NewMemoryStorage()
	h := handler.NewHandler(store)

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", h)
}
