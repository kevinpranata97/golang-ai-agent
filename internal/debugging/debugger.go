package debugging

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Debugger struct {
	projectPath string
	logLevel    string
}

type DebugSession struct {
	ID          string                 `json:"id"`
	ProjectPath string                 `json:"project_path"`
	StartTime   time.Time              `json:"start_time"`
	Status      string                 `json:"status"`
	Breakpoints []Breakpoint           `json:"breakpoints"`
	Variables   map[string]interface{} `json:"variables"`
	StackTrace  []StackFrame           `json:"stack_trace"`
	Logs        []LogEntry             `json:"logs"`
}

type Breakpoint struct {
	File string `json:"file"`
	Line int    `json:"line"`
	ID   int    `json:"id"`
}

type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	File      string    `json:"file,omitempty"`
	Line      int       `json:"line,omitempty"`
}

type DebugResult struct {
	Success      bool                   `json:"success"`
	Issues       []DebugIssue           `json:"issues"`
	Suggestions  []Suggestion           `json:"suggestions"`
	Performance  PerformanceAnalysis    `json:"performance"`
	MemoryLeaks  []MemoryLeak           `json:"memory_leaks"`
	Metrics      map[string]interface{} `json:"metrics"`
	Duration     time.Duration          `json:"duration"`
}

type DebugIssue struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	File        string `json:"file"`
	Line        int    `json:"line"`
	Function    string `json:"function,omitempty"`
	Description string `json:"description"`
	Context     string `json:"context,omitempty"`
}

type Suggestion struct {
	Type        string `json:"type"`
	Priority    string `json:"priority"`
	Description string `json:"description"`
	Code        string `json:"code,omitempty"`
	File        string `json:"file,omitempty"`
	Line        int    `json:"line,omitempty"`
}

type PerformanceAnalysis struct {
	HotSpots        []HotSpot `json:"hot_spots"`
	SlowFunctions   []SlowFunction `json:"slow_functions"`
	MemoryUsage     MemoryUsage `json:"memory_usage"`
	GoroutineLeaks  int `json:"goroutine_leaks"`
}

type HotSpot struct {
	Function    string        `json:"function"`
	File        string        `json:"file"`
	Line        int           `json:"line"`
	CPUPercent  float64       `json:"cpu_percent"`
	CallCount   int           `json:"call_count"`
	TotalTime   time.Duration `json:"total_time"`
	AverageTime time.Duration `json:"average_time"`
}

type SlowFunction struct {
	Function    string        `json:"function"`
	File        string        `json:"file"`
	Line        int           `json:"line"`
	Duration    time.Duration `json:"duration"`
	Calls       int           `json:"calls"`
}

type MemoryUsage struct {
	HeapSize     int64 `json:"heap_size"`
	StackSize    int64 `json:"stack_size"`
	Allocations  int64 `json:"allocations"`
	Deallocations int64 `json:"deallocations"`
}

type MemoryLeak struct {
	Type        string `json:"type"`
	Location    string `json:"location"`
	Size        int64  `json:"size"`
	Description string `json:"description"`
}

func NewDebugger(projectPath string) *Debugger {
	return &Debugger{
		projectPath: projectPath,
		logLevel:    "info",
	}
}

func (d *Debugger) StartDebugSession() (*DebugSession, error) {
	session := &DebugSession{
		ID:          fmt.Sprintf("debug_%d", time.Now().Unix()),
		ProjectPath: d.projectPath,
		StartTime:   time.Now(),
		Status:      "active",
		Breakpoints: []Breakpoint{},
		Variables:   make(map[string]interface{}),
		StackTrace:  []StackFrame{},
		Logs:        []LogEntry{},
	}
	
	return session, nil
}

