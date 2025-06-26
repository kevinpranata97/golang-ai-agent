package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kevinpranata97/golang-ai-agent/internal/agent"
	"github.com/kevinpranata97/golang-ai-agent/internal/github"
	"github.com/kevinpranata97/golang-ai-agent/internal/storage"
	"github.com/kevinpranata97/golang-ai-agent/internal/testing"
	"github.com/kevinpranata97/golang-ai-agent/internal/workflow"
)

func main() {
	// Initialize components
	storage := storage.NewFileStorage("./data")
	githubClient := github.NewClient(os.Getenv("GITHUB_TOKEN"))
	testRunner := testing.NewTestRunner()
	workflowEngine := workflow.NewEngine()
	
	// Initialize AI agent
	aiAgent := agent.NewAgent(storage, githubClient, testRunner, workflowEngine)
	
	// Setup HTTP server for webhooks
	http.HandleFunc("/webhook", aiAgent.HandleWebhook)
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/status", aiAgent.GetStatus)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("AI Agent starting on port %s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	})
}

