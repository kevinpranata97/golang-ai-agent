package workflow

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

type Engine struct {
	workflows   map[string]Workflow
	activeJobs  int
	totalJobs   int
	mutex       sync.RWMutex
}

type Workflow struct {
	Name  string
	Steps []Step
}

type Step struct {
	Name    string
	Command string
	Args    []string
	WorkDir string
	Timeout time.Duration
}

type Context struct {
	Repository string
	CloneURL   string
	Ref        string
	Commits    []Commit
	WorkDir    string
}

type Commit struct {
	ID      string
	Message string
	Author  string
}

type Result struct {
	Success   bool                   `json:"success"`
	Error     string                 `json:"error,omitempty"`
	Steps     []StepResult           `json:"steps"`
	Duration  time.Duration          `json:"duration"`
	Context   Context                `json:"context"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type StepResult struct {
	Name     string        `json:"name"`
	Success  bool          `json:"success"`
	Output   string        `json:"output"`
	Error    string        `json:"error,omitempty"`
	Duration time.Duration `json:"duration"`
}

func NewEngine() *Engine {
	engine := &Engine{
		workflows: make(map[string]Workflow),
	}
	
	// Register default workflows
	engine.registerDefaultWorkflows()
	
	return engine
}

func (e *Engine) registerDefaultWorkflows() {
	// CI/CD Workflow
	cicdWorkflow := Workflow{
		Name: "ci_cd",
		Steps: []Step{
			{
				Name:    "clone",
				Command: "git",
				Args:    []string{"clone", "", ""},
				Timeout: 5 * time.Minute,
			},
			{
				Name:    "analyze",
				Command: "echo",
				Args:    []string{"Analyzing repository structure..."},
				Timeout: 1 * time.Minute,
			},
			{
				Name:    "build",
				Command: "echo",
				Args:    []string{"Building application..."},
				Timeout: 10 * time.Minute,
			},
			{
				Name:    "test",
				Command: "echo",
				Args:    []string{"Running tests..."},
				Timeout: 15 * time.Minute,
			},
			{
				Name:    "security_scan",
				Command: "echo",
				Args:    []string{"Running security scan..."},
				Timeout: 5 * time.Minute,
			},
		},
	}
	
	e.workflows["ci_cd"] = cicdWorkflow
}

func (e *Engine) ExecuteWorkflow(name string, ctx Context) Result {
	e.mutex.Lock()
	e.activeJobs++
	e.totalJobs++
	e.mutex.Unlock()
	
	defer func() {
		e.mutex.Lock()
		e.activeJobs--
		e.mutex.Unlock()
	}()
	
	workflow, exists := e.workflows[name]
	if !exists {
		return Result{
			Success: false,
			Error:   fmt.Sprintf("workflow '%s' not found", name),
			Context: ctx,
		}
	}
	
	log.Printf("Executing workflow '%s' for repository '%s'", name, ctx.Repository)
	
	startTime := time.Now()
	result := Result{
		Success:  true,
		Steps:    make([]StepResult, 0, len(workflow.Steps)),
		Context:  ctx,
		Metadata: make(map[string]interface{}),
	}
	
	// Create temporary working directory
	tempDir, err := os.MkdirTemp("", "workflow_"+name+"_")
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("failed to create temp directory: %v", err)
		result.Duration = time.Since(startTime)
		return result
	}
	defer os.RemoveAll(tempDir)
	
	ctx.WorkDir = tempDir
	
	// Execute each step
	for _, step := range workflow.Steps {
		stepResult := e.executeStep(step, ctx)
		result.Steps = append(result.Steps, stepResult)
		
		if !stepResult.Success {
			result.Success = false
			result.Error = fmt.Sprintf("step '%s' failed: %s", step.Name, stepResult.Error)
			break
		}
	}
	
	result.Duration = time.Since(startTime)
	log.Printf("Workflow '%s' completed in %v, success: %v", name, result.Duration, result.Success)
	
	return result
}

func (e *Engine) executeStep(step Step, ctx Context) StepResult {
	log.Printf("Executing step: %s", step.Name)
	
	startTime := time.Now()
	stepResult := StepResult{
		Name:    step.Name,
		Success: true,
	}
	
	// Prepare command and arguments
	command := step.Command
	args := make([]string, len(step.Args))
	copy(args, step.Args)
	
	// Handle special cases
	switch step.Name {
	case "clone":
		if len(args) >= 2 {
			args[1] = ctx.CloneURL
			args[2] = filepath.Join(ctx.WorkDir, "repo")
		}
	case "build":
		// Detect project type and use appropriate build command
		repoPath := filepath.Join(ctx.WorkDir, "repo")
		if e.fileExists(filepath.Join(repoPath, "go.mod")) {
			command = "go"
			args = []string{"build", "./..."}
		} else if e.fileExists(filepath.Join(repoPath, "package.json")) {
			command = "npm"
			args = []string{"install"}
		} else if e.fileExists(filepath.Join(repoPath, "Makefile")) {
			command = "make"
			args = []string{}
		}
	case "test":
		// Detect test framework and run tests
		repoPath := filepath.Join(ctx.WorkDir, "repo")
		if e.fileExists(filepath.Join(repoPath, "go.mod")) {
			command = "go"
			args = []string{"test", "./..."}
		} else if e.fileExists(filepath.Join(repoPath, "package.json")) {
			command = "npm"
			args = []string{"test"}
		}
	}
	
	// Set working directory
	workDir := ctx.WorkDir
	if step.WorkDir != "" {
		workDir = step.WorkDir
	} else if step.Name != "clone" {
		workDir = filepath.Join(ctx.WorkDir, "repo")
	}
	
	// Execute command
	cmd := exec.Command(command, args...)
	cmd.Dir = workDir
	
	output, err := cmd.CombinedOutput()
	stepResult.Output = string(output)
	stepResult.Duration = time.Since(startTime)
	
	if err != nil {
		stepResult.Success = false
		stepResult.Error = err.Error()
		log.Printf("Step '%s' failed: %v", step.Name, err)
	} else {
		log.Printf("Step '%s' completed successfully", step.Name)
	}
	
	return stepResult
}

func (e *Engine) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (e *Engine) GetActiveJobs() int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.activeJobs
}

func (e *Engine) GetTotalJobs() int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.totalJobs
}

func (e *Engine) RegisterWorkflow(workflow Workflow) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.workflows[workflow.Name] = workflow
}

func (e *Engine) ListWorkflows() []string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	
	names := make([]string, 0, len(e.workflows))
	for name := range e.workflows {
		names = append(names, name)
	}
	return names
}

