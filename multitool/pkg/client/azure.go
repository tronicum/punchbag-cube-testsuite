package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"punchbag-cube-testsuite/multitool/pkg/models"
)

// LogAnalyticsClient provides methods for Log Analytics operations
type LogAnalyticsClient struct {
	client *APIClient
}

func NewLogAnalyticsClient(client *APIClient) *LogAnalyticsClient {
	return &LogAnalyticsClient{client: client}
}

func (c *LogAnalyticsClient) Create(workspace *models.LogAnalyticsWorkspace) (*models.LogAnalyticsWorkspace, error) {
	url := fmt.Sprintf("%s/api/azure/loganalytics", c.client.baseURL)
	data, err := json.Marshal(workspace)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.httpClient.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	var result models.LogAnalyticsWorkspace
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *LogAnalyticsClient) Get(id string) (*models.LogAnalyticsWorkspace, error) {
	url := fmt.Sprintf("%s/api/azure/loganalytics/%s", c.client.baseURL, id)
	resp, err := c.client.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	var result models.LogAnalyticsWorkspace
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *LogAnalyticsClient) List() ([]*models.LogAnalyticsWorkspace, error) {
	url := fmt.Sprintf("%s/api/azure/loganalytics", c.client.baseURL)
	resp, err := c.client.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	var results []*models.LogAnalyticsWorkspace
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func (c *LogAnalyticsClient) Delete(id string) error {
	url := fmt.Sprintf("%s/api/azure/loganalytics/%s", c.client.baseURL, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found")
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	return nil
}

// AzureBudgetClient provides methods for Azure Budget operations
type AzureBudgetClient struct {
	client *APIClient
}

func NewAzureBudgetClient(client *APIClient) *AzureBudgetClient {
	return &AzureBudgetClient{client: client}
}

func (c *AzureBudgetClient) Create(budget *models.AzureBudget) (*models.AzureBudget, error) {
	url := fmt.Sprintf("%s/api/azure/budget", c.client.baseURL)
	data, err := json.Marshal(budget)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.httpClient.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	var result models.AzureBudget
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *AzureBudgetClient) Get(id string) (*models.AzureBudget, error) {
	url := fmt.Sprintf("%s/api/azure/budget/%s", c.client.baseURL, id)
	resp, err := c.client.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	var result models.AzureBudget
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *AzureBudgetClient) List() ([]*models.AzureBudget, error) {
	url := fmt.Sprintf("%s/api/azure/budget", c.client.baseURL)
	resp, err := c.client.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	var results []*models.AzureBudget
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func (c *AzureBudgetClient) Delete(id string) error {
	url := fmt.Sprintf("%s/api/azure/budget/%s", c.client.baseURL, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found")
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	return nil
}

// AppInsightsClient provides methods for Application Insights operations
type AppInsightsClient struct {
	client *APIClient
}

func NewAppInsightsClient(client *APIClient) *AppInsightsClient {
	return &AppInsightsClient{client: client}
}

func (c *AppInsightsClient) Create(app *models.AppInsightsResource) (*models.AppInsightsResource, error) {
	url := fmt.Sprintf("%s/api/azure/appinsights", c.client.baseURL)
	data, err := json.Marshal(app)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.httpClient.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	var result models.AppInsightsResource
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *AppInsightsClient) Get(id string) (*models.AppInsightsResource, error) {
	url := fmt.Sprintf("%s/api/azure/appinsights/%s", c.client.baseURL, id)
	resp, err := c.client.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	var result models.AppInsightsResource
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *AppInsightsClient) List() ([]*models.AppInsightsResource, error) {
	url := fmt.Sprintf("%s/api/azure/appinsights", c.client.baseURL)
	resp, err := c.client.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	var results []*models.AppInsightsResource
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}

func (c *AppInsightsClient) Delete(id string) error {
	url := fmt.Sprintf("%s/api/azure/appinsights/%s", c.client.baseURL, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found")
	}
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}
	return nil
}
