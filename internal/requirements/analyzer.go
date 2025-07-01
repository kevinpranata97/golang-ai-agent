package requirements

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// ApplicationRequirement represents the parsed requirements for an application
type ApplicationRequirement struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Type         string                 `json:"type"` // web, api, cli, etc.
	Language     string                 `json:"language"`
	Framework    string                 `json:"framework"`
	Database     string                 `json:"database"`
	Features     []string               `json:"features"`
	Entities     []Entity               `json:"entities"`
	Endpoints    []APIEndpoint          `json:"endpoints"`
	Pages        []UIPage               `json:"pages"`
	Dependencies []string               `json:"dependencies"`
	Config       map[string]interface{} `json:"config"`
}

// Entity represents a data entity in the application
type Entity struct {
	Name       string            `json:"name"`
	Fields     []EntityField     `json:"fields"`
	Relations  []EntityRelation  `json:"relations"`
	Operations []string          `json:"operations"` // CRUD operations
}

// EntityField represents a field in an entity
type EntityField struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Required   bool   `json:"required"`
	Validation string `json:"validation"`
}

// EntityRelation represents relationships between entities
type EntityRelation struct {
	Type   string `json:"type"` // one-to-one, one-to-many, many-to-many
	Target string `json:"target"`
}

// APIEndpoint represents an API endpoint
type APIEndpoint struct {
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Description string            `json:"description"`
	Parameters  []EndpointParam   `json:"parameters"`
	Response    map[string]string `json:"response"`
}

// EndpointParam represents a parameter for an API endpoint
type EndpointParam struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Source   string `json:"source"` // query, path, body
}

// UIPage represents a UI page/component
type UIPage struct {
	Name        string   `json:"name"`
	Route       string   `json:"route"`
	Description string   `json:"description"`
	Components  []string `json:"components"`
}

// RequirementAnalyzer handles the analysis of user requirements
type RequirementAnalyzer struct {
	geminiAPIKey string
	httpClient   *http.Client
}

