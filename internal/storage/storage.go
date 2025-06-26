package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kevinpranata97/golang-ai-agent/internal/apptesting"
	"github.com/kevinpranata97/golang-ai-agent/internal/requirements"
)

// ProjectData represents stored project data
type ProjectData struct {
	ID           string                              `json:"id"`
	Name         string                              `json:"name"`
	Description  string                              `json:"description"`
	Requirements *requirements.ApplicationRequirement `json:"requirements"`
	GeneratedAt  time.Time                           `json:"generated_at"`
	AppPath      string                              `json:"app_path"`
	TestResults  *apptesting.TestSuite               `json:"test_results,omitempty"`
	Status       string                              `json:"status"` // generating, testing, completed, failed
	Iterations   []IterationData                     `json:"iterations"`
	Metadata     map[string]interface{}              `json:"metadata"`
}

// IterationData represents data for each iteration/improvement
type IterationData struct {
	ID          string                    `json:"id"`
	Timestamp   time.Time                 `json:"timestamp"`
	Changes     []string                  `json:"changes"`
	TestResults *apptesting.TestSuite     `json:"test_results,omitempty"`
	Improvements []string                 `json:"improvements"`
	Status      string                    `json:"status"`
}

// AnalysisData represents analysis results
type AnalysisData struct {
	ProjectID     string                 `json:"project_id"`
	Timestamp     time.Time              `json:"timestamp"`
	CodeQuality   CodeQualityMetrics     `json:"code_quality"`
	Performance   PerformanceMetrics     `json:"performance"`
	Security      SecurityMetrics        `json:"security"`
	Suggestions   []ImprovementSuggestion `json:"suggestions"`
}

// CodeQualityMetrics represents code quality metrics
type CodeQualityMetrics struct {
	LinesOfCode       int     `json:"lines_of_code"`
	CyclomaticComplexity int  `json:"cyclomatic_complexity"`
	TestCoverage      float64 `json:"test_coverage"`
	DuplicationRatio  float64 `json:"duplication_ratio"`
	TechnicalDebt     string  `json:"technical_debt"`
	Maintainability   string  `json:"maintainability"`
}

// PerformanceMetrics represents performance metrics
type PerformanceMetrics struct {
	BinarySize       int64   `json:"binary_size"`
	BuildTime        float64 `json:"build_time"`
	StartupTime      float64 `json:"startup_time"`
	MemoryUsage      int64   `json:"memory_usage"`
	ResponseTime     float64 `json:"response_time"`
}

// SecurityMetrics represents security metrics
type SecurityMetrics struct {
	Vulnerabilities   int      `json:"vulnerabilities"`
	SecurityScore     float64  `json:"security_score"`
	HardcodedSecrets  int      `json:"hardcoded_secrets"`
	SecurityIssues    []string `json:"security_issues"`
}

// ImprovementSuggestion represents a suggestion for improvement
type ImprovementSuggestion struct {
	Type        string `json:"type"` // performance, security, quality, functionality
	Priority    string `json:"priority"` // high, medium, low
	Description string `json:"description"`
	Impact      string `json:"impact"`
	Effort      string `json:"effort"`
	Code        string `json:"code,omitempty"`
}

// Storage interface defines storage operations
type Storage interface {
	SaveProject(project *ProjectData) error
	GetProject(id string) (*ProjectData, error)
	ListProjects() ([]*ProjectData, error)
	UpdateProject(project *ProjectData) error
	DeleteProject(id string) error
	SaveAnalysis(analysis *AnalysisData) error
	GetAnalysis(projectID string) ([]*AnalysisData, error)
	GetProjectStats() (*ProjectStats, error)
	Cleanup(olderThan time.Duration) error
}

// ProjectStats represents overall project statistics
type ProjectStats struct {
	TotalProjects     int                    `json:"total_projects"`
	CompletedProjects int                    `json:"completed_projects"`
	FailedProjects    int                    `json:"failed_projects"`
	AvgTestCoverage   float64                `json:"avg_test_coverage"`
	AvgBuildTime      float64                `json:"avg_build_time"`
	PopularLanguages  map[string]int         `json:"popular_languages"`
	PopularFrameworks map[string]int         `json:"popular_frameworks"`
	RecentActivity    []ProjectData          `json:"recent_activity"`
}

// FileStorage implements Storage interface using file system
type FileStorage struct {
	baseDir string
}

// NewFileStorage creates a new file storage instance
func NewFileStorage(baseDir string) *FileStorage {
	return &FileStorage{
		baseDir: baseDir,
	}
}

// Initialize initializes the storage
func (fs *FileStorage) Initialize() error {
	dirs := []string{
		filepath.Join(fs.baseDir, "projects"),
		filepath.Join(fs.baseDir, "analysis"),
		filepath.Join(fs.baseDir, "backups"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	return nil
}

// SaveProject saves project data to storage
func (fs *FileStorage) SaveProject(project *ProjectData) error {
	if err := fs.Initialize(); err != nil {
		return err
	}

	projectPath := filepath.Join(fs.baseDir, "projects", project.ID+".json")
	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal project data: %v", err)
	}

	if err := os.WriteFile(projectPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write project file: %v", err)
	}

	return nil
}

// GetProject retrieves project data from storage
func (fs *FileStorage) GetProject(id string) (*ProjectData, error) {
	projectPath := filepath.Join(fs.baseDir, "projects", id+".json")
	data, err := os.ReadFile(projectPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("project not found: %s", id)
		}
		return nil, fmt.Errorf("failed to read project file: %v", err)
	}

	var project ProjectData
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("failed to unmarshal project data: %v", err)
	}

	return &project, nil
}

