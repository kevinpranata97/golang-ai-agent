package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Client struct {
	token      string
	httpClient *http.Client
}

type CommitStatus struct {
	State       string `json:"state"`
	Description string `json:"description"`
	Context     string `json:"context"`
}

type Repository struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	CloneURL    string `json:"clone_url"`
	Language    string `json:"language"`
	Description string `json:"description"`
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{},
	}
}

func (c *Client) SetCommitStatus(repo, sha, state, description string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/statuses/%s", repo, sha)
	
	status := CommitStatus{
		State:       state,
		Description: description,
		Context:     "golang-ai-agent",
	}
	
	jsonData, err := json.Marshal(status)
	if err != nil {
		return err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to set commit status: %s", string(body))
	}
	
	return nil
}

func (c *Client) GetRepository(repo string) (*Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", repo)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get repository: %d", resp.StatusCode)
	}
	
	var repository Repository
	if err := json.NewDecoder(resp.Body).Decode(&repository); err != nil {
		return nil, err
	}
	
	return &repository, nil
}

func (c *Client) CloneRepository(cloneURL, destination string) error {
	// Add token to clone URL for authentication
	authenticatedURL := strings.Replace(cloneURL, "https://", fmt.Sprintf("https://%s@", c.token), 1)
	
	cmd := exec.Command("git", "clone", authenticatedURL, destination)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %s, output: %s", err, string(output))
	}
	
	return nil
}

func (c *Client) AnalyzeRepository(repoPath string) (*RepositoryAnalysis, error) {
	analysis := &RepositoryAnalysis{
		Languages:    make(map[string]int),
		Files:        []string{},
		HasTests:     false,
		HasDockerfile: false,
		HasMakefile:  false,
	}
	
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}
		
		relPath, _ := filepath.Rel(repoPath, path)
		analysis.Files = append(analysis.Files, relPath)
		
		// Analyze file types
		ext := filepath.Ext(path)
		switch ext {
		case ".go":
			analysis.Languages["Go"]++
			if strings.Contains(relPath, "_test.go") {
				analysis.HasTests = true
			}
		case ".js", ".jsx":
			analysis.Languages["JavaScript"]++
		case ".ts", ".tsx":
			analysis.Languages["TypeScript"]++
		case ".py":
			analysis.Languages["Python"]++
		case ".java":
			analysis.Languages["Java"]++
		case ".cpp", ".cc", ".cxx":
			analysis.Languages["C++"]++
		case ".c":
			analysis.Languages["C"]++
		}
		
		// Check for special files
		filename := filepath.Base(path)
		switch filename {
		case "Dockerfile":
			analysis.HasDockerfile = true
		case "Makefile":
			analysis.HasMakefile = true
		case "package.json":
			analysis.HasPackageJSON = true
		case "go.mod":
			analysis.HasGoMod = true
		}
		
		return nil
	})
	
	return analysis, err
}

type RepositoryAnalysis struct {
	Languages      map[string]int `json:"languages"`
	Files          []string       `json:"files"`
	HasTests       bool           `json:"has_tests"`
	HasDockerfile  bool           `json:"has_dockerfile"`
	HasMakefile    bool           `json:"has_makefile"`
	HasPackageJSON bool           `json:"has_package_json"`
	HasGoMod       bool           `json:"has_go_mod"`
}