// NewRequirementAnalyzer creates a new requirement analyzer
func NewRequirementAnalyzer(geminiAPIKey string) *RequirementAnalyzer {
	return &RequirementAnalyzer{
		geminiAPIKey: geminiAPIKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AnalyzeRequirements analyzes user requirements and returns structured application requirements
func (ra *RequirementAnalyzer) AnalyzeRequirements(userDescription string) (*ApplicationRequirement, error) {
	// First, try to use Gemini API for analysis
	if ra.geminiAPIKey != "" {
		result, err := ra.analyzeWithGemini(userDescription)
		if err == nil {
			return result, nil
		}
		fmt.Printf("Gemini API failed, falling back to rule-based analysis: %v\n", err)
	}

	// Fallback to rule-based analysis
	return ra.analyzeWithRules(userDescription)
}

// analyzeWithGemini uses Google Gemini API for requirement analysis
func (ra *RequirementAnalyzer) analyzeWithGemini(userDescription string) (*ApplicationRequirement, error) {
	prompt := fmt.Sprintf(`
Analyze the following application requirements and return a structured JSON response:

User Description: %s

Please analyze this description and return a JSON object with the following structure:
{
  "name": "application name",
  "description": "detailed description",
  "type": "web|api|cli|desktop",
  "language": "go|javascript|python|java",
  "framework": "gin|echo|react|vue|flask|spring",
  "database": "postgresql|mysql|sqlite|mongodb",
  "features": ["list of main features"],
  "entities": [
    {
      "name": "entity name",
      "fields": [
        {
          "name": "field name",
          "type": "string|int|bool|date|email",
          "required": true|false,
          "validation": "validation rules"
        }
      ],
      "relations": [
        {
          "type": "one-to-one|one-to-many|many-to-many",
          "target": "related entity name"
        }
      ],
      "operations": ["create", "read", "update", "delete"]
    }
  ],
  "endpoints": [
    {
      "method": "GET|POST|PUT|DELETE",
      "path": "/api/path",
      "description": "endpoint description",
      "parameters": [
        {
          "name": "param name",
          "type": "string|int|bool",
          "required": true|false,
          "source": "query|path|body"
        }
      ],
      "response": {"field": "type"}
    }
  ],
  "pages": [
    {
      "name": "page name",
      "route": "/route",
      "description": "page description",
      "components": ["list of components"]
    }
  ],
  "dependencies": ["list of required dependencies"],
  "config": {
    "port": 8080,
    "other_config": "values"
  }
}

Focus on extracting entities, relationships, and required functionality. Make reasonable assumptions for missing details.
`, userDescription)

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.1,
			"maxOutputTokens": 2048,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=%s", ra.geminiAPIKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := ra.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	responseText := geminiResp.Candidates[0].Content.Parts[0].Text
	
	// Extract JSON from the response (it might be wrapped in markdown)
	jsonStart := strings.Index(responseText, "{")
	jsonEnd := strings.LastIndex(responseText, "}")
	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("no JSON found in response")
	}

	jsonStr := responseText[jsonStart : jsonEnd+1]
	
	var appReq ApplicationRequirement
	if err := json.Unmarshal([]byte(jsonStr), &appReq); err != nil {
		return nil, fmt.Errorf("failed to unmarshal application requirements: %v", err)
	}

	return &appReq, nil
}

// analyzeWithRules provides rule-based analysis as fallback
func (ra *RequirementAnalyzer) analyzeWithRules(userDescription string) (*ApplicationRequirement, error) {
	desc := strings.ToLower(userDescription)
	
	appReq := &ApplicationRequirement{
		Name:        "Generated Application",
		Description: userDescription,
		Language:    "go", // default
		Framework:   "gin", // default
		Database:    "sqlite",
		Features:    []string{},
		Entities:    []Entity{},
		Endpoints:   []APIEndpoint{},
		Pages:       []UIPage{},
		Dependencies: []string{},
		Config: map[string]interface{}{
			"port": 8080,
		},
	}

	// Detect programming language from description
	if strings.Contains(desc, "node") || strings.Contains(desc, "nodejs") || strings.Contains(desc, "node.js") || 
	   strings.Contains(desc, "javascript") || strings.Contains(desc, "js") || strings.Contains(desc, "express") {
		appReq.Language = "javascript"
		appReq.Framework = "express"
		appReq.Dependencies = []string{"express", "cors", "helmet", "morgan"}
	} else if strings.Contains(desc, "python") || strings.Contains(desc, "flask") || strings.Contains(desc, "django") || strings.Contains(desc, "fastapi") {
		appReq.Language = "python"
		if strings.Contains(desc, "django") {
			appReq.Framework = "django"
			appReq.Dependencies = []string{"django", "djangorestframework", "django-cors-headers"}
		} else if strings.Contains(desc, "fastapi") {
			appReq.Framework = "fastapi"
			appReq.Dependencies = []string{"fastapi", "uvicorn", "pydantic", "sqlalchemy"}
		} else {
			appReq.Framework = "flask"
			appReq.Dependencies = []string{"flask", "flask-cors", "flask-sqlalchemy", "flask-migrate"}
		}
	} else if strings.Contains(desc, "java") || strings.Contains(desc, "spring") || strings.Contains(desc, "springboot") {
		appReq.Language = "java"
		appReq.Framework = "spring"
		appReq.Dependencies = []string{"spring-boot-starter-web", "spring-boot-starter-data-jpa", "spring-boot-starter-security"}
	} else if strings.Contains(desc, "php") || strings.Contains(desc, "laravel") || strings.Contains(desc, "symfony") {
		appReq.Language = "php"
		if strings.Contains(desc, "laravel") {
			appReq.Framework = "laravel"
			appReq.Dependencies = []string{"laravel/framework", "laravel/sanctum", "laravel/tinker"}
		} else {
			appReq.Framework = "symfony"
			appReq.Dependencies = []string{"symfony/framework-bundle", "symfony/console", "symfony/dotenv"}
		}
	} else if strings.Contains(desc, "ruby") || strings.Contains(desc, "rails") || strings.Contains(desc, "sinatra") {
		appReq.Language = "ruby"
		if strings.Contains(desc, "rails") {
			appReq.Framework = "rails"
			appReq.Dependencies = []string{"rails", "pg", "puma", "bootsnap"}
		} else {
			appReq.Framework = "sinatra"
			appReq.Dependencies = []string{"sinatra", "sinatra-contrib", "rack-cors"}
		}
	} else if strings.Contains(desc, "go") || strings.Contains(desc, "golang") || strings.Contains(desc, "gin") || strings.Contains(desc, "echo") {
		appReq.Language = "go"
		if strings.Contains(desc, "echo") {
			appReq.Framework = "echo"
			appReq.Dependencies = []string{"github.com/labstack/echo/v4", "github.com/labstack/echo/v4/middleware"}
		} else if strings.Contains(desc, "fiber") {
			appReq.Framework = "fiber"
			appReq.Dependencies = []string{"github.com/gofiber/fiber/v2", "github.com/gofiber/fiber/v2/middleware/cors"}
		} else {
			appReq.Framework = "gin"
			appReq.Dependencies = []string{"github.com/gin-gonic/gin", "github.com/gin-contrib/cors"}
		}
	}

	// Determine application type
	if strings.Contains(desc, "web") || strings.Contains(desc, "website") || strings.Contains(desc, "frontend") {
		appReq.Type = "web"
	} else if strings.Contains(desc, "api") || strings.Contains(desc, "rest") || strings.Contains(desc, "service") {
		appReq.Type = "api"
	} else if strings.Contains(desc, "cli") || strings.Contains(desc, "command") {
		appReq.Type = "cli"
	} else {
		appReq.Type = "api" // default to API for most cases
	}

	// Extract common entities
	if strings.Contains(desc, "user") || strings.Contains(desc, "account") || strings.Contains(desc, "login") {
		userEntity := Entity{
			Name: "User",
			Fields: []EntityField{
				{Name: "id", Type: "int", Required: true},
				{Name: "username", Type: "string", Required: true, Validation: "min=3,max=50"},
				{Name: "email", Type: "email", Required: true},
				{Name: "password", Type: "string", Required: true, Validation: "min=8"},
				{Name: "created_at", Type: "date", Required: true},
			},
			Operations: []string{"create", "read", "update", "delete"},
		}
		appReq.Entities = append(appReq.Entities, userEntity)
		appReq.Features = append(appReq.Features, "user_management", "authentication")
	}

	if strings.Contains(desc, "product") || strings.Contains(desc, "item") || strings.Contains(desc, "catalog") {
		productEntity := Entity{
			Name: "Product",
			Fields: []EntityField{
				{Name: "id", Type: "int", Required: true},
				{Name: "name", Type: "string", Required: true, Validation: "min=1,max=200"},
				{Name: "description", Type: "string", Required: false},
				{Name: "price", Type: "float", Required: true, Validation: "min=0"},
				{Name: "created_at", Type: "date", Required: true},
			},
			Operations: []string{"create", "read", "update", "delete"},
		}
		appReq.Entities = append(appReq.Entities, productEntity)
		appReq.Features = append(appReq.Features, "product_management")
	}

	if strings.Contains(desc, "blog") || strings.Contains(desc, "post") || strings.Contains(desc, "article") {
		postEntity := Entity{
			Name: "Post",
			Fields: []EntityField{
				{Name: "id", Type: "int", Required: true},
				{Name: "title", Type: "string", Required: true, Validation: "min=1,max=200"},
				{Name: "content", Type: "string", Required: true},
				{Name: "author_id", Type: "int", Required: true},
				{Name: "published", Type: "bool", Required: true},
				{Name: "created_at", Type: "date", Required: true},
			},
			Relations: []EntityRelation{
				{Type: "many-to-one", Target: "User"},
			},
			Operations: []string{"create", "read", "update", "delete"},
		}
		appReq.Entities = append(appReq.Entities, postEntity)
		appReq.Features = append(appReq.Features, "content_management", "blog")
	}

	// Generate basic CRUD endpoints for each entity
	for _, entity := range appReq.Entities {
		entityLower := strings.ToLower(entity.Name)
		
		// GET all
		appReq.Endpoints = append(appReq.Endpoints, APIEndpoint{
			Method:      "GET",
			Path:        fmt.Sprintf("/api/%ss", entityLower),
			Description: fmt.Sprintf("Get all %ss", entityLower),
			Response:    map[string]string{"data": fmt.Sprintf("[]%s", entity.Name)},
		})

		// GET by ID
		appReq.Endpoints = append(appReq.Endpoints, APIEndpoint{
			Method:      "GET",
			Path:        fmt.Sprintf("/api/%ss/{id}", entityLower),
			Description: fmt.Sprintf("Get %s by ID", entityLower),
			Parameters: []EndpointParam{
				{Name: "id", Type: "int", Required: true, Source: "path"},
			},
			Response: map[string]string{"data": entity.Name},
		})

		// POST create
		appReq.Endpoints = append(appReq.Endpoints, APIEndpoint{
			Method:      "POST",
			Path:        fmt.Sprintf("/api/%ss", entityLower),
			Description: fmt.Sprintf("Create new %s", entityLower),
			Parameters: []EndpointParam{
				{Name: "body", Type: entity.Name, Required: true, Source: "body"},
			},
			Response: map[string]string{"data": entity.Name},
		})

		// PUT update
		appReq.Endpoints = append(appReq.Endpoints, APIEndpoint{
			Method:      "PUT",
			Path:        fmt.Sprintf("/api/%ss/{id}", entityLower),
			Description: fmt.Sprintf("Update %s", entityLower),
			Parameters: []EndpointParam{
				{Name: "id", Type: "int", Required: true, Source: "path"},
				{Name: "body", Type: "string", Required: true, Source: "body"},
			},
			Response: map[string]string{"data": entity.Name},
		})

		// DELETE
		appReq.Endpoints = append(appReq.Endpoints, APIEndpoint{
			Method:      "DELETE",
			Path:        fmt.Sprintf("/api/%ss/{id}", entityLower),
			Description: fmt.Sprintf("Delete %s", entityLower),
			Parameters: []EndpointParam{
				{Name: "id", Type: "int", Required: true, Source: "path"},
			},
			Response: map[string]string{"message": "string"},
		})
	}

	// Add basic pages if it's a web application
	if appReq.Type == "web" {
		appReq.Pages = append(appReq.Pages, UIPage{
			Name:        "Home",
			Route:       "/",
			Description: "Home page",
			Components:  []string{"Header", "Navigation", "Content", "Footer"},
		})

		for _, entity := range appReq.Entities {
			entityLower := strings.ToLower(entity.Name)
			appReq.Pages = append(appReq.Pages, UIPage{
				Name:        fmt.Sprintf("%s List", entity.Name),
				Route:       fmt.Sprintf("/%ss", entityLower),
				Description: fmt.Sprintf("List all %ss", entityLower),
				Components:  []string{"Header", "Navigation", fmt.Sprintf("%sList", entity.Name), "Footer"},
			})
		}
	}

	return appReq, nil
}

// ValidateRequirements validates the parsed requirements
func (ra *RequirementAnalyzer) ValidateRequirements(appReq *ApplicationRequirement) error {
	if appReq.Name == "" {
		return fmt.Errorf("application name is required")
	}

	if appReq.Type == "" {
		return fmt.Errorf("application type is required")
	}

	if appReq.Language == "" {
		return fmt.Errorf("programming language is required")
	}

	// Validate entities
	for _, entity := range appReq.Entities {
		if entity.Name == "" {
			return fmt.Errorf("entity name is required")
		}
		if len(entity.Fields) == 0 {
			return fmt.Errorf("entity %s must have at least one field", entity.Name)
		}
	}

	return nil
}

// GetGeminiAPIKey gets the Gemini API key from environment
func GetGeminiAPIKey() string {
	return os.Getenv("GEMINI_API_KEY")
}

