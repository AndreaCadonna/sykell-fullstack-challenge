package main

import (
	"log"
	"os"
)

func main() {
	log.Println("Starting Web Crawler API...")

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server will start on port %s", port)
	log.Println("Initial scaffolding - full implementation coming in next commits")

	// Placeholder for now - will be replaced with Gin server setup
	select {}
}