func (d *Debugger) AnalyzeProject() DebugResult {
	startTime := time.Now()
	
	result := DebugResult{
		Success:     true,
		Issues:      []DebugIssue{},
		Suggestions: []Suggestion{},
		Metrics:     make(map[string]interface{}),
	}
	
	// Analyze code for common issues
	d.analyzeCodeIssues(&result)
	
	// Analyze logs for errors
	d.analyzeLogs(&result)
	
	// Check for performance issues
	d.analyzePerformance(&result)
	
	// Check for memory leaks
	d.analyzeMemoryLeaks(&result)
	
	// Generate suggestions
	d.generateSuggestions(&result)
	
	result.Duration = time.Since(startTime)
	
	return result
}

func (d *Debugger) analyzeCodeIssues(result *DebugResult) {
	filepath.Walk(d.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !d.isSourceFile(path) {
			return nil
		}
		
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		
		contentStr := string(content)
		lines := strings.Split(contentStr, "\n")
		
		// Check for common Go issues
		if strings.HasSuffix(path, ".go") {
			d.analyzeGoIssues(path, lines, result)
		}
		
		// Check for common JavaScript issues
		if strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".ts") {
			d.analyzeJavaScriptIssues(path, lines, result)
		}
		
		return nil
	})
}

func (d *Debugger) analyzeGoIssues(filePath string, lines []string, result *DebugResult) {
	for i, line := range lines {
		lineNum := i + 1
		trimmedLine := strings.TrimSpace(line)
		
		// Check for potential nil pointer dereferences
		if strings.Contains(trimmedLine, ".") && !strings.Contains(trimmedLine, "if") && !strings.Contains(trimmedLine, "!=") {
			if matched, _ := regexp.MatchString(`\w+\.\w+`, trimmedLine); matched {
				result.Issues = append(result.Issues, DebugIssue{
					Type:        "nil_pointer_risk",
					Severity:    "warning",
					File:        filePath,
					Line:        lineNum,
					Description: "Potential nil pointer dereference - consider nil check",
					Context:     trimmedLine,
				})
			}
		}
		
		// Check for missing error handling
		if strings.Contains(trimmedLine, ":=") && strings.Contains(trimmedLine, "err") && !strings.Contains(trimmedLine, "if err") {
			nextLine := ""
			if i+1 < len(lines) {
				nextLine = strings.TrimSpace(lines[i+1])
			}
			if !strings.Contains(nextLine, "if err") {
				result.Issues = append(result.Issues, DebugIssue{
					Type:        "missing_error_handling",
					Severity:    "error",
					File:        filePath,
					Line:        lineNum,
					Description: "Error not handled - this could cause runtime panics",
					Context:     trimmedLine,
				})
			}
		}
		
		// Check for goroutine leaks
		if strings.Contains(trimmedLine, "go func") && !strings.Contains(trimmedLine, "defer") {
			result.Issues = append(result.Issues, DebugIssue{
				Type:        "goroutine_leak_risk",
				Severity:    "warning",
				File:        filePath,
				Line:        lineNum,
				Description: "Goroutine without proper cleanup - potential leak",
				Context:     trimmedLine,
			})
		}
		
		// Check for hardcoded values
		if matched, _ := regexp.MatchString(`"(http://|https://|localhost|127\.0\.0\.1)"`, trimmedLine); matched {
			result.Issues = append(result.Issues, DebugIssue{
				Type:        "hardcoded_url",
				Severity:    "info",
				File:        filePath,
				Line:        lineNum,
				Description: "Hardcoded URL found - consider using configuration",
				Context:     trimmedLine,
			})
		}
	}
}

