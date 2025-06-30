package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/shyam078/internal-transfer-system/db"
	"github.com/shyam078/internal-transfer-system/handlers"
)

func main() {
	log.Println("Initializing database...")
	db.Init()
	log.Println("Database initialized.")

	r := chi.NewRouter()
	log.Println("Setting up routes...")
	r.Post("/accounts", handlers.CreateAccount)
	r.Get("/accounts/{account_id}", handlers.GetAccount)
	r.Post("/transactions", handlers.CreateTransaction)
	log.Println("Routes set up.")
	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
