package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	simulation "github.com/tronicum/punchbag-cube-testsuite/shared/simulation"
)

// SimulatorClient allows generator to call the cube-server simulation endpoints
type SimulatorClient struct {
	BaseURL string
}

// SimulateAKSCluster calls the /simulator/azure/aks endpoint
func (c *SimulatorClient) SimulateAKSCluster(params map[string]interface{}) (*simulation.SimulationResult, error) {
	url := c.BaseURL + "/api/v1/simulator/azure/aks"
	request := simulation.SimulationRequest{
		Provider:   "azure",
		Operation:  "create_cluster",
		Parameters: params,
	}
	body, _ := json.Marshal(request)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("simulator error: %s", resp.Status)
	}
	var result simulation.SimulationResult
	data, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SimulateAzureBudget calls the /simulator/azure/budget endpoint
func (c *SimulatorClient) SimulateAzureBudget(params map[string]interface{}) (*simulation.SimulationResult, error) {
	url := c.BaseURL + "/api/v1/simulator/azure/budget"
	request := simulation.SimulationRequest{
		Provider:   "azure",
		Operation:  "create_budget",
		Parameters: params,
	}
	body, _ := json.Marshal(request)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("simulator error: %s", resp.Status)
	}
	var result simulation.SimulationResult
	data, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// NOTE: All files in this directory should use the same package name for Go compatibility.
// For generator integration, use 'package generator' in all files (main.go, simulator_client.go, etc.).
//
// Example usage in main.go:
//
// import (
//   "punchbag-cube-testsuite/generator"
// )
//
// client := generator.SimulatorClient{BaseURL: "http://localhost:8080"}
// result, err := client.SimulateAKSCluster(params)
// if err == nil {
//     // Use result to generate Terraform template
// }
