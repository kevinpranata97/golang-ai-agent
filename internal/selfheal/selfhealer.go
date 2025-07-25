package selfheal

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/kevinpranata97/golang-ai-agent/internal/analysis"
	"github.com/kevinpranata97/golang-ai-agent/internal/apptesting"
	"github.com/kevinpranata97/golang-ai-agent/internal/database"
	"github.com/kevinpranata97/golang-ai-agent/internal/requirements"
	"github.com/kevinpranata97/golang-ai-agent/internal/storage"
)

type SelfHealer struct {
	analyzer *analysis.CodeAnalyzer
	tester *apptesting.ApplicationTester
	db *database.DB
	// Mungkin perlu referensi ke CodeModifier di sini
}

func NewSelfHealer(analyzer *analysis.CodeAnalyzer, tester *apptesting.ApplicationTester, db *database.DB) *SelfHealer {
	return &SelfHealer{
		analyzer: analyzer,
		tester: tester,
		db: db,
	}
}

// AttemptSelfFix mencoba memperbaiki aplikasi berdasarkan hasil analisis dan tes.
func (sh *SelfHealer) AttemptSelfFix(projectID, appPath string, appReq *requirements.ApplicationRequirement) error {
	log.Printf("Attempting self-fix for project %s at %s", projectID, appPath)

	initialTestResults, err := sh.tester.TestApplication(appPath, appReq)
	if err != nil {
		return fmt.Errorf("failed to run initial tests for self-fix: %w", err)
	}

	if initialTestResults.OverallStatus == "success" {
		log.Printf("Initial tests passed for project %s. No self-fix needed.", projectID)
		return nil
	}

	log.Printf("Initial tests failed for project %s. Analyzing for suggestions...", projectID)

	analysisData, err := sh.analyzer.AnalyzeProject(projectID, appPath, appReq, initialTestResults)
	if err != nil {
		return fmt.Errorf("failed to analyze project for self-fix: %w", err)
	}

	for i, suggestion := range analysisData.Suggestions {
		log.Printf("Applying suggestion %d: %s (Type: %s, Priority: %s)", i+1, suggestion.Description, suggestion.Type, suggestion.Priority)

		// Log the attempt
		logEntry := database.InteractionLog{
			ID:                     fmt.Sprintf("%s-%d-%d", projectID, time.Now().UnixNano(), i),
			Timestamp:              time.Now(),
			Endpoint:               "self-fix",
			AppName:                appReq.Name,
			AppPath:                appPath,
			AnalysisResultsJSON:    "", // TODO: Marshal analysisData
			FeedbackJSON:           fmt.Sprintf("Attempting fix for: %s", suggestion.Description),
			Status:                 "attempting",
			ProcessedForFinetuning: false,
		}
		// sh.db.InsertInteractionLog(logEntry) // Uncomment after implementing CodeModifier and marshaling

		// Apply the fix based on suggestion type and \'Code\' field
		fixApplied := false
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
				fixApplied = true
				logEntry.Status = "applied_fix_command"
				logEntry.FeedbackJSON = fmt.Sprintf("Applied command: %s, Output: %s", suggestion.Code, string(output))
			}
		} else {
			// This is where CodeModifier would be called for more complex code changes
			log.Printf("Suggestion requires complex code modification, skipping for now: %s", suggestion.Description)
			logEntry.Status = "skipped_complex_fix"
			logEntry.FeedbackJSON = fmt.Sprintf("Skipped complex fix: %s", suggestion.Description)
		}

		// sh.db.InsertInteractionLog(logEntry) // Uncomment after implementing CodeModifier and marshaling

		if fixApplied {
			// Re-test after applying a fix
			log.Printf("Re-testing project %s after applying fix...", projectID)
			retestResults, err := sh.tester.TestApplication(appPath, appReq)
			if err != nil {
				log.Printf("Failed to re-test after fix: %v", err)
				// Update logEntry status to indicate retest failure
			} else {
				// Update logEntry with retest results
				log.Printf("Retest results for project %s: %s", projectID, retestResults.OverallStatus)
				if retestResults.OverallStatus == "success" {
					log.Printf("Self-fix successful for project %s!", projectID)
					return nil // Fix successful, exit
				}
			}
		}
	}

	log.Printf("Self-fix attempts completed for project %s. Project still has issues.", projectID)
	return fmt.Errorf("self-fix attempts failed to resolve all issues for project %s", projectID)
}

// PerformUpgrade mencoba melakukan upgrade pada aplikasi berdasarkan hasil analisis.
func (sh *SelfHealer) PerformUpgrade(projectID, appPath string, appReq *requirements.ApplicationRequirement) error {
	log.Printf("Attempting upgrade for project %s at %s", projectID, appPath)

	// Run analysis to find upgrade opportunities
	analysisData, err := sh.analyzer.AnalyzeProject(projectID, appPath, appReq, nil) // No test results needed for initial upgrade analysis
	if err != nil {
		return fmt.Errorf("failed to analyze project for upgrade: %w", err)
	}

	for i, suggestion := range analysisData.Suggestions {
		// Filter for upgrade-specific suggestions (e.g., performance, quality, non-critical functionality)
		if suggestion.Type == "performance" || suggestion.Type == "quality" || (suggestion.Type == "functionality" && suggestion.Priority != "high") {
			log.Printf("Applying upgrade suggestion %d: %s (Type: %s, Priority: %s)", i+1, suggestion.Description, suggestion.Type, suggestion.Priority)

			// Apply the upgrade based on \'Code\' field
			upgradeApplied := false
			if suggestion.Code != "" {
				cmd := exec.Command("bash", "-c", suggestion.Code)
				cmd.Dir = appPath
				output, cmdErr := cmd.CombinedOutput()
				if cmdErr != nil {
					log.Printf("Failed to apply upgrade (command): %s, Error: %v, Output: %s", suggestion.Code, cmdErr, string(output))
				} else {
					log.Printf("Successfully applied upgrade (command): %s, Output: %s", suggestion.Code, string(output))
					upgradeApplied = true
				}
			} else {
				log.Printf("Upgrade suggestion requires complex code modification, skipping for now: %s", suggestion.Description)
			}

			if upgradeApplied {
				// Re-test after applying an upgrade to ensure no regressions
				log.Printf("Re-testing project %s after applying upgrade...", projectID)
				retestResults, err := sh.tester.TestApplication(appPath, appReq)
				if err != nil {
					log.Printf("Failed to re-test after upgrade: %v", err)
				} else {
					log.Printf("Retest results for project %s after upgrade: %s", projectID, retestResults.OverallStatus)
					// Log upgrade success/failure and new test results
				}
			}
		}
	}

	log.Printf("Upgrade attempts completed for project %s.", projectID)
	return nil
}


