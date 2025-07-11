package main

import (
	"testing"

	"github.com/kevinpranata97/golang-ai-agent/internal/agent"
	"github.com/kevinpranata97/golang-ai-agent/internal/github"
	"github.com/kevinpranata97/golang-ai-agent/internal/storage"
	testingpkg "github.com/kevinpranata97/golang-ai-agent/internal/testing"
	"github.com/kevinpranata97/golang-ai-agent/internal/workflow"
)

func TestAgentInitialization(t *testing.T) {
	// Test storage initialization
	storage := storage.NewFileStorage("./test_data")
	if storage == nil {
		t.Fatal("Failed to initialize storage")
	}

	// Test GitHub client initialization
	githubClient := github.NewClient("test_token")
	if githubClient == nil {
		t.Fatal("Failed to initialize GitHub client")
	}

	// Test test runner initialization
	testRunner := testingpkg.NewTestRunner()
	if testRunner == nil {
		t.Fatal("Failed to initialize test runner")
	}

	// Test workflow engine initialization
	workflowEngine := workflow.NewEngine()
	if workflowEngine == nil {
		t.Fatal("Failed to initialize workflow engine")
	}

	// Test agent initialization
	aiAgent := agent.NewAgent(storage, githubClient, testRunner, workflowEngine)
	if aiAgent == nil {
		t.Fatal("Failed to initialize AI agent")
	}
}

func TestWorkflowEngine(t *testing.T) {
	engine := workflow.NewEngine()
	
	// Test workflow execution
	ctx := workflow.Context{
		Repository: "test/repo",
		CloneURL:   "https://github.com/test/repo.git",
		Ref:        "refs/heads/main",
		Commits:    []workflow.Commit{},
	}
	
	result := engine.ExecuteWorkflow("ci_cd", ctx)
	if result.Duration == 0 {
		t.Error("Workflow execution should have a duration")
	}
}

func TestStorage(t *testing.T) {
	storage := storage.NewFileStorage("./test_data")
	
	// Test storing data
	testData := map[string]interface{}{
		"test_key": "test_value",
	}
	
	err := storage.Store("test_item", testData)
	if err != nil {
		t.Fatalf("Failed to store data: %v", err)
	}
	
	// Test retrieving data
	var retrieved map[string]interface{}
	err = storage.Retrieve("test_item", &retrieved)
	if err != nil {
		t.Fatalf("Failed to retrieve data: %v", err)
	}
	
	if retrieved["test_key"] != "test_value" {
		t.Error("Retrieved data doesn't match stored data")
	}
	
	// Cleanup
	storage.Delete("test_item")
}


