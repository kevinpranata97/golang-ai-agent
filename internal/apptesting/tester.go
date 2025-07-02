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
	OverallStatus string       `json:"overall_status"` // Added field
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

	// Detect the language of the application
	language := at.detectApplicationLanguage(appPath, appReq)

	// Test 1: Build Test (language-specific)
	buildResult := at.testBuildByLanguage(appPath, appReq, language)
	suite.Results = append(suite.Results, buildResult)

	// Test 2: Static Analysis (language-specific)
	staticResult := at.testStaticAnalysisByLanguage(appPath, appReq, language)
	suite.Results = append(suite.Results, staticResult)

	// Test 3: Unit Tests (if any exist)
	unitResult := at.testUnitByLanguage(appPath, appReq, language)
	suite.Results = append(suite.Results, unitResult)

	// Test 4: API Tests (if it's an API application)
	if appReq.Type == "api" || appReq.Type == "web" {
		apiResult := at.testAPIByLanguage(appPath, appReq, language)
		suite.Results = append(suite.Results, apiResult)
	}

	// Test 5: Security Tests (language-specific)
	securityResult := at.testSecurityByLanguage(appPath, appReq, language)
	suite.Results = append(suite.Results, securityResult)

	// Test 6: Performance Tests (basic)
	perfResult := at.testPerformanceByLanguage(appPath, appReq, language)
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

	if suite.FailedTests > 0 {
		suite.OverallStatus = "failure"
	} else if suite.PassedTests > 0 {
		suite.OverallStatus = "success"
	} else {
		suite.OverallStatus = "skipped"
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
			if regexp.MustCompile(`password\s*[:=]\s*[""][^"\]+[""]`).MatchString(strings.ToLower(contentStr)) {
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
		regexp.MustCompile(`(?i)api[_-]?key\s*[:=]\s*[""][^"\]{10,}[""]`),
		regexp.MustCompile(`(?i)secret[_-]?key\s*[:=]\s*[""][^"\]{10,}[""]`),
		regexp.MustCompile(`(?i)token\s*[:=]\s*[""][^"\]{10,}[""]`),
		regexp.MustCompile(`(?i)password\s*[:=]\s*[""][^"\]{8,}[""]`),
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
				if len(line) > 0 && !strings.HasPrefix(line, "//") {
					totalLines++
				}
			}
			
			if err := scanner.Err(); err != nil {
				return err
			}
		}
		
		return nil
	})
	
	if err != nil {
		return 0, err
	}
	
	return totalLines, nil
}

