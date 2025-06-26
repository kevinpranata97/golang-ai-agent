package testing

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type TestRunner struct {
	httpClient *http.Client
}

type TestResult struct {
	Success      bool                   `json:"success"`
	TestsPassed  int                    `json:"tests_passed"`
	TestsFailed  int                    `json:"tests_failed"`
	Coverage     float64                `json:"coverage"`
	Duration     time.Duration          `json:"duration"`
	Details      []TestDetail           `json:"details"`
	Analysis     CodeAnalysis           `json:"analysis"`
	Metrics      PerformanceMetrics     `json:"metrics"`
	SecurityScan SecurityScanResult     `json:"security_scan"`
}

type TestDetail struct {
	Name     string        `json:"name"`
	Status   string        `json:"status"`
	Duration time.Duration `json:"duration"`
	Error    string        `json:"error,omitempty"`
	Output   string        `json:"output,omitempty"`
}

type CodeAnalysis struct {
	LinesOfCode      int                    `json:"lines_of_code"`
	Functions        int                    `json:"functions"`
	Complexity       int                    `json:"complexity"`
	Duplicates       int                    `json:"duplicates"`
	Issues           []CodeIssue            `json:"issues"`
	Dependencies     []Dependency           `json:"dependencies"`
	Architecture     ArchitectureAnalysis   `json:"architecture"`
	QualityScore     float64                `json:"quality_score"`
}

type CodeIssue struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion,omitempty"`
}

type Dependency struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Type    string `json:"type"`
	License string `json:"license,omitempty"`
}

type ArchitectureAnalysis struct {
	Packages     []string `json:"packages"`
	Interfaces   int      `json:"interfaces"`
	Structs      int      `json:"structs"`
	Methods      int      `json:"methods"`
	Coupling     float64  `json:"coupling"`
	Cohesion     float64  `json:"cohesion"`
}

type PerformanceMetrics struct {
	ResponseTime    time.Duration `json:"response_time"`
	Throughput      float64       `json:"throughput"`
	MemoryUsage     int64         `json:"memory_usage"`
	CPUUsage        float64       `json:"cpu_usage"`
	LoadTestResult  LoadTestResult `json:"load_test"`
}

type LoadTestResult struct {
	TotalRequests   int           `json:"total_requests"`
	SuccessfulReqs  int           `json:"successful_requests"`
	FailedRequests  int           `json:"failed_requests"`
	AverageResponse time.Duration `json:"average_response"`
	MaxResponse     time.Duration `json:"max_response"`
	MinResponse     time.Duration `json:"min_response"`
}

type SecurityScanResult struct {
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	Score           float64         `json:"score"`
	Recommendations []string        `json:"recommendations"`
}

type Vulnerability struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Description string `json:"description"`
	Fix         string `json:"fix,omitempty"`
}

func NewTestRunner() *TestRunner {
	return &TestRunner{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (tr *TestRunner) RunTests(projectPath string) TestResult {
	startTime := time.Now()
	
	result := TestResult{
		Success: true,
		Details: []TestDetail{},
	}
	
	// Run unit tests
	unitTestResult := tr.runUnitTests(projectPath)
	result.Details = append(result.Details, unitTestResult...)
	
	// Run integration tests
	integrationTestResult := tr.runIntegrationTests(projectPath)
	result.Details = append(result.Details, integrationTestResult...)
	
	// Analyze code
	result.Analysis = tr.analyzeCode(projectPath)
	
	// Run performance tests
	result.Metrics = tr.runPerformanceTests(projectPath)
	
	// Run security scan
	result.SecurityScan = tr.runSecurityScan(projectPath)
	
	// Calculate overall results
	for _, detail := range result.Details {
		if detail.Status == "PASS" {
			result.TestsPassed++
		} else {
			result.TestsFailed++
			result.Success = false
		}
	}
	
	result.Duration = time.Since(startTime)
	
	return result
}

func (tr *TestRunner) runUnitTests(projectPath string) []TestDetail {
	var details []TestDetail
	
	// Check if it's a Go project
	if tr.fileExists(filepath.Join(projectPath, "go.mod")) {
		details = append(details, tr.runGoTests(projectPath)...)
	}
	
	// Check if it's a Node.js project
	if tr.fileExists(filepath.Join(projectPath, "package.json")) {
		details = append(details, tr.runNodeTests(projectPath)...)
	}
	
	// Check if it's a Python project
	if tr.fileExists(filepath.Join(projectPath, "requirements.txt")) || tr.fileExists(filepath.Join(projectPath, "setup.py")) {
		details = append(details, tr.runPythonTests(projectPath)...)
	}
	
	return details
}

func (tr *TestRunner) runGoTests(projectPath string) []TestDetail {
	var details []TestDetail
	
	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Dir = projectPath
	
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	
	if err != nil {
		details = append(details, TestDetail{
			Name:   "Go Unit Tests",
			Status: "FAIL",
			Error:  err.Error(),
			Output: outputStr,
		})
		return details
	}
	
	// Parse go test output
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, "RUN") {
			testName := strings.TrimSpace(strings.Split(line, "RUN")[1])
			details = append(details, TestDetail{
				Name:   testName,
				Status: "PASS",
				Output: line,
			})
		} else if strings.Contains(line, "FAIL") {
			testName := strings.TrimSpace(strings.Split(line, "FAIL")[1])
			details = append(details, TestDetail{
				Name:   testName,
				Status: "FAIL",
				Output: line,
			})
		}
	}
	
	return details
}

