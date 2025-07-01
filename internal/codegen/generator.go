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

	// Generate application based on language and type
	switch appReq.Language {
	case "javascript":
		return cg.generateJavaScriptApplication(appDir, appReq)
	case "python":
		return cg.generatePythonApplication(appDir, appReq)
	case "java":
		return cg.generateJavaApplication(appDir, appReq)
	case "php":
		return cg.generatePHPApplication(appDir, appReq)
	case "ruby":
		return cg.generateRubyApplication(appDir, appReq)
	case "go":
		fallthrough
	default:
		return cg.generateGoApplication(appDir, appReq)
	}
}

// generateGoApplication generates a Go application
func (cg *CodeGenerator) generateGoApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// Generate different components based on application type
	switch appReq.Type {
	case "api":
		return cg.generateGoAPIApplication(appDir, appReq)
	case "web":
		return cg.generateGoWebApplication(appDir, appReq)
	case "cli":
		return cg.generateGoCLIApplication(appDir, appReq)
	default:
		return cg.generateGoAPIApplication(appDir, appReq) // default to API
	}
}

// generateJavaScriptApplication generates a Node.js/JavaScript application
func (cg *CodeGenerator) generateJavaScriptApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// Generate different components based on application type
	switch appReq.Type {
	case "api":
		return cg.generateJavaScriptAPIApplication(appDir, appReq)
	case "web":
		return cg.generateJavaScriptWebApplication(appDir, appReq)
	default:
		return cg.generateJavaScriptAPIApplication(appDir, appReq) // default to API
	}
}

// generatePythonApplication generates a Python application
func (cg *CodeGenerator) generatePythonApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// Generate different components based on application type
	switch appReq.Type {
	case "api":
		return cg.generatePythonAPIApplication(appDir, appReq)
	case "web":
		return cg.generatePythonWebApplication(appDir, appReq)
	default:
		return cg.generatePythonAPIApplication(appDir, appReq) // default to API
	}
}

// generateJavaApplication generates a Java application
func (cg *CodeGenerator) generateJavaApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// Generate different components based on application type
	switch appReq.Type {
	case "api":
		return cg.generateJavaAPIApplication(appDir, appReq)
	case "web":
		return cg.generateJavaWebApplication(appDir, appReq)
	default:
		return cg.generateJavaAPIApplication(appDir, appReq) // default to API
	}
}

// generatePHPApplication generates a PHP application
func (cg *CodeGenerator) generatePHPApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// Generate different components based on application type
	switch appReq.Type {
	case "api":
		return cg.generatePHPAPIApplication(appDir, appReq)
	case "web":
		return cg.generatePHPWebApplication(appDir, appReq)
	default:
		return cg.generatePHPAPIApplication(appDir, appReq) // default to API
	}
}

// generateRubyApplication generates a Ruby application
func (cg *CodeGenerator) generateRubyApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// Generate different components based on application type
	switch appReq.Type {
	case "api":
		return cg.generateRubyAPIApplication(appDir, appReq)
	case "web":
		return cg.generateRubyWebApplication(appDir, appReq)
	default:
		return cg.generateRubyAPIApplication(appDir, appReq) // default to API
	}
}

