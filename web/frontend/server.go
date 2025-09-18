package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Serve static files from current directory
	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)

	port := "3000"
	if p := os.Getenv("FRONTEND_PORT"); p != "" {
		port = p
	}

	fmt.Printf("Frontend server starting on http://localhost:%s\n", port)
	fmt.Println("Open your browser to view the React app!")

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
