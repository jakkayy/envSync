package main

import (
	"fmt"
	"os"

	"github.com/jakkayy/envSync/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := server.NewServer(port)
	if err := srv.Run(); err != nil {
		fmt.Printf("Server failure: %v\n", err)
		os.Exit(1)
	}
}
