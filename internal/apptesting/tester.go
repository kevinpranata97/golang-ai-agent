package apptesting

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kevinpranata97/golang-ai-agent/internal/requirements"
)

// TestResult represents the result of a single test
type TestResult struct {
	TestName string `json:"test_name"`
	Status   string `json:"status"`
	Error    string `json:"error,omitempty"`
}
// Tester provides functionality to test generated applications
type Tester struct{}

// NewTester creates a new Tester instance
func NewTester() *Tester {
	return &Tester{}
}

// RunTests runs a series of tests on the generated application
func (t *Tester) RunTests(appPath string, appReq *requirements.ApplicationRequirement) ([]TestResult, error) {
	var results []TestResult

	// Example: Basic build test
	buildCmd := exec.Command("go", "build", "-o", appReq.Name, ".")
	buildCmd.Dir = appPath
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		results = append(results, TestResult{
			TestName: "Build Application",
			Status:   "Failed",
			Error:    fmt.Sprintf("Build failed: %s\n%s", err.Error(), string(output)),
		})
		return results, nil // Return early if build fails
	}
	results = append(results, TestResult{
		TestName: "Build Application",
		Status:   "Passed",
	})

	// Example: Run application and check if it starts (for API/Web apps)
	if appReq.Type == "api" || appReq.Type == "web" {
		cmd := exec.Command(filepath.Join(appPath, appReq.Name))
		cmd.Dir = appPath

		// Start the application in a goroutine
		if err := cmd.Start(); err != nil {
			results = append(results, TestResult{
				TestName: "Start Application",
				Status:   "Failed",
				Error:    fmt.Sprintf("Failed to start application: %v", err),
			})
		} else {
			results = append(results, TestResult{
				TestName: "Start Application",
				Status:   "Passed",
			})
			// Give the application some time to start
			time.Sleep(2 * time.Second)

			// Kill the process after testing
			if err := cmd.Process.Kill(); err != nil {
				fmt.Printf("Failed to kill application process: %v\n", err)
			}
		}
	}

	// Add more specific tests based on appReq (e.g., API endpoint tests, UI tests)

	return results, nil
}

// SaveTestResults saves the test results to a JSON file
func (t *Tester) SaveTestResults(appPath string, results []TestResult) error {
	resultsFile := filepath.Join(appPath, "test_results.json")
	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal test results: %w", err)
	}

	if err := ioutil.WriteFile(resultsFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write test results to file: %w", err)
	}

	return nil
}


