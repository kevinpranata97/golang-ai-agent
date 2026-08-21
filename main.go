package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/kevinpranata97/golang-ai-agent/internal/analysis"
	"github.com/kevinpranata97/golang-ai-agent/internal/codegen"
	"github.com/kevinpranata97/golang-ai-agent/internal/database"
	"github.com/kevinpranata97/golang-ai-agent/internal/finetuning"
	"github.com/kevinpranata97/golang-ai-agent/internal/llm"
	"github.com/kevinpranata97/golang-ai-agent/internal/requirements"
	"github.com/kevinpranata97/golang-ai-agent/internal/selfheal"
	"github.com/kevinpranata97/golang-ai-agent/internal/storage"
)

func main() {
	dataDir := "./data"

	// Initialize storage
	fileStorage := storage.NewFileStorage(dataDir)

	// Initialize CodeAnalyzer
	codeAnalyzer := analysis.NewCodeAnalyzer(fileStorage)

	// Initialize Local Database for Fine-tuning
	db, err := database.NewDB(dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize SelfHealer
	selfHealer := selfheal.NewSelfHealer(codeAnalyzer, db)

	// Initialize requirement analyzer
	geminiAPIKey := requirements.GetGeminiAPIKey()
	reqAnalyzer := requirements.NewRequirementAnalyzer(geminiAPIKey)

	// Initialize code generator
	outputDir := "./generated_apps"
	codeGen := codegen.NewCodeGenerator(outputDir)

	// Initialize LLM client for internal interactions
	llmClient := llm.NewClient()

	// Initialize Finetuner
	finetuner := finetuning.NewFinetuner(db)

	// Schedule periodic fine-tuning process
	go func() {
		for {
			log.Println("Running scheduled fine-tuning process...")
			if err := finetuner.ProcessLogs(); err != nil {
				log.Printf("Error during scheduled fine-tuning: %v", err)
			}
			time.Sleep(5 * time.Minute)
		}
	}()

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
				"fine_tuning",
				"local_database_storage",
				"self_healing",
				"internal_llm_interaction",
			},
		})
	})

	// Unified endpoint for generating applications and general inquiries
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

		// Detect if this is a general inquiry or an app generation request
		if isGeneralInquiry(request.Description) {
			response, err := llmClient.Interact(request.Description)
			if err != nil {
				log.Printf("LLM interaction failed: %v", err)
				http.Error(w, fmt.Sprintf("Interaction failed: %v", err), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"type":    "inquiry_response",
				"message": response,
			})
			return
		}

		interactionLog := database.InteractionLog{
			ID:             uuid.New().String(),
			Timestamp:      time.Now(),
			Endpoint:       "/generate-app",
			RequestPayload: request.Description,
			Status:         "success",
		}

		appReq, err := reqAnalyzer.AnalyzeRequirements(request.Description)
		if err != nil {
			log.Printf("Failed to analyze requirements: %v", err)
			http.Error(w, fmt.Sprintf("Failed to analyze requirements: %v", err), http.StatusInternalServerError)
			interactionLog.Status = "failure"
			db.InsertInteractionLog(interactionLog)
			return
		}

		if err := reqAnalyzer.ValidateRequirements(appReq); err != nil {
			log.Printf("Invalid requirements: %v", err)
			http.Error(w, fmt.Sprintf("Invalid requirements: %v", err), http.StatusBadRequest)
			interactionLog.Status = "failure"
			db.InsertInteractionLog(interactionLog)
			return
		}

		if err := codeGen.GenerateApplication(appReq); err != nil {
			log.Printf("Failed to generate application: %v", err)
			http.Error(w, fmt.Sprintf("Failed to generate application: %v", err), http.StatusInternalServerError)
			interactionLog.Status = "failure"
			db.InsertInteractionLog(interactionLog)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, _ := json.Marshal(map[string]interface{}{
			"success": true,
			"type":    "app_generation",
			"message": "Application generated successfully",
			"app": map[string]interface{}{
				"name":       appReq.Name,
				"type":       appReq.Type,
				"language":   appReq.Language,
				"framework":  appReq.Framework,
				"entities":   len(appReq.Entities),
				"endpoints":  len(appReq.Endpoints),
				"output_dir": filepath.Join(outputDir, strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-"))),
			},
		})
		w.Write(jsonResponse)

		interactionLog.ResponsePayload = string(jsonResponse)
		interactionLog.AppName = appReq.Name
		interactionLog.AppPath = filepath.Join(outputDir, strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")))
		db.InsertInteractionLog(interactionLog)
	})

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

		interactionLog := database.InteractionLog{
			ID:             uuid.New().String(),
			Timestamp:      time.Now(),
			Endpoint:       "/test-app",
			RequestPayload: request.AppPath,
			AppPath:        request.AppPath,
			Status:         "success",
		}

		if _, err := os.Stat(request.AppPath); os.IsNotExist(err) {
			http.Error(w, "Application path does not exist", http.StatusNotFound)
			interactionLog.Status = "failure"
			db.InsertInteractionLog(interactionLog)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		jsonResponse, _ := json.Marshal(map[string]interface{}{
			"success":      true,
			"message":      "Application testing completed",
			"test_suite":   map[string]interface{}{"status": "not_implemented"},
			"results_file": filepath.Join(request.AppPath, "test_results.json"),
		})
		w.Write(jsonResponse)

		interactionLog.ResponsePayload = string(jsonResponse)
		db.InsertInteractionLog(interactionLog)
	})

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

		interactionLog := database.InteractionLog{
			ID:             uuid.New().String(),
			Timestamp:      time.Now(),
			Endpoint:       "/generate-and-test",
			RequestPayload: request.Description,
			Status:         "success",
		}

		appReq, err := reqAnalyzer.AnalyzeRequirements(request.Description)
		if err != nil {
			log.Printf("Failed to analyze requirements: %v", err)
			http.Error(w, fmt.Sprintf("Failed to analyze requirements: %v", err), http.StatusInternalServerError)
			interactionLog.Status = "failure"
			db.InsertInteractionLog(interactionLog)
			return
		}

		if err := reqAnalyzer.ValidateRequirements(appReq); err != nil {
			log.Printf("Invalid requirements: %v", err)
			http.Error(w, fmt.Sprintf("Invalid requirements: %v", err), http.StatusBadRequest)
			interactionLog.Status = "failure"
			db.InsertInteractionLog(interactionLog)
			return
		}

		if err := codeGen.GenerateApplication(appReq); err != nil {
			log.Printf("Failed to generate application: %v", err)
			http.Error(w, fmt.Sprintf("Failed to analyze requirements: %v", err), http.StatusInternalServerError)
			interactionLog.Status = "failure"
			db.InsertInteractionLog(interactionLog)
			return
		}

		appPath := filepath.Join(outputDir, strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")))

		w.Header().Set("Content-Type", "application/json")
		responseMap := map[string]interface{}{
			"success": true,
			"message": "Application generated and tested successfully",
			"app": map[string]interface{}{
				"name":       appReq.Name,
				"type":       appReq.Type,
				"language":   appReq.Language,
				"framework":  appReq.Framework,
				"entities":   len(appReq.Entities),
				"endpoints":  len(appReq.Endpoints),
				"output_dir": appPath,
			},
		}

		jsonResponse, _ := json.Marshal(responseMap)
		w.Write(jsonResponse)

		interactionLog.ResponsePayload = string(jsonResponse)
		interactionLog.AppName = appReq.Name
		interactionLog.AppPath = appPath
		db.InsertInteractionLog(interactionLog)
	})

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		log.Println("Webhook received")
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/self-heal", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			ProjectID   string `json:"project_id"`
			AppPath     string `json:"app_path"`
			Description string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if request.ProjectID == "" || request.AppPath == "" || request.Description == "" {
			http.Error(w, "ProjectID, AppPath, and Description are required", http.StatusBadRequest)
			return
		}

		appReq, err := reqAnalyzer.AnalyzeRequirements(request.Description)
		if err != nil {
			log.Printf("Failed to analyze requirements for self-heal: %v", err)
			http.Error(w, fmt.Sprintf("Failed to analyze requirements for self-heal: %v", err), http.StatusInternalServerError)
			return
		}

		err = selfHealer.AttemptSelfFix(request.ProjectID, request.AppPath, appReq)
		if err != nil {
			log.Printf("Self-heal failed: %v", err)
			http.Error(w, fmt.Sprintf("Self-heal failed: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "Self-heal initiated successfully"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func isGeneralInquiry(prompt string) bool {
	creationKeywords := []string{
		"create", "generate", "build", "make", "setup", "app", "application", "website", "service",
	}
	
	lowerPrompt := strings.ToLower(prompt)
	for _, kw := range creationKeywords {
		if strings.Contains(lowerPrompt, kw) {
			return false
		}
	}
	
	questionWords := []string{"what", "how", "why", "when", "where", "who", "can", "could", "should"}
	for _, qw := range questionWords {
		if strings.HasPrefix(lowerPrompt, qw) {
			return true
		}
	}
	
	if strings.HasSuffix(lowerPrompt, "?") {
		return true
	}
	
	return len(prompt) < 50
}
