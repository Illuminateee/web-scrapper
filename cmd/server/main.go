package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Illuminateee/web-scrapper.git/internal/api"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Initialize router
	router := mux.NewRouter()

	// Setup API routes
	api.SetupRoutes(router)

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Wrap router with CORS
	handler := c.Handler(router)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
