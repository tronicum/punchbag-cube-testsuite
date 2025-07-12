package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

// In-memory mock storage (thread-safe)
var (
	mockDataStore = make(map[string]json.RawMessage)
	mockDataMutex sync.RWMutex
)

// SaveMockData saves mock data by key
func SaveMockData(c *gin.Context) {
	key := c.Param("key")
	var data json.RawMessage
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	mockDataMutex.Lock()
	mockDataStore[key] = data
	mockDataMutex.Unlock()
	c.JSON(http.StatusOK, gin.H{"status": "saved", "key": key})
}

// GetMockData retrieves mock data by key
func GetMockData(c *gin.Context) {
	key := c.Param("key")
	mockDataMutex.RLock()
	data, ok := mockDataStore[key]
	mockDataMutex.RUnlock()
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.Data(http.StatusOK, "application/json", data)
}

// SaveMockDataToFile saves mock data to a local file
func SaveMockDataToFile(c *gin.Context) {
	key := c.Param("key")
	var data json.RawMessage
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	filePath := "mock_" + key + ".json"
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "saved to file", "file": filePath})
}
