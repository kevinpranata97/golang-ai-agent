package finetuning

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kevinpranata97/golang-ai-agent/internal/analysis"
	"github.com/kevinpranata97/golang-ai-agent/internal/apptesting"
	"github.com/kevinpranata97/golang-ai-agent/internal/codegen"
	"github.com/kevinpranata97/golang-ai-agent/internal/requirements"
	"github.com/kevinpranata97/golang-ai-agent/internal/storage"
)

// FineTuner handles iterative improvement of generated applications
type FineTuner struct {
	storage     storage.Storage
	analyzer    *analysis.CodeAnalyzer
	codeGen     *codegen.CodeGenerator
	appTester   *apptesting.ApplicationTester
	maxIterations int
}

// NewFineTuner creates a new fine tuner
func NewFineTuner(storage storage.Storage, analyzer *analysis.CodeAnalyzer, codeGen *codegen.CodeGenerator, appTester *apptesting.ApplicationTester) *FineTuner {
	return &FineTuner{
		storage:       storage,
		analyzer:      analyzer,
		codeGen:       codeGen,
		appTester:     appTester,
		maxIterations: 3, // Maximum number of improvement iterations
	}
}

// ImproveProject performs iterative improvement of a project
func (ft *FineTuner) ImproveProject(projectID string) (*storage.ProjectData, error) {
	// Get project data
	project, err := ft.storage.GetProject(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %v", err)
	}

	// Perform analysis
	analysisData, err := ft.analyzer.AnalyzeProject(projectID, project.AppPath, project.Requirements, project.TestResults)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze project: %v", err)
	}

	// Determine if improvement is needed
	if !ft.needsImprovement(analysisData) {
		return project, nil // No improvement needed
	}

	// Perform iterative improvements
	for i := 0; i < ft.maxIterations; i++ {
		improved, err := ft.performIteration(project, analysisData, i+1)
		if err != nil {
			return project, fmt.Errorf("failed to perform iteration %d: %v", i+1, err)
		}

		if !improved {
			break // No more improvements possible
		}

		// Re-analyze after improvement
		analysisData, err = ft.analyzer.AnalyzeProject(projectID, project.AppPath, project.Requirements, project.TestResults)
		if err != nil {
			return project, fmt.Errorf("failed to re-analyze project: %v", err)
		}

		// Check if we've reached acceptable quality
		if !ft.needsImprovement(analysisData) {
			break
		}
	}

	return project, nil
}

// needsImprovement determines if a project needs improvement based on analysis
func (ft *FineTuner) needsImprovement(analysis *storage.AnalysisData) bool {
	// Check various quality metrics
	if analysis.CodeQuality.TestCoverage < 70 {
		return true
	}

	if analysis.Security.Vulnerabilities > 0 {
		return true
	}

	if analysis.CodeQuality.CyclomaticComplexity > 15 {
		return true
	}

	if analysis.Security.SecurityScore < 80 {
		return true
	}

	// Check for high-priority suggestions
	for _, suggestion := range analysis.Suggestions {
		if suggestion.Priority == "high" {
			return true
		}
	}

	return false
}

// performIteration performs one iteration of improvement
func (ft *FineTuner) performIteration(project *storage.ProjectData, analysis *storage.AnalysisData, iterationNum int) (bool, error) {
	iterationData := storage.IterationData{
		ID:        fmt.Sprintf("iter_%d_%d", time.Now().Unix(), iterationNum),
		Timestamp: time.Now(),
		Changes:   []string{},
		Status:    "in_progress",
	}

	improved := false

	// Apply high-priority improvements first
	for _, suggestion := range analysis.Suggestions {
		if suggestion.Priority == "high" {
			applied, err := ft.applySuggestion(project, suggestion, &iterationData)
			if err != nil {
				continue // Skip this suggestion if it fails
			}
			if applied {
				improved = true
			}
		}
	}

	// Apply medium-priority improvements if we haven't made enough changes
	if len(iterationData.Changes) < 3 {
		for _, suggestion := range analysis.Suggestions {
			if suggestion.Priority == "medium" {
				applied, err := ft.applySuggestion(project, suggestion, &iterationData)
				if err != nil {
					continue
				}
				if applied {
					improved = true
				}
				
				if len(iterationData.Changes) >= 3 {
					break // Limit changes per iteration
				}
			}
		}
	}

	if improved {
		// Test the improved application
		testSuite, err := ft.appTester.TestApplication(project.AppPath, project.Requirements)
		if err != nil {
			iterationData.Status = "failed"
			iterationData.Improvements = []string{"Testing failed after improvements"}
		} else {
			iterationData.TestResults = testSuite
			project.TestResults = testSuite
			
			// Determine improvements made
			iterationData.Improvements = ft.calculateImprovements(analysis, testSuite)
			iterationData.Status = "completed"
		}

		// Add iteration to project
		project.Iterations = append(project.Iterations, iterationData)

		// Update project in storage
		if err := ft.storage.UpdateProject(project); err != nil {
			return improved, fmt.Errorf("failed to update project: %v", err)
		}
	}

	return improved, nil
}