func (d *Debugger) analyzeJavaScriptIssues(filePath string, lines []string, result *DebugResult) {
	for i, line := range lines {
		lineNum := i + 1
		trimmedLine := strings.TrimSpace(line)
		
		// Check for console.log statements
		if strings.Contains(trimmedLine, "console.log") {
			result.Issues = append(result.Issues, DebugIssue{
				Type:        "debug_statement",
				Severity:    "info",
				File:        filePath,
				Line:        lineNum,
				Description: "Debug console.log statement found - remove before production",
				Context:     trimmedLine,
			})
		}
		
		// Check for == instead of ===
		if strings.Contains(trimmedLine, "==") && !strings.Contains(trimmedLine, "===") {
			result.Issues = append(result.Issues, DebugIssue{
				Type:        "loose_equality",
				Severity:    "warning",
				File:        filePath,
				Line:        lineNum,
				Description: "Use strict equality (===) instead of loose equality (==)",
				Context:     trimmedLine,
			})
		}
		
		// Check for var declarations
		if matched, _ := regexp.MatchString(`^\s*var\s+`, trimmedLine); matched {
			result.Issues = append(result.Issues, DebugIssue{
				Type:        "var_declaration",
				Severity:    "info",
				File:        filePath,
				Line:        lineNum,
				Description: "Consider using 'let' or 'const' instead of 'var'",
				Context:     trimmedLine,
			})
		}
	}
}

func (d *Debugger) analyzeLogs(result *DebugResult) {
	logFiles := d.findLogFiles()
	
	for _, logFile := range logFiles {
		d.analyzeLogFile(logFile, result)
	}
}

func (d *Debugger) findLogFiles() []string {
	var logFiles []string
	
	filepath.Walk(d.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if strings.HasSuffix(path, ".log") || strings.Contains(path, "log") {
			logFiles = append(logFiles, path)
		}
		
		return nil
	})
	
	return logFiles
}

func (d *Debugger) analyzeLogFile(logFile string, result *DebugResult) {
	file, err := os.Open(logFile)
	if err != nil {
		return
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	lineNum := 0
	
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		
		// Check for error patterns
		if d.containsErrorPattern(line) {
			result.Issues = append(result.Issues, DebugIssue{
				Type:        "log_error",
				Severity:    "error",
				File:        logFile,
				Line:        lineNum,
				Description: "Error found in log file",
				Context:     line,
			})
		}
		
		// Check for warning patterns
		if d.containsWarningPattern(line) {
			result.Issues = append(result.Issues, DebugIssue{
				Type:        "log_warning",
				Severity:    "warning",
				File:        logFile,
				Line:        lineNum,
				Description: "Warning found in log file",
				Context:     line,
			})
		}
	}
}

func (d *Debugger) containsErrorPattern(line string) bool {
	errorPatterns := []string{
		"ERROR", "error", "Error",
		"FATAL", "fatal", "Fatal",
		"PANIC", "panic", "Panic",
		"exception", "Exception",
		"failed", "Failed",
	}
	
	for _, pattern := range errorPatterns {
		if strings.Contains(line, pattern) {
			return true
		}
	}
	
	return false
}

func (d *Debugger) containsWarningPattern(line string) bool {
	warningPatterns := []string{
		"WARN", "warn", "Warning", "WARNING",
		"deprecated", "Deprecated",
		"timeout", "Timeout",
	}
	
	for _, pattern := range warningPatterns {
		if strings.Contains(line, pattern) {
			return true
		}
	}
	
	return false
}

func (d *Debugger) analyzePerformance(result *DebugResult) {
	// This is a simplified performance analysis
	// In a real implementation, you would use profiling tools
	
	result.Performance = PerformanceAnalysis{
		HotSpots:       []HotSpot{},
		SlowFunctions:  []SlowFunction{},
		MemoryUsage:    MemoryUsage{},
		GoroutineLeaks: 0,
	}
	
	// Check for potential performance issues in code
	filepath.Walk(d.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return err
		}
		
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		
		contentStr := string(content)
		lines := strings.Split(contentStr, "\n")
		
		for i, line := range lines {
			lineNum := i + 1
			trimmedLine := strings.TrimSpace(line)
			
			// Check for loops that might be inefficient
			if strings.Contains(trimmedLine, "for") && strings.Contains(trimmedLine, "range") {
				if strings.Contains(trimmedLine, "append") {
					result.Issues = append(result.Issues, DebugIssue{
						Type:        "performance_issue",
						Severity:    "warning",
						File:        path,
						Line:        lineNum,
						Description: "Potential performance issue: appending in loop without pre-allocation",
						Context:     trimmedLine,
					})
				}
			}
			
			// Check for string concatenation in loops
			if strings.Contains(trimmedLine, "for") && strings.Contains(trimmedLine, "+") && strings.Contains(trimmedLine, "string") {
				result.Issues = append(result.Issues, DebugIssue{
					Type:        "performance_issue",
					Severity:    "warning",
					File:        path,
					Line:        lineNum,
					Description: "String concatenation in loop - consider using strings.Builder",
					Context:     trimmedLine,
				})
			}
		}
		
		return nil
	})
}

