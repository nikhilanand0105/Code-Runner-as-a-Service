package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nikhilanand0105/Code-Runner-as-a-Service/pkg/api"
	"github.com/nikhilanand0105/Code-Runner-as-a-Service/pkg/runner"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	sandboxRunner := runner.NewLocalSandboxRunner()
	server := api.NewServer(sandboxRunner)

	log.Printf("Starting Code Runner as a Service on port :%s", port)
	if err := http.ListenAndServe(":"+port, server.Routes()); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
