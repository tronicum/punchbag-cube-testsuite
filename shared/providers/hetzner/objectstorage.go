package hetzner

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

// getHetznerS3Endpoint returns the correct endpoint for Hetzner S3 for a given region
func getHetznerS3Endpoint(region string) string {
	if region == "" {
		region = "fsn1"
	}
	return "https://" + region + ".your-objectstorage.com"
}

// maskToken returns a masked version of a token, showing only the last 4 characters
func maskToken(token string) string {
	if len(token) <= 4 {
		return "****"
	}
	return "****" + token[len(token)-4:]
}

type HetznerS3Client struct {
	accessKey string
	secretKey string
	region    string
	// Simulation persistence
	simBuckets     []models.ObjectStorageBucket
	simMu          sync.Mutex
	simPersistPath string
	simEnabled     bool
}

// NewHetznerS3Client creates a Hetzner S3-compatible client using environment variables or defaults.
// Looks for HETZNER_S3_ACCESS_KEY, HETZNER_S3_SECRET_KEY, HETZNER_S3_REGION (default: fsn1)
// If CUBE_SERVER_SIM_PERSIST is set, enables general simulation persistence for all providers.
func NewHetznerS3Client() *HetznerS3Client {
	accessKey := os.Getenv("HETZNER_S3_ACCESS_KEY")
	secretKey := os.Getenv("HETZNER_S3_SECRET_KEY")
	if accessKey == "" {
		accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	}
	if secretKey == "" {
		secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}
	region := os.Getenv("HETZNER_S3_REGION")
	if region == "" {
		region = "fsn1"
	}
	simEnabled := os.Getenv("CUBE_SERVER_SIM_PERSIST") != ""
	simPersistPath := os.Getenv("CUBE_SERVER_SIM_PERSIST")
	if simPersistPath == "" {
		simPersistPath = "testdata/cube_server_sim_buckets.json"
	}
	client := &HetznerS3Client{accessKey: accessKey, secretKey: secretKey, region: region, simEnabled: simEnabled, simPersistPath: simPersistPath}
	if simEnabled {
		client.simLoad()
	}
	return client
}

// simLoad loads buckets from the general simulation persistence file (CUBE_SERVER_SIM_PERSIST).
func (c *HetznerS3Client) simLoad() {
	c.simMu.Lock()
	defer c.simMu.Unlock()
	data, err := ioutil.ReadFile(c.simPersistPath)
	if err == nil {
		_ = json.Unmarshal(data, &c.simBuckets)
	}
}

// simSave saves buckets to the general simulation persistence file (CUBE_SERVER_SIM_PERSIST).
func (c *HetznerS3Client) simSave() {
	c.simMu.Lock()
	defer c.simMu.Unlock()
	data, err := json.MarshalIndent(c.simBuckets, "", "  ")
	if err != nil {
		fmt.Printf("[DEBUG] simSave marshal error: %v\n", err)
		return
	}
	if len(data) == 0 {
		fmt.Printf("[DEBUG] simSave: marshaled data is empty! simBuckets: %+v\n", c.simBuckets)
	}
	err = ioutil.WriteFile(c.simPersistPath, data, 0644)
	if err != nil {
		fmt.Printf("[DEBUG] simSave write error: %v\n", err)
	} else {
		fmt.Printf("[DEBUG] simSave wrote %d bytes to %s\n", len(data), c.simPersistPath)
	}
}

// NewHetznerS3ClientWithKeys creates a Hetzner S3-compatible client with explicit credentials.
func NewHetznerS3ClientWithKeys(accessKey, secretKey, region string) *HetznerS3Client {
	if region == "" {
		region = "fsn1"
	}
	return &HetznerS3Client{accessKey: accessKey, secretKey: secretKey, region: region}
}

