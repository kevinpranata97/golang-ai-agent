package codegen

import (
	"fmt"
	"os"
	"os/exec"
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
	switch appReq.Type {
	case "api":
		return cg.generateJavaScriptAPIApplication(appDir, appReq)
	case "web":
		return cg.generateJavaScriptWebApplication(appDir, appReq)
	default:
		return cg.generateJavaScriptAPIApplication(appDir, appReq)
	}
}

// generatePythonApplication generates a Python application
func (cg *CodeGenerator) generatePythonApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	switch appReq.Type {
	case "api":
		return cg.generatePythonAPIApplication(appDir, appReq)
	case "web":
		return cg.generatePythonWebApplication(appDir, appReq)
	default:
		return cg.generatePythonAPIApplication(appDir, appReq)
	}
}

// generateJavaApplication generates a Java application
func (cg *CodeGenerator) generateJavaApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	switch appReq.Type {
	case "api":
		return cg.generateJavaAPIApplication(appDir, appReq)
	case "web":
		return cg.generateJavaWebApplication(appDir, appReq)
	default:
		return cg.generateJavaAPIApplication(appDir, appReq)
	}
}

// generatePHPApplication generates a PHP application
func (cg *CodeGenerator) generatePHPApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	switch appReq.Type {
	case "api":
		return cg.generatePHPAPIApplication(appDir, appReq)
	case "web":
		return cg.generatePHPWebApplication(appDir, appReq)
	default:
		return cg.generatePHPAPIApplication(appDir, appReq)
	}
}

// generateRubyApplication generates a Ruby application
func (cg *CodeGenerator) generateRubyApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	switch appReq.Type {
	case "api":
		return cg.generateRubyAPIApplication(appDir, appReq)
	case "web":
		return cg.generateRubyWebApplication(appDir, appReq)
	default:
		return cg.generateRubyAPIApplication(appDir, appReq)
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

	// Run go mod tidy
	if _, err := exec.LookPath("go"); err == nil {
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Dir = appDir
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("Warning: failed to run go mod tidy: %v, output: %s\n", err, string(output))
		}
	}

	return nil
}

// generateGoWebApplication generates a web application with frontend and backend
func (cg *CodeGenerator) generateGoWebApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	if err := cg.generateGoAPIApplication(appDir, appReq); err != nil {
		return err
	}

	staticDir := filepath.Join(appDir, "static")
	if err := os.MkdirAll(staticDir, 0755); err != nil {
		return fmt.Errorf("failed to create static directory: %v", err)
	}

	if err := cg.generateHTMLTemplates(staticDir, appReq); err != nil {
		return err
	}

	if err := cg.generateCSS(staticDir, appReq); err != nil {
		return err
	}

	if err := cg.generateJavaScript(staticDir, appReq); err != nil {
		return err
	}

	return nil
}

