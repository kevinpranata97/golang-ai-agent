package apptesting

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/kevinpranata97/golang-ai-agent/internal/requirements"
)

// TestResult represents the result of a test
type TestResult struct {
	Name        string        `json:"name"`
	Type        string        `json:"type"` // unit, integration, build, api
	Status      string        `json:"status"` // pass, fail, skip
	Duration    time.Duration `json:"duration"`
	Output      string        `json:"output"`
	Error       string        `json:"error,omitempty"`
	Coverage    float64       `json:"coverage,omitempty"`
	Details     interface{}   `json:"details,omitempty"`
}

// TestSuite represents a collection of test results
type TestSuite struct {
	Name         string       `json:"name"`
	AppPath      string       `json:"app_path"`
	StartTime    time.Time    `json:"start_time"`
	EndTime      time.Time    `json:"end_time"`
	Duration     time.Duration `json:"duration"`
	TotalTests   int          `json:"total_tests"`
	PassedTests  int          `json:"passed_tests"`
	FailedTests  int          `json:"failed_tests"`
	SkippedTests int          `json:"skipped_tests"`
	Coverage     float64      `json:"coverage"`
	Results      []TestResult `json:"results"`
	Summary      string       `json:"summary"`
}

// ApplicationTester handles testing of generated applications
type ApplicationTester struct {
	workingDir string
	timeout    time.Duration
}

// NewApplicationTester creates a new application tester
func NewApplicationTester(workingDir string) *ApplicationTester {
	return &ApplicationTester{
		workingDir: workingDir,
		timeout:    5 * time.Minute,
	}
}

// TestApplication runs comprehensive tests on a generated application
func (at *ApplicationTester) TestApplication(appPath string, appReq *requirements.ApplicationRequirement) (*TestSuite, error) {
	suite := &TestSuite{
		Name:      appReq.Name,
		AppPath:   appPath,
		StartTime: time.Now(),
		Results:   []TestResult{},
	}

	// Test 1: Build Test
	buildResult := at.testBuild(appPath, appReq)
	suite.Results = append(suite.Results, buildResult)

	// Test 2: Static Analysis
	staticResult := at.testStaticAnalysis(appPath, appReq)
	suite.Results = append(suite.Results, staticResult)

	// Test 3: Unit Tests (if any exist)
	unitResult := at.testUnit(appPath, appReq)
	suite.Results = append(suite.Results, unitResult)

	// Test 4: API Tests (if it's an API application)
	if appReq.Type == "api" || appReq.Type == "web" {
		apiResult := at.testAPI(appPath, appReq)
		suite.Results = append(suite.Results, apiResult)
	}

	// Test 5: Security Tests
	securityResult := at.testSecurity(appPath, appReq)
	suite.Results = append(suite.Results, securityResult)

	// Test 6: Performance Tests (basic)
	perfResult := at.testPerformance(appPath, appReq)
	suite.Results = append(suite.Results, perfResult)

	// Calculate summary
	suite.EndTime = time.Now()
	suite.Duration = suite.EndTime.Sub(suite.StartTime)
	suite.TotalTests = len(suite.Results)

	for _, result := range suite.Results {
		switch result.Status {
		case "pass":
			suite.PassedTests++
		case "fail":
			suite.FailedTests++
		case "skip":
			suite.SkippedTests++
		}
	}

	// Calculate overall coverage (average of all coverage results)
	var totalCoverage float64
	var coverageCount int
	for _, result := range suite.Results {
		if result.Coverage > 0 {
			totalCoverage += result.Coverage
			coverageCount++
		}
	}
	if coverageCount > 0 {
		suite.Coverage = totalCoverage / float64(coverageCount)
	}

	// Generate summary
	suite.Summary = at.generateSummary(suite)

	return suite, nil
}

// testBuild tests if the application builds successfully
func (at *ApplicationTester) testBuild(appPath string, appReq *requirements.ApplicationRequirement) TestResult {
	result := TestResult{
		Name: "Build Test",
		Type: "build",
	}

	startTime := time.Now()

	// Change to app directory and run go build
	cmd := exec.Command("go", "build", "-v", ".")
	cmd.Dir = appPath
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH")+":/usr/local/go/bin")

	output, err := cmd.CombinedOutput()
	result.Duration = time.Since(startTime)
	result.Output = string(output)

	if err != nil {
		result.Status = "fail"
		result.Error = err.Error()
	} else {
		result.Status = "pass"
	}

	return result
}