// ListProjects lists all projects in storage
func (fs *FileStorage) ListProjects() ([]*ProjectData, error) {
	projectsDir := filepath.Join(fs.baseDir, "projects")
	if _, err := os.Stat(projectsDir); os.IsNotExist(err) {
		return []*ProjectData{}, nil
	}

	files, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read projects directory: %v", err)
	}

	var projects []*ProjectData
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			id := strings.TrimSuffix(file.Name(), ".json")
			project, err := fs.GetProject(id)
			if err != nil {
				continue // Skip corrupted files
			}
			projects = append(projects, project)
		}
	}

	return projects, nil
}

// UpdateProject updates existing project data
func (fs *FileStorage) UpdateProject(project *ProjectData) error {
	return fs.SaveProject(project) // Same as save for file storage
}

// DeleteProject deletes project data from storage
func (fs *FileStorage) DeleteProject(id string) error {
	projectPath := filepath.Join(fs.baseDir, "projects", id+".json")
	if err := os.Remove(projectPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("project not found: %s", id)
		}
		return fmt.Errorf("failed to delete project file: %v", err)
	}

	// Also delete associated analysis data
	analysisDir := filepath.Join(fs.baseDir, "analysis", id)
	if _, err := os.Stat(analysisDir); err == nil {
		os.RemoveAll(analysisDir)
	}

	return nil
}

// SaveAnalysis saves analysis data to storage
func (fs *FileStorage) SaveAnalysis(analysis *AnalysisData) error {
	if err := fs.Initialize(); err != nil {
		return err
	}

	analysisDir := filepath.Join(fs.baseDir, "analysis", analysis.ProjectID)
	if err := os.MkdirAll(analysisDir, 0755); err != nil {
		return fmt.Errorf("failed to create analysis directory: %v", err)
	}

	filename := fmt.Sprintf("%d.json", analysis.Timestamp.Unix())
	analysisPath := filepath.Join(analysisDir, filename)

	data, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal analysis data: %v", err)
	}

	if err := os.WriteFile(analysisPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write analysis file: %v", err)
	}

	return nil
}

// GetAnalysis retrieves analysis data for a project
func (fs *FileStorage) GetAnalysis(projectID string) ([]*AnalysisData, error) {
	analysisDir := filepath.Join(fs.baseDir, "analysis", projectID)
	if _, err := os.Stat(analysisDir); os.IsNotExist(err) {
		return []*AnalysisData{}, nil
	}

	files, err := os.ReadDir(analysisDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read analysis directory: %v", err)
	}

	var analyses []*AnalysisData
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			analysisPath := filepath.Join(analysisDir, file.Name())
			data, err := os.ReadFile(analysisPath)
			if err != nil {
				continue // Skip corrupted files
			}

			var analysis AnalysisData
			if err := json.Unmarshal(data, &analysis); err != nil {
				continue // Skip corrupted files
			}

			analyses = append(analyses, &analysis)
		}
	}

	return analyses, nil
}

// GetProjectStats calculates and returns project statistics
func (fs *FileStorage) GetProjectStats() (*ProjectStats, error) {
	projects, err := fs.ListProjects()
	if err != nil {
		return nil, err
	}

	stats := &ProjectStats{
		TotalProjects:     len(projects),
		PopularLanguages:  make(map[string]int),
		PopularFrameworks: make(map[string]int),
		RecentActivity:    []ProjectData{},
	}

	var totalCoverage float64
	var totalBuildTime float64
	var coverageCount int
	var buildTimeCount int

	for _, project := range projects {
		// Count by status
		switch project.Status {
		case "completed":
			stats.CompletedProjects++
		case "failed":
			stats.FailedProjects++
		}

		// Count languages and frameworks
		if project.Requirements != nil {
			stats.PopularLanguages[project.Requirements.Language]++
			stats.PopularFrameworks[project.Requirements.Framework]++
		}

		// Calculate averages
		if project.TestResults != nil {
			if project.TestResults.Coverage > 0 {
				totalCoverage += project.TestResults.Coverage
				coverageCount++
			}
			
			buildTimeCount++
			totalBuildTime += project.TestResults.Duration.Seconds()
		}

		// Recent activity (last 10 projects)
		if len(stats.RecentActivity) < 10 {
			stats.RecentActivity = append(stats.RecentActivity, *project)
		}
	}

	// Calculate averages
	if coverageCount > 0 {
		stats.AvgTestCoverage = totalCoverage / float64(coverageCount)
	}
	if buildTimeCount > 0 {
		stats.AvgBuildTime = totalBuildTime / float64(buildTimeCount)
	}

	return stats, nil
}

// Cleanup removes old data based on age
func (fs *FileStorage) Cleanup(olderThan time.Duration) error {
	cutoffTime := time.Now().Add(-olderThan)

	// Clean up old projects
	projectsDir := filepath.Join(fs.baseDir, "projects")
	if _, err := os.Stat(projectsDir); err == nil {
		files, err := os.ReadDir(projectsDir)
		if err != nil {
			return err
		}

		for _, file := range files {
			filePath := filepath.Join(projectsDir, file.Name())
			info, err := file.Info()
			if err != nil {
				continue
			}

			if info.ModTime().Before(cutoffTime) {
				os.Remove(filePath)
			}
		}
	}

	// Clean up old analysis data
	analysisDir := filepath.Join(fs.baseDir, "analysis")
	if _, err := os.Stat(analysisDir); err == nil {
		err := filepath.Walk(analysisDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && info.ModTime().Before(cutoffTime) {
				os.Remove(path)
			}

			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

