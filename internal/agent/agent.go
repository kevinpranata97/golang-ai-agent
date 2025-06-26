package agent

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kevinpranata97/golang-ai-agent/internal/github"
	"github.com/kevinpranata97/golang-ai-agent/internal/storage"
	"github.com/kevinpranata97/golang-ai-agent/internal/testing"
	"github.com/kevinpranata97/golang-ai-agent/internal/workflow"
)

type Agent struct {
	storage        storage.Storage
	githubClient   *github.Client
	testRunner     *testing.TestRunner
	workflowEngine *workflow.Engine
	webhookSecret  string
}

type WebhookPayload struct {
	Action     string `json:"action"`
	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		CloneURL string `json:"clone_url"`
	} `json:"repository"`
	Ref    string `json:"ref"`
	Commits []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"commits"`
}

type Status struct {
	LastActivity time.Time `json:"last_activity"`
	ActiveJobs   int       `json:"active_jobs"`
	TotalJobs    int       `json:"total_jobs"`
	Health       string    `json:"health"`
}

func NewAgent(storage storage.Storage, githubClient *github.Client, testRunner *testing.TestRunner, workflowEngine *workflow.Engine) *Agent {
	return &Agent{
		storage:        storage,
		githubClient:   githubClient,
		testRunner:     testRunner,
		workflowEngine: workflowEngine,
		webhookSecret:  os.Getenv("WEBHOOK_SECRET"),
	}
}

func (a *Agent) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading webhook body: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Verify webhook signature if secret is configured
	if a.webhookSecret != "" {
		signature := r.Header.Get("X-Hub-Signature-256")
		if !a.verifySignature(body, signature) {
			log.Printf("Invalid webhook signature")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}

	var payload WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("Error parsing webhook payload: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("Received webhook for repository: %s, ref: %s", payload.Repository.FullName, payload.Ref)

	// Process webhook asynchronously
	go a.processWebhook(payload)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}

func (a *Agent) processWebhook(payload WebhookPayload) {
	// Only process push events to main/master branch
	if !strings.HasSuffix(payload.Ref, "/main") && !strings.HasSuffix(payload.Ref, "/master") {
		log.Printf("Ignoring push to branch: %s", payload.Ref)
		return
	}

	// Create workflow context
	ctx := workflow.Context{
		Repository: payload.Repository.FullName,
		CloneURL:   payload.Repository.CloneURL,
		Ref:        payload.Ref,
		Commits:    make([]workflow.Commit, len(payload.Commits)),
	}

	for i, commit := range payload.Commits {
		ctx.Commits[i] = workflow.Commit{
			ID:      commit.ID,
			Message: commit.Message,
			Author:  commit.Author.Name,
		}
	}

	// Execute CI/CD workflow
	result := a.workflowEngine.ExecuteWorkflow("ci_cd", ctx)
	
	// Store result
	a.storage.Store(fmt.Sprintf("workflow_%s_%d", payload.Repository.Name, time.Now().Unix()), result)
	
	// Send status back to GitHub
	if len(payload.Commits) > 0 {
		status := "success"
		description := "All checks passed"
		if !result.Success {
			status = "failure"
			description = result.Error
		}
		
		a.githubClient.SetCommitStatus(payload.Repository.FullName, payload.Commits[0].ID, status, description)
	}
}

func (a *Agent) verifySignature(body []byte, signature string) bool {
	if !strings.HasPrefix(signature, "sha256=") {
		return false
	}

	expectedMAC := hmac.New(sha256.New, []byte(a.webhookSecret))
	expectedMAC.Write(body)
	expectedSignature := "sha256=" + hex.EncodeToString(expectedMAC.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (a *Agent) GetStatus(w http.ResponseWriter, r *http.Request) {
	status := Status{
		LastActivity: time.Now(),
		ActiveJobs:   a.workflowEngine.GetActiveJobs(),
		TotalJobs:    a.workflowEngine.GetTotalJobs(),
		Health:       "healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

