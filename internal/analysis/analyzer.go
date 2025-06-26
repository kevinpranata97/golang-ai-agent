package analysis

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/kevinpranata97/golang-ai-agent/internal/apptesting"
	"github.com/kevinpranata97/golang-ai-agent/internal/requirements"
	"github.com/kevinpranata97/golang-ai-agent/internal/storage"
)

// CodeAnalyzer handles code analysis and improvement suggestions
type CodeAnalyzer struct {
	storage storage.Storage
}

// NewCodeAnalyzer creates a new code analyzer
func NewCodeAnalyzer(storage storage.Storage) *CodeAnalyzer {
	return &CodeAnalyzer{
		storage: storage,
	}
}

// AnalyzeProject performs comprehensive analysis of a generated project
func (ca *CodeAnalyzer) AnalyzeProject(projectID, appPath string, appReq *requirements.ApplicationRequirement, testResults *apptesting.TestSuite) (*storage.AnalysisData, error) {
	analysis := &storage.AnalysisData{
		ProjectID: projectID,
		Timestamp: time.Now(),
	}

	// Analyze code quality
	codeQuality, err := ca.analyzeCodeQuality(appPath)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze code quality: %v", err)
	}
	analysis.CodeQuality = *codeQuality

	// Analyze performance
	performance, err := ca.analyzePerformance(appPath, testResults)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze performance: %v", err)
	}
	analysis.Performance = *performance

	// Analyze security
	security, err := ca.analyzeSecurity(appPath)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze security: %v", err)
	}
	analysis.Security = *security

	// Generate improvement suggestions
	suggestions := ca.generateImprovementSuggestions(analysis, appReq, testResults)
	analysis.Suggestions = suggestions

	// Save analysis to storage
	if err := ca.storage.SaveAnalysis(analysis); err != nil {
		return nil, fmt.Errorf("failed to save analysis: %v", err)
	}

	return analysis, nil
}

// analyzeCodeQuality analyzes code quality metrics
func (ca *CodeAnalyzer) analyzeCodeQuality(appPath string) (*storage.CodeQualityMetrics, error) {
	metrics := &storage.CodeQualityMetrics{}

	// Count lines of code
	loc, err := ca.countLinesOfCode(appPath)
	if err != nil {
		return nil, err
	}
	metrics.LinesOfCode = loc

	// Calculate cyclomatic complexity
	complexity, err := ca.calculateCyclomaticComplexity(appPath)
	if err != nil {
		return nil, err
	}
	metrics.CyclomaticComplexity = complexity

	// Extract test coverage from test results (if available)
	// This would typically be extracted from go test -cover output
	metrics.TestCoverage = 0.0 // Placeholder

	// Calculate code duplication
	duplication, err := ca.calculateDuplication(appPath)
	if err != nil {
		return nil, err
	}
	metrics.DuplicationRatio = duplication

	// Assess technical debt and maintainability
	metrics.TechnicalDebt = ca.assessTechnicalDebt(metrics)
	metrics.Maintainability = ca.assessMaintainability(metrics)

	return metrics, nil
}

// analyzePerformance analyzes performance metrics
func (ca *CodeAnalyzer) analyzePerformance(appPath string, testResults *apptesting.TestSuite) (*storage.PerformanceMetrics, error) {
	metrics := &storage.PerformanceMetrics{}

	// Get binary size
	binaryPath := filepath.Join(appPath, filepath.Base(appPath))
	if info, err := os.Stat(binaryPath); err == nil {
		metrics.BinarySize = info.Size()
	}

	// Extract build time from test results
	if testResults != nil {
		for _, result := range testResults.Results {
			if result.Type == "build" {
				metrics.BuildTime = result.Duration.Seconds()
				break
			}
		}
	}

	// Estimate startup time (placeholder - would need actual measurement)
	metrics.StartupTime = 0.5 // seconds

	// Estimate memory usage (placeholder - would need actual profiling)
	metrics.MemoryUsage = 10 * 1024 * 1024 // 10MB

	// Extract response time from API tests
	if testResults != nil {
		for _, result := range testResults.Results {
			if result.Type == "api" && result.Details != nil {
				// Extract average response time from API test details
				metrics.ResponseTime = 0.1 // seconds (placeholder)
				break
			}
		}
	}

	return metrics, nil
}