func (d *Debugger) analyzeMemoryLeaks(result *DebugResult) {
	// Simplified memory leak detection
	result.MemoryLeaks = []MemoryLeak{}
	
	// This would typically involve running the application with memory profiling
	// For now, we'll just check for common patterns that can cause leaks
	
	filepath.Walk(d.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return err
		}
		
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		
		contentStr := string(content)
		
		// Check for unclosed resources
		if strings.Contains(contentStr, "os.Open") && !strings.Contains(contentStr, "defer") {
			result.MemoryLeaks = append(result.MemoryLeaks, MemoryLeak{
				Type:        "unclosed_file",
				Location:    path,
				Description: "File opened without defer close - potential resource leak",
			})
		}
		
		return nil
	})
}

func (d *Debugger) generateSuggestions(result *DebugResult) {
	// Generate suggestions based on found issues
	for _, issue := range result.Issues {
		switch issue.Type {
		case "missing_error_handling":
			result.Suggestions = append(result.Suggestions, Suggestion{
				Type:        "error_handling",
				Priority:    "high",
				Description: "Add proper error handling",
				Code:        "if err != nil {\n    return err\n}",
				File:        issue.File,
				Line:        issue.Line,
			})
		case "nil_pointer_risk":
			result.Suggestions = append(result.Suggestions, Suggestion{
				Type:        "nil_check",
				Priority:    "medium",
				Description: "Add nil check before dereferencing",
				Code:        "if obj != nil {\n    // safe to use obj\n}",
				File:        issue.File,
				Line:        issue.Line,
			})
		case "performance_issue":
			result.Suggestions = append(result.Suggestions, Suggestion{
				Type:        "performance",
				Priority:    "medium",
				Description: "Optimize for better performance",
				File:        issue.File,
				Line:        issue.Line,
			})
		}
	}
}

func (d *Debugger) isSourceFile(path string) bool {
	ext := filepath.Ext(path)
	sourceExts := []string{".go", ".js", ".ts", ".py", ".java", ".cpp", ".c", ".h"}
	
	for _, sourceExt := range sourceExts {
		if ext == sourceExt {
			return true
		}
	}
	
	return false
}

func (d *Debugger) RunProfiler(duration time.Duration) (map[string]interface{}, error) {
	// This would run a profiler for the specified duration
	// For Go applications, this could use pprof
	
	result := make(map[string]interface{})
	
	// Simulate profiling
	time.Sleep(duration)
	
	result["cpu_profile"] = "cpu_profile_data"
	result["memory_profile"] = "memory_profile_data"
	result["goroutine_profile"] = "goroutine_profile_data"
	
	return result, nil
}

func (d *Debugger) AttachDebugger(processID int) error {
	// This would attach a debugger to a running process
	// For Go, this could use delve
	
	cmd := exec.Command("dlv", "attach", fmt.Sprintf("%d", processID))
	return cmd.Start()
}

func (d *Debugger) SetBreakpoint(file string, line int) error {
	// This would set a breakpoint in the debugger
	// Implementation would depend on the specific debugger being used
	
	return nil
}

func (d *Debugger) GetStackTrace() ([]StackFrame, error) {
	// This would get the current stack trace from the debugger
	
	return []StackFrame{}, nil
}

