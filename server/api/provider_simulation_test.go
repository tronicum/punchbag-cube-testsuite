package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	// replace with your actual import path
)

func NewMockLogger() *zap.Logger {
	logger, _ := zap.NewDevelopment()
	return logger
}

func TestValidateProvider(t *testing.T) {
	// Set up Gin router
	router := gin.Default()
	logger := NewMockLogger() // Mock logger for testing
	providerHandlers := NewProviderSimulationHandlers(nil, logger)

	// Add validation endpoint
	router.GET("/validate/:provider", providerHandlers.ValidateProvider)

	// Test valid provider
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/validate/azure", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "azure")

	// Test invalid provider
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/validate/unknown", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "unsupported provider")
}

func TestSimulateProviderOperation(t *testing.T) {
	// Set up Gin router
	router := gin.Default()
	logger := NewMockLogger() // Mock logger for testing
	providerHandlers := NewProviderSimulationHandlers(nil, logger)

	// Add simulation endpoint
	router.POST("/providers/:provider/operations/:operation", providerHandlers.SimulateProviderOperation)

	// Test valid operation
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/providers/azure/operations/create-cluster", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "create-cluster")

	// Test invalid operation
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("POST", "/providers/unknown/operations/invalid", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code) // This should be OK as it still simulates the operation
	assert.Contains(t, w.Body.String(), "unknown")
}

// Test all provider validations
func TestValidateAllProviders(t *testing.T) {
	router := gin.Default()
	logger := NewMockLogger()
	providerHandlers := NewProviderSimulationHandlers(nil, logger)
	router.GET("/validate/:provider", providerHandlers.ValidateProvider)

	providers := []string{"azure", "hetzner-hcloud", "united-ionos", "schwarz-stackit", "aws", "gcp"}

	for _, provider := range providers {
		t.Run(provider, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/validate/"+provider, nil)
			router.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Contains(t, w.Body.String(), provider)
			assert.Contains(t, w.Body.String(), "valid")
		})
	}
}

// Test GetProviderInfo endpoint
func TestGetProviderInfo(t *testing.T) {
	router := gin.Default()
	logger := NewMockLogger()
	providerHandlers := NewProviderSimulationHandlers(nil, logger)
	router.GET("/providers/:provider/info", providerHandlers.GetProviderInfo)

	providers := []string{"azure", "hetzner-hcloud", "united-ionos", "schwarz-stackit", "aws", "gcp"}

	for _, provider := range providers {
		t.Run(provider, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("GET", "/providers/"+provider+"/info", nil)
			router.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, provider, response["provider"])
			assert.Contains(t, response, "name")
			assert.Contains(t, response, "description")
		})
	}

	// Test invalid provider
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/providers/invalid/info", nil)
	router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Test provider operations with request body
func TestSimulateProviderOperationWithBody(t *testing.T) {
	router := gin.Default()
	logger := NewMockLogger()
	providerHandlers := NewProviderSimulationHandlers(nil, logger)
	router.POST("/providers/:provider/operations/:operation", providerHandlers.SimulateProviderOperation)

	// Test with request body
	requestBody := map[string]interface{}{
		"cluster_name": "test-cluster",
		"node_count":   3,
		"machine_type": "Standard_B2s",
	}
	bodyBytes, _ := json.Marshal(requestBody)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/providers/azure/operations/create-cluster", bytes.NewBuffer(bodyBytes))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "azure", response["provider"])
	assert.Equal(t, "create-cluster", response["operation"])
	assert.Equal(t, "success", response["status"])
	assert.Contains(t, response, "operation_id")
	assert.Contains(t, response, "azure_details")
}

// Test Hetzner-specific operation details
func TestSimulateHetznerOperation(t *testing.T) {
	router := gin.Default()
	logger := NewMockLogger()
	providerHandlers := NewProviderSimulationHandlers(nil, logger)
	router.POST("/providers/:provider/operations/:operation", providerHandlers.SimulateProviderOperation)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/providers/hetzner-hcloud/operations/create-cluster", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "hetzner_details")
	hetznerDetails := response["hetzner_details"].(map[string]interface{})
	assert.Contains(t, hetznerDetails, "server_id")
	assert.Contains(t, hetznerDetails, "datacenter")
}

// Test IONOS-specific operation details
func TestSimulateIONOSOperation(t *testing.T) {
	router := gin.Default()
	logger := NewMockLogger()
	providerHandlers := NewProviderSimulationHandlers(nil, logger)
	router.POST("/providers/:provider/operations/:operation", providerHandlers.SimulateProviderOperation)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/providers/united-ionos/operations/create-cluster", nil)
	router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "ionos_details")
	ionosDetails := response["ionos_details"].(map[string]interface{})
	assert.Contains(t, ionosDetails, "datacenter_id")
	assert.Contains(t, ionosDetails, "contract_number")
}

// Test different operations
func TestSimulateVariousOperations(t *testing.T) {
	router := gin.Default()
	logger := NewMockLogger()
	providerHandlers := NewProviderSimulationHandlers(nil, logger)
	router.POST("/providers/:provider/operations/:operation", providerHandlers.SimulateProviderOperation)

	operations := []string{"create-cluster", "delete-cluster", "scale-cluster", "upgrade-cluster"}

	for _, operation := range operations {
		t.Run(operation, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest("POST", "/providers/azure/operations/"+operation, nil)
			router.ServeHTTP(w, r)

			assert.Equal(t, http.StatusOK, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, operation, response["operation"])
		})
	}
}
