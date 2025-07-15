package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/kevinpranata97/golang-ai-agent/internal/requirements"
)// CodeGenerator generates application code based on requirements
type CodeGenerator struct{}

// NewCodeGenerator creates a new CodeGenerator instance
func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

// GenerateApplication generates the application code
func (cg *CodeGenerator) GenerateApplication(appReq *requirements.ApplicationRequirement) (string, error) {
	appDir := filepath.Join("generated_apps", strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")))
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create application directory: %v", err)
	}

	switch appReq.Language {
	case "Go":
		switch appReq.Type {
		case "api":
			if err := cg.generateGoAPIApplication(appDir, appReq); err != nil {
				return "", fmt.Errorf("failed to generate Go API application: %v", err)
			}
		case "web":
			return "", fmt.Errorf("Go web application generation not yet implemented")
		case "cli":
			return "", fmt.Errorf("Go CLI application generation not yet implemented")
		default:
			return "", fmt.Errorf("unsupported Go application type: %s", appReq.Type)
		}
	case "JavaScript":
		switch appReq.Type {
		case "api":
			if err := cg.generateJavaScriptAPIApplication(appDir, appReq); err != nil {
				return "", fmt.Errorf("failed to generate JavaScript API application: %v", err)
			}
		case "web":
			if err := cg.generateJavaScriptWebApplication(appDir, appReq); err != nil {
				return "", fmt.Errorf("failed to generate JavaScript web application: %v", err)
			}
		default:
			return "", fmt.Errorf("unsupported JavaScript application type: %s", appReq.Type)
		}
	case "Python":
		return "", fmt.Errorf("Python application generation not yet implemented")
	case "Java":
		return "", fmt.Errorf("Java application generation not yet implemented")
	case "PHP":
		return "", fmt.Errorf("PHP application generation not yet implemented")
	case "Ruby":
		return "", fmt.Errorf("Ruby application generation not yet implemented")
	default:
		return "", fmt.Errorf("unsupported language: %s", appReq.Language)
	}

	return appDir, nil
}

// generateGoAPIApplication generates a REST API application in Go
func (cg *CodeGenerator) generateGoAPIApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// Generate main.go
	if err := cg.generateGoMain(appDir, appReq); err != nil {
		return err
	}

	// Generate models
	if err := cg.generateModels(appDir, appReq); err != nil {
		return err
	}

	// Generate handlers
	if err := cg.generateHandlers(appDir, appReq); err != nil {
		return err
	}

	// Generate database setup
	if err := cg.generateDatabase(appDir, appReq); err != nil {
		return err
	}

	// Generate routes
	if err := cg.generateRoutes(appDir, appReq); err != nil {
		return err
	}

	// Generate config
	if err := cg.generateConfig(appDir, appReq); err != nil {
		return err
	}

	// Generate Dockerfile
	if err := cg.generateDockerfile(appDir, appReq); err != nil {
		return err
	}

	// Generate README
	if err := cg.generateReadme(appDir, appReq); err != nil {
		return err
	}

	return nil
}