// generateGoAPIApplication generates a REST API application in Go
func (cg *CodeGenerator) generateGoAPIApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
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
	if err := cg.generateGoAPIApplication(appDir, appReq); err != nil {
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



// generateJavaScriptAPIApplication generates a REST API application in Node.js/JavaScript
func (cg *CodeGenerator) generateJavaScriptAPIApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// Generate package.json
	if err := cg.generatePackageJSON(appDir, appReq); err != nil {
		return err
	}

	// Generate main server file (app.js or server.js)
	if err := cg.generateJavaScriptMainFile(appDir, appReq); err != nil {
		return err
	}

	// Generate models
	if err := cg.generateJavaScriptModels(appDir, appReq); err != nil {
		return err
	}

	// Generate routes
	if err := cg.generateJavaScriptRoutes(appDir, appReq); err != nil {
		return err
	}

	// Generate controllers
	if err := cg.generateJavaScriptControllers(appDir, appReq); err != nil {
		return err
	}

	// Generate middleware
	if err := cg.generateJavaScriptMiddleware(appDir, appReq); err != nil {
		return err
	}

	// Generate database configuration
	if err := cg.generateJavaScriptDatabase(appDir, appReq); err != nil {
		return err
	}

	// Generate environment configuration
	if err := cg.generateJavaScriptEnvConfig(appDir, appReq); err != nil {
		return err
	}

	// Generate Dockerfile
	if err := cg.generateJavaScriptDockerfile(appDir, appReq); err != nil {
		return err
	}

	// Generate README
	if err := cg.generateJavaScriptReadme(appDir, appReq); err != nil {
		return err
	}

	return nil
}