// testStaticAnalysis runs static analysis tools
func (at *ApplicationTester) testStaticAnalysis(appPath string, appReq *requirements.ApplicationRequirement) TestResult {
	result := TestResult{
		Name: "Static Analysis",
		Type: "static",
	}

	startTime := time.Now()

	var outputs []string
	var errors []string

	// Run go vet
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = appPath
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH")+":/usr/local/go/bin")

	output, err := cmd.CombinedOutput()
	outputs = append(outputs, "=== go vet ===\n"+string(output))
	if err != nil {
		errors = append(errors, "go vet: "+err.Error())
	}

	// Run go fmt check
	cmd = exec.Command("go", "fmt", "./...")
	cmd.Dir = appPath
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH")+":/usr/local/go/bin")

	output, err = cmd.CombinedOutput()
	outputs = append(outputs, "=== go fmt ===\n"+string(output))
	if err != nil {
		errors = append(errors, "go fmt: "+err.Error())
	}

	result.Duration = time.Since(startTime)
	result.Output = strings.Join(outputs, "\n\n")

	if len(errors) > 0 {
		result.Status = "fail"
		result.Error = strings.Join(errors, "; ")
	} else {
		result.Status = "pass"
	}

	return result
}

// testUnit runs unit tests
func (at *ApplicationTester) testUnit(appPath string, appReq *requirements.ApplicationRequirement) TestResult {
	result := TestResult{
		Name: "Unit Tests",
		Type: "unit",
	}

	startTime := time.Now()

	// Check if there are any test files
	hasTests, err := at.hasTestFiles(appPath)
	if err != nil {
		result.Status = "fail"
		result.Error = err.Error()
		result.Duration = time.Since(startTime)
		return result
	}

	if !hasTests {
		result.Status = "skip"
		result.Output = "No test files found"
		result.Duration = time.Since(startTime)
		return result
	}

	// Run go test with coverage
	cmd := exec.Command("go", "test", "-v", "-cover", "./...")
	cmd.Dir = appPath
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH")+":/usr/local/go/bin")

	output, err := cmd.CombinedOutput()
	result.Duration = time.Since(startTime)
	result.Output = string(output)

	if err != nil {
		result.Status = "fail"
		result.Error = err.Error()
	} else {
		result.Status = "pass"
		// Extract coverage from output
		result.Coverage = at.extractCoverage(string(output))
	}

	return result
}

// testAPI tests API endpoints
func (at *ApplicationTester) testAPI(appPath string, appReq *requirements.ApplicationRequirement) TestResult {
	result := TestResult{
		Name: "API Tests",
		Type: "api",
	}

	startTime := time.Now()

	// Start the application
	cmd := exec.Command("./"+filepath.Base(appPath))
	cmd.Dir = appPath
	cmd.Env = append(os.Environ(), "PORT=8081") // Use different port for testing

	// Start the server
	err := cmd.Start()
	if err != nil {
		result.Status = "fail"
		result.Error = "Failed to start application: " + err.Error()
		result.Duration = time.Since(startTime)
		return result
	}

	// Ensure we kill the process when done
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	// Wait a moment for server to start
	time.Sleep(2 * time.Second)

	var testResults []map[string]interface{}
	var errors []string

	// Test health endpoint
	healthResult := at.testEndpoint("GET", "http://localhost:8081/health", nil)
	testResults = append(testResults, map[string]interface{}{
		"endpoint": "/health",
		"method":   "GET",
		"result":   healthResult,
	})

	// Test each API endpoint
	for _, endpoint := range appReq.Endpoints {
		url := "http://localhost:8081" + endpoint.Path
		
		// Replace path parameters with test values
		url = strings.ReplaceAll(url, "{id}", "1")
		
		var body []byte
		if endpoint.Method == "POST" || endpoint.Method == "PUT" {
			// Create test data based on the first entity
			if len(appReq.Entities) > 0 {
				testData := at.generateTestData(appReq.Entities[0])
				body, _ = json.Marshal(testData)
			}
		}

		endpointResult := at.testEndpoint(endpoint.Method, url, body)
		testResults = append(testResults, map[string]interface{}{
			"endpoint": endpoint.Path,
			"method":   endpoint.Method,
			"result":   endpointResult,
		})

		if !endpointResult["success"].(bool) {
			errors = append(errors, fmt.Sprintf("%s %s: %s", endpoint.Method, endpoint.Path, endpointResult["error"]))
		}
	}

	result.Duration = time.Since(startTime)
	result.Details = testResults

	if len(errors) > 0 {
		result.Status = "fail"
		result.Error = strings.Join(errors, "; ")
	} else {
		result.Status = "pass"
	}

	result.Output = fmt.Sprintf("Tested %d endpoints", len(testResults))

	return result
}

