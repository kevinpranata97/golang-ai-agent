package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kevinpranata97/golang-ai-agent/internal/apptesting"
	"github.com/kevinpranata97/golang-ai-agent/internal/codegen"
	"github.com/kevinpranata97/golang-ai-agent/internal/requirements"
)

func main() {
	// Initialize requirement analyzer
	geminiAPIKey := requirements.GetGeminiAPIKey()
	reqAnalyzer := requirements.NewRequirementAnalyzer(geminiAPIKey)
	
	// Initialize code generator
	outputDir := "./generated_apps"
	codeGen := codegen.NewCodeGenerator(outputDir)
	
	// Initialize application tester
	appTester := apptesting.NewApplicationTester(outputDir)

	// Setup HTTP routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "running",
			"agent":  "golang-ai-agent",
			"features": []string{
				"application_generation",
				"code_testing",
				"requirement_analysis",
				"github_integration",
			},
		})
	})

	// New endpoint for generating applications
	http.HandleFunc("/generate-app", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			Description string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if request.Description == "" {
			http.Error(w, "Description is required", http.StatusBadRequest)
			return
		}

		// Analyze requirements
		appReq, err := reqAnalyzer.AnalyzeRequirements(request.Description)
		if err != nil {
			log.Printf("Failed to analyze requirements: %v", err)
			http.Error(w, fmt.Sprintf("Failed to analyze requirements: %v", err), http.StatusInternalServerError)
			return
		}

		// Validate requirements
		if err := reqAnalyzer.ValidateRequirements(appReq); err != nil {
			log.Printf("Invalid requirements: %v", err)
			http.Error(w, fmt.Sprintf("Invalid requirements: %v", err), http.StatusBadRequest)
			return
		}

		// Generate application
		if err := codeGen.GenerateApplication(appReq); err != nil {
			log.Printf("Failed to generate application: %v", err)
			http.Error(w, fmt.Sprintf("Failed to generate application: %v", err), http.StatusInternalServerError)
			return
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Application generated successfully",
			"app": map[string]interface{}{
				"name":        appReq.Name,
				"type":        appReq.Type,
				"language":    appReq.Language,
				"framework":   appReq.Framework,
				"entities":    len(appReq.Entities),
				"endpoints":   len(appReq.Endpoints),
				"output_dir":  filepath.Join(outputDir, strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-"))),
			},
		})
	})

	// New endpoint for testing generated applications
	http.HandleFunc("/test-app", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			AppPath string `json:"app_path"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if request.AppPath == "" {
			http.Error(w, "App path is required", http.StatusBadRequest)
			return
		}

		// Check if app path exists
		if _, err := os.Stat(request.AppPath); os.IsNotExist(err) {
			http.Error(w, "Application path does not exist", http.StatusNotFound)
			return
		}

		// Load application requirements (this would typically be saved during generation)
		// For now, we'll create a basic requirement structure
		appReq := &requirements.ApplicationRequirement{
			Name:     filepath.Base(request.AppPath),
			Type:     "api", // Default assumption
			Language: "go",
		}

		// Run tests
		testSuite, err := appTester.TestApplication(request.AppPath, appReq)
		if err != nil {
			log.Printf("Failed to test application: %v", err)
			http.Error(w, fmt.Sprintf("Failed to test application: %v", err), http.StatusInternalServerError)
			return
		}

		// Save test results
		resultsPath := filepath.Join(request.AppPath, "test_results.json")
		if err := appTester.SaveTestResults(testSuite, resultsPath); err != nil {
			log.Printf("Failed to save test results: %v", err)
		}

		// Return test results
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":      true,
			"message":      "Application testing completed",
			"test_suite":   testSuite,
			"results_file": resultsPath,
		})
	})

	// Combined endpoint for generating and testing applications
	http.HandleFunc("/generate-and-test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			Description string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if request.Description == "" {
			http.Error(w, "Description is required", http.StatusBadRequest)
			return
		}

		// Analyze requirements
		appReq, err := reqAnalyzer.AnalyzeRequirements(request.Description)
		if err != nil {
			log.Printf("Failed to analyze requirements: %v", err)
			http.Error(w, fmt.Sprintf("Failed to analyze requirements: %v", err), http.StatusInternalServerError)
			return
		}

		// Validate requirements
		if err := reqAnalyzer.ValidateRequirements(appReq); err != nil {
			log.Printf("Invalid requirements: %v", err)
			http.Error(w, fmt.Sprintf("Invalid requirements: %v", err), http.StatusBadRequest)
			return
		}

		// Generate application
		if err := codeGen.GenerateApplication(appReq); err != nil {
			log.Printf("Failed to generate application: %v", err)
			http.Error(w, fmt.Sprintf("Failed to generate application: %v", err), http.StatusInternalServerError)
			return
		}

		appPath := filepath.Join(outputDir, strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")))

		// Test the generated application
		testSuite, err := appTester.TestApplication(appPath, appReq)
		if err != nil {
			log.Printf("Failed to test application: %v", err)
			// Don't fail the entire request if testing fails
		}

		// Save test results if testing was successful
		var resultsPath string
		if testSuite != nil {
			resultsPath = filepath.Join(appPath, "test_results.json")
			if err := appTester.SaveTestResults(testSuite, resultsPath); err != nil {
				log.Printf("Failed to save test results: %v", err)
			}
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"success": true,
			"message": "Application generated and tested successfully",
			"app": map[string]interface{}{
				"name":        appReq.Name,
				"type":        appReq.Type,
				"language":    appReq.Language,
				"framework":   appReq.Framework,
				"entities":    len(appReq.Entities),
				"endpoints":   len(appReq.Endpoints),
				"output_dir":  appPath,
			},
		}

		if testSuite != nil {
			response["test_results"] = map[string]interface{}{
				"total_tests":    testSuite.TotalTests,
				"passed_tests":   testSuite.PassedTests,
				"failed_tests":   testSuite.FailedTests,
				"skipped_tests":  testSuite.SkippedTests,
				"coverage":       testSuite.Coverage,
				"duration":       testSuite.Duration.String(),
				"results_file":   resultsPath,
				"summary":        testSuite.Summary,
			}
		}

		json.NewEncoder(w).Encode(response)
	})

	// Webhook endpoint (existing functionality)
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Process webhook (existing logic)
		log.Println("Webhook received")
		w.WriteHeader(http.StatusOK)
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET  /health - Health check")
	log.Printf("  GET  /status - Agent status")
	log.Printf("  POST /generate-app - Generate application from description")
	log.Printf("  POST /test-app - Test generated application")
	log.Printf("  POST /generate-and-test - Generate and test application")
	log.Printf("  POST /webhook - GitHub webhook")
	
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

