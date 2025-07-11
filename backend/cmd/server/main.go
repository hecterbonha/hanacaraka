package main

import (
	"log"
	"net/http"

	"hanacaraka/internal/api"
	"hanacaraka/internal/middleware"
)

func main() {
	mux := http.NewServeMux()
	api.RegisterRoutes(mux)

	handler := middleware.Apply(mux)

	log.Println("Server running at http://localhost:8080")
	port := ":8080"
	err := http.ListenAndServe(port, handler)
	if err != nil {
		log.Fatal(err)
	}
}