// testSecurity runs basic security tests
func (at *ApplicationTester) testSecurity(appPath string, appReq *requirements.ApplicationRequirement) TestResult {
	result := TestResult{
		Name: "Security Tests",
		Type: "security",
	}

	startTime := time.Now()

	var issues []string
	var outputs []string

	// Check for common security issues in code
	securityIssues := at.scanForSecurityIssues(appPath)
	if len(securityIssues) > 0 {
		issues = append(issues, securityIssues...)
	}

	// Check for hardcoded secrets
	secrets := at.scanForHardcodedSecrets(appPath)
	if len(secrets) > 0 {
		issues = append(issues, secrets...)
	}

	outputs = append(outputs, fmt.Sprintf("Security scan completed. Found %d potential issues.", len(issues)))

	result.Duration = time.Since(startTime)
	result.Output = strings.Join(outputs, "\n")

	if len(issues) > 0 {
		result.Status = "fail"
		result.Error = strings.Join(issues, "; ")
		result.Details = map[string]interface{}{
			"issues": issues,
		}
	} else {
		result.Status = "pass"
	}

	return result
}

// testPerformance runs basic performance tests
func (at *ApplicationTester) testPerformance(appPath string, appReq *requirements.ApplicationRequirement) TestResult {
	result := TestResult{
		Name: "Performance Tests",
		Type: "performance",
	}

	startTime := time.Now()

	// For now, just check binary size and basic metrics
	var metrics []string

	// Check binary size
	binaryPath := filepath.Join(appPath, filepath.Base(appPath))
	if info, err := os.Stat(binaryPath); err == nil {
		size := info.Size()
		metrics = append(metrics, fmt.Sprintf("Binary size: %d bytes (%.2f MB)", size, float64(size)/1024/1024))
	}

	// Count lines of code
	loc, err := at.countLinesOfCode(appPath)
	if err == nil {
		metrics = append(metrics, fmt.Sprintf("Lines of code: %d", loc))
	}

	result.Duration = time.Since(startTime)
	result.Output = strings.Join(metrics, "\n")
	result.Status = "pass"

	result.Details = map[string]interface{}{
		"metrics": metrics,
	}

	return result
}

// Helper methods

// hasTestFiles checks if there are any test files in the project
func (at *ApplicationTester) hasTestFiles(appPath string) (bool, error) {
	return filepath.Walk(appPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(info.Name(), "_test.go") {
			return fmt.Errorf("found test file") // Use error to break out of walk
		}
		return nil
	}) != nil, nil
}

// extractCoverage extracts coverage percentage from go test output
func (at *ApplicationTester) extractCoverage(output string) float64 {
	re := regexp.MustCompile(`coverage: ([\d.]+)% of statements`)
	matches := re.FindStringSubmatch(output)
	if len(matches) > 1 {
		var coverage float64
		fmt.Sscanf(matches[1], "%f", &coverage)
		return coverage
	}
	return 0
}

// testEndpoint tests a single API endpoint
func (at *ApplicationTester) testEndpoint(method, url string, body []byte) map[string]interface{} {
	client := &http.Client{Timeout: 10 * time.Second}
	
	var req *http.Request
	var err error
	
	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	return map[string]interface{}{
		"success":     resp.StatusCode < 400,
		"status_code": resp.StatusCode,
		"response":    string(responseBody),
	}
}

