package selfheal

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/kevinpranata97/golang-ai-agent/internal/analysis"
	"github.com/kevinpranata97/golang-ai-agent/internal/apptesting"
	"github.com/kevinpranata97/golang-ai-agent/internal/codemod"
	"github.com/kevinpranata97/golang-ai-agent/internal/database"
	"github.com/kevinpranata97/golang-ai-agent/internal/requirements"
	"github.com/kevinpranata97/golang-ai-agent/internal/storage"
))

type SelfHealer struct {
	analyzer *analysis.CodeAnalyzer
	tester *apptesting.ApplicationTester
	db *database.DB
	codeModifier *codemod.CodeModifier
}
func NewSelfHealer(analyzer *analysis.CodeAnalyzer, tester *apptesting.ApplicationTester, db *database.DB) *SelfHealer {
	return &SelfHealer{
		analyzer: analyzer,
		tester: tester,
		db: db,
		codeModifier: codemod.NewCodeModifier(),
	}
}dasarkan hasil analisis dan tes.
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
				fixApplied = true
				logEntry.Status = "applied_complex_fix"
				logEntry.FeedbackJSON = "Successfully applied complex fix."
			}
		} else {
			log.Printf("Suggestion requires complex code modification or missing target file, skipping for now: %s", suggestion.Description)
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
					logEntry.FixStatus = "failed_retest"
					logEntry.FixOutput = fmt.Sprintf("Retest failed: %v", err)
				} else {
					log.Printf("Retest results for project %s: %s", projectID, retestResults.OverallStatus)
					retestResultsJSON, _ := json.Marshal(retestResults)
					logEntry.RetestResultsJSON = string(retestResultsJSON)
					if retestResults.OverallStatus == "success" {
						log.Printf("Self-fix successful for project %s!", projectID)
						logEntry.FixStatus = "success"
						sh.db.InsertInteractionLog(logEntry)
						return nil // Fix successful, exit
					} else {
						log.Printf("Self-fix did not fully resolve issues for project %s.", projectID)
						logEntry.FixStatus = "failed_retest_still_issues"
					}
				}
				sh.db.InsertInteractionLog(logEntry)
			}

	}

	log.Printf("Self-fix attempts completed for project %s. Project still has issues.", projectID)
	return fmt.Errorf("self-fix attempts failed to resolve all issues for project %s", projectID)
}

// PerformUpgrade mencoba melakukan upgrade pada aplikasi berdasarkan hasil analisis.
func (sh *SelfHealer) PerformUpgrade(projectID, appPath string, appReq *requirements.ApplicationRequirement) error {
	log.Printf("Attempting upgrade for project %s at %s", projectID, appPath)

	// Run analysis to find upgrade opportunities
	analysisData, err := sh.analyzer.AnalyzeProject(projectID, appPath, appReq, nil)
	if err != nil {
		return fmt.Errorf("failed to analyze project for upgrade: %w", err)
	}

	analysisDataJSON, _ := json.Marshal(analysisData)

	for i, suggestion := range analysisData.Suggestions {
		// Filter for upgrade-specific suggestions (e.g., performance, quality, non-critical functionality)
		if suggestion.Type == "performance" || suggestion.Type == "quality" || (suggestion.Type == "functionality" && suggestion.Priority != "high") {
			log.Printf("Applying upgrade suggestion %d: %s (Type: %s, Priority: %s)", i+1, suggestion.Description, suggestion.Type, suggestion.Priority)

			suggestionJSON, _ := json.Marshal(suggestion)

			logEntry := database.InteractionLog{
				ID:                     fmt.Sprintf("%s-%d-%d", projectID, time.Now().UnixNano(), i),
				Timestamp:              time.Now(),
				Endpoint:               "self-upgrade",
				AppName:                appReq.Name,
				AppPath:                appPath,
				AnalysisResultsJSON:    string(analysisDataJSON),
				FeedbackJSON:           fmt.Sprintf("Attempting upgrade for: %s", suggestion.Description),
				Status:                 "attempting",
				ProcessedForFinetuning: false,
				AppliedSuggestionJSON:  string(suggestionJSON),
			}
			sh.db.InsertInteractionLog(logEntry)

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
				} else if suggestion.TargetFile != "" && suggestion.Code != "" {
					log.Printf("Attempting complex code modification in %s: %s", suggestion.TargetFile, suggestion.Description)
					err := sh.codeModifier.ApplyModification(appPath, suggestion.TargetFile, suggestion.Code)
					if err != nil {
						log.Printf("Failed to apply complex upgrade: %v", err)
						logEntry.FixStatus = "failed_complex_upgrade"
						logEntry.FixOutput = fmt.Sprintf("Failed complex upgrade: %v", err)
					} else {
						log.Printf("Successfully applied complex upgrade.")
						upgradeApplied = true
						logEntry.FixStatus = "applied_complex_upgrade"
						logEntry.FixOutput = "Successfully applied complex upgrade."
					}
				} else {
					log.Printf("Upgrade suggestion requires complex code modification or missing target file, skipping for now: %s", suggestion.Description)
					logEntry.FixStatus = "skipped_complex_upgrade"
					logEntry.FixOutput = fmt.Sprintf("Skipped complex upgrade: %s", suggestion.Description)
				}


			if upgradeApplied {
				// Re-test after applying an upgrade to ensure no regressions
				log.Printf("Re-testing project %s after applying upgrade...", projectID)
				retestResults, err := sh.tester.TestApplication(appPath, appReq)
				if err != nil {
					log.Printf("Failed to re-test after upgrade: %v", err)
					logEntry.FixStatus = "failed_retest_upgrade"
					logEntry.FixOutput = fmt.Sprintf("Retest failed: %v", err)
				} else {
					log.Printf("Retest results for project %s after upgrade: %s", projectID, retestResults.OverallStatus)
					retestResultsJSON, _ := json.Marshal(retestResults)
					logEntry.RetestResultsJSON = string(retestResultsJSON)
					if retestResults.OverallStatus == "success" {
						logEntry.FixStatus = "success_upgrade"
					} else {
						logEntry.FixStatus = "failed_retest_upgrade_still_issues"
					}
				}
				sh.db.InsertInteractionLog(logEntry)
			}
		}
	}

	log.Printf("Upgrade attempts completed for project %s.", projectID)
	return nil
}