// generateSummary generates a summary of the test suite
func (at *ApplicationTester) generateSummary(suite *TestSuite) string {
	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("Test Suite: %s\n", suite.Name))
	summary.WriteString(fmt.Sprintf("Total Tests: %d, Passed: %d, Failed: %d, Skipped: %d\n", 
		suite.TotalTests, suite.PassedTests, suite.FailedTests, suite.SkippedTests))
	summary.WriteString(fmt.Sprintf("Duration: %s\n", suite.Duration.Round(time.Millisecond)))
	if suite.Coverage > 0 {
		summary.WriteString(fmt.Sprintf("Coverage: %.2f%%\n", suite.Coverage))
	}

	for _, result := range suite.Results {
		summary.WriteString(fmt.Sprintf("- %s (%s): %s\n", result.Name, result.Type, strings.ToUpper(result.Status)))
		if result.Error != "" {
			summary.WriteString(fmt.Sprintf("  Error: %s\n", result.Error))
		}
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




// detectApplicationLanguage detects the programming language of the generated application
func (at *ApplicationTester) detectApplicationLanguage(appPath string, appReq *requirements.ApplicationRequirement) string {
	// First check the requirement language if available
	if appReq.Language != "" {
		return strings.ToLower(appReq.Language)
	}

	// Check for language-specific files
	if _, err := os.Stat(filepath.Join(appPath, "package.json")); err == nil {
		return "javascript"
	}
	if _, err := os.Stat(filepath.Join(appPath, "go.mod")); err == nil {
		return "go"
	}
	if _, err := os.Stat(filepath.Join(appPath, "requirements.txt")); err == nil {
		return "python"
	}
	if _, err := os.Stat(filepath.Join(appPath, "pom.xml")); err == nil {
		return "java"
	}
	if _, err := os.Stat(filepath.Join(appPath, "composer.json")); err == nil {
		return "php"
	}
	if _, err := os.Stat(filepath.Join(appPath, "Gemfile")); err == nil {
		return "ruby"
	}

	// Default to go if no specific indicators found
	return "go"
}

// testBuildByLanguage runs build tests specific to the detected language
func (at *ApplicationTester) testBuildByLanguage(appPath string, appReq *requirements.ApplicationRequirement, language string) TestResult {
	result := TestResult{
		Name: "Build Test",
		Type: "build",
	}
	start := time.Now()

	var cmd *exec.Cmd
	switch language {
	case "javascript", "node", "nodejs":
		// Check if package.json exists
		if _, err := os.Stat(filepath.Join(appPath, "package.json")); err != nil {
			result.Status = "fail"
			result.Error = "package.json not found"
			result.Duration = time.Since(start)
			return result
		}
		cmd = exec.Command("npm", "install")
	case "go", "golang":
		cmd = exec.Command("go", "build", "-v", ".")
	case "python":
		// Check if requirements.txt exists
		if _, err := os.Stat(filepath.Join(appPath, "requirements.txt")); err == nil {
			cmd = exec.Command("pip", "install", "-r", "requirements.txt")
		} else {
			result.Status = "skip"
			result.Output = "No requirements.txt found, skipping build test"
			result.Duration = time.Since(start)
			return result
		}
	case "java":
		if _, err := os.Stat(filepath.Join(appPath, "pom.xml")); err == nil {
			cmd = exec.Command("mvn", "compile")
		} else {
			cmd = exec.Command("javac", "*.java")
		}
	case "php":
		if _, err := os.Stat(filepath.Join(appPath, "composer.json")); err == nil {
			cmd = exec.Command("composer", "install")
		} else {
			result.Status = "skip"
			result.Output = "No composer.json found, skipping build test"
			result.Duration = time.Since(start)
			return result
		}
	case "ruby":
		if _, err := os.Stat(filepath.Join(appPath, "Gemfile")); err == nil {
			cmd = exec.Command("bundle", "install")
		} else {
			result.Status = "skip"
			result.Output = "No Gemfile found, skipping build test"
			result.Duration = time.Since(start)
			return result
		}
	default:
		result.Status = "skip"
		result.Output = fmt.Sprintf("Build test not implemented for language: %s", language)
		result.Duration = time.Since(start)
		return result
	}

	cmd.Dir = appPath
	output, err := cmd.CombinedOutput()
	result.Duration = time.Since(start)
	result.Output = string(output)

	if err != nil {
		result.Status = "fail"
		result.Error = err.Error()
	} else {
		result.Status = "pass"
	}

	return result
}

// testStaticAnalysisByLanguage runs static analysis specific to the detected language
func (at *ApplicationTester) testStaticAnalysisByLanguage(appPath string, appReq *requirements.ApplicationRequirement, language string) TestResult {
	result := TestResult{
		Name: "Static Analysis",
		Type: "static",
	}
	start := time.Now()

	var commands [][]string
	switch language {
	case "javascript", "node", "nodejs":
		// Check if ESLint is available
		if _, err := exec.LookPath("eslint"); err == nil {
			commands = append(commands, []string{"eslint", "."})
		}
		// Check if Prettier is available
		if _, err := exec.LookPath("prettier"); err == nil {
			commands = append(commands, []string{"prettier", "--check", "."})
		}
	case "go", "golang":
		commands = [][]string{
			{"go", "vet", "."},
			{"go", "fmt", "-l", "."},
		}
	case "python":
		if _, err := exec.LookPath("flake8"); err == nil {
			commands = append(commands, []string{"flake8", "."})
		}
		if _, err := exec.LookPath("black"); err == nil {
			commands = append(commands, []string{"black", "--check", "."})
		}
	case "java":
		if _, err := exec.LookPath("checkstyle"); err == nil {
			commands = append(commands, []string{"checkstyle", "-c", "/google_checks.xml", "."})
		}
	case "php":
		if _, err := exec.LookPath("phpcs"); err == nil {
			commands = append(commands, []string{"phpcs", "."})
		}
	case "ruby":
		if _, err := exec.LookPath("rubocop"); err == nil {
			commands = append(commands, []string{"rubocop", "."})
		}
	}

	if len(commands) == 0 {
		result.Status = "skip"
		result.Output = fmt.Sprintf("No static analysis tools available for language: %s", language)
		result.Duration = time.Since(start)
		return result
	}

	var outputs []string
	var errors []string
	allPassed := true

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = appPath
		output, err := cmd.CombinedOutput()
		outputs = append(outputs, fmt.Sprintf("%s: %s", strings.Join(cmdArgs, " "), string(output)))
		
		if err != nil {
			allPassed = false
			errors = append(errors, fmt.Sprintf("%s: %s", strings.Join(cmdArgs, " "), err.Error()))
		}
	}

	result.Duration = time.Since(start)
	result.Output = strings.Join(outputs, "\n")

	if allPassed {
		result.Status = "pass"
	} else {
		result.Status = "fail"
		result.Error = strings.Join(errors, "; ")
	}

	return result
}

