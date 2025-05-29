package main

import (
    "encoding/json"
    "net/http"
)

// Plugin metadata - these will be read by KubeStellar
var (
    PluginName    = "sample-plugin"
    PluginVersion = "1.0.0"
    PluginAuthor  = "Your Name"
)

// Example endpoint handler
func healthCheck(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{
        "status": "healthy",
        "plugin": PluginName,
        "version": PluginVersion,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// Initialize function - called when plugin loads
func Initialize() error {
    // Plugin initialization logic
    return nil
}

// GetEndpoints - returns available endpoints
func GetEndpoints() []string {
    return []string{"/health"}
}

// HandleRequest - main request handler
func HandleRequest(endpoint string, w http.ResponseWriter, r *http.Request) {
    switch endpoint {
    case "/health":
        healthCheck(w, r)
    default:
        http.NotFound(w, r)
    }
}

func main() {
    // This is required for Go plugins
}