// generateGoMain generates the main.go file for a Go application
func (cg *CodeGenerator) generateGoMain(appDir string, appReq *requirements.ApplicationRequirement) error {
	mainTemplate := `package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"{{.ModuleName}}/internal/config"
	"{{.ModuleName}}/internal/database"
	"{{.ModuleName}}/internal/handlers"
	"{{.ModuleName}}/internal/routes"
)

func main() {
	cfg := config.Load()

	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	defer db.Close()

	h := handlers.New(db)

	r := gin.Default()
	routes.Setup(r, h)

	port := fmt.Sprintf(":%v", cfg.Port)
	log.Printf("Server listening on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
`

	tmpl, err := template.New("main").Parse(mainTemplate)
	if err != nil {
		return err
	}

	data := struct {
		ModuleName string
	}{
		ModuleName: strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")),
	}

	file, err := os.Create(filepath.Join(appDir, "main.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateModels generates model files
func (cg *CodeGenerator) generateModels(appDir string, appReq *requirements.ApplicationRequirement) error {
	modelsDir := filepath.Join(appDir, "internal", "models")
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return err
	}

	for _, entity := range appReq.Entities {
		if err := cg.generateModelFile(modelsDir, entity); err != nil {
			return err
		}
	}

	return nil
}

// generateModelFile generates a model file for a specific entity
func (cg *CodeGenerator) generateModelFile(modelsDir string, entity requirements.Entity) error {
	modelTemplate := `package models\n\n` +
		`import (\n` +
		`\t"database/sql"\n` +
		`\t"time"\n` +
		`)\n\n` +
		`// {{.Name}} represents the {{.Name}} entity\n` +
		`type {{.Name}} struct {\n` +
		`\tID int ` + "`json:\"id\"`" + `\n` +
		`{{range .Fields}}\t{{.GoName}} {{.GoType}} ` + "`json:\"{{.JSONName}}\"{{if .Required}} validate:\"required\"{{end}}`" + `\n` +
		`{{end}}}\n\n` +
		`// Create{{.Name}} creates a new {{.Name}} in the database\n` +
		`func Create{{.Name}}(db *sql.DB, {{.LowerName}} *{{.Name}}) error {\n` +
		`\tquery := ` + "`INSERT INTO {{.TableName}} ({{.InsertFields}}) VALUES ({{.InsertPlaceholders}})`" + `\n` +
		`\tresult, err := db.Exec(query, {{range .InsertValues}} {{.}},{{end}})\n` +
		`\tif err != nil {\n` +
		`\t\treturn err\n` +
		`\t}\n` +
		`\tid, err := result.LastInsertId()\n` +
		`\tif err != nil {\n` +
		`\t\treturn err\n` +
		`\t}\n` +
		`\t{{$.LowerName}}.ID = int(id)\n` +
		`\treturn nil\n` +
		`}\n\n` +
		`// Get{{.Name}}ByID retrieves a {{.Name}} by ID\n` +
		`func Get{{.Name}}ByID(db *sql.DB, id int) (*{{.Name}}, error) {\n` +
		`\t{{.LowerName}} := &{{.Name}}{}\n` +
		`\tquery := ` + "`SELECT {{.SelectFields}} FROM {{.TableName}} WHERE id = ?`" + `\n` +
		`\terr := db.QueryRow(query, id).Scan({{range .ScanFields}}&{{$.LowerName}}.{{.}},{{end}})\n` +
		`\tif err != nil {\n` +
		`\t\treturn nil, err\n` +
		`\t}\n` +
		`\treturn {{.LowerName}}, nil\n` +
		`}\n\n` +
		`// GetAll{{.Name}}s retrieves all {{.LowerName}}s\n` +
		`func GetAll{{.Name}}s(db *sql.DB) ([]*{{.Name}}, error) {\n` +
		`\tquery := ` + "`SELECT {{.SelectFields}} FROM {{.TableName}}`" + `\n` +
		`\trows, err := db.Query(query)\n` +
		`\tif err != nil {\n` +
		`\t\treturn nil, err\n` +
		`\t}\n` +
		`\tdefer rows.Close()\n\n` +
		`\tvar {{.LowerName}}s []*{{.Name}}\n` +
		`\tfor rows.Next() {\n` +
		`\t\t{{.LowerName}} := &{{.Name}}{}\n` +
	`\t\tif err := rows.Scan({{range .ScanFields}}&{{$.LowerName}}.{{.}},{{end}}); err != nil {\n` +
		`\t\t\treturn nil, err\n` +
		`\t\t}\n` +
		`\t\t{{.LowerName}}s = append({{.LowerName}}s, {{.LowerName}})\n` +
		`\t}\n\n` +
		`\treturn {{.LowerName}}s, nil\n` +
		`}\n\n` +
		`// Update{{.Name}} updates a {{.LowerName}} in the database\n` +
		`func Update{{.Name}}(db *sql.DB, {{.LowerName}} *{{.Name}}) error {\n` +
		`\tquery := ` + "`UPDATE {{.TableName}} SET {{.UpdateFields}} WHERE id = ?`" + `\n` +
		`\t_, err := db.Exec(query, {{range .UpdateValues}} {{.}},{{end}} {{.LowerName}}.ID)\n` +
		`\treturn err\n` +
		`}\n\n` +
		`// Delete{{.Name}} deletes a {{.Name}} from the database\n` +
		`func Delete{{.Name}}(db *sql.DB, id int) error {\n` +
		`\tquery := ` + "`DELETE FROM {{.TableName}} WHERE id = ?`" + `\n` +
		`\t_, err := db.Exec(query, id)\n` +
		`\treturn err\n` +
		`}\n`

	// Prepare template data
	data := cg.prepareModelData(entity)

	tmpl, err := template.New("model").Parse(modelTemplate)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s.go", strings.ToLower(entity.Name))
	file, err := os.Create(filepath.Join(modelsDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// prepareModelData prepares template data for model generation
func (cg *CodeGenerator) prepareModelData(entity requirements.Entity) map[string]interface{} {
	data := map[string]interface{}{
		"Name":      entity.Name,
		"LowerName": strings.ToLower(entity.Name),
		"TableName": strings.ToLower(entity.Name) + "s",
	}

	var fields []map[string]interface{}
	var insertFields []string
	var insertPlaceholders []string
	var insertValues []string
	var selectFields []string
	var scanFields []string
	var updateFields []string
	var updateValues []string

	// Fix template execution issue by ensuring all fields are properly set
	for _, field := range entity.Fields {
		goType := cg.mapFieldTypeToGo(field.Type)
		goName := strings.Title(field.Name)
		jsonName := strings.ToLower(field.Name)

		fields = append(fields, map[string]interface{}{
			"GoName":   goName,
			"GoType":   goType,
			"JSONName": jsonName,
			"Required": field.Required,
		})

		if field.Name != "id" && field.Name != "created_at" {
			insertFields = append(insertFields, field.Name)
			insertPlaceholders = append(insertPlaceholders, "?")
			insertValues = append(insertValues, goName)
			updateFields = append(updateFields, field.Name+" = ?")
			updateValues = append(updateValues, goName)
		}

		selectFields = append(selectFields, field.Name)
		scanFields = append(scanFields, goName)
	}

	data["Fields"] = fields
	data["InsertFields"] = strings.Join(insertFields, ", ")
	data["InsertPlaceholders"] = strings.Join(insertPlaceholders, ", ")
	data["InsertValues"] = insertValues
	data["SelectFields"] = strings.Join(selectFields, ", ")
	data["ScanFields"] = scanFields
	data["UpdateFields"] = strings.Join(updateFields, ", ")
	data["UpdateValues"] = updateValues

	return data
}

// mapFieldTypeToGo maps field types to Go types
func (cg *CodeGenerator) mapFieldTypeToGo(fieldType string) string {
	switch fieldType {
	case "string", "email":
		return "string"
	case "int":
		return "int"
	case "float":
		return "float64"
	case "bool":
		return "bool"
	case "date":
		return "time.Time"
	default:
		return "string"
	}
}

// generateHandlers generates handler files
func (cg *CodeGenerator) generateHandlers(appDir string, appReq *requirements.ApplicationRequirement) error {
	handlersDir := filepath.Join(appDir, "internal", "handlers")
	if err := os.MkdirAll(handlersDir, 0755); err != nil {
		return err
	}

	// Generate base handler
	if err := cg.generateBaseHandler(handlersDir); err != nil {
		return err
	}

	// Generate handlers for each entity
	for _, entity := range appReq.Entities {
		if err := cg.generateEntityHandler(handlersDir, entity, appReq.Name); err != nil {
			return err
		}
	}

	return nil
}

// generateBaseHandler generates the base handler file
func (cg *CodeGenerator) generateBaseHandler(handlersDir string) error {
	handlerTemplate := `package handlers\n\n` +
		`import (\n` +
		`\t"database/sql"\n` +
		`)\n\n` +
		`// Handler contains the database connection and other dependencies\n` +
		`type Handler struct {\n` +
		`\tDB *sql.DB\n` +
		`}\n\n` +
		`// New creates a new handler instance\n` +
		`func New(db *sql.DB) *Handler {\n` +
		`\treturn &Handler{\n` +
		`\t\tDB: db,\n` +
		`\t}\n` +
		`}\n\n` +
		`// ErrorResponse represents an error response\n` +
		`type ErrorResponse struct {\n` +
		`\tError string ` + "`json:\"error\"`" + `\n` +
		`}\n\n` +
		`// SuccessResponse represents a success response\n` +
		`type SuccessResponse struct {\n` +
		`\tMessage string      ` + "`json:\"message\"`" + `\n` +
		`\tData    interface{} ` + "`json:\"data,omitempty\"`" + `\n` +
		`}\n`

	file, err := os.Create(filepath.Join(handlersDir, "handler.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(handlerTemplate)
	return err
}

// generateEntityHandler generates a handler file for a specific entity
func (cg *CodeGenerator) generateEntityHandler(handlersDir string, entity requirements.Entity, appName string) error {
	controllerTemplate := `package handlers\n\n` +
		`import (\n` +
		`\t"net/http"\n` +
		`\t"strconv"\n\n` +
		`\t"github.com/gin-gonic/gin"\n` +
		`\t"{{.ModuleName}}/internal/models"\n` +
		`)\n\n` +
		`// Create{{.Name}} creates a new {{.LowerName}}\n` +
		`func (h *Handler) Create{{.Name}}(c *gin.Context) {\n` +
		`\tvar {{.LowerName}} models.{{.Name}}\n` +
		`\tif err := c.ShouldBindJSON(&{{.LowerName}}); err != nil {\n` +
		`\t\tc.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})\n` +
		`\t\treturn\n` +
		`\t}\n\n` +
		`\tif err := models.Create{{.Name}}(h.DB, &{{.LowerName}}); err != nil {\n` +
		`\t\tc.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})\n` +
		`\t\treturn\n` +
		`\t}\n\n` +
		`\tc.JSON(http.StatusCreated, SuccessResponse{Message: "{{.Name}} created successfully", Data: {{.LowerName}}})\n` +
		`}\n\n` +
		`// Get{{.Name}} retrieves a {{.LowerName}} by ID\n` +
		`func (h *Handler) Get{{.Name}}(c *gin.Context) {\n` +
		`\tid, err := strconv.Atoi(c.Param("id"))\n` +
		`\tif err != nil {\n` +
		`\t\tc.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})\n` +
		`\t\treturn\n` +
		`\t}\n\n` +
		`\t{{.LowerName}}, err := models.Get{{.Name}}ByID(h.DB, id)\n` +
		`\tif err != nil {\n` +
		`\t\tc.JSON(http.StatusNotFound, ErrorResponse{Error: "{{.Name}} not found"})\n` +
		`\t\treturn\n` +
		`\t}\n\n` +
		`\tc.JSON(http.StatusOK, SuccessResponse{Data: {{.LowerName}}})\n` +
		`}\n\n` +
		`// GetAll{{.Name}}s retrieves all {{.LowerName}}s\n` +
		`func (h *Handler) GetAll{{.Name}}s(c *gin.Context) {\n` +
		`\t{{.LowerName}}s, err := models.GetAll{{.Name}}s(h.DB)\n` +
		`\tif err != nil {\n` +
		`\t\tc.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})\n` +
		`\t\treturn\n` +
		`\t}\n\n` +
		`\tc.JSON(http.StatusOK, SuccessResponse{Data: {{.LowerName}}s})\n` +
		`}\n\n` +
		`// Update{{.Name}} updates a {{.LowerName}} in the database\n` +
		`func (h *Handler) Update{{.Name}}(c *gin.Context) {\n` +
		`\tid, err := strconv.Atoi(c.Param("id"))\n` +
		`\tif err != nil {\n` +
		`\t\tc.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})\n` +
		`\t\treturn\n` +
		`\t}\n\n` +
		`\tvar {{.LowerName}} models.{{.Name}}\n` +
		`\tif err := c.ShouldBindJSON(&{{.LowerName}}); err != nil {\n` +
		`\t\tc.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})\n` +
		`\t\treturn\n` +
		`\t}\n\n` +
		`\t{{.LowerName}}.ID = id\n` +
		`\tif err := models.Update{{.Name}}(h.DB, &{{.LowerName}}); err != nil {\n` +
		`\t\tc.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})\n` +
		`\t\treturn\n` +
		`\t}\n\n` +
		`\tc.JSON(http.StatusOK, SuccessResponse{Message: "{{.Name}} updated successfully", Data: {{.LowerName}}})\n` +
		`}\n\n` +
		`// Delete{{.Name}} deletes a {{.Name}}\n` +
		`func (h *Handler) Delete{{.Name}}(c *gin.Context) {\n` +
		`\tid, err := strconv.Atoi(c.Param("id"))\n` +
		`\tif err != nil {\n` +
		`\t\tc.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})\n` +
		`\t\treturn\n` +
		`\t}\n\n` +
		`\tif err := models.Delete{{.Name}}(h.DB, id); err != nil {\n` +
		`\t\tc.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})\n` +
		`\t\treturn\n` +
		`\t}\n\n` +
		`\tc.JSON(http.StatusOK, SuccessResponse{Message: "{{.Name}} deleted successfully"})\n` +
		`}\n`

	tmpl, err := template.New("controller").Parse(controllerTemplate)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s.go", strings.ToLower(entity.Name)+"_handler")
	file, err := os.Create(filepath.Join(handlersDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	data := struct {
		Name       string
		LowerName  string
		ModuleName string
	}{
		Name:      entity.Name,
		LowerName: strings.ToLower(entity.Name),
		ModuleName: strings.ToLower(strings.ReplaceAll(appName, " ", "-")),
	}

	return tmpl.Execute(file, data)
}

func (cg *CodeGenerator) generateDatabase(appDir string, appReq *requirements.ApplicationRequirement) error {
	dbDir := filepath.Join(appDir, "internal", "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return err
	}

	data := struct {
		Entities []map[string]interface{}
	}{}

	for _, entity := range appReq.Entities {
		data.Entities = append(data.Entities, cg.prepareModelData(entity))
	}

	tmpl, err := template.New("database").Parse(dbTemplate)
	if err != nil {
		return err
	}

	fileName := "database.go"
	file, err := os.Create(filepath.Join(dbDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateConfig generates the config file
func (cg *CodeGenerator) generateConfig(appDir string, appReq *requirements.ApplicationRequirement) error {
	configDir := filepath.Join(appDir, "internal", "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configTemplate := `package config\n\n` +
		`import (\n` +
		`\t"log"\n` +
		`\t"os"\n\n` +
		`\t"github.com/joho/godotenv"\n` +
		`)\n\n` +
		`// Config holds the application configuration\n` +
		`type Config struct {\n` +
		`\tDatabaseURL string\n` +
		`\tPort        string\n` +
		`}\n\n` +
		`// Load loads the configuration from environment variables or .env file\n` +
		`func Load() *Config {\n` +
		`\tif err := godotenv.Load(); err != nil {\n` +
		`\t\tlog.Println("No .env file found, using environment variables")\n` +
		`\t}\n\n` +
		`\treturn &Config{\n` +
		`\t\tDatabaseURL: getEnv("DATABASE_URL", "./test.db"),\n` +
		`\t\tPort:        getEnv("PORT", "8080"),\n` +
		`}\n` +
		`}\n\n` +
		`func getEnv(key, defaultValue string) string {\n` +
		`\tvalue := os.Getenv(key)\n` +
		`\tif value == "" {\n` +
		`\t\treturn defaultValue\n` +
		`\t}\n` +
		`\treturn value\n` +
		`}\n`

	file, err := os.Create(filepath.Join(configDir, "config.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(configTemplate)
	return err
}




// generateRoutes generates the routes file
func (cg *CodeGenerator) generateRoutes(appDir string, appReq *requirements.ApplicationRequirement) error {
	routesDir := filepath.Join(appDir, "internal", "routes")
	if err := os.MkdirAll(routesDir, 0755); err != nil {
		return err
	}

	routesTemplate := `package routes\n\n` +
		`import (\n` +
		`\t"github.com/gin-gonic/gin"\n` +
		`\t"{{.ModuleName}}/internal/handlers"\n` +
		`)\n\n` +
		`// Setup sets up the routes for the application\n` +
		`func Setup(r *gin.Engine, h *handlers.Handler) {\n` +
		`\tapi := r.Group("/api")\n` +
		`\t{\n` +
		`{{range .Entities}}\t\tapi.POST("/{{.LowerName}}", h.Create{{.Name}})\n` +
		`\t\tapi.GET("/{{.LowerName}}/:id", h.Get{{.Name}})\n` +
		`\t\tapi.GET("/{{.LowerName}}s", h.GetAll{{.Name}}s)\n` +
		`\t\tapi.PUT("/{{.LowerName}}/:id", h.Update{{.Name}})\n` +
		`\t\tapi.DELETE("/{{.LowerName}}/:id", h.Delete{{.Name}})\n` +
		`{{end}}\t}\n` +
		`}\n`

	tmpl, err := template.New("routes").Parse(routesTemplate)
	if err != nil {
		return err
	}

	data := struct {
		ModuleName string
		Entities   []requirements.Entity
	}{
		ModuleName: strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")),
		Entities:   appReq.Entities,
	}

	file, err := os.Create(filepath.Join(routesDir, "routes.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}




// generateDockerfile generates the Dockerfile
func (cg *CodeGenerator) generateDockerfile(appDir string, appReq *requirements.ApplicationRequirement) error {
	dockerfileTemplate := `FROM golang:1.22-alpine AS builder\n\n` +
		`WORKDIR /app\n\n` +
		`COPY go.mod ./\n` +
		`COPY go.sum ./\n` +
		`RUN go mod download\n\n` +
		`COPY . .\n\n` +
		`RUN go build -o main main.go\n\n` +
		`FROM alpine:latest\n\n` +
		`WORKDIR /root/\n\n` +
		`COPY --from=builder /app/main .\n` +
		`COPY --from=builder /app/.env .\n\n` +
		`EXPOSE 8080\n\n` +
		`CMD ["./main"]\n`

	file, err := os.Create(filepath.Join(appDir, "Dockerfile"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(dockerfileTemplate)
	return err
}




// generateReadme generates the README.md file
func (cg *CodeGenerator) generateReadme(appDir string, appReq *requirements.ApplicationRequirement) error {
	readmeTemplate := `
# {{.Name}}

This is a generated {{.Language}} {{.Type}} application.
`

	tmpl, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(appDir, "README.md"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, appReq)
}




// generateJavaScriptAPIApplication generates a JavaScript API application
func (cg *CodeGenerator) generateJavaScriptAPIApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// Implement JavaScript API application generation logic here
	return fmt.Errorf("JavaScript API application generation not yet implemented")
}




// generateJavaScriptWebApplication generates a JavaScript web application
func (cg *CodeGenerator) generateJavaScriptWebApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// Implement JavaScript web application generation logic here
	return fmt.Errorf("JavaScript web application generation not yet implemented")
}




var dbTemplate = `package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Initialize initializes the database connection
func Initialize(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")

	// Create tables
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return db, nil
}

// createTables creates necessary tables in the database
func createTables(db *sql.DB) error {
	// {{range .Entities}}
	// Create {{.Name}} table
	create{{.Name}}TableSQL := `
	CREATE TABLE IF NOT EXISTS {{.TableName}} (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		{{range .Fields}}{{.Name}} {{.SQLType}}{{if .Required}} NOT NULL{{end}},
		{{end}}
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(create{{.Name}}TableSQL)
	if err != nil {
		return fmt.Errorf("failed to create {{.Name}} table: %v", err)
	}
	// {{end}}

	log.Println("Tables created successfully")
	return nil
}

func (cg *CodeGenerator) mapFieldTypeToSQL(fieldType string) string {
	switch fieldType {
	case "string", "email":
		return "TEXT"
	case "int":
		return "INTEGER"
	case "float":
		return "REAL"
	case "bool":
		return "BOOLEAN"
	case "date":
		return "DATETIME"
	default:
		return "TEXT"
	}
}

var dbTemplate = `package database

