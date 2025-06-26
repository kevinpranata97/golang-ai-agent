package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server struct {
		Port         string `json:"port"`
		Host         string `json:"host"`
		ReadTimeout  int    `json:"read_timeout"`
		WriteTimeout int    `json:"write_timeout"`
	} `json:"server"`
	
	GitHub struct {
		Token         string `json:"token"`
		WebhookSecret string `json:"webhook_secret"`
		BaseURL       string `json:"base_url"`
	} `json:"github"`
	
	Storage struct {
		Type string `json:"type"`
		Path string `json:"path"`
	} `json:"storage"`
	
	Testing struct {
		Timeout       int  `json:"timeout"`
		Parallel      bool `json:"parallel"`
		Coverage      bool `json:"coverage"`
		SecurityScan  bool `json:"security_scan"`
	} `json:"testing"`
	
	Debugging struct {
		LogLevel    string `json:"log_level"`
		ProfileMode bool   `json:"profile_mode"`
		MaxSessions int    `json:"max_sessions"`
	} `json:"debugging"`
	
	Workflow struct {
		MaxConcurrent int `json:"max_concurrent"`
		RetryAttempts int `json:"retry_attempts"`
		CleanupAfter  int `json:"cleanup_after"`
	} `json:"workflow"`
}

func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}
	
	// Set defaults
	config.Server.Port = "8080"
	config.Server.Host = "0.0.0.0"
	config.Server.ReadTimeout = 30
	config.Server.WriteTimeout = 30
	
	config.GitHub.BaseURL = "https://api.github.com"
	
	config.Storage.Type = "file"
	config.Storage.Path = "./data"
	
	config.Testing.Timeout = 300
	config.Testing.Parallel = true
	config.Testing.Coverage = true
	config.Testing.SecurityScan = true
	
	config.Debugging.LogLevel = "info"
	config.Debugging.ProfileMode = false
	config.Debugging.MaxSessions = 5
	
	config.Workflow.MaxConcurrent = 3
	config.Workflow.RetryAttempts = 3
	config.Workflow.CleanupAfter = 24
	
	// Load from file if exists
	if configPath != "" {
		if _, err := os.Stat(configPath); err == nil {
			data, err := os.ReadFile(configPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
			
			if err := json.Unmarshal(data, config); err != nil {
				return nil, fmt.Errorf("failed to parse config file: %w", err)
			}
		}
	}
	
	// Override with environment variables
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		config.GitHub.Token = token
	}
	
	if secret := os.Getenv("WEBHOOK_SECRET"); secret != "" {
		config.GitHub.WebhookSecret = secret
	}
	
	if port := os.Getenv("PORT"); port != "" {
		config.Server.Port = port
	}
	
	return config, nil
}

func (c *Config) Save(configPath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

