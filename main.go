package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
)

// PluginMetadata contains plugin information
type PluginMetadata struct {
    ID           string            `yaml:"id"`
    Name         string            `yaml:"name"`
    Version      string            `yaml:"version"`
    Description  string            `yaml:"description"`
    Author       string            `yaml:"author"`
    Endpoints    []EndpointConfig  `yaml:"endpoints"`
    Dependencies []string          `yaml:"dependencies"`
    Permissions  []string          `yaml:"permissions"`
}

// EndpointConfig defines plugin endpoint configuration
type EndpointConfig struct {
    Path    string `yaml:"path"`
    Method  string `yaml:"method"`
    Handler string `yaml:"handler"`
}

// KubestellarPlugin defines the interface that all dynamic plugins must implement
type KubestellarPlugin interface {
    Initialize(config map[string]interface{}) error
    GetMetadata() PluginMetadata
    GetHandlers() map[string]gin.HandlerFunc
    Health() error
    Cleanup() error
}

// SimpleTestPlugin implements the KubestellarPlugin interface for testing
type SimpleTestPlugin struct {
    initialized bool
    mutex       sync.RWMutex
}

// Initialize initializes the test plugin
func (p *SimpleTestPlugin) Initialize(config map[string]interface{}) error {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    
    if p.initialized {
        return fmt.Errorf("plugin already initialized")
    }
    
    p.initialized = true
    log.Println("✅ Simple Test Plugin initialized successfully")
    return nil
}

// GetMetadata returns plugin metadata
func (p *SimpleTestPlugin) GetMetadata() PluginMetadata {
    return PluginMetadata{
        ID:          "simple-test-plugin",
        Name:        "Simple Test Plugin",
        Version:     "1.0.0",
        Description: "A simple test plugin for GitHub installation testing",
        Author:      "Your Name",
        Endpoints: []EndpointConfig{
            {Path: "/health", Method: "GET", Handler: "HealthHandler"},
            {Path: "/info", Method: "GET", Handler: "InfoHandler"},
        },
        Dependencies: []string{},
        Permissions:  []string{},
    }
}

// GetHandlers returns the plugin's HTTP handlers
func (p *SimpleTestPlugin) GetHandlers() map[string]gin.HandlerFunc {
    return map[string]gin.HandlerFunc{
        "HealthHandler": p.HealthHandler,
        "InfoHandler":  p.InfoHandler,
    }
}

// Health performs a health check
func (p *SimpleTestPlugin) Health() error {
    if !p.initialized {
        return fmt.Errorf("plugin not initialized")
    }
    return nil
}

// Cleanup performs cleanup operations
func (p *SimpleTestPlugin) Cleanup() error {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    p.initialized = false
    log.Println("🧹 Simple Test Plugin cleaned up")
    return nil
}

// HealthHandler handles health check requests
func (p *SimpleTestPlugin) HealthHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":    "healthy",
        "plugin":    "simple-test-plugin",
        "timestamp": time.Now().Format(time.RFC3339),
        "message":   "Plugin is running correctly",
    })
}

// InfoHandler handles info requests
func (p *SimpleTestPlugin) InfoHandler(c *gin.Context) {
    metadata := p.GetMetadata()
    c.JSON(http.StatusOK, gin.H{
        "plugin":      metadata,
        "initialized": p.initialized,
        "timestamp":   time.Now().Format(time.RFC3339),
        "endpoints":   []string{"/health", "/info"},
    })
}

// NewPlugin creates a new instance of the test plugin
// This is the required symbol that will be looked up when loading the plugin
func NewPlugin() interface{} {
    log.Println("🏗️ Creating new SimpleTestPlugin instance")
    return &SimpleTestPlugin{}
}

func main() {
    // This is required for Go plugins but won't be executed
    fmt.Println("This is a Go plugin, not a standalone executable")
}