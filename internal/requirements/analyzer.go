package requirements

import (
	"fmt"
	"strings"
)

// Analyzer provides functionality to analyze and validate application requirements
type Analyzer struct{}

// NewAnalyzer creates a new Analyzer instance
func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

// ValidateRequirements validates the given application requirements
func (a *Analyzer) ValidateRequirements(appReq *ApplicationRequirement) error {
	if appReq.Name == "" {
		return fmt.Errorf("application name cannot be empty")
	}

	if appReq.Language == "" {
		return fmt.Errorf("application language cannot be empty")
	}

	if appReq.Type == "" {
		return fmt.Errorf("application type cannot be empty")
	}

	// Validate entities
	for _, entity := range appReq.Entities {
		if entity.Name == "" {
			return fmt.Errorf("entity name cannot be empty")
		}

		// Validate fields
		for _, field := range entity.Fields {
			if field.Name == "" {
				return fmt.Errorf("field name cannot be empty for entity %s", entity.Name)
			}
			if field.Type == "" {
				return fmt.Errorf("field type cannot be empty for field %s in entity %s", field.Name, entity.Name)
			}
		}
	}

	// Validate API endpoints
	for _, endpoint := range appReq.APIEndpoints {
		if endpoint.Entity == "" {
			return fmt.Errorf("API endpoint entity cannot be empty")
		}
		if len(endpoint.Operations) == 0 {
			return fmt.Errorf("API endpoint operations cannot be empty for entity %s", endpoint.Entity)
		}
	}

	return nil
}

// ApplicationRequirement represents the overall requirements for an application
type ApplicationRequirement struct {
	Name        string
	Description string
	Language    string
	Type        string
	Framework   string
	Entities    []Entity
	APIEndpoints []APIEndpoint
}

// Entity represents a data entity in the application
type Entity struct {
	Name   string
	Fields []Field
}

// Field represents a field within an entity
type Field struct {
	Name          string
	Type          string
	PrimaryKey    bool
	AutoIncrement bool
	Required      bool
	Unique        bool
	Default       string
}

// APIEndpoint represents an API endpoint for an entity
type APIEndpoint struct {
	Entity     string
	Operations []string // e.g., "create", "read", "update", "delete"
}

// EntityField represents a field within an entity for template processing
type EntityField struct {
	Name     string
	Type     string
	GoName   string
	GoType   string
	JSONName string
	SQLType  string
	Required bool
	Unique   bool
	Default  string
}

// mapFieldTypeToSQL maps field types to SQL types
func MapFieldTypeToSQL(fieldType string) string {
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

// mapFieldTypeToGo maps field types to Go types
func MapFieldTypeToGo(fieldType string) string {
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

// toSnakeCase converts a string to snake_case
func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}


