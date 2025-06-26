package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/kevinpranata97/golang-ai-agent/internal/requirements"
)

// CodeGenerator handles the generation of application code
type CodeGenerator struct {
	outputDir string
	templates map[string]*template.Template
}

// NewCodeGenerator creates a new code generator
func NewCodeGenerator(outputDir string) *CodeGenerator {
	return &CodeGenerator{
		outputDir: outputDir,
		templates: make(map[string]*template.Template),
	}
}

// GenerateApplication generates a complete application based on requirements
func (cg *CodeGenerator) GenerateApplication(appReq *requirements.ApplicationRequirement) error {
	// Create output directory
	appDir := filepath.Join(cg.outputDir, strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")))
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return fmt.Errorf("failed to create app directory: %v", err)
	}

	// Generate different components based on application type
	switch appReq.Type {
	case "api":
		return cg.generateAPIApplication(appDir, appReq)
	case "web":
		return cg.generateWebApplication(appDir, appReq)
	case "cli":
		return cg.generateCLIApplication(appDir, appReq)
	default:
		return cg.generateAPIApplication(appDir, appReq) // default to API
	}
}

// generateAPIApplication generates a REST API application
func (cg *CodeGenerator) generateAPIApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// Generate main.go
	if err := cg.generateMainFile(appDir, appReq); err != nil {
		return err
	}

	// Generate go.mod
	if err := cg.generateGoMod(appDir, appReq); err != nil {
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

// generateWebApplication generates a web application with frontend and backend
func (cg *CodeGenerator) generateWebApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// First generate the API backend
	if err := cg.generateAPIApplication(appDir, appReq); err != nil {
		return err
	}

	// Generate static files directory
	staticDir := filepath.Join(appDir, "static")
	if err := os.MkdirAll(staticDir, 0755); err != nil {
		return fmt.Errorf("failed to create static directory: %v", err)
	}

	// Generate basic HTML templates
	if err := cg.generateHTMLTemplates(staticDir, appReq); err != nil {
		return err
	}

	// Generate CSS
	if err := cg.generateCSS(staticDir, appReq); err != nil {
		return err
	}

	// Generate JavaScript
	if err := cg.generateJavaScript(staticDir, appReq); err != nil {
		return err
	}

	return nil
}

// generateCLIApplication generates a CLI application
func (cg *CodeGenerator) generateCLIApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// Generate main.go for CLI
	if err := cg.generateCLIMain(appDir, appReq); err != nil {
		return err
	}

	// Generate go.mod
	if err := cg.generateGoMod(appDir, appReq); err != nil {
		return err
	}

	// Generate commands
	if err := cg.generateCLICommands(appDir, appReq); err != nil {
		return err
	}

	return nil
}

