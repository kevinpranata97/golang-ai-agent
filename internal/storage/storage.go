package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Storage interface {
	Store(key string, data interface{}) error
	Retrieve(key string, dest interface{}) error
	List(prefix string) ([]string, error)
	Delete(key string) error
}

type FileStorage struct {
	basePath string
}

type StoredItem struct {
	Key       string      `json:"key"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewFileStorage(basePath string) *FileStorage {
	// Create base directory if it doesn't exist
	os.MkdirAll(basePath, 0755)
	
	return &FileStorage{
		basePath: basePath,
	}
}

func (fs *FileStorage) Store(key string, data interface{}) error {
	item := StoredItem{
		Key:       key,
		Data:      data,
		Timestamp: time.Now(),
	}
	
	jsonData, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	
	filePath := filepath.Join(fs.basePath, key+".json")
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	if err := os.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

func (fs *FileStorage) Retrieve(key string, dest interface{}) error {
	filePath := filepath.Join(fs.basePath, key+".json")
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to read file: %w", err)
	}
	
	var item StoredItem
	if err := json.Unmarshal(data, &item); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}
	
	// Convert the data back to the destination type
	jsonData, err := json.Marshal(item.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal item data: %w", err)
	}
	
	if err := json.Unmarshal(jsonData, dest); err != nil {
		return fmt.Errorf("failed to unmarshal to destination: %w", err)
	}
	
	return nil
}

func (fs *FileStorage) List(prefix string) ([]string, error) {
	var keys []string
	
	err := filepath.Walk(fs.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}
		
		if filepath.Ext(path) != ".json" {
			return nil
		}
		
		relPath, err := filepath.Rel(fs.basePath, path)
		if err != nil {
			return err
		}
		
		// Remove .json extension
		key := relPath[:len(relPath)-5]
		
		if prefix == "" || filepath.HasPrefix(key, prefix) {
			keys = append(keys, key)
		}
		
		return nil
	})
	
	return keys, err
}

func (fs *FileStorage) Delete(key string) error {
	filePath := filepath.Join(fs.basePath, key+".json")
	
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("key not found: %s", key)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}
	
	return nil
}

// Additional utility methods
func (fs *FileStorage) GetStorageStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	var totalFiles int
	var totalSize int64
	
	err := filepath.Walk(fs.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			totalFiles++
			totalSize += info.Size()
		}
		
		return nil
	})
	
	if err != nil {
		return nil, err
	}
	
	stats["total_files"] = totalFiles
	stats["total_size_bytes"] = totalSize
	stats["base_path"] = fs.basePath
	
	return stats, nil
}

func (fs *FileStorage) Cleanup(olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	
	return filepath.Walk(fs.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		
		if info.ModTime().Before(cutoff) {
			return os.Remove(path)
		}
		
		return nil
	})
}