// generateGoCLIApplication generates a CLI application
func (cg *CodeGenerator) generateGoCLIApplication(appDir string, appReq *requirements.ApplicationRequirement) error {
	if err := cg.generateCLIMain(appDir, appReq); err != nil {
		return err
	}

	if err := cg.generateGoMod(appDir, appReq); err != nil {
		return err
	}

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

	port := "8080"
	if p, ok := appReq.Config["port"]; ok {
		port = fmt.Sprintf("%v", p)
	}

	data := struct {
		ModuleName string
		Port       string
	}{
		ModuleName: strings.ToLower(strings.ReplaceAll(appReq.Name, " ", "-")),
		Port:       port,
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

go 1.18

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
	
	err := db.QueryRow(query, id).Scan({{range $i, $field := .ScanFields}}{{if $i}}, {{end}}&{{$.LowerName}}.{{$field}}{{end}})
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
		err := rows.Scan({{range $i, $field := .ScanFields}}{{if $i}}, {{end}}&{{$.LowerName}}.{{$field}}{{end}})
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

	for _, field := range entity.Fields {
		goType := cg.mapFieldTypeToGo(field.Type)
		
		// Consistently handle ID naming
		goName := strings.ToUpper(field.Name[:1]) + field.Name[1:]
		if strings.ToLower(field.Name) == "id" {
			goName = "ID"
		}
		
		jsonName := strings.ToLower(field.Name)

		fields = append(fields, map[string]interface{}{
			"GoName":   goName,
			"GoType":   goType,
			"JSONName": jsonName,
			"Required": field.Required,
		})

		if strings.ToLower(field.Name) != "id" && strings.ToLower(field.Name) != "created_at" {
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

	if err := cg.generateBaseHandler(handlersDir); err != nil {
		return err
	}

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

type Handler struct {
	DB *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

type ErrorResponse struct {
	Error string ` + "`json:\"error\"`" + `
}

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

func (h *Handler) GetAll{{.Name}}s(c *gin.Context) {
	{{.LowerName}}s, err := models.GetAll{{.Name}}s(h.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{Data: {{.LowerName}}s})
}

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

	if err := cg.generateDatabaseInit(dbDir, appReq); err != nil {
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
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func Initialize(databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		databaseURL = "./app.db"
	}

	// Ensure directory exists for SQLite
	if !strings.HasPrefix(databaseURL, "file:") {
		dir := filepath.Dir(databaseURL)
		if dir != "." && dir != "/" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create database directory: %v", err)
			}
		}
	}

	db, err := sql.Open("sqlite3", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}

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

		if strings.ToLower(field.Name) == "id" {
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

func Setup(r *gin.Engine, h *handlers.Handler) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	_ = api // avoid unused variable error
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

type Config struct {
	Port        string
	DatabaseURL string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "{{.Port}}"),
		DatabaseURL: getEnv("DATABASE_URL", "{{.DatabaseURL}}"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
`

	port := "8080"
	if p, ok := appReq.Config["port"]; ok {
		port = fmt.Sprintf("%v", p)
	}

	data := map[string]interface{}{
		"Port":        port,
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

	port := "8080"
	if p, ok := appReq.Config["port"]; ok {
		port = fmt.Sprintf("%v", p)
	}

	data := map[string]interface{}{
		"Port": port,
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
`

	data := struct {
		Name        string
		Description string
		Features    []string
		Endpoints   []requirements.APIEndpoint
	}{
		Name:        appReq.Name,
		Description: appReq.Description,
		Features:    appReq.Features,
		Endpoints:   appReq.Endpoints,
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

// Helper methods for other languages (placeholders)
func (cg *CodeGenerator) generateJavaScriptAPIApplication(dir string, req *requirements.ApplicationRequirement) error { return nil }
func (cg *CodeGenerator) generateJavaScriptWebApplication(dir string, req *requirements.ApplicationRequirement) error { return nil }
func (cg *CodeGenerator) generatePythonAPIApplication(dir string, req *requirements.ApplicationRequirement) error { return nil }
func (cg *CodeGenerator) generatePythonWebApplication(dir string, req *requirements.ApplicationRequirement) error { return nil }
func (cg *CodeGenerator) generateJavaAPIApplication(dir string, req *requirements.ApplicationRequirement) error { return nil }
func (cg *CodeGenerator) generateJavaWebApplication(dir string, req *requirements.ApplicationRequirement) error { return nil }
func (cg *CodeGenerator) generatePHPAPIApplication(dir string, req *requirements.ApplicationRequirement) error { return nil }
func (cg *CodeGenerator) generatePHPWebApplication(dir string, req *requirements.ApplicationRequirement) error { return nil }
func (cg *CodeGenerator) generateRubyAPIApplication(dir string, req *requirements.ApplicationRequirement) error { return nil }
func (cg *CodeGenerator) generateRubyWebApplication(dir string, req *requirements.ApplicationRequirement) error { return nil }
func (cg *CodeGenerator) generateCLIMain(dir string, req *requirements.ApplicationRequirement) error { return nil }
func (cg *CodeGenerator) generateCLICommands(dir string, req *requirements.ApplicationRequirement) error { return nil }
func (cg *CodeGenerator) generateHTMLTemplates(dir string, req *requirements.ApplicationRequirement) error { return nil }
func (cg *CodeGenerator) generateCSS(dir string, req *requirements.ApplicationRequirement) error { return nil }
func (cg *CodeGenerator) generateJavaScript(dir string, req *requirements.ApplicationRequirement) error { return nil }