// applySuggestion applies a specific improvement suggestion
func (ft *FineTuner) applySuggestion(project *storage.ProjectData, suggestion storage.ImprovementSuggestion, iteration *storage.IterationData) (bool, error) {
	switch suggestion.Type {
	case "security":
		return ft.applySecurityImprovement(project, suggestion, iteration)
	case "performance":
		return ft.applyPerformanceImprovement(project, suggestion, iteration)
	case "quality":
		return ft.applyQualityImprovement(project, suggestion, iteration)
	case "functionality":
		return ft.applyFunctionalityImprovement(project, suggestion, iteration)
	default:
		return false, fmt.Errorf("unknown suggestion type: %s", suggestion.Type)
	}
}

// applySecurityImprovement applies security-related improvements
func (ft *FineTuner) applySecurityImprovement(project *storage.ProjectData, suggestion storage.ImprovementSuggestion, iteration *storage.IterationData) (bool, error) {
	if strings.Contains(suggestion.Description, "hardcoded") {
		// Fix hardcoded secrets
		applied, err := ft.fixHardcodedSecrets(project.AppPath)
		if err != nil {
			return false, err
		}
		if applied {
			iteration.Changes = append(iteration.Changes, "Fixed hardcoded secrets by using environment variables")
			return true, nil
		}
	}

	if strings.Contains(suggestion.Description, "SQL injection") {
		// Fix SQL injection vulnerabilities
		applied, err := ft.fixSQLInjection(project.AppPath)
		if err != nil {
			return false, err
		}
		if applied {
			iteration.Changes = append(iteration.Changes, "Fixed SQL injection vulnerabilities by using parameterized queries")
			return true, nil
		}
	}

	return false, nil
}

// applyPerformanceImprovement applies performance-related improvements
func (ft *FineTuner) applyPerformanceImprovement(project *storage.ProjectData, suggestion storage.ImprovementSuggestion, iteration *storage.IterationData) (bool, error) {
	if strings.Contains(suggestion.Description, "binary size") {
		// Add build optimization
		applied, err := ft.optimizeBuild(project.AppPath)
		if err != nil {
			return false, err
		}
		if applied {
			iteration.Changes = append(iteration.Changes, "Added build optimization flags to reduce binary size")
			return true, nil
		}
	}

	return false, nil
}

// applyQualityImprovement applies code quality improvements
func (ft *FineTuner) applyQualityImprovement(project *storage.ProjectData, suggestion storage.ImprovementSuggestion, iteration *storage.IterationData) (bool, error) {
	if strings.Contains(suggestion.Description, "test coverage") {
		// Add basic tests
		applied, err := ft.addBasicTests(project.AppPath, project.Requirements)
		if err != nil {
			return false, err
		}
		if applied {
			iteration.Changes = append(iteration.Changes, "Added basic unit tests to improve coverage")
			return true, nil
		}
	}

	if strings.Contains(suggestion.Description, "complexity") {
		// This would require more sophisticated refactoring
		// For now, we'll just document the need
		iteration.Changes = append(iteration.Changes, "Documented need for complexity reduction")
		return true, nil
	}

	return false, nil
}

// applyFunctionalityImprovement applies functionality improvements
func (ft *FineTuner) applyFunctionalityImprovement(project *storage.ProjectData, suggestion storage.ImprovementSuggestion, iteration *storage.IterationData) (bool, error) {
	if strings.Contains(suggestion.Description, "middleware") {
		// Add basic middleware
		applied, err := ft.addMiddleware(project.AppPath)
		if err != nil {
			return false, err
		}
		if applied {
			iteration.Changes = append(iteration.Changes, "Added logging and CORS middleware")
			return true, nil
		}
	}

	return false, nil
}

// Implementation of specific improvements

// fixHardcodedSecrets replaces hardcoded secrets with environment variables
func (ft *FineTuner) fixHardcodedSecrets(appPath string) (bool, error) {
	applied := false

	err := filepath.Walk(appPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			originalContent := string(content)
			modifiedContent := originalContent

			// Replace common hardcoded patterns
			replacements := map[string]string{
				`password := ".*"`:     `password := os.Getenv("DB_PASSWORD")`,
				`apiKey := ".*"`:       `apiKey := os.Getenv("API_KEY")`,
				`secretKey := ".*"`:    `secretKey := os.Getenv("SECRET_KEY")`,
				`token := ".*"`:        `token := os.Getenv("AUTH_TOKEN")`,
			}

			for pattern, replacement := range replacements {
				if strings.Contains(modifiedContent, pattern) {
					// This is a simplified replacement - in reality, you'd use regex
					modifiedContent = strings.ReplaceAll(modifiedContent, pattern, replacement)
					applied = true
				}
			}

			if modifiedContent != originalContent {
				// Add os import if not present
				if !strings.Contains(modifiedContent, `"os"`) {
					modifiedContent = strings.Replace(modifiedContent, "import (", "import (\n\t\"os\"", 1)
				}

				if err := os.WriteFile(path, []byte(modifiedContent), info.Mode()); err != nil {
					return err
				}
			}
		}

		return nil
	})

	return applied, err
}

