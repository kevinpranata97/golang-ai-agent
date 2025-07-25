package codemod

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CodeModifier handles programmatic modification of code files.
// This is a simplified placeholder. A real implementation would involve AST manipulation,
// regex-based replacements, or even LLM-driven code generation for complex changes.
type CodeModifier struct {
	// Configuration or tools for code modification
}

func NewCodeModifier() *CodeModifier {
	return &CodeModifier{}
}

// ApplyModification applies a specific code modification to a file.
// This function would be called by SelfHealer for suggestions that don't have direct shell commands.
func (cm *CodeModifier) ApplyModification(filePath, modificationType, targetContent, newContent string) error {
	// Example: Replace targetContent with newContent in the file
	// This is a very basic string replacement and needs to be much more robust for real use.

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	originalContent := string(content)
	modifiedContent := strings.ReplaceAll(originalContent, targetContent, newContent)

	if originalContent == modifiedContent {
		return fmt.Errorf("modification did not find target content in %s", filePath)
	}

	if err := os.WriteFile(filePath, []byte(modifiedContent), 0644); err != nil {
		return fmt.Errorf("failed to write modified file %s: %w", filePath, err)
	}

	return nil
}

// AddLineToFile adds a line of code to a specific position in a file.
func (cm *CodeModifier) AddLineToFile(filePath, lineContent string, afterLine string) error {
	// Simplified: finds 'afterLine' and inserts 'lineContent' after it.
	// More advanced logic would handle indentation, context, etc.

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	inserted := false

	for _, line := range lines {
		newLines = append(newLines, line)
		if strings.Contains(line, afterLine) && !inserted {
			newLines = append(newLines, lineContent)
			inserted = true
		}
	}

	if !inserted {
		return fmt.Errorf("could not find target line '%s' in %s", afterLine, filePath)
	}

	if err := os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write modified file %s: %w", filePath, err)
	}

	return nil
}

// DeleteLineFromFile deletes a specific line of code from a file.
func (cm *CodeModifier) DeleteLineFromFile(filePath, lineContent string) error {
	// Simplified: removes 'lineContent' from the file.

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	deleted := false

	for _, line := range lines {
		if strings.Contains(line, lineContent) && !deleted {
			deleted = true
			continue // Skip this line
		}
		newLines = append(newLines, line)
	}

	if !deleted {
		return fmt.Errorf("could not find target line '%s' in %s", lineContent, filePath)
	}

	if err := os.WriteFile(filePath, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return fmt.Errorf("failed to write modified file %s: %w", filePath, err)
	}

	return nil
}


