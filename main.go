package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kevinpranata97/golang-ai-agent/internal/apptesting"
	"github.com/kevinpranata97/golang-ai-agent/internal/codegen"
	"github.com/kevinpranata97/golang-ai-agent/internal/requirements"
)
func main() {
	// Initialize components
	reqAnalyzer := requirements.NewAnalyzer()
	codeGenerator := codegen.NewCodeGenerator()
	appTester := apptesting.NewTester()

	// Define the application requirements (example)
	appReq := &requirements.ApplicationRequirement{
		Name:        "My Awesome App",
		Description: "A simple application to manage users and products.",
		Language:    "Go",
		Type:        "api",
		Framework:   "Gin Gonic",
		Entities: []requirements.Entity{
			{
				Name: "User",
				Fields: []requirements.Field{
					{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
					{Name: "name", Type: "string", Required: true},
					{Name: "email", Type: "email", Required: true, Unique: true},
					{Name: "created_at", Type: "date"},
				},
			},
			{
				Name: "Product",
				Fields: []requirements.Field{
					{Name: "id", Type: "int", PrimaryKey: true, AutoIncrement: true},
					{Name: "name", Type: "string", Required: true},
					{Name: "price", Type: "float", Required: true},
					{Name: "created_at", Type: "date"},
				},
			},
		},
		APIEndpoints: []requirements.APIEndpoint{
			{
				Entity:     "User",
				Operations: []string{"create", "read", "read_all", "update", "delete"},
			},
			{
				Entity:     "Product",
				Operations: []string{"create", "read", "read_all", "update", "delete"},
			},
		},
	}

	// Validate requirements
	if err := reqAnalyzer.ValidateRequirements(appReq); err != nil {
		log.Fatalf("Requirement validation failed: %v", err)
	}

	// Generate application
	appPath, err := codeGenerator.GenerateApplication(appReq)
	if err != nil {
		log.Fatalf("Failed to generate application: %v", err)
	}

	fmt.Printf("Application generated at: %s\n", appPath)

	// Test application
	testResults, err := appTester.RunTests(appPath, appReq)
	if err != nil {
		log.Fatalf("Failed to run tests: %v", err)
	}

	fmt.Println("Test Results:")
	for _, result := range testResults {
		fmt.Printf("- %s: %s\n", result.TestName, result.Status)
		if result.Error != "" {
			fmt.Printf("  Error: %s\n", result.Error)
		}
	}

	// Save test results
	if err := appTester.SaveTestResults(appPath, testResults); err != nil {
		log.Fatalf("Failed to save test results: %v", err)
	}

	fmt.Printf("Test results saved to: %s\n", filepath.Join(appPath, "test_results.json"))

	// Clean up (optional)
	// defer os.RemoveAll(appPath)
}

// Helper function to create directories if they don't exist
func createDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory: %v", err)
		}
	}
}


