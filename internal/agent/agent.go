package agent

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// Agent represents the core AI agent
type Agent struct {
	ID        string
	Status    string
	CreatedAt time.Time
	Jobs      map[string]*Job
}

// Job represents a processing job
type Job struct {
	ID          string
	Type        string
	Status      string
	Description string
	CreatedAt   time.Time
	CompletedAt *time.Time
	Result      interface{}
	Error       string
}

// NewAgent creates a new AI agent instance
func NewAgent() *Agent {
	return &Agent{
		ID:        generateID(),
		Status:    "idle",
		CreatedAt: time.Now(),
		Jobs:      make(map[string]*Job),
	}
}

// generateID generates a random ID
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// GetStatus returns the current status of the agent
func (a *Agent) GetStatus() map[string]interface{} {
	activeJobs := 0
	for _, job := range a.Jobs {
		if job.Status == "running" {
			activeJobs++
		}
	}

	return map[string]interface{}{
		"id":            a.ID,
		"status":        a.Status,
		"created_at":    a.CreatedAt,
		"active_jobs":   activeJobs,
		"total_jobs":    len(a.Jobs),
		"capabilities": []string{
			"application_generation",
			"requirement_analysis",
			"code_generation",
			"testing",
			"github_integration",
		},
	}
}

// CreateJob creates a new job
func (a *Agent) CreateJob(jobType, description string) *Job {
	job := &Job{
		ID:          generateID(),
		Type:        jobType,
		Status:      "created",
		Description: description,
		CreatedAt:   time.Now(),
	}

	a.Jobs[job.ID] = job
	return job
}

// StartJob starts a job
func (a *Agent) StartJob(jobID string) error {
	job, exists := a.Jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	job.Status = "running"
	return nil
}

// CompleteJob completes a job
func (a *Agent) CompleteJob(jobID string, result interface{}) error {
	job, exists := a.Jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	now := time.Now()
	job.Status = "completed"
	job.CompletedAt = &now
	job.Result = result
	return nil
}

// FailJob marks a job as failed
func (a *Agent) FailJob(jobID string, err error) error {
	job, exists := a.Jobs[jobID]
	if !exists {
		return fmt.Errorf("job not found: %s", jobID)
	}

	now := time.Now()
	job.Status = "failed"
	job.CompletedAt = &now
	job.Error = err.Error()
	return nil
}

