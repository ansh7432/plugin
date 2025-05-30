package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
)

// This should match your backend's PluginMetadata exactly
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

type EndpointConfig struct {
    Path    string `yaml:"path" json:"path"`
    Method  string `yaml:"method" json:"method"`
    Handler string `yaml:"handler" json:"handler"`
}

// This interface should match what your backend expects
type KubestellarPlugin interface {
    Initialize(config map[string]interface{}) error
    GetMetadata() PluginMetadata
    GetHandlers() map[string]gin.HandlerFunc
    Health() error
    Cleanup() error
}

// TestClusterPlugin implements the KubestellarPlugin interface
type TestClusterPlugin struct {
    initialized bool
    mutex       sync.RWMutex
}

// Initialize initializes the plugin
func (p *TestClusterPlugin) Initialize(config map[string]interface{}) error {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    
    if p.initialized {
        return fmt.Errorf("plugin already initialized")
    }
    
    p.initialized = true
    log.Println("✅ TestClusterPlugin initialized successfully")
    return nil
}

// GetMetadata returns plugin metadata
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
    p.mutex.RLock()
    defer p.mutex.RUnlock()
    
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
    log.Println("🧹 TestClusterPlugin cleaned up")
    return nil
}

// GetClusterStatusHandler handles cluster status requests
func (p *TestClusterPlugin) GetClusterStatusHandler(c *gin.Context) {
    log.Printf("📊 GetClusterStatusHandler called")
    
    // Mock cluster data for testing
    clusters := []map[string]interface{}{
        {
            "clusterName":  "test-cluster-1",
            "status":       "failed",
            "message":      "niii bdlunga",
            "lastUpdated":  time.Now().Format(time.RFC3339),
        },
        {
            "clusterName":  "gya", 
            "status":       "failed",  // ✅ CHANGE THIS LINE
            "message":      "Cluster onboarding completed successfully",  // ✅ UPDATE MESSAGE TOO
            "lastUpdated":  time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
        },
        {
            "clusterName":  "prod-cluster-1",
            "status":       "failed",  // ✅ ALSO FIX THIS (was "pending" but summary says "failed")
            "message":      "Connection timeout during onboarding",
            "lastUpdated":  time.Now().Add(-10 * time.Minute).Format(time.RFC3339),
        },
    }

    summary := map[string]int{
        "total":     3,
        "ready":     3,  // ✅ UPDATE: test-cluster-1 + test-cluster-2
        "pending":   0,  // ✅ UPDATE: none pending now
        "failed":    0,  // ✅ UPDATE: prod-cluster-1
        "detaching": 0,
    }

    response := map[string]interface{}{
        "clusters": clusters,
        "summary":  summary,
        "timestamp": time.Now().Format(time.RFC3339),
        "plugin": "GitHub Test Plugin v2", // ✅ VERSION BUMP TO VERIFY UPDATE
    }

    log.Printf("✅ Returning cluster status: %d clusters", len(clusters))
    c.JSON(http.StatusOK, response)
}

// OnboardClusterHandler handles cluster onboarding requests
func (p *TestClusterPlugin) OnboardClusterHandler(c *gin.Context) {
    log.Printf("🚀 OnboardClusterHandler called")
    
    var request map[string]interface{}
    if err := c.ShouldBindJSON(&request); err != nil {
        log.Printf("❌ Invalid request format: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request format",
            "details": err.Error(),
        })
        return
    }

    clusterName, exists := request["clusterName"]
    if !exists || clusterName == "" {
        log.Printf("❌ Missing clusterName in request")
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "clusterName is required",
        })
        return
    }

    log.Printf("🚀 Mock onboarding cluster: %s", clusterName)

    response := gin.H{
        "message":     fmt.Sprintf("Cluster '%s' onboarding started successfully", clusterName),
        "clusterName": clusterName,
        "status":      "pending",
        "timestamp":   time.Now().Format(time.RFC3339),
        "plugin":      "GitHub Test Plugin",
    }

    log.Printf("✅ Onboarding request processed for cluster: %s", clusterName)
    c.JSON(http.StatusOK, response)
}

// DetachClusterHandler handles cluster detachment requests  
func (p *TestClusterPlugin) DetachClusterHandler(c *gin.Context) {
    log.Printf("🗑️ DetachClusterHandler called")
    
    var request map[string]interface{}
    if err := c.ShouldBindJSON(&request); err != nil {
        log.Printf("❌ Invalid request format: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request format",
            "details": err.Error(),
        })
        return
    }

    clusterName, exists := request["clusterName"]
    if !exists || clusterName == "" {
        log.Printf("❌ Missing clusterName in request")
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "clusterName is required",
        })
        return
    }

    log.Printf("🗑️ Mock detaching cluster: %s", clusterName)

    response := gin.H{
        "message":     fmt.Sprintf("Cluster '%s' detachment started successfully", clusterName),
        "clusterName": clusterName,
        "status":      "detaching",
        "timestamp":   time.Now().Format(time.RFC3339),
        "plugin":      "GitHub Test Plugin",
    }

    log.Printf("✅ Detachment request processed for cluster: %s", clusterName)
    c.JSON(http.StatusOK, response)
}

// NewPlugin creates a new instance of the plugin
// This is the EXACT symbol name that your plugin manager will look for
func NewPlugin() interface{} {
    log.Println("🏗️ Creating new TestClusterPlugin instance")
    plugin := &TestClusterPlugin{}
    
    // Initialize the plugin immediately
    if err := plugin.Initialize(nil); err != nil {
        log.Printf("❌ Failed to initialize plugin: %v", err)
        return nil
    }
    
    return plugin
}

// Alternative symbol names in case your plugin manager looks for different ones
var (
    // These are common plugin symbol names
    Plugin   = NewPlugin
    GetPlugin = NewPlugin
    CreatePlugin = NewPlugin
    NewKubestellarPlugin = NewPlugin
)

func main() {
    fmt.Println("This is a Go plugin, not a standalone executable")
}