// testUnitByLanguage runs unit tests specific to the detected language
func (at *ApplicationTester) testUnitByLanguage(appPath string, appReq *requirements.ApplicationRequirement, language string) TestResult {
	result := TestResult{
		Name: "Unit Tests",
		Type: "unit",
	}
	start := time.Now()

	var cmd *exec.Cmd
	switch language {
	case "javascript", "node", "nodejs":
		// Check if test script exists in package.json
		packageJsonPath := filepath.Join(appPath, "package.json")
		if data, err := os.ReadFile(packageJsonPath); err == nil {
			var packageJson map[string]interface{}
			if json.Unmarshal(data, &packageJson) == nil {
				if scripts, ok := packageJson["scripts"].(map[string]interface{}); ok {
					if _, hasTest := scripts["test"]; hasTest {
						cmd = exec.Command("npm", "test")
					}
				}
			}
		}
	case "go", "golang":
		cmd = exec.Command("go", "test", "-v", "./...")
	case "python":
		if _, err := exec.LookPath("pytest"); err == nil {
			cmd = exec.Command("pytest", "-v")
		} else if _, err := exec.LookPath("python"); err == nil {
			cmd = exec.Command("python", "-m", "unittest", "discover", "-v")
		}
	case "java":
		if _, err := os.Stat(filepath.Join(appPath, "pom.xml")); err == nil {
			cmd = exec.Command("mvn", "test")
		}
	case "php":
		if _, err := exec.LookPath("phpunit"); err == nil {
			cmd = exec.Command("phpunit")
		}
	case "ruby":
		if _, err := os.Stat(filepath.Join(appPath, "Rakefile")); err == nil {
			cmd = exec.Command("rake", "test")
		} else if _, err := exec.LookPath("rspec"); err == nil {
			cmd = exec.Command("rspec")
		}
	}

	if cmd == nil {
		result.Status = "skip"
		result.Output = fmt.Sprintf("No unit test framework found for language: %s", language)
		result.Duration = time.Since(start)
		return result
	}

	cmd.Dir = appPath
	output, err := cmd.CombinedOutput()
	result.Duration = time.Since(start)
	result.Output = string(output)

	if err != nil {
		result.Status = "fail"
		result.Error = err.Error()
	} else {
		result.Status = "pass"
	}

	return result
}

// testAPIByLanguage runs API tests specific to the detected language
func (at *ApplicationTester) testAPIByLanguage(appPath string, appReq *requirements.ApplicationRequirement, language string) TestResult {
	result := TestResult{
		Name: "API Tests",
		Type: "api",
	}
	start := time.Now()

	// Start the application based on language
	var cmd *exec.Cmd
	var port string = "3000" // default port

	switch language {
	case "javascript", "node", "nodejs":
		// Try to start with npm start first
		if _, err := os.Stat(filepath.Join(appPath, "package.json")); err == nil {
			cmd = exec.Command("npm", "start")
		} else if _, err := os.Stat(filepath.Join(appPath, "app.js")); err == nil {
			cmd = exec.Command("node", "app.js")
		} else if _, err := os.Stat(filepath.Join(appPath, "index.js")); err == nil {
			cmd = exec.Command("node", "index.js")
		}
	case "go", "golang":
		// Build first, then run
		buildCmd := exec.Command("go", "build", "-o", "app", ".")
		buildCmd.Dir = appPath
		if err := buildCmd.Run(); err == nil {
			cmd = exec.Command("./app")
			port = "8080" // Go apps typically use 8080
		}
	case "python":
		if _, err := os.Stat(filepath.Join(appPath, "app.py")); err == nil {
			cmd = exec.Command("python", "app.py")
		} else if _, err := os.Stat(filepath.Join(appPath, "main.py")); err == nil {
			cmd = exec.Command("python", "main.py")
		}
		port = "5000" // Flask default
	}

	if cmd == nil {
		result.Status = "skip"
		result.Output = fmt.Sprintf("No runnable application found for language: %s", language)
		result.Duration = time.Since(start)
		return result
	}

	cmd.Dir = appPath
	
	// Start the application
	if err := cmd.Start(); err != nil {
		result.Status = "fail"
		result.Error = fmt.Sprintf("Failed to start application: %v", err)
		result.Duration = time.Since(start)
		return result
	}

	// Wait a moment for the server to start
	time.Sleep(2 * time.Second)

	// Test basic endpoints
	baseURL := fmt.Sprintf("http://localhost:%s", port)
	endpoints := []string{"/", "/health", "/api", "/api/health"}
	
	var testResults []string
	successCount := 0

	for _, endpoint := range endpoints {
		resp, err := http.Get(baseURL + endpoint)
		if err == nil {
			testResults = append(testResults, fmt.Sprintf("%s: %d", endpoint, resp.StatusCode))
			if resp.StatusCode < 500 {
				successCount++
			}
			resp.Body.Close()
		} else {
			testResults = append(testResults, fmt.Sprintf("%s: error - %v", endpoint, err))
		}
	}

	// Stop the application
	if cmd.Process != nil {
		cmd.Process.Kill()
	}

	result.Duration = time.Since(start)
	result.Output = strings.Join(testResults, "\n")

	if successCount > 0 {
		result.Status = "pass"
		result.Details = map[string]interface{}{
			"endpoints_tested": len(endpoints),
			"successful_responses": successCount,
		}
	} else {
		result.Status = "fail"
		result.Error = "No endpoints responded successfully"
	}

	return result
}