// analyzeSecurity analyzes security metrics
func (ca *CodeAnalyzer) analyzeSecurity(appPath string) (*storage.SecurityMetrics, error) {
	metrics := &storage.SecurityMetrics{}

	var vulnerabilities []string
	var hardcodedSecrets int

	// Scan for security issues
	err := filepath.Walk(appPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			contentStr := string(content)

			// Check for SQL injection vulnerabilities
			if ca.hasSQLInjectionRisk(contentStr) {
				vulnerabilities = append(vulnerabilities, fmt.Sprintf("Potential SQL injection in %s", path))
			}

			// Check for hardcoded secrets
			if ca.hasHardcodedSecrets(contentStr) {
				hardcodedSecrets++
				vulnerabilities = append(vulnerabilities, fmt.Sprintf("Hardcoded secret in %s", path))
			}

			// Check for insecure HTTP usage
			if ca.hasInsecureHTTP(contentStr) {
				vulnerabilities = append(vulnerabilities, fmt.Sprintf("Insecure HTTP usage in %s", path))
			}

			// Check for weak cryptography
			if ca.hasWeakCryptography(contentStr) {
				vulnerabilities = append(vulnerabilities, fmt.Sprintf("Weak cryptography in %s", path))
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	metrics.Vulnerabilities = len(vulnerabilities)
	metrics.HardcodedSecrets = hardcodedSecrets
	metrics.SecurityIssues = vulnerabilities

	// Calculate security score (0-100)
	maxScore := 100.0
	deductions := float64(len(vulnerabilities)) * 10.0
	metrics.SecurityScore = maxScore - deductions
	if metrics.SecurityScore < 0 {
		metrics.SecurityScore = 0
	}

	return metrics, nil
}

// generateImprovementSuggestions generates suggestions for improving the code
func (ca *CodeAnalyzer) generateImprovementSuggestions(analysis *storage.AnalysisData, appReq *requirements.ApplicationRequirement, testResults *apptesting.TestSuite) []storage.ImprovementSuggestion {
	var suggestions []storage.ImprovementSuggestion

	// Code quality suggestions
	if analysis.CodeQuality.CyclomaticComplexity > 10 {
		suggestions = append(suggestions, storage.ImprovementSuggestion{
			Type:        "quality",
			Priority:    "medium",
			Description: "High cyclomatic complexity detected. Consider breaking down complex functions into smaller, more manageable pieces.",
			Impact:      "Improved code readability and maintainability",
			Effort:      "medium",
		})
	}

	if analysis.CodeQuality.TestCoverage < 80 {
		suggestions = append(suggestions, storage.ImprovementSuggestion{
			Type:        "quality",
			Priority:    "high",
			Description: "Low test coverage. Add more unit tests to improve code reliability.",
			Impact:      "Better bug detection and code confidence",
			Effort:      "high",
		})
	}

	if analysis.CodeQuality.DuplicationRatio > 0.1 {
		suggestions = append(suggestions, storage.ImprovementSuggestion{
			Type:        "quality",
			Priority:    "medium",
			Description: "Code duplication detected. Extract common functionality into shared functions or modules.",
			Impact:      "Reduced maintenance burden and improved consistency",
			Effort:      "medium",
		})
	}

	// Performance suggestions
	if analysis.Performance.BinarySize > 50*1024*1024 { // 50MB
		suggestions = append(suggestions, storage.ImprovementSuggestion{
			Type:        "performance",
			Priority:    "low",
			Description: "Large binary size. Consider optimizing dependencies and build flags.",
			Impact:      "Faster deployment and reduced resource usage",
			Effort:      "low",
			Code:        "go build -ldflags='-s -w' .",
		})
	}

	if analysis.Performance.BuildTime > 60 { // 60 seconds
		suggestions = append(suggestions, storage.ImprovementSuggestion{
			Type:        "performance",
			Priority:    "medium",
			Description: "Slow build time. Consider optimizing imports and using build cache.",
			Impact:      "Faster development cycle",
			Effort:      "medium",
		})
	}

	// Security suggestions
	if analysis.Security.Vulnerabilities > 0 {
		suggestions = append(suggestions, storage.ImprovementSuggestion{
			Type:        "security",
			Priority:    "high",
			Description: fmt.Sprintf("Found %d security vulnerabilities. Review and fix security issues.", analysis.Security.Vulnerabilities),
			Impact:      "Improved application security",
			Effort:      "high",
		})
	}

	if analysis.Security.HardcodedSecrets > 0 {
		suggestions = append(suggestions, storage.ImprovementSuggestion{
			Type:        "security",
			Priority:    "high",
			Description: "Hardcoded secrets detected. Move sensitive data to environment variables or secure configuration.",
			Impact:      "Improved security and configuration management",
			Effort:      "medium",
			Code:        "Use os.Getenv() for sensitive configuration",
		})
	}

	// Functionality suggestions based on test results
	if testResults != nil {
		for _, result := range testResults.Results {
			if result.Status == "fail" {
				suggestions = append(suggestions, storage.ImprovementSuggestion{
					Type:        "functionality",
					Priority:    "high",
					Description: fmt.Sprintf("Test failure in %s: %s", result.Name, result.Error),
					Impact:      "Fixed functionality and improved reliability",
					Effort:      "medium",
				})
			}
		}
	}

	// Framework-specific suggestions
	if appReq != nil {
		switch appReq.Framework {
		case "gin":
			suggestions = append(suggestions, storage.ImprovementSuggestion{
				Type:        "functionality",
				Priority:    "low",
				Description: "Consider adding middleware for logging, rate limiting, and request validation.",
				Impact:      "Better observability and security",
				Effort:      "low",
			})
		}

		// Database-specific suggestions
		switch appReq.Database {
		case "sqlite":
			suggestions = append(suggestions, storage.ImprovementSuggestion{
				Type:        "performance",
				Priority:    "medium",
				Description: "Consider migrating to PostgreSQL or MySQL for production use.",
				Impact:      "Better performance and scalability",
				Effort:      "high",
			})
		}
	}

	return suggestions
}

// Helper methods for analysis

// countLinesOfCode counts non-empty, non-comment lines of code
func (ca *CodeAnalyzer) countLinesOfCode(appPath string) (int, error) {
	totalLines := 0

	err := filepath.Walk(appPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".go") {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" && !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "/*") {
					totalLines++
				}
			}
		}

		return nil
	})

	return totalLines, err
}

// calculateCyclomaticComplexity calculates cyclomatic complexity
func (ca *CodeAnalyzer) calculateCyclomaticComplexity(appPath string) (int, error) {
	totalComplexity := 0

	err := filepath.Walk(appPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".go") {
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err != nil {
				return err
			}

			ast.Inspect(node, func(n ast.Node) bool {
				switch n.(type) {
				case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt:
					totalComplexity++
				case *ast.CaseClause:
					totalComplexity++
				}
				return true
			})
		}

		return nil
	})

	return totalComplexity, err
}

// calculateDuplication calculates code duplication ratio
func (ca *CodeAnalyzer) calculateDuplication(appPath string) (float64, error) {
	// This is a simplified implementation
	// In a real system, you'd use more sophisticated algorithms
	
	var allLines []string
	lineCount := make(map[string]int)

	err := filepath.Walk(appPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".go") {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" && !strings.HasPrefix(line, "//") && len(line) > 10 {
					allLines = append(allLines, line)
					lineCount[line]++
				}
			}
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	duplicatedLines := 0
	for _, count := range lineCount {
		if count > 1 {
			duplicatedLines += count - 1
		}
	}

	if len(allLines) == 0 {
		return 0, nil
	}

	return float64(duplicatedLines) / float64(len(allLines)), nil
}

// assessTechnicalDebt assesses technical debt level
func (ca *CodeAnalyzer) assessTechnicalDebt(metrics *storage.CodeQualityMetrics) string {
	score := 0

	if metrics.CyclomaticComplexity > 20 {
		score += 3
	} else if metrics.CyclomaticComplexity > 10 {
		score += 2
	} else if metrics.CyclomaticComplexity > 5 {
		score += 1
	}

	if metrics.TestCoverage < 50 {
		score += 3
	} else if metrics.TestCoverage < 70 {
		score += 2
	} else if metrics.TestCoverage < 80 {
		score += 1
	}

	if metrics.DuplicationRatio > 0.2 {
		score += 3
	} else if metrics.DuplicationRatio > 0.1 {
		score += 2
	} else if metrics.DuplicationRatio > 0.05 {
		score += 1
	}

	switch {
	case score >= 7:
		return "high"
	case score >= 4:
		return "medium"
	case score >= 2:
		return "low"
	default:
		return "minimal"
	}
}

// assessMaintainability assesses code maintainability
func (ca *CodeAnalyzer) assessMaintainability(metrics *storage.CodeQualityMetrics) string {
	score := 100

	// Deduct points for complexity
	score -= metrics.CyclomaticComplexity * 2

	// Deduct points for low test coverage
	if metrics.TestCoverage < 80 {
		score -= int((80 - metrics.TestCoverage) / 2)
	}

	// Deduct points for duplication
	score -= int(metrics.DuplicationRatio * 100)

	switch {
	case score >= 80:
		return "excellent"
	case score >= 60:
		return "good"
	case score >= 40:
		return "fair"
	case score >= 20:
		return "poor"
	default:
		return "very poor"
	}
}

// Security analysis helper methods

func (ca *CodeAnalyzer) hasSQLInjectionRisk(content string) bool {
	// Look for string concatenation in SQL queries
	patterns := []string{
		`db\.Exec\([^)]*\+`,
		`db\.Query\([^)]*\+`,
		`fmt\.Sprintf.*SELECT`,
		`fmt\.Sprintf.*INSERT`,
		`fmt\.Sprintf.*UPDATE`,
		`fmt\.Sprintf.*DELETE`,
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, content); matched {
			return true
		}
	}

	return false
}

func (ca *CodeAnalyzer) hasHardcodedSecrets(content string) bool {
	patterns := []string{
		`(?i)password\s*[:=]\s*["'][^"']{8,}["']`,
		`(?i)api[_-]?key\s*[:=]\s*["'][^"']{10,}["']`,
		`(?i)secret[_-]?key\s*[:=]\s*["'][^"']{10,}["']`,
		`(?i)token\s*[:=]\s*["'][^"']{10,}["']`,
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, content); matched {
			return true
		}
	}

	return false
}

func (ca *CodeAnalyzer) hasInsecureHTTP(content string) bool {
	patterns := []string{
		`http://`,
		`InsecureSkipVerify:\s*true`,
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, content); matched {
			return true
		}
	}

	return false
}

func (ca *CodeAnalyzer) hasWeakCryptography(content string) bool {
	patterns := []string{
		`md5\.`,
		`sha1\.`,
		`des\.`,
		`rc4\.`,
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, content); matched {
			return true
		}
	}

	return false
}

