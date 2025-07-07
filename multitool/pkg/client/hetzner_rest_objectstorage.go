package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

type HetznerRESTObjectStorageClient struct {
	APIToken string
}

func NewHetznerRESTObjectStorageClient(token string) *HetznerRESTObjectStorageClient {
	return &HetznerRESTObjectStorageClient{APIToken: token}
}

func (c *HetznerRESTObjectStorageClient) CreateBucket(bucket *sharedmodels.ObjectStorageBucket) (*sharedmodels.ObjectStorageBucket, error) {
	fmt.Println("[Hetzner REST] CreateBucket called")
	url := "https://api.hetzner.cloud/v1/object_storage" // Adjust if endpoint differs
	body := map[string]interface{}{
		"name": bucket.Name,
		"location": bucket.Region,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("[Hetzner REST] CreateBucket response: %s\n", string(respBody))
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("Hetzner REST API error: %s", resp.Status)
	}
	// Optionally parse response for bucket ID
	bucket.ID = bucket.Name + "-hetzner-rest"
	return bucket, nil
}

func (c *HetznerRESTObjectStorageClient) ListBuckets() ([]*sharedmodels.ObjectStorageBucket, error) {
	fmt.Println("[Hetzner REST] ListBuckets called")
	url := "https://api.hetzner.cloud/v1/object_storage" // Adjust if endpoint differs
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("[Hetzner REST] ListBuckets response: %s\n", string(respBody))
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Hetzner REST API error: %s", resp.Status)
	}
	// Optionally parse response for bucket list
	return []*sharedmodels.ObjectStorageBucket{}, nil
}
