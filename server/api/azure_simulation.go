package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tronicum/punchbag-cube-testsuite/shared/simulation"
	"go.uber.org/zap"
)

// AzureHandlers provides endpoints for Azure resource simulation
type AzureHandlers struct {
	simulator *simulation.SimulationService
	logger    *zap.Logger
}

func NewAzureHandlers(logger *zap.Logger) *AzureHandlers {
	return &AzureHandlers{
		simulator: simulation.NewSimulationService(),
		logger:    logger,
	}
}

// SimulateAKS handles POST /api/v1/azure/aks/simulate
func (h *AzureHandlers) SimulateAKS(c *gin.Context) {
	var req simulation.SimulationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	// Force provider/operation for AKS
	req.Provider = "azure"
	if req.Operation == "" {
		req.Operation = "create_cluster"
	}
	result := h.simulator.SimulateOperation(&req)
	c.JSON(http.StatusOK, result)
}

// SimulateAzureBudget handles POST /api/v1/azure/budget/simulate
func (h *AzureHandlers) SimulateAzureBudget(c *gin.Context) {
	var req simulation.SimulationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	req.Provider = "azure"
	if req.Operation == "" {
		req.Operation = "create_budget"
	}
	result := h.simulator.SimulateOperation(&req)
	c.JSON(http.StatusOK, result)
}