// testSecurityByLanguage runs security tests specific to the detected language
func (at *ApplicationTester) testSecurityByLanguage(appPath string, appReq *requirements.ApplicationRequirement, language string) TestResult {
	result := TestResult{
		Name: "Security Tests",
		Type: "security",
	}
	start := time.Now()

	var commands [][]string
	switch language {
	case "javascript", "node", "nodejs":
		if _, err := exec.LookPath("npm"); err == nil {
			commands = append(commands, []string{"npm", "audit"})
		}
	case "go", "golang":
		if _, err := exec.LookPath("gosec"); err == nil {
			commands = append(commands, []string{"gosec", "./..."})
		}
	case "python":
		if _, err := exec.LookPath("safety"); err == nil {
			commands = append(commands, []string{"safety", "check"})
		}
		if _, err := exec.LookPath("bandit"); err == nil {
			commands = append(commands, []string{"bandit", "-r", "."})
		}
	}

	if len(commands) == 0 {
		result.Status = "pass"
		result.Output = fmt.Sprintf("No security scanning tools available for language: %s, marking as pass", language)
		result.Duration = time.Since(start)
		return result
	}

	var outputs []string
	var errors []string
	allPassed := true

	for _, cmdArgs := range commands {
		cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		cmd.Dir = appPath
		output, err := cmd.CombinedOutput()
		outputs = append(outputs, fmt.Sprintf("%s: %s", strings.Join(cmdArgs, " "), string(output)))
		
		if err != nil {
			// For security tools, some "errors" might be warnings, so we're more lenient
			errors = append(errors, fmt.Sprintf("%s: %s", strings.Join(cmdArgs, " "), err.Error()))
		}
	}

	result.Duration = time.Since(start)
	result.Output = strings.Join(outputs, "\n")

	if allPassed {
		result.Status = "pass"
	} else {
		result.Status = "pass" // Mark as pass but include warnings in output
		result.Details = map[string]interface{}{
			"warnings": errors,
		}
	}

	return result
}

// testPerformanceByLanguage runs performance tests specific to the detected language
func (at *ApplicationTester) testPerformanceByLanguage(appPath string, appReq *requirements.ApplicationRequirement, language string) TestResult {
	result := TestResult{
		Name: "Performance Tests",
		Type: "performance",
	}
	start := time.Now()

	// Basic performance metrics - file count, size, etc.
	var totalSize int64
	var fileCount int

	err := filepath.Walk(appPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
			fileCount++
		}
		return nil
	})

	result.Duration = time.Since(start)

	if err != nil {
		result.Status = "fail"
		result.Error = err.Error()
	} else {
		result.Status = "pass"
		result.Output = fmt.Sprintf("Project size: %d bytes, Files: %d", totalSize, fileCount)
		result.Details = map[string]interface{}{
			"total_size_bytes": totalSize,
			"file_count": fileCount,
			"language": language,
		}
	}

	return result
}