func (tr *TestRunner) runNodeTests(projectPath string) []TestDetail {
	var details []TestDetail
	
	cmd := exec.Command("npm", "test")
	cmd.Dir = projectPath
	
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	
	status := "PASS"
	if err != nil {
		status = "FAIL"
	}
	
	details = append(details, TestDetail{
		Name:   "Node.js Tests",
		Status: status,
		Output: outputStr,
		Error:  func() string { if err != nil { return err.Error() }; return "" }(),
	})
	
	return details
}

func (tr *TestRunner) runPythonTests(projectPath string) []TestDetail {
	var details []TestDetail
	
	// Try pytest first, then unittest
	cmd := exec.Command("python", "-m", "pytest", "-v")
	cmd.Dir = projectPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Fallback to unittest
		cmd = exec.Command("python", "-m", "unittest", "discover", "-v")
		cmd.Dir = projectPath
		output, err = cmd.CombinedOutput()
	}
	
	outputStr := string(output)
	status := "PASS"
	if err != nil {
		status = "FAIL"
	}
	
	details = append(details, TestDetail{
		Name:   "Python Tests",
		Status: status,
		Output: outputStr,
		Error:  func() string { if err != nil { return err.Error() }; return "" }(),
	})
	
	return details
}

func (tr *TestRunner) runIntegrationTests(projectPath string) []TestDetail {
	var details []TestDetail
	
	// Look for integration test files
	integrationTestFiles := tr.findIntegrationTests(projectPath)
	
	for _, testFile := range integrationTestFiles {
		detail := TestDetail{
			Name:   fmt.Sprintf("Integration Test: %s", filepath.Base(testFile)),
			Status: "PASS",
		}
		
		// Run the integration test
		if strings.HasSuffix(testFile, "_test.go") {
			cmd := exec.Command("go", "test", "-v", testFile)
			cmd.Dir = filepath.Dir(testFile)
			
			output, err := cmd.CombinedOutput()
			detail.Output = string(output)
			
			if err != nil {
				detail.Status = "FAIL"
				detail.Error = err.Error()
			}
		}
		
		details = append(details, detail)
	}
	
	return details
}

func (tr *TestRunner) findIntegrationTests(projectPath string) []string {
	var testFiles []string
	
	filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if strings.Contains(path, "integration") && strings.HasSuffix(path, "_test.go") {
			testFiles = append(testFiles, path)
		}
		
		return nil
	})
	
	return testFiles
}

func (tr *TestRunner) analyzeCode(projectPath string) CodeAnalysis {
	analysis := CodeAnalysis{
		Issues:       []CodeIssue{},
		Dependencies: []Dependency{},
	}
	
	// Analyze Go code
	if tr.fileExists(filepath.Join(projectPath, "go.mod")) {
		tr.analyzeGoCode(projectPath, &analysis)
	}
	
	// Count lines of code
	analysis.LinesOfCode = tr.countLinesOfCode(projectPath)
	
	// Calculate quality score
	analysis.QualityScore = tr.calculateQualityScore(analysis)
	
	return analysis
}

func (tr *TestRunner) analyzeGoCode(projectPath string, analysis *CodeAnalysis) {
	fset := token.NewFileSet()
	
	filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") {
			return err
		}
		
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		
		file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
		if err != nil {
			analysis.Issues = append(analysis.Issues, CodeIssue{
				Type:        "syntax",
				Severity:    "error",
				File:        path,
				Description: fmt.Sprintf("Parse error: %v", err),
			})
			return nil
		}
		
		// Count functions and analyze complexity
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				analysis.Functions++
				complexity := tr.calculateCyclomaticComplexity(x)
				analysis.Complexity += complexity
				
				if complexity > 10 {
					analysis.Issues = append(analysis.Issues, CodeIssue{
						Type:        "complexity",
						Severity:    "warning",
						File:        path,
						Line:        fset.Position(x.Pos()).Line,
						Description: fmt.Sprintf("Function %s has high cyclomatic complexity: %d", x.Name.Name, complexity),
						Suggestion:  "Consider breaking this function into smaller functions",
					})
				}
			case *ast.GenDecl:
				if x.Tok == token.TYPE {
					for _, spec := range x.Specs {
						if ts, ok := spec.(*ast.TypeSpec); ok {
							if _, ok := ts.Type.(*ast.StructType); ok {
								analysis.Architecture.Structs++
							} else if _, ok := ts.Type.(*ast.InterfaceType); ok {
								analysis.Architecture.Interfaces++
							}
						}
					}
				}
			}
			return true
		})
		
		return nil
	})
}

