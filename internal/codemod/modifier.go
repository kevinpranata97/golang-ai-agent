package codemod

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// CodeModifier handles programmatic modification of code files.
type CodeModifier struct{}

// NewCodeModifier creates a new CodeModifier instance.
func NewCodeModifier() *CodeModifier {
	return &CodeModifier{}
}

// ApplyModification applies a given code modification to a target file.
// The 'modification' string can be a simple replacement, or a more complex instruction
// that the LLM (or a more sophisticated parser) would interpret.
// For simplicity, this example assumes 'modification' is a direct string replacement
// in the format "OLD_CODE:::NEW_CODE".
func (cm *CodeModifier) ApplyModification(appPath, targetFile, modification string) error {
	filePath := filepath.Join(appPath, targetFile)

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	parts := strings.SplitN(modification, ":::", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid modification format. Expected 'OLD_CODE:::NEW_CODE'")
	}
	oldCode := parts[0]
	newCode := parts[1]

	modifiedContent := strings.ReplaceAll(string(content), oldCode, newCode)

	if err := ioutil.WriteFile(filePath, []byte(modifiedContent), 0644); err != nil {
		return fmt.Errorf("failed to write modified file %s: %w", filePath, err)
	}

	return nil
}