// generateMainFile generates the main.go file
func (cg *CodeGenerator) generateMainFile(appDir string, appReq *requirements.ApplicationRequirement) error {
	mainTemplate := `package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"{{.ModuleName}}/internal/config"
	"{{.ModuleName}}/internal/database"
	"{{.ModuleName}}/internal/handlers"
	"{{.ModuleName}}/internal/routes"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize Gin router
	r := gin.Default()

	// Setup CORS
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Initialize handlers
	h := handlers.New(db)

	// Setup routes
	routes.Setup(r, h)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "{{.Port}}"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, r))
}
`

	tmpl, err := template.New("main").Parse(mainTemplate)
	if err != nil {
		return err
	}

	data := struct {
		ModuleName string
		Port       string
	}{
		ModuleName: strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")),
		Port:       fmt.Sprintf("%v", appReq.Config["port"]),
	}

	file, err := os.Create(filepath.Join(appDir, "main.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateGoMod generates the go.mod file
func (cg *CodeGenerator) generateGoMod(appDir string, appReq *requirements.ApplicationRequirement) error {
	modTemplate := `module {{.ModuleName}}

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/mattn/go-sqlite3 v1.14.17
{{range .Dependencies}}	{{.}}
{{end}})
`

	tmpl, err := template.New("gomod").Parse(modTemplate)
	if err != nil {
		return err
	}

	data := struct {
		ModuleName   string
		Dependencies []string
	}{
		ModuleName:   strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")),
		Dependencies: appReq.Dependencies,
	}

	file, err := os.Create(filepath.Join(appDir, "go.mod"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateModels generates model files for each entity
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

// generateModelFile generates a single model file
func (cg *CodeGenerator) generateModelFile(modelsDir string, entity requirements.Entity) error {
	modelTemplate := `package models

import (
	"time"
	"database/sql"
)

// {{.Name}} represents the {{.Name}} entity
type {{.Name}} struct {
{{range .Fields}}	{{.GoName}} {{.GoType}} ` + "`json:\"{{.JSONName}}\"{{if .Required}} validate:\"required\"{{end}}`" + `
{{end}}}

// Create{{.Name}} creates a new {{.Name}} in the database
func Create{{.Name}}(db *sql.DB, {{.LowerName}} *{{.Name}}) error {
	query := ` + "`INSERT INTO {{.TableName}} ({{.InsertFields}}) VALUES ({{.InsertPlaceholders}})`" + `
	
	result, err := db.Exec(query{{range .InsertValues}}, {{$.LowerName}}.{{.}}{{end}})
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	{{.LowerName}}.ID = int(id)
	return nil
}

// Get{{.Name}}ByID retrieves a {{.Name}} by ID
func Get{{.Name}}ByID(db *sql.DB, id int) (*{{.Name}}, error) {
	{{.LowerName}} := &{{.Name}}{}
	query := ` + "`SELECT {{.SelectFields}} FROM {{.TableName}} WHERE id = ?`" + `
	
	err := db.QueryRow(query, id).Scan({{range .ScanFields}}&{{$.LowerName}}.{{.}}{{end}})
	if err != nil {
		return nil, err
	}

	return {{.LowerName}}, nil
}

// GetAll{{.Name}}s retrieves all {{.Name}}s
func GetAll{{.Name}}s(db *sql.DB) ([]{{.Name}}, error) {
	query := ` + "`SELECT {{.SelectFields}} FROM {{.TableName}}`" + `
	
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var {{.LowerName}}s []{{.Name}}
	for rows.Next() {
		{{.LowerName}} := {{.Name}}{}
		err := rows.Scan({{range .ScanFields}}&{{$.LowerName}}.{{.}}{{end}})
		if err != nil {
			return nil, err
		}
		{{.LowerName}}s = append({{.LowerName}}s, {{.LowerName}})
	}

	return {{.LowerName}}s, nil
}

// Update{{.Name}} updates a {{.Name}} in the database
func Update{{.Name}}(db *sql.DB, {{.LowerName}} *{{.Name}}) error {
	query := ` + "`UPDATE {{.TableName}} SET {{.UpdateFields}} WHERE id = ?`" + `
	
	_, err := db.Exec(query{{range .UpdateValues}}, {{$.LowerName}}.{{.}}{{end}}, {{.LowerName}}.ID)
	return err
}

// Delete{{.Name}} deletes a {{.Name}} from the database
func Delete{{.Name}}(db *sql.DB, id int) error {
	query := ` + "`DELETE FROM {{.TableName}} WHERE id = ?`" + `
	
	_, err := db.Exec(query, id)
	return err
}
`

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
	handlerTemplate := `package handlers

import (
	"database/sql"
)

// Handler contains the database connection and other dependencies
type Handler struct {
	DB *sql.DB
}

// New creates a new handler instance
func New(db *sql.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string ` + "`json:\"error\"`" + `
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Message string      ` + "`json:\"message\"`" + `
	Data    interface{} ` + "`json:\"data,omitempty\"`" + `
}
`

	file, err := os.Create(filepath.Join(handlersDir, "handler.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(handlerTemplate)
	return err
}

// generateEntityHandler generates handler for a specific entity
func (cg *CodeGenerator) generateEntityHandler(handlersDir string, entity requirements.Entity, appName string) error {
	handlerTemplate := `package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"{{.ModuleName}}/internal/models"
)

// Create{{.Name}} creates a new {{.Name}}
func (h *Handler) Create{{.Name}}(c *gin.Context) {
	var {{.LowerName}} models.{{.Name}}
	
	if err := c.ShouldBindJSON(&{{.LowerName}}); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := models.Create{{.Name}}(h.DB, &{{.LowerName}}); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Message: "{{.Name}} created successfully",
		Data:    {{.LowerName}},
	})
}

// Get{{.Name}} retrieves a {{.Name}} by ID
func (h *Handler) Get{{.Name}}(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	{{.LowerName}}, err := models.Get{{.Name}}ByID(h.DB, id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "{{.Name}} not found"})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: {{.LowerName}}})
}

// GetAll{{.Name}}s retrieves all {{.Name}}s
func (h *Handler) GetAll{{.Name}}s(c *gin.Context) {
	{{.LowerName}}s, err := models.GetAll{{.Name}}s(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: {{.LowerName}}s})
}

// Update{{.Name}} updates a {{.Name}}
func (h *Handler) Update{{.Name}}(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	var {{.LowerName}} models.{{.Name}}
	if err := c.ShouldBindJSON(&{{.LowerName}}); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	{{.LowerName}}.ID = id
	if err := models.Update{{.Name}}(h.DB, &{{.LowerName}}); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "{{.Name}} updated successfully",
		Data:    {{.LowerName}},
	})
}

// Delete{{.Name}} deletes a {{.Name}}
func (h *Handler) Delete{{.Name}}(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	if err := models.Delete{{.Name}}(h.DB, id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Message: "{{.Name}} deleted successfully"})
}
`

	data := map[string]interface{}{
		"Name":       entity.Name,
		"LowerName":  strings.ToLower(entity.Name),
		"ModuleName": strings.ToLower(strings.ReplaceAll(appName, " ", "-")),
	}

	tmpl, err := template.New("handler").Parse(handlerTemplate)
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s_handler.go", strings.ToLower(entity.Name))
	file, err := os.Create(filepath.Join(handlersDir, fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateDatabase generates database setup files
func (cg *CodeGenerator) generateDatabase(appDir string, appReq *requirements.ApplicationRequirement) error {
	dbDir := filepath.Join(appDir, "internal", "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return err
	}

	// Generate database initialization
	if err := cg.generateDatabaseInit(dbDir, appReq); err != nil {
		return err
	}

	// Generate migrations
	if err := cg.generateMigrations(dbDir, appReq); err != nil {
		return err
	}

	return nil
}

// generateDatabaseInit generates database initialization file
func (cg *CodeGenerator) generateDatabaseInit(dbDir string, appReq *requirements.ApplicationRequirement) error {
	dbTemplate := `package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Initialize initializes the database connection and runs migrations
func Initialize(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		databaseURL = "./app.db"
	}

	db, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}

// runMigrations runs database migrations
func runMigrations(db *sql.DB) error {
	migrations := []string{
{{range .Migrations}}		` + "`{{.}}`" + `,
{{end}}	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("failed to execute migration: %v", err)
		}
	}

	return nil
}
`

	var migrations []string
	for _, entity := range appReq.Entities {
		migration := cg.generateCreateTableSQL(entity)
		migrations = append(migrations, migration)
	}

	data := map[string]interface{}{
		"Migrations": migrations,
	}

	tmpl, err := template.New("database").Parse(dbTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(dbDir, "database.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateCreateTableSQL generates CREATE TABLE SQL for an entity
func (cg *CodeGenerator) generateCreateTableSQL(entity requirements.Entity) string {
	tableName := strings.ToLower(entity.Name) + "s"
	var fields []string

	for _, field := range entity.Fields {
		sqlType := cg.mapFieldTypeToSQL(field.Type)
		fieldDef := fmt.Sprintf("%s %s", field.Name, sqlType)
		
		if field.Name == "id" {
			fieldDef += " PRIMARY KEY AUTOINCREMENT"
		} else if field.Required {
			fieldDef += " NOT NULL"
		}

		fields = append(fields, fieldDef)
	}

	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(fields, ", "))
}

// mapFieldTypeToSQL maps field types to SQL types
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
		return "DATETIME DEFAULT CURRENT_TIMESTAMP"
	default:
		return "TEXT"
	}
}

// generateMigrations generates migration files
func (cg *CodeGenerator) generateMigrations(dbDir string, appReq *requirements.ApplicationRequirement) error {
	// For now, we'll keep it simple and not generate separate migration files
	// In a more complex system, we would generate individual migration files
	return nil
}

// generateRoutes generates route setup
func (cg *CodeGenerator) generateRoutes(appDir string, appReq *requirements.ApplicationRequirement) error {
	routesDir := filepath.Join(appDir, "internal", "routes")
	if err := os.MkdirAll(routesDir, 0755); err != nil {
		return err
	}

	routesTemplate := `package routes

import (
	"github.com/gin-gonic/gin"
	"{{.ModuleName}}/internal/handlers"
)

// Setup configures all routes
func Setup(r *gin.Engine, h *handlers.Handler) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
{{range .Entities}}		// {{.Name}} routes
		api.GET("/{{.LowerPlural}}", h.GetAll{{.Name}}s)
		api.GET("/{{.LowerPlural}}/:id", h.Get{{.Name}})
		api.POST("/{{.LowerPlural}}", h.Create{{.Name}})
		api.PUT("/{{.LowerPlural}}/:id", h.Update{{.Name}})
		api.DELETE("/{{.LowerPlural}}/:id", h.Delete{{.Name}})

{{end}}	}
}
`

	var entities []map[string]interface{}
	for _, entity := range appReq.Entities {
		entities = append(entities, map[string]interface{}{
			"Name":        entity.Name,
			"LowerPlural": strings.ToLower(entity.Name) + "s",
		})
	}

	data := map[string]interface{}{
		"ModuleName": strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")),
		"Entities":   entities,
	}

	tmpl, err := template.New("routes").Parse(routesTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(routesDir, "routes.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateConfig generates configuration files
func (cg *CodeGenerator) generateConfig(appDir string, appReq *requirements.ApplicationRequirement) error {
	configDir := filepath.Join(appDir, "internal", "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configTemplate := `package config

import (
	"os"
)

// Config holds application configuration
type Config struct {
	Port        string
	DatabaseURL string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "{{.Port}}"),
		DatabaseURL: getEnv("DATABASE_URL", "{{.DatabaseURL}}"),
	}
}

// getEnv gets an environment variable with a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
`

	data := map[string]interface{}{
		"Port":        fmt.Sprintf("%v", appReq.Config["port"]),
		"DatabaseURL": "./app.db",
	}

	tmpl, err := template.New("config").Parse(configTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(configDir, "config.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateDockerfile generates Dockerfile
func (cg *CodeGenerator) generateDockerfile(appDir string, appReq *requirements.ApplicationRequirement) error {
	dockerfileTemplate := `# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE {{.Port}}

# Run the application
CMD ["./main"]
`

	data := map[string]interface{}{
		"Port": fmt.Sprintf("%v", appReq.Config["port"]),
	}

	tmpl, err := template.New("dockerfile").Parse(dockerfileTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(appDir, "Dockerfile"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateReadme generates README.md
func (cg *CodeGenerator) generateReadme(appDir string, appReq *requirements.ApplicationRequirement) error {
	readmeTemplate := `# {{.Name}}

{{.Description}}

## Features

{{range .Features}}- {{.}}
{{end}}

## API Endpoints

{{range .Endpoints}}### {{.Method}} {{.Path}}
{{.Description}}

{{if .Parameters}}**Parameters:**
{{range .Parameters}}- {{.Name}} ({{.Type}}) - {{if .Required}}Required{{else}}Optional{{end}} - {{.Source}}
{{end}}{{end}}

{{end}}

## Getting Started

### Prerequisites

- Go 1.21 or higher
- SQLite (for development)

### Installation

1. Clone the repository
2. Install dependencies:
   ` + "```bash" + `
   go mod tidy
   ` + "```" + `

3. Run the application:
   ` + "```bash" + `
   go run main.go
   ` + "```" + `

The server will start on port {{.Port}}.

### Docker

Build and run with Docker:

` + "```bash" + `
docker build -t {{.DockerName}} .
docker run -p {{.Port}}:{{.Port}} {{.DockerName}}
` + "```" + `

## Configuration

Environment variables:

- ` + "`PORT`" + ` - Server port (default: {{.Port}})
- ` + "`DATABASE_URL`" + ` - Database connection string (default: ./app.db)

## Testing

Run tests:

` + "```bash" + `
go test ./...
` + "```" + `

## License

This project is generated by Golang AI Agent.
`

	data := map[string]interface{}{
		"Name":        appReq.Name,
		"Description": appReq.Description,
		"Features":    appReq.Features,
		"Endpoints":   appReq.Endpoints,
		"Port":        fmt.Sprintf("%v", appReq.Config["port"]),
		"DockerName":  strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")),
	}

	tmpl, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(appDir, "README.md"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateHTMLTemplates generates basic HTML templates for web applications
func (cg *CodeGenerator) generateHTMLTemplates(staticDir string, appReq *requirements.ApplicationRequirement) error {
	// Create templates directory
	templatesDir := filepath.Join(staticDir, "templates")
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return err
	}

	// Generate index.html
	indexTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Name}}</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <header>
        <nav>
            <h1>{{.Name}}</h1>
            <ul>
                <li><a href="/">Home</a></li>
{{range .Pages}}                <li><a href="{{.Route}}">{{.Name}}</a></li>
{{end}}            </ul>
        </nav>
    </header>

    <main>
        <h2>Welcome to {{.Name}}</h2>
        <p>{{.Description}}</p>
        
        <div class="features">
            <h3>Features:</h3>
            <ul>
{{range .Features}}                <li>{{.}}</li>
{{end}}            </ul>
        </div>
    </main>

    <script src="/static/js/app.js"></script>
</body>
</html>
`

	data := map[string]interface{}{
		"Name":        appReq.Name,
		"Description": appReq.Description,
		"Features":    appReq.Features,
		"Pages":       appReq.Pages,
	}

	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(templatesDir, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateCSS generates basic CSS
func (cg *CodeGenerator) generateCSS(staticDir string, appReq *requirements.ApplicationRequirement) error {
	cssDir := filepath.Join(staticDir, "css")
	if err := os.MkdirAll(cssDir, 0755); err != nil {
		return err
	}

	css := `/* Reset and base styles */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    line-height: 1.6;
    color: #333;
    background-color: #f4f4f4;
}

/* Header and navigation */
header {
    background: #2c3e50;
    color: white;
    padding: 1rem 0;
    box-shadow: 0 2px 5px rgba(0,0,0,0.1);
}

nav {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

nav h1 {
    font-size: 1.5rem;
}

nav ul {
    display: flex;
    list-style: none;
    gap: 2rem;
}

nav a {
    color: white;
    text-decoration: none;
    transition: color 0.3s;
}

nav a:hover {
    color: #3498db;
}

/* Main content */
main {
    max-width: 1200px;
    margin: 2rem auto;
    padding: 0 2rem;
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 10px rgba(0,0,0,0.1);
    padding: 2rem;
}

h2 {
    color: #2c3e50;
    margin-bottom: 1rem;
}

.features {
    margin-top: 2rem;
}

.features ul {
    list-style-type: disc;
    margin-left: 2rem;
}

.features li {
    margin-bottom: 0.5rem;
}

/* Responsive design */
@media (max-width: 768px) {
    nav {
        flex-direction: column;
        gap: 1rem;
    }
    
    nav ul {
        gap: 1rem;
    }
    
    main {
        margin: 1rem;
        padding: 1rem;
    }
}
`

	file, err := os.Create(filepath.Join(cssDir, "style.css"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(css)
	return err
}

// generateJavaScript generates basic JavaScript
func (cg *CodeGenerator) generateJavaScript(staticDir string, appReq *requirements.ApplicationRequirement) error {
	jsDir := filepath.Join(staticDir, "js")
	if err := os.MkdirAll(jsDir, 0755); err != nil {
		return err
	}

	js := `// Basic JavaScript for the application
console.log('Application loaded successfully');

// API base URL
const API_BASE = '/api';

// Utility function to make API calls
async function apiCall(endpoint, options = {}) {
    try {
        const response = await fetch(API_BASE + endpoint, {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            ...options
        });
        
        if (!response.ok) {
            throw new Error(` + "`HTTP error! status: ${response.status}`" + `);
        }
        
        return await response.json();
    } catch (error) {
        console.error('API call failed:', error);
        throw error;
    }
}

// Example functions for each entity
` + cg.generateEntityJSFunctions(appReq.Entities) + `

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM loaded, initializing application...');
    
    // Add any initialization code here
});
`

	file, err := os.Create(filepath.Join(jsDir, "app.js"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(js)
	return err
}

// generateEntityJSFunctions generates JavaScript functions for entities
func (cg *CodeGenerator) generateEntityJSFunctions(entities []requirements.Entity) string {
	var functions []string

	for _, entity := range entities {
		entityLower := strings.ToLower(entity.Name)
		entityPlural := entityLower + "s"

		functions = append(functions, fmt.Sprintf(`
// %s functions
async function getAll%ss() {
    return await apiCall('/%s');
}

async function get%s(id) {
    return await apiCall('/%s/' + id);
}

async function create%s(data) {
    return await apiCall('/%s', {
        method: 'POST',
        body: JSON.stringify(data)
    });
}

async function update%s(id, data) {
    return await apiCall('/%s/' + id, {
        method: 'PUT',
        body: JSON.stringify(data)
    });
}

async function delete%s(id) {
    return await apiCall('/%s/' + id, {
        method: 'DELETE'
    });
}`, entity.Name, entityPlural, entityPlural, entity.Name, entityPlural, entity.Name, entityPlural, entity.Name, entityPlural, entity.Name, entityPlural))
	}

	return strings.Join(functions, "\n")
}

// generateCLIMain generates main.go for CLI applications
func (cg *CodeGenerator) generateCLIMain(appDir string, appReq *requirements.ApplicationRequirement) error {
	cliTemplate := `package main

import (
	"fmt"
	"os"

	"{{.ModuleName}}/internal/commands"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: {{.AppName}} <command> [args...]")
		fmt.Println("Available commands:")
{{range .Commands}}		fmt.Println("  {{.}}")
{{end}}		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	switch command {
{{range .Commands}}	case "{{.}}":
		commands.{{.Title}}(args)
{{end}}	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
`

	var commands []string
	for _, entity := range appReq.Entities {
		entityLower := strings.ToLower(entity.Name)
		commands = append(commands, "list-"+entityLower+"s")
		commands = append(commands, "create-"+entityLower)
	}

	data := map[string]interface{}{
		"ModuleName": strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")),
		"AppName":    strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")),
		"Commands":   commands,
	}

	tmpl, err := template.New("climain").Parse(cliTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(appDir, "main.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateCLICommands generates CLI command files
func (cg *CodeGenerator) generateCLICommands(appDir string, appReq *requirements.ApplicationRequirement) error {
	commandsDir := filepath.Join(appDir, "internal", "commands")
	if err := os.MkdirAll(commandsDir, 0755); err != nil {
		return err
	}

	// Generate basic command structure
	commandTemplate := `package commands

import (
	"fmt"
)

// Example command functions
{{range .Commands}}
func {{.Function}}(args []string) {
	fmt.Println("Executing {{.Name}} command with args:", args)
	// TODO: Implement {{.Name}} logic
}
{{end}}
`

	var commands []map[string]string
	for _, entity := range appReq.Entities {
		entityLower := strings.ToLower(entity.Name)
		commands = append(commands, map[string]string{
			"Name":     "list-" + entityLower + "s",
			"Function": "List" + entity.Name + "s",
		})
		commands = append(commands, map[string]string{
			"Name":     "create-" + entityLower,
			"Function": "Create" + entity.Name,
		})
	}

	data := map[string]interface{}{
		"Commands": commands,
	}

	tmpl, err := template.New("commands").Parse(commandTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(commandsDir, "commands.go"))
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