// generatePackageJSON generates package.json for Node.js application
func (cg *CodeGenerator) generatePackageJSON(appDir string, appReq *requirements.ApplicationRequirement) error {
	packageJSON := `{
  "name": "{{.AppName}}",
  "version": "1.0.0",
  "description": "{{.Description}}",
  "main": "app.js",
  "scripts": {
    "start": "node app.js",
    "dev": "nodemon app.js",
    "test": "jest"
  },
  "dependencies": {
{{range $i, $dep := .Dependencies}}    "{{$dep}}": "latest"{{if ne $i (sub (len $.Dependencies) 1)}},{{end}}
{{end}}  },
  "devDependencies": {
    "nodemon": "^3.0.0",
    "jest": "^29.0.0"
  },
  "keywords": [
    "api",
    "{{.Framework}}",
    "rest"
  ],
  "author": "",
  "license": "MIT"
}`

	tmpl, err := template.New("package.json").Funcs(template.FuncMap{
		"sub": func(a, b int) int { return a - b },
	}).Parse(packageJSON)
	if err != nil {
		return fmt.Errorf("failed to parse package.json template: %v", err)
	}

	data := struct {
		AppName      string
		Description  string
		Framework    string
		Dependencies []string
	}{
		AppName:      strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")),
		Description:  appReq.Description,
		Framework:    appReq.Framework,
		Dependencies: appReq.Dependencies,
	}

	file, err := os.Create(filepath.Join(appDir, "package.json"))
	if err != nil {
		return fmt.Errorf("failed to create package.json: %v", err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateJavaScriptMainFile generates the main server file (app.js)
func (cg *CodeGenerator) generateJavaScriptMainFile(appDir string, appReq *requirements.ApplicationRequirement) error {
	mainFile := `const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const morgan = require('morgan');
{{if .HasDatabase}}const db = require('./config/database');{{end}}

// Import routes
{{range .Entities}}const {{.LowerName}}Routes = require('./routes/{{.LowerName}}Routes');
{{end}}

const app = express();
const PORT = process.env.PORT || {{.Port}};

// Middleware
app.use(helmet());
app.use(cors());
app.use(morgan('combined'));
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Routes
app.get('/', (req, res) => {
  res.json({
    message: 'Welcome to {{.AppName}} API',
    version: '1.0.0',
    endpoints: [
{{range .Endpoints}}      '{{.Method}} {{.Path}}',
{{end}}    ]
  });
});

{{range .Entities}}app.use('/api/{{.LowerName}}s', {{.LowerName}}Routes);
{{end}}

// Error handling middleware
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({
    error: 'Something went wrong!',
    message: err.message
  });
});

// 404 handler
app.use('*', (req, res) => {
  res.status(404).json({
    error: 'Route not found'
  });
});

{{if .HasDatabase}}// Initialize database connection
db.connect().then(() => {
  console.log('Database connected successfully');
  
  app.listen(PORT, '0.0.0.0', () => {
    console.log('Server is running on port ' + PORT);
    console.log('API Documentation: http://localhost:' + PORT);
  });
}).catch(err => {
  console.error('Database connection failed:', err);
  process.exit(1);
});{{else}}app.listen(PORT, '0.0.0.0', () => {
  console.log('Server is running on port ' + PORT);
  console.log('API Documentation: http://localhost:' + PORT);
});{{end}}`

	tmpl, err := template.New("app.js").Parse(mainFile)
	if err != nil {
		return fmt.Errorf("failed to parse app.js template: %v", err)
	}

	// Prepare entities with lowercase names
	var entities []map[string]interface{}
	for _, entity := range appReq.Entities {
		entities = append(entities, map[string]interface{}{
			"Name":      entity.Name,
			"LowerName": strings.ToLower(entity.Name),
		})
	}

	data := struct {
		AppName     string
		Port        interface{}
		HasDatabase bool
		Entities    []map[string]interface{}
		Endpoints   []requirements.APIEndpoint
	}{
		AppName:     appReq.Name,
		Port:        appReq.Config["port"],
		HasDatabase: appReq.Database != "",
		Entities:    entities,
		Endpoints:   appReq.Endpoints,
	}

	file, err := os.Create(filepath.Join(appDir, "app.js"))
	if err != nil {
		return fmt.Errorf("failed to create app.js: %v", err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateJavaScriptModels generates model files for JavaScript application
func (cg *CodeGenerator) generateJavaScriptModels(appDir string, appReq *requirements.ApplicationRequirement) error {
	modelsDir := filepath.Join(appDir, "models")
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %v", err)
	}

	for _, entity := range appReq.Entities {
		if err := cg.generateJavaScriptModel(modelsDir, entity); err != nil {
			return err
		}
	}

	return nil
}

// generateJavaScriptModel generates a single model file
func (cg *CodeGenerator) generateJavaScriptModel(modelsDir string, entity requirements.Entity) error {
	modelTemplate := `class {{.Name}} {
  constructor(data = {}) {
{{range .Fields}}    this.{{.Name}} = data.{{.Name}} || {{.DefaultValue}};
{{end}}  }

  // Validation method
  validate() {
    const errors = [];
{{range .Fields}}{{if .Required}}
    if (!this.{{.Name}}) {
      errors.push('{{.Name}} is required');
    }{{end}}{{if .Validation}}
    // Add validation for {{.Name}}: {{.Validation}}{{end}}
{{end}}
    return errors;
  }

  // Convert to JSON
  toJSON() {
    return {
{{range .Fields}}      {{.Name}}: this.{{.Name}},
{{end}}    };
  }

  // Create from database row
  static fromRow(row) {
    return new {{.Name}}(row);
  }
}

module.exports = {{.Name}};`

	tmpl, err := template.New("model").Parse(modelTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse model template: %v", err)
	}

	// Prepare fields with default values
	var fields []map[string]interface{}
	for _, field := range entity.Fields {
		defaultValue := "null"
		switch field.Type {
		case "string", "email":
			defaultValue = "''"
		case "int", "float":
			defaultValue = "0"
		case "bool":
			defaultValue = "false"
		case "date":
			defaultValue = "new Date()"
		}

		fields = append(fields, map[string]interface{}{
			"Name":         field.Name,
			"Type":         field.Type,
			"Required":     field.Required,
			"Validation":   field.Validation,
			"DefaultValue": defaultValue,
		})
	}

	data := struct {
		Name   string
		Fields []map[string]interface{}
	}{
		Name:   entity.Name,
		Fields: fields,
	}

	filename := filepath.Join(modelsDir, fmt.Sprintf("%s.js", entity.Name))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create model file %s: %v", filename, err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateJavaScriptRoutes generates route files for JavaScript application
func (cg *CodeGenerator) generateJavaScriptRoutes(appDir string, appReq *requirements.ApplicationRequirement) error {
	routesDir := filepath.Join(appDir, "routes")
	if err := os.MkdirAll(routesDir, 0755); err != nil {
		return fmt.Errorf("failed to create routes directory: %v", err)
	}

	for _, entity := range appReq.Entities {
		if err := cg.generateJavaScriptRoute(routesDir, entity); err != nil {
			return err
		}
	}

	return nil
}

// generateJavaScriptRoute generates a single route file
func (cg *CodeGenerator) generateJavaScriptRoute(routesDir string, entity requirements.Entity) error {
	routeTemplate := `const express = require('express');
const router = express.Router();
const {{.LowerName}}Controller = require('../controllers/{{.LowerName}}Controller');

// GET /api/{{.LowerName}}s - Get all {{.LowerName}}s
router.get('/', {{.LowerName}}Controller.getAll);

// GET /api/{{.LowerName}}s/:id - Get {{.LowerName}} by ID
router.get('/:id', {{.LowerName}}Controller.getById);

// POST /api/{{.LowerName}}s - Create new {{.LowerName}}
router.post('/', {{.LowerName}}Controller.create);

// PUT /api/{{.LowerName}}s/:id - Update {{.LowerName}}
router.put('/:id', {{.LowerName}}Controller.update);

// DELETE /api/{{.LowerName}}s/:id - Delete {{.LowerName}}
router.delete('/:id', {{.LowerName}}Controller.delete);

module.exports = router;`

	tmpl, err := template.New("route").Parse(routeTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse route template: %v", err)
	}

	data := struct {
		Name      string
		LowerName string
	}{
		Name:      entity.Name,
		LowerName: strings.ToLower(entity.Name),
	}

	filename := filepath.Join(routesDir, fmt.Sprintf("%sRoutes.js", strings.ToLower(entity.Name)))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create route file %s: %v", filename, err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateJavaScriptControllers generates controller files for JavaScript application
func (cg *CodeGenerator) generateJavaScriptControllers(appDir string, appReq *requirements.ApplicationRequirement) error {
	controllersDir := filepath.Join(appDir, "controllers")
	if err := os.MkdirAll(controllersDir, 0755); err != nil {
		return fmt.Errorf("failed to create controllers directory: %v", err)
	}

	for _, entity := range appReq.Entities {
		if err := cg.generateJavaScriptController(controllersDir, entity); err != nil {
			return err
		}
	}

	return nil
}

// generateJavaScriptController generates a single controller file
func (cg *CodeGenerator) generateJavaScriptController(controllersDir string, entity requirements.Entity) error {
	controllerTemplate := `const {{.Name}} = require('../models/{{.Name}}');

class {{.Name}}Controller {
  // Get all {{.LowerName}}s
  static async getAll(req, res) {
    try {
      // TODO: Implement database query to get all {{.LowerName}}s
      const {{.LowerName}}s = [];
      
      res.json({
        success: true,
        data: {{.LowerName}}s,
        count: {{.LowerName}}s.length
      });
    } catch (error) {
      console.error('Error getting {{.LowerName}}s:', error);
      res.status(500).json({
        success: false,
        error: 'Failed to retrieve {{.LowerName}}s'
      });
    }
  }

  // Get {{.LowerName}} by ID
  static async getById(req, res) {
    try {
      const { id } = req.params;
      
      // TODO: Implement database query to get {{.LowerName}} by ID
      const {{.LowerName}} = null;
      
      if (!{{.LowerName}}) {
        return res.status(404).json({
          success: false,
          error: '{{.Name}} not found'
        });
      }

      res.json({
        success: true,
        data: {{.LowerName}}
      });
    } catch (error) {
      console.error('Error getting {{.LowerName}}:', error);
      res.status(500).json({
        success: false,
        error: 'Failed to retrieve {{.LowerName}}'
      });
    }
  }

  // Create new {{.LowerName}}
  static async create(req, res) {
    try {
      const {{.LowerName}}Data = req.body;
      const {{.LowerName}} = new {{.Name}}({{.LowerName}}Data);
      
      // Validate {{.LowerName}} data
      const validationErrors = {{.LowerName}}.validate();
      if (validationErrors.length > 0) {
        return res.status(400).json({
          success: false,
          error: 'Validation failed',
          details: validationErrors
        });
      }

      // TODO: Implement database insert
      const created{{.Name}} = {{.LowerName}};
      
      res.status(201).json({
        success: true,
        data: created{{.Name}},
        message: '{{.Name}} created successfully'
      });
    } catch (error) {
      console.error('Error creating {{.LowerName}}:', error);
      res.status(500).json({
        success: false,
        error: 'Failed to create {{.LowerName}}'
      });
    }
  }

  // Update {{.LowerName}}
  static async update(req, res) {
    try {
      const { id } = req.params;
      const updateData = req.body;
      
      // TODO: Implement database update
      const updated{{.Name}} = null;
      
      if (!updated{{.Name}}) {
        return res.status(404).json({
          success: false,
          error: '{{.Name}} not found'
        });
      }

      res.json({
        success: true,
        data: updated{{.Name}},
        message: '{{.Name}} updated successfully'
      });
    } catch (error) {
      console.error('Error updating {{.LowerName}}:', error);
      res.status(500).json({
        success: false,
        error: 'Failed to update {{.LowerName}}'
      });
    }
  }

  // Delete {{.LowerName}}
  static async delete(req, res) {
    try {
      const { id } = req.params;
      
      // TODO: Implement database delete
      const deleted = false;
      
      if (!deleted) {
        return res.status(404).json({
          success: false,
          error: '{{.Name}} not found'
        });
      }

      res.json({
        success: true,
        message: '{{.Name}} deleted successfully'
      });
    } catch (error) {
      console.error('Error deleting {{.LowerName}}:', error);
      res.status(500).json({
        success: false,
        error: 'Failed to delete {{.LowerName}}'
      });
    }
  }
}

module.exports = {{.Name}}Controller;`

	tmpl, err := template.New("controller").Parse(controllerTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse controller template: %v", err)
	}

	data := struct {
		Name      string
		LowerName string
	}{
		Name:      entity.Name,
		LowerName: strings.ToLower(entity.Name),
	}

	filename := filepath.Join(controllersDir, fmt.Sprintf("%sController.js", strings.ToLower(entity.Name)))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create controller file %s: %v", filename, err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateJavaScriptMiddleware generates middleware files
func (cg *CodeGenerator) generateJavaScriptMiddleware(appDir string, appReq *requirements.ApplicationRequirement) error {
	middlewareDir := filepath.Join(appDir, "middleware")
	if err := os.MkdirAll(middlewareDir, 0755); err != nil {
		return fmt.Errorf("failed to create middleware directory: %v", err)
	}

	// Generate auth middleware
	authMiddleware := `// Authentication middleware
const auth = (req, res, next) => {
  try {
    const token = req.header('Authorization')?.replace('Bearer ', '');
    
    if (!token) {
      return res.status(401).json({
        success: false,
        error: 'Access denied. No token provided.'
      });
    }

    // TODO: Implement JWT token verification
    // const decoded = jwt.verify(token, process.env.JWT_SECRET);
    // req.user = decoded;
    
    next();
  } catch (error) {
    res.status(400).json({
      success: false,
      error: 'Invalid token.'
    });
  }
};

module.exports = auth;`

	authFile, err := os.Create(filepath.Join(middlewareDir, "auth.js"))
	if err != nil {
		return fmt.Errorf("failed to create auth middleware: %v", err)
	}
	defer authFile.Close()

	if _, err := authFile.WriteString(authMiddleware); err != nil {
		return fmt.Errorf("failed to write auth middleware: %v", err)
	}

	return nil
}

// generateJavaScriptDatabase generates database configuration
func (cg *CodeGenerator) generateJavaScriptDatabase(appDir string, appReq *requirements.ApplicationRequirement) error {
	configDir := filepath.Join(appDir, "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	dbConfig := `// Database configuration
const config = {
  development: {
    host: process.env.DB_HOST || 'localhost',
    port: process.env.DB_PORT || 5432,
    database: process.env.DB_NAME || '{{.AppName}}_dev',
    username: process.env.DB_USER || 'postgres',
    password: process.env.DB_PASSWORD || 'password',
    dialect: '{{.Database}}',
    logging: console.log
  },
  production: {
    host: process.env.DB_HOST,
    port: process.env.DB_PORT,
    database: process.env.DB_NAME,
    username: process.env.DB_USER,
    password: process.env.DB_PASSWORD,
    dialect: '{{.Database}}',
    logging: false
  }
};

const env = process.env.NODE_ENV || 'development';
const dbConfig = config[env];

// Simple database connection (placeholder)
const db = {
  connect: async () => {
    console.log('Connecting to ' + dbConfig.dialect + ' database...');
    // TODO: Implement actual database connection
    return Promise.resolve();
  },
  
  disconnect: async () => {
    console.log('Disconnecting from database...');
    // TODO: Implement actual database disconnection
    return Promise.resolve();
  }
};

module.exports = db;`

	tmpl, err := template.New("database").Parse(dbConfig)
	if err != nil {
		return fmt.Errorf("failed to parse database template: %v", err)
	}

	data := struct {
		AppName  string
		Database string
	}{
		AppName:  strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")),
		Database: appReq.Database,
	}

	file, err := os.Create(filepath.Join(configDir, "database.js"))
	if err != nil {
		return fmt.Errorf("failed to create database config: %v", err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateJavaScriptEnvConfig generates environment configuration
func (cg *CodeGenerator) generateJavaScriptEnvConfig(appDir string, appReq *requirements.ApplicationRequirement) error {
	envContent := `# Environment Configuration
NODE_ENV=development
PORT={{.Port}}

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME={{.AppName}}_dev
DB_USER=postgres
DB_PASSWORD=password

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_EXPIRES_IN=24h

# CORS Configuration
CORS_ORIGIN=*

# Logging
LOG_LEVEL=info`

	tmpl, err := template.New("env").Parse(envContent)
	if err != nil {
		return fmt.Errorf("failed to parse env template: %v", err)
	}

	data := struct {
		AppName string
		Port    interface{}
	}{
		AppName: strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")),
		Port:    appReq.Config["port"],
	}

	file, err := os.Create(filepath.Join(appDir, ".env.example"))
	if err != nil {
		return fmt.Errorf("failed to create .env.example: %v", err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateJavaScriptDockerfile generates Dockerfile for JavaScript application
func (cg *CodeGenerator) generateJavaScriptDockerfile(appDir string, appReq *requirements.ApplicationRequirement) error {
	dockerfile := `# Use official Node.js runtime as base image
FROM node:18-alpine

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production

# Copy application code
COPY . .

# Create non-root user
RUN addgroup -g 1001 -S nodejs
RUN adduser -S nodejs -u 1001

# Change ownership of the app directory
RUN chown -R nodejs:nodejs /app
USER nodejs

# Expose port
EXPOSE {{.Port}}

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD node healthcheck.js

# Start the application
CMD ["npm", "start"]`

	tmpl, err := template.New("dockerfile").Parse(dockerfile)
	if err != nil {
		return fmt.Errorf("failed to parse dockerfile template: %v", err)
	}

	data := struct {
		Port interface{}
	}{
		Port: appReq.Config["port"],
	}

	file, err := os.Create(filepath.Join(appDir, "Dockerfile"))
	if err != nil {
		return fmt.Errorf("failed to create Dockerfile: %v", err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// generateJavaScriptReadme generates README for JavaScript application
func (cg *CodeGenerator) generateJavaScriptReadme(appDir string, appReq *requirements.ApplicationRequirement) error {
	readme := `# {{.AppName}}

{{.Description}}

## Features

{{range .Features}}- {{.}}
{{end}}

## Prerequisites

- Node.js 18+ 
- npm or yarn
{{if .HasDatabase}}- {{.Database}} database{{end}}

## Installation

1. Clone the repository
2. Install dependencies:
   ` + "`" + `bash
   npm install
   ` + "`" + `

3. Copy environment configuration:
   ` + "`" + `bash
   cp .env.example .env
   ` + "`" + `

4. Update the ` + "`" + `.env` + "`" + ` file with your configuration

{{if .HasDatabase}}5. Set up your {{.Database}} database

6. Start the application:{{else}}5. Start the application:{{end}}
   ` + "`" + `bash
   npm run dev
   ` + "`" + `

## API Endpoints

{{range .Endpoints}}- ` + "`" + `{{.Method}} {{.Path}}` + "`" + ` - {{.Description}}
{{end}}

## Project Structure

` + "`" + `
{{.AppName}}/
├── app.js              # Main application file
├── package.json        # Dependencies and scripts
├── .env.example        # Environment configuration template
├── Dockerfile          # Docker configuration
├── controllers/        # Request handlers
├── models/            # Data models
├── routes/            # API routes
├── middleware/        # Custom middleware
└── config/            # Configuration files
` + "`" + `

## Development

- ` + "`" + `npm run dev` + "`" + ` - Start development server with auto-reload
- ` + "`" + `npm start` + "`" + ` - Start production server
- ` + "`" + `npm test` + "`" + ` - Run tests

## Docker

Build and run with Docker:

` + "`" + `bash
docker build -t {{.AppName}} .
docker run -p {{.Port}}:{{.Port}} {{.AppName}}
` + "`" + `

## License

MIT`

	tmpl, err := template.New("readme").Parse(readme)
	if err != nil {
		return fmt.Errorf("failed to parse readme template: %v", err)
	}

	data := struct {
		AppName     string
		Description string
		Features    []string
		HasDatabase bool
		Database    string
		Endpoints   []requirements.APIEndpoint
		Port        interface{}
	}{
		AppName:     appReq.Name,
		Description: appReq.Description,
		Features:    appReq.Features,
		HasDatabase: appReq.Database != "",
		Database:    appReq.Database,
		Endpoints:   appReq.Endpoints,
		Port:        appReq.Config["port"],
	}

	file, err := os.Create(filepath.Join(appDir, "README.md"))
	if err != nil {
		return fmt.Errorf("failed to create README.md: %v", err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// Placeholder methods for other language implementations
func (cg *CodeGenerator) generateJavaScriptWebApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// TODO: Implement JavaScript web application generation
	return fmt.Errorf("JavaScript web application generation not yet implemented")
}

func (cg *CodeGenerator) generatePythonAPIApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// TODO: Implement Python API application generation
	return fmt.Errorf("Python API application generation not yet implemented")
}

func (cg *CodeGenerator) generatePythonWebApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// TODO: Implement Python web application generation
	return fmt.Errorf("Python web application generation not yet implemented")
}

func (cg *CodeGenerator) generateJavaAPIApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// TODO: Implement Java API application generation
	return fmt.Errorf("Java API application generation not yet implemented")
}

func (cg *CodeGenerator) generateJavaWebApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// TODO: Implement Java web application generation
	return fmt.Errorf("Java web application generation not yet implemented")
}

func (cg *CodeGenerator) generatePHPAPIApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// TODO: Implement PHP API application generation
	return fmt.Errorf("PHP API application generation not yet implemented")
}

func (cg *CodeGenerator) generatePHPWebApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// TODO: Implement PHP web application generation
	return fmt.Errorf("PHP web application generation not yet implemented")
}

func (cg *CodeGenerator) generateRubyAPIApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// TODO: Implement Ruby API application generation
	return fmt.Errorf("Ruby API application generation not yet implemented")
}

func (cg *CodeGenerator) generateRubyWebApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// TODO: Implement Ruby web application generation
	return fmt.Errorf("Ruby web application generation not yet implemented")
}

func (cg *CodeGenerator) generateGoWebApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// TODO: Implement Go web application generation
	return fmt.Errorf("Go web application generation not yet implemented")
}

func (cg *CodeGenerator) generateGoCLIApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	// TODO: Implement Go CLI application generation
	return fmt.Errorf("Go CLI application generation not yet implemented")
}