// generateTestData generates test data for an entity
func (at *ApplicationTester) generateTestData(entity requirements.Entity) map[string]interface{} {
	data := make(map[string]interface{})
	
	for _, field := range entity.Fields {
		if field.Name == "id" || field.Name == "created_at" {
			continue // Skip auto-generated fields
		}
		
		switch field.Type {
		case "string":
			data[field.Name] = "test_" + field.Name
		case "email":
			data[field.Name] = "test@example.com"
		case "int":
			data[field.Name] = 1
		case "float":
			data[field.Name] = 1.0
		case "bool":
			data[field.Name] = true
		default:
			data[field.Name] = "test_value"
		}
	}
	
	return data
}

// scanForSecurityIssues scans code for common security issues
func (at *ApplicationTester) scanForSecurityIssues(appPath string) []string {
	var issues []string
	
	// This is a basic implementation - in a real system, you'd use tools like gosec
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
			if strings.Contains(contentStr, "db.Exec(") && strings.Contains(contentStr, "+") {
				issues = append(issues, fmt.Sprintf("Potential SQL injection in %s", path))
			}
			
			// Check for hardcoded passwords
			if regexp.MustCompile(`password\s*[:=]\s*["'][^"']+["']`).MatchString(strings.ToLower(contentStr)) {
				issues = append(issues, fmt.Sprintf("Potential hardcoded password in %s", path))
			}
		}
		
		return nil
	})
	
	if err != nil {
		issues = append(issues, "Error scanning for security issues: "+err.Error())
	}
	
	return issues
}

// scanForHardcodedSecrets scans for hardcoded secrets
func (at *ApplicationTester) scanForHardcodedSecrets(appPath string) []string {
	var secrets []string
	
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)api[_-]?key\s*[:=]\s*["'][^"']{10,}["']`),
		regexp.MustCompile(`(?i)secret[_-]?key\s*[:=]\s*["'][^"']{10,}["']`),
		regexp.MustCompile(`(?i)token\s*[:=]\s*["'][^"']{10,}["']`),
		regexp.MustCompile(`(?i)password\s*[:=]\s*["'][^"']{8,}["']`),
	}
	
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
			
			for _, pattern := range patterns {
				if pattern.MatchString(contentStr) {
					secrets = append(secrets, fmt.Sprintf("Potential hardcoded secret in %s", path))
				}
			}
		}
		
		return nil
	})
	
	if err != nil {
		secrets = append(secrets, "Error scanning for secrets: "+err.Error())
	}
	
	return secrets
}

// countLinesOfCode counts lines of code in the project
func (at *ApplicationTester) countLinesOfCode(appPath string) (int, error) {
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
				if line != "" && !strings.HasPrefix(line, "//") {
					totalLines++
				}
			}
		}
		
		return nil
	})
	
	return totalLines, err
}

// generateSummary generates a summary of the test results
func (at *ApplicationTester) generateSummary(suite *TestSuite) string {
	var summary strings.Builder
	
	summary.WriteString(fmt.Sprintf("Test Summary for %s:\n", suite.Name))
	summary.WriteString(fmt.Sprintf("Duration: %v\n", suite.Duration))
	summary.WriteString(fmt.Sprintf("Total Tests: %d\n", suite.TotalTests))
	summary.WriteString(fmt.Sprintf("Passed: %d\n", suite.PassedTests))
	summary.WriteString(fmt.Sprintf("Failed: %d\n", suite.FailedTests))
	summary.WriteString(fmt.Sprintf("Skipped: %d\n", suite.SkippedTests))
	
	if suite.Coverage > 0 {
		summary.WriteString(fmt.Sprintf("Coverage: %.2f%%\n", suite.Coverage))
	}
	
	if suite.FailedTests > 0 {
		summary.WriteString("\nFailed Tests:\n")
		for _, result := range suite.Results {
			if result.Status == "fail" {
				summary.WriteString(fmt.Sprintf("- %s: %s\n", result.Name, result.Error))
			}
		}
	}
	
	// Overall status
	if suite.FailedTests == 0 {
		summary.WriteString("\n✅ All tests passed!")
	} else {
		summary.WriteString(fmt.Sprintf("\n❌ %d test(s) failed", suite.FailedTests))
	}
	
	return summary.String()
}

// SaveTestResults saves test results to a file
func (at *ApplicationTester) SaveTestResults(suite *TestSuite, outputPath string) error {
	data, err := json.MarshalIndent(suite, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(outputPath, data, 0644)
}

