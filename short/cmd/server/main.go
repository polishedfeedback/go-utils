package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/polishedfeedback/go-utils/short/internal/handler"
	"github.com/polishedfeedback/go-utils/short/internal/storage"
)

func main() {
	s, err := storage.NewSQLiteStorage("urls.db")
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	h := handler.NewHandler(s)
	fmt.Println("Starting server on port 8080...")
	http.ListenAndServe(":8080", h)
}
