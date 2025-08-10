package client

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// APIClient represents a minimal client for interacting with the cube-server API (health, login, SSO, etc).
type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAPIClient creates a new API client. Returns nil if baseURL is empty.
func NewAPIClient(baseURL string) *APIClient {
	if strings.TrimSpace(baseURL) == "" {
		return nil
	}
	return &APIClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// buildURL helps construct endpoint URLs.
func (c *APIClient) buildURL(path string) string {
	return c.baseURL + path
}

// Ping checks if the API server is reachable and healthy.
func (c *APIClient) Ping() error {
	url := c.buildURL("/ping")
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API ping failed: status %d", resp.StatusCode)
	}
	return nil
}

// TODO (medium priority):
// - Only keep endpoints in this client for:
//   - Health checks (Ping)
//   - Authentication (login, logout, SSO, token refresh)
//   - User/session info (current user, session diagnostics)
//   - Server info/config (version, capabilities, status)
// - All cloud resource logic (clusters, tests, etc.) must be handled by the shared Go library, not by this client.
// - Add mock/test code for login, SSO, and other cube-server interaction endpoints as they are implemented.
// - Remove all cluster/test logic from this client; use only for server/user/session endpoints.
// - Proxy/Simulation mode:
//   - Add a way to switch between proxy and simulation mode (flag/config)
//   - In proxy mode, forward API requests to cube-server (configurable URL)
//   - In simulation mode, use local mocks/stubs for endpoints
//   - Pass through authentication/session tokens if required
//   - Unified interface for proxy/simulation so rest of mt is agnostic
//   - Document configuration for users (mode, server URL, etc)
// Next milestone:
//   - Emulate a Hetzner S3 session in simulation mode (mock S3 API)