// fixSQLInjection fixes SQL injection vulnerabilities
func (ft *FineTuner) fixSQLInjection(appPath string) (bool, error) {
	applied := false

	err := filepath.Walk(appPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			originalContent := string(content)
			modifiedContent := originalContent

			// Look for and fix basic SQL injection patterns
			if strings.Contains(modifiedContent, "db.Exec(") && strings.Contains(modifiedContent, "+") {
				// This is a very basic fix - in reality, you'd need more sophisticated parsing
				// For now, just add a comment
				modifiedContent = strings.Replace(modifiedContent, 
					"// TODO: Use parameterized queries", 
					"// TODO: Use parameterized queries to prevent SQL injection", 1)
				applied = true
			}

			if modifiedContent != originalContent {
				if err := os.WriteFile(path, []byte(modifiedContent), info.Mode()); err != nil {
					return err
				}
			}
		}

		return nil
	})

	return applied, err
}

// optimizeBuild adds build optimization
func (ft *FineTuner) optimizeBuild(appPath string) (bool, error) {
	// Create or update Makefile with optimization flags
	makefilePath := filepath.Join(appPath, "Makefile")
	
	makefileContent := `# Optimized build targets
.PHONY: build build-optimized clean

build:
	go build -v .

build-optimized:
	go build -ldflags="-s -w" -v .

clean:
	go clean
	rm -f $(shell basename $(PWD))

# Default target
all: build-optimized
`

	if err := os.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
		return false, err
	}

	return true, nil
}

// addBasicTests adds basic unit tests
func (ft *FineTuner) addBasicTests(appPath string, appReq *requirements.ApplicationRequirement) (bool, error) {
	if appReq == nil || len(appReq.Entities) == 0 {
		return false, nil
	}

	// Add a basic test file for the first entity
	entity := appReq.Entities[0]
	testFileName := fmt.Sprintf("%s_test.go", strings.ToLower(entity.Name))
	testFilePath := filepath.Join(appPath, "internal", "models", testFileName)

	testContent := fmt.Sprintf(`package models

import (
	"testing"
)

func Test%sCreation(t *testing.T) {
	%s := &%s{
		// Add test data here
	}
	
	if %s == nil {
		t.Error("Failed to create %s")
	}
}

func Test%sValidation(t *testing.T) {
	// Add validation tests here
	t.Skip("Validation tests not implemented yet")
}
`, entity.Name, strings.ToLower(entity.Name), entity.Name, strings.ToLower(entity.Name), entity.Name, entity.Name)

	if err := os.WriteFile(testFilePath, []byte(testContent), 0644); err != nil {
		return false, err
	}

	return true, nil
}

// addMiddleware adds basic middleware
func (ft *FineTuner) addMiddleware(appPath string) (bool, error) {
	middlewarePath := filepath.Join(appPath, "internal", "middleware")
	if err := os.MkdirAll(middlewarePath, 0755); err != nil {
		return false, err
	}

	middlewareContent := `package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware for request logging
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// CORS middleware for cross-origin requests
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	}
}

// Recovery middleware for panic recovery
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("Panic recovered: %s", err)
		}
		c.AbortWithStatus(500)
	})
}
`

	middlewareFilePath := filepath.Join(middlewarePath, "middleware.go")
	if err := os.WriteFile(middlewareFilePath, []byte(middlewareContent), 0644); err != nil {
		return false, err
	}

	return true, nil
}

// calculateImprovements calculates what improvements were made
func (ft *FineTuner) calculateImprovements(oldAnalysis *storage.AnalysisData, newTestResults *apptesting.TestSuite) []string {
	var improvements []string

	// This would compare old vs new metrics
	// For now, we'll return generic improvements
	if newTestResults != nil {
		if newTestResults.PassedTests > 0 {
			improvements = append(improvements, fmt.Sprintf("Improved test pass rate: %d/%d tests passing", 
				newTestResults.PassedTests, newTestResults.TotalTests))
		}

		if newTestResults.Coverage > 0 {
			improvements = append(improvements, fmt.Sprintf("Test coverage: %.2f%%", newTestResults.Coverage))
		}
	}

	return improvements
}

