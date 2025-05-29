package main

import (

    "fmt"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
)

// PluginMetadata contains plugin information - matches your backend's expected structure
type PluginMetadata struct {
    ID             string                    `yaml:"id" json:"id"`
    Name           string                    `yaml:"name" json:"name"`
    Version        string                    `yaml:"version" json:"version"`
    Description    string                    `yaml:"description" json:"description"`
    Author         string                    `yaml:"author" json:"author"`
    Endpoints      []EndpointConfig          `yaml:"endpoints" json:"endpoints"`
    Dependencies   []string                  `yaml:"dependencies" json:"dependencies"`
    Permissions    []string                  `yaml:"permissions" json:"permissions"`
    Compatibility  map[string]string         `yaml:"compatibility" json:"compatibility"`
}

// EndpointConfig defines plugin endpoint configuration
type EndpointConfig struct {
    Path    string `yaml:"path" json:"path"`
    Method  string `yaml:"method" json:"method"`
    Handler string `yaml:"handler" json:"handler"`
}

// KubestellarPlugin defines the interface that matches your backend's expectations
type KubestellarPlugin interface {
    Initialize(config map[string]interface{}) error
    GetMetadata() PluginMetadata
    GetHandlers() map[string]gin.HandlerFunc
    Health() error
    Cleanup() error
}

// TestClusterPlugin implements the KubestellarPlugin interface for cluster management testing
type TestClusterPlugin struct {
    initialized bool
    mutex       sync.RWMutex
}

// Initialize initializes the test cluster plugin
func (p *TestClusterPlugin) Initialize(config map[string]interface{}) error {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    
    if p.initialized {
        return fmt.Errorf("plugin already initialized")
    }
    
    p.initialized = true
    log.Println("✅ Test Cluster Plugin initialized successfully")
    return nil
}

// GetMetadata returns plugin metadata that matches your backend's expected structure
func (p *TestClusterPlugin) GetMetadata() PluginMetadata {
    return PluginMetadata{
        ID:          "kubestellar-cluster-plugin",
        Name:        "KubeStellar Cluster Management",
        Version:     "1.0.0",
        Description: "Plugin for cluster onboarding and detachment operations with real functionality",
        Author:      "CNCF LFX Mentee",
        Endpoints: []EndpointConfig{
            {Path: "/onboard", Method: "POST", Handler: "OnboardClusterHandler"},
            {Path: "/detach", Method: "POST", Handler: "DetachClusterHandler"},
            {Path: "/status", Method: "GET", Handler: "GetClusterStatusHandler"},
        },
        Dependencies: []string{"kubectl", "clusteradm"},
        Permissions:  []string{"cluster.read", "cluster.write"},
        Compatibility: map[string]string{
            "go":          ">=1.21",
            "kubestellar": ">=0.21.0",
        },
    }
}

// GetHandlers returns the plugin's HTTP handlers
func (p *TestClusterPlugin) GetHandlers() map[string]gin.HandlerFunc {
    return map[string]gin.HandlerFunc{
        "GetClusterStatusHandler": p.GetClusterStatusHandler,
        "OnboardClusterHandler":   p.OnboardClusterHandler,
        "DetachClusterHandler":    p.DetachClusterHandler,
    }
}

// Health performs a health check
func (p *TestClusterPlugin) Health() error {
    if !p.initialized {
        return fmt.Errorf("plugin not initialized")
    }
    return nil
}

// Cleanup performs cleanup operations
func (p *TestClusterPlugin) Cleanup() error {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    p.initialized = false
    log.Println("🧹 Test Cluster Plugin cleaned up")
    return nil
}

// GetClusterStatusHandler handles cluster status requests
func (p *TestClusterPlugin) GetClusterStatusHandler(c *gin.Context) {
    // Mock cluster data for testing
    clusters := []map[string]interface{}{
        {
            "clusterName":  "test-cluster-1",
            "status":       "ready",
            "message":      "Cluster is healthy and ready",
            "lastUpdated":  time.Now().Format(time.RFC3339),
        },
        {
            "clusterName":  "test-cluster-2", 
            "status":       "pending",
            "message":      "Cluster onboarding in progress",
            "lastUpdated":  time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
        },
    }

    summary := map[string]int{
        "total":     2,
        "ready":     1,
        "pending":   1,
        "failed":    0,
        "detaching": 0,
    }

    response := map[string]interface{}{
        "clusters": clusters,
        "summary":  summary,
        "timestamp": time.Now().Format(time.RFC3339),
    }

    c.JSON(http.StatusOK, response)
}

// OnboardClusterHandler handles cluster onboarding requests
func (p *TestClusterPlugin) OnboardClusterHandler(c *gin.Context) {
    var request map[string]interface{}
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request format",
            "details": err.Error(),
        })
        return
    }

    clusterName, exists := request["clusterName"]
    if !exists || clusterName == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "clusterName is required",
        })
        return
    }

    log.Printf("🚀 Mock onboarding cluster: %s", clusterName)

    c.JSON(http.StatusOK, gin.H{
        "message":     fmt.Sprintf("Cluster '%s' onboarding started successfully", clusterName),
        "clusterName": clusterName,
        "status":      "pending",
        "timestamp":   time.Now().Format(time.RFC3339),
    })
}

// DetachClusterHandler handles cluster detachment requests  
func (p *TestClusterPlugin) DetachClusterHandler(c *gin.Context) {
    var request map[string]interface{}
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request format",
            "details": err.Error(),
        })
        return
    }

    clusterName, exists := request["clusterName"]
    if !exists || clusterName == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "clusterName is required",
        })
        return
    }

    log.Printf("🗑️ Mock detaching cluster: %s", clusterName)

    c.JSON(http.StatusOK, gin.H{
        "message":     fmt.Sprintf("Cluster '%s' detachment started successfully", clusterName),
        "clusterName": clusterName,
        "status":      "detaching",
        "timestamp":   time.Now().Format(time.RFC3339),
    })
}

// NewPlugin creates a new instance of the test cluster plugin
// This is the required symbol that will be looked up when loading the plugin
func NewPlugin() interface{} {
    log.Println("🏗️ Creating new TestClusterPlugin instance")
    return &TestClusterPlugin{}
}

func main() {
    // This is required for Go plugins but won't be executed
    fmt.Println("This is a Go plugin, not a standalone executable")
}