package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/polishedfeedback/paste/internal/handler"
	"github.com/polishedfeedback/paste/internal/storage"
)

func main() {
	store, err := storage.NewSQLiteStorage("pastes.db")
	if err != nil {
		log.Fatal(err)
	}
	h := handler.NewHandler(store)

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", h)
}
