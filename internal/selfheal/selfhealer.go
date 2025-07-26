package selfheal
import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"
	"github.com/kevinpranata97/golang-ai-agent/internal/analysis"
	"github.com/kevinpranata97/golang-ai-agent/internal/codemod"
	"github.com/kevinpranata97/golang-ai-agent/internal/database"
	"github.com/kevinpranata97/golang-ai-agent/internal/requirements"

)

type SelfHealer struct {
	analyzer *analysis.CodeAnalyzer
	db *database.DB
	codeModifier *codemod.CodeModifier
}
func NewSelfHealer(analyzer *analysis.CodeAnalyzer, db *database.DB) *SelfHealer {
	return &SelfHealer{
		analyzer: analyzer,
		db: db,
		codeModifier: codemod.NewCodeModifier(),
	}
}
func (sh *SelfHealer) AttemptSelfFix(projectID, appPath string, appReq *requirements.ApplicationRequirement) error {
	log.Printf("Attempting self-fix for project %s at %s", projectID, appPath)

	analysisData, err := sh.analyzer.AnalyzeProject(projectID, appPath, appReq)
	if err != nil {
		return fmt.Errorf("failed to analyze project for self-fix: %w", err)
	}
	analysisDataJSON, _ := json.Marshal(analysisData)
	for i, suggestion := range analysisData.Suggestions {
		log.Printf("Applying suggestion %d: %s (Type: %s, Priority: %s)", i+1, suggestion.Description, suggestion.Type, suggestion.Priority)
		suggestionJSON, _ := json.Marshal(suggestion)
		// Log the attempt
		logEntry := database.InteractionLog{
			ID:                     fmt.Sprintf("%s-%d-%d", projectID, time.Now().UnixNano(), i),
			Timestamp:              time.Now(),
			Endpoint:               "self-fix",
			AppName:                appReq.Name,
			AppPath:                appPath,
			AnalysisResultsJSON:    string(analysisDataJSON),
			FeedbackJSON:           fmt.Sprintf("Attempting fix for: %s", suggestion.Description),
			Status:                 "attempting",
			ProcessedForFinetuning: false,
			AppliedSuggestionJSON:  string(suggestionJSON),
		}
		sh.db.InsertInteractionLog(logEntry)

		// Apply the fix based on suggestion type and \'Code\' field
		if suggestion.Code != "" {
			// Try to execute the code as a shell command
			cmd := exec.Command("bash", "-c", suggestion.Code)
			cmd.Dir = appPath // Execute in the app\'s directory
			output, cmdErr := cmd.CombinedOutput()
			if cmdErr != nil {
				log.Printf("Failed to apply fix (command): %s, Error: %v, Output: %s", suggestion.Code, cmdErr, string(output))
				logEntry.Status = "failed_fix_command"
				logEntry.FeedbackJSON = fmt.Sprintf("Failed command: %s, Error: %v, Output: %s", suggestion.Code, cmdErr, string(output))
			} else {
				log.Printf("Successfully applied fix (command): %s, Output: %s", suggestion.Code, string(output))

				logEntry.Status = "applied_fix_command"
				logEntry.FeedbackJSON = fmt.Sprintf("Applied command: %s, Output: %s", suggestion.Code, string(output))
			}
		} else if suggestion.TargetFile != "" && suggestion.Code != "" {
			// Attempt complex code modification using CodeModifier
			log.Printf("Attempting complex code modification in %s: %s", suggestion.TargetFile, suggestion.Description)
			err := sh.codeModifier.ApplyModification(appPath, suggestion.TargetFile, suggestion.Code)
			if err != nil {
				log.Printf("Failed to apply complex fix: %v", err)
				logEntry.Status = "failed_complex_fix"
				logEntry.FeedbackJSON = fmt.Sprintf("Failed complex fix: %v", err)
			} else {
				log.Printf("Successfully applied complex fix.")

				logEntry.Status = "applied_complex_fix"
				logEntry.FeedbackJSON = "Successfully applied complex fix."
			}
		} else {
			log.Printf("Suggestion requires complex code modification or missing target file, skipping for now: %s", suggestion.Description)
			logEntry.Status = "skipped_complex_fix"
			logEntry.FeedbackJSON = fmt.Sprintf("Skipped complex fix: %s", suggestion.Description)
		}
		// sh.db.InsertInteractionLog(logEntry) // Uncomment after implementing CodeModifier and marshaling

	}
	log.Printf("Self-fix attempts completed for project %s. Project still has issues.", projectID)
	return fmt.Errorf("self-fix attempts failed to resolve all issues for project %s", projectID)
}
// PerformUpgrade mencoba melakukan upgrade pada aplikasi berdasarkan hasil analisis.