func (c *HetznerS3Client) CreateBucket(ctx context.Context, bucket *models.ObjectStorageBucket) error {
	if c.simEnabled {
		c.simMu.Lock()
		defer c.simMu.Unlock()
		for _, b := range c.simBuckets {
			if b.Name == bucket.Name {
				return fmt.Errorf("bucket %s already exists", bucket.Name)
			}
		}
		bucket.Provider = models.CloudProviderHetzner
		bucket.Region = c.region
		bucket.CreatedAt = time.Now().UTC()
		c.simBuckets = append(c.simBuckets, *bucket)
		fmt.Printf("[DEBUG] simBuckets before save: %+v\n", c.simBuckets)
		c.simSave()
		// Double-check file size after save
		fi, err := os.Stat(c.simPersistPath)
		if err == nil {
			fmt.Printf("[DEBUG] simSave file size: %d bytes\n", fi.Size())
		} else {
			fmt.Printf("[DEBUG] simSave stat error: %v\n", err)
		}
		return nil
	}
	// Normal (real) mode
	endpoint := getHetznerS3Endpoint(c.region)
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(c.region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.accessKey, c.secretKey, "")),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint, SigningRegion: c.region}, nil
			}),
		),
	)
	if err != nil {
		return err
	}
	s3Client := s3.NewFromConfig(cfg)
	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucket.Name,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *HetznerS3Client) ListBuckets(ctx context.Context) ([]models.ObjectStorageBucket, error) {
	if c.simEnabled {
		c.simMu.Lock()
		defer c.simMu.Unlock()
		return append([]models.ObjectStorageBucket(nil), c.simBuckets...), nil
	}
	// If S3 credentials are present, use S3 API
	if c.accessKey != "" && c.secretKey != "" {
		fmt.Printf("[DEBUG] Using Hetzner S3 credentials: accessKey=%s secretKey=%s region=%s\n", maskToken(c.accessKey), maskToken(c.secretKey), c.region)
		endpoint := getHetznerS3Endpoint(c.region)
		fmt.Printf("[DEBUG] S3 endpoint: %s\n", endpoint)
		cfg, err := config.LoadDefaultConfig(ctx,
			config.WithRegion(c.region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.accessKey, c.secretKey, "")),
			config.WithEndpointResolver(
				aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{URL: endpoint, SigningRegion: c.region}, nil
				}),
			),
		)
		if err != nil {
			fmt.Printf("[DEBUG] AWS config error: %v\n", err)
			return nil, err
		}
		s3Client := s3.NewFromConfig(cfg)
		resp, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
		if err != nil {
			fmt.Printf("[DEBUG] S3 ListBuckets error: %v\n", err)
			return nil, err
		}
		fmt.Printf("[DEBUG] S3 ListBuckets response: %+v\n", resp)
		var buckets []models.ObjectStorageBucket
		for _, b := range resp.Buckets {
			created := time.Now().UTC()
			if b.CreationDate != nil && b.CreationDate.Year() > 2000 {
				created = *b.CreationDate
			}
			buckets = append(buckets, models.ObjectStorageBucket{
				Name:      *b.Name,
				Provider:  models.CloudProviderHetzner,
				Region:    c.region,
				CreatedAt: created,
			})
		}
		return buckets, nil
	}
	// Otherwise, use REST API with HCLOUD_TOKEN
	token := os.Getenv("HCLOUD_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("HCLOUD_TOKEN not set. Cannot list Hetzner object storages.")
	}
	url := "https://api.hetzner.cloud/v1/object_storages"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Hetzner REST API error: %s\n%s", resp.Status, string(respBody))
	}
	type hcloudObjectStorage struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Location string `json:"location"`
	}
	type hcloudObjectStorageListResp struct {
		ObjectStorages []hcloudObjectStorage `json:"object_storages"`
	}
	var result hcloudObjectStorageListResp
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}
	buckets := make([]models.ObjectStorageBucket, 0, len(result.ObjectStorages))
	for _, b := range result.ObjectStorages {
		buckets = append(buckets, models.ObjectStorageBucket{
			Name:     b.Name,
			Provider: models.CloudProviderHetzner,
			Region:   b.Location,
		})
	}
	return buckets, nil
}

// DeleteBucket deletes a bucket in Hetzner S3-compatible object storage.
func (c *HetznerS3Client) DeleteBucket(ctx context.Context, bucketName string) error {
	if c.simEnabled {
		c.simMu.Lock()
		idx := -1
		for i, b := range c.simBuckets {
			if b.Name == bucketName {
				idx = i
				break
			}
		}
		if idx == -1 {
			c.simMu.Unlock()
			return fmt.Errorf("bucket %s not found", bucketName)
		}
		c.simBuckets = append(c.simBuckets[:idx], c.simBuckets[idx+1:]...)
		c.simSave()
		c.simMu.Unlock()
		return nil
	}
	endpoint := "https://" + c.region + ".your-objectstorage.com"
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(c.region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.accessKey, c.secretKey, "")),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint, SigningRegion: c.region}, nil
			}),
		),
	)
	if err != nil {
		return err
	}
	s3Client := s3.NewFromConfig(cfg)
	_, err = s3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		return err
	}
	// Verify deletion by listing buckets and checking if the bucket still exists
	resp, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return err
	}
	for _, b := range resp.Buckets {
		if aws.ToString(b.Name) == bucketName {
			return fmt.Errorf("bucket %s still exists after deletion", bucketName)
		}
	}
	return nil
}