func (tr *TestRunner) calculateCyclomaticComplexity(fn *ast.FuncDecl) int {
	complexity := 1 // Base complexity
	
	ast.Inspect(fn, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.TypeSwitchStmt:
			complexity++
		case *ast.CaseClause:
			complexity++
		}
		return true
	})
	
	return complexity
}

func (tr *TestRunner) countLinesOfCode(projectPath string) int {
	var totalLines int
	
	filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if tr.isSourceFile(path) {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" && !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "#") {
					totalLines++
				}
			}
		}
		
		return nil
	})
	
	return totalLines
}

func (tr *TestRunner) isSourceFile(path string) bool {
	ext := filepath.Ext(path)
	sourceExts := []string{".go", ".js", ".ts", ".py", ".java", ".cpp", ".c", ".h"}
	
	for _, sourceExt := range sourceExts {
		if ext == sourceExt {
			return true
		}
	}
	
	return false
}

func (tr *TestRunner) calculateQualityScore(analysis CodeAnalysis) float64 {
	score := 100.0
	
	// Deduct points for issues
	for _, issue := range analysis.Issues {
		switch issue.Severity {
		case "error":
			score -= 10
		case "warning":
			score -= 5
		case "info":
			score -= 1
		}
	}
	
	// Deduct points for high complexity
	if analysis.Functions > 0 {
		avgComplexity := float64(analysis.Complexity) / float64(analysis.Functions)
		if avgComplexity > 5 {
			score -= (avgComplexity - 5) * 2
		}
	}
	
	if score < 0 {
		score = 0
	}
	
	return score
}

func (tr *TestRunner) runPerformanceTests(projectPath string) PerformanceMetrics {
	metrics := PerformanceMetrics{}
	
	// Simple performance test - this would be more sophisticated in a real implementation
	startTime := time.Now()
	
	// Simulate some performance testing
	time.Sleep(100 * time.Millisecond)
	
	metrics.ResponseTime = time.Since(startTime)
	metrics.Throughput = 100.0 // requests per second
	metrics.MemoryUsage = 1024 * 1024 // 1MB
	metrics.CPUUsage = 15.5 // 15.5%
	
	return metrics
}

func (tr *TestRunner) runSecurityScan(projectPath string) SecurityScanResult {
	result := SecurityScanResult{
		Vulnerabilities: []Vulnerability{},
		Score:           85.0,
		Recommendations: []string{
			"Update dependencies to latest versions",
			"Add input validation",
			"Implement proper error handling",
		},
	}
	
	// Scan for common security issues
	filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || !tr.isSourceFile(path) {
			return err
		}
		
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		
		contentStr := string(content)
		
		// Check for hardcoded secrets
		secretPatterns := []string{
			`password\s*=\s*["'][^"']+["']`,
			`api_key\s*=\s*["'][^"']+["']`,
			`secret\s*=\s*["'][^"']+["']`,
		}
		
		for _, pattern := range secretPatterns {
			re := regexp.MustCompile(pattern)
			if re.MatchString(contentStr) {
				result.Vulnerabilities = append(result.Vulnerabilities, Vulnerability{
					Type:        "hardcoded_secret",
					Severity:    "high",
					File:        path,
					Description: "Potential hardcoded secret found",
					Fix:         "Use environment variables or secure configuration",
				})
			}
		}
		
		return nil
	})
	
	return result
}

func (tr *TestRunner) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Web-specific testing methods
func (tr *TestRunner) TestWebEndpoint(url string) TestDetail {
	startTime := time.Now()
	
	resp, err := tr.httpClient.Get(url)
	if err != nil {
		return TestDetail{
			Name:     fmt.Sprintf("GET %s", url),
			Status:   "FAIL",
			Duration: time.Since(startTime),
			Error:    err.Error(),
		}
	}
	defer resp.Body.Close()
	
	status := "PASS"
	if resp.StatusCode >= 400 {
		status = "FAIL"
	}
	
	return TestDetail{
		Name:     fmt.Sprintf("GET %s", url),
		Status:   status,
		Duration: time.Since(startTime),
		Output:   fmt.Sprintf("Status: %d %s", resp.StatusCode, resp.Status),
	}
}

func (tr *TestRunner) RunLoadTest(url string, requests int, concurrency int) LoadTestResult {
	// Simple load test implementation
	result := LoadTestResult{
		TotalRequests: requests,
	}
	
	// This is a simplified implementation
	// In a real scenario, you'd use a proper load testing tool
	
	return result
}

