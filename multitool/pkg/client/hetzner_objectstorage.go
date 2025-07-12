package client

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

// Hetzner Object Storage S3-compatible API usage:
// - All bucket operations (create/list/delete) must use the S3-compatible API.
// - Use endpoint: https://<region>.your-objectstorage.com
// - Use S3 access key and secret key generated in the Hetzner Cloud Console (not hcloud API token).
// - Example access: https://<bucket-name>.<region>.your-objectstorage.com/<file-name>
// - REST API does NOT support object storage management.
// See: https://docs.hetzner.com/cloud/object-storage/overview/

type HetznerObjectStorageClient struct {
	accessKey string
	secretKey string
}

func NewHetznerObjectStorageClientFromKeys(accessKey, secretKey string) *HetznerObjectStorageClient {
	return &HetznerObjectStorageClient{
		accessKey: accessKey,
		secretKey: secretKey,
	}
}

// CreateBucket creates a new Hetzner object storage bucket (Spaces API compatible)
func (c *HetznerObjectStorageClient) CreateBucket(bucket *sharedmodels.ObjectStorageBucket) (*sharedmodels.ObjectStorageBucket, error) {
	fmt.Println("[Hetzner] CreateBucket called")
	region := bucket.Region
	if region == "" || region == "eu-central" {
		region = "fsn1" // Hetzner's default region is fsn1
	}
	endpoint := fmt.Sprintf("https://%s.your-objectstorage.com", region)
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		fmt.Printf("[Hetzner] ERROR: Invalid endpoint URL: %s\n", endpoint)
		return nil, fmt.Errorf("invalid endpoint URL: %s", endpoint)
	}
	urlStr := parsed.String()
	fmt.Printf("[Hetzner] Using endpoint: %s\n", urlStr)
	fmt.Printf("[Hetzner] Using region: %s\n", region)
	fmt.Printf("[Hetzner] Using access key: %s...\n", maskToken(c.accessKey))

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.accessKey, c.secretKey, "")),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: urlStr, SigningRegion: region}, nil
			}),
		),
	)
	if err != nil {
		fmt.Printf("[Hetzner] config error: %v\n", err)
		return nil, err
	}

	s3client := s3.NewFromConfig(cfg)
	resp, err := s3client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: &bucket.Name,
	})
	if err != nil {
		fmt.Printf("[Hetzner] CreateBucket error: %v\n", err)
		if resp != nil {
			fmt.Printf("[Hetzner] CreateBucket response: %+v\n", resp)
		}
		if ae, ok := err.(interface {
			ErrorCode() string
			ErrorMessage() string
		}); ok {
			fmt.Printf("[Hetzner] API ErrorCode: %s, Message: %s\n", ae.ErrorCode(), ae.ErrorMessage())
		}
		return nil, err
	}
	fmt.Printf("[Hetzner] CreateBucket response: %+v\n", resp)
	bucket.ID = bucket.Name + "-hetzner-real"
	// Validate CreatedAt timestamp
	if bucket.CreatedAt.IsZero() || bucket.CreatedAt.Year() < 2000 {
		fmt.Println("[Hetzner] WARNING: Invalid or missing CreatedAt timestamp in response. Setting to current time.")
		bucket.CreatedAt = time.Now().UTC()
	}
	return bucket, nil
}

// ListBuckets lists all buckets in the Hetzner project/region
func (c *HetznerObjectStorageClient) ListBuckets() ([]*sharedmodels.ObjectStorageBucket, error) {
	fmt.Println("[Hetzner] ListBuckets called")
	region := "fsn1"
	endpoint := fmt.Sprintf("https://%s.your-objectstorage.com", region)
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		fmt.Printf("[Hetzner] ERROR: Invalid endpoint URL: %s\n", endpoint)
		return nil, fmt.Errorf("invalid endpoint URL: %s", endpoint)
	}
	urlStr := parsed.String()
	fmt.Printf("[Hetzner] Using endpoint: %s\n", urlStr)
	fmt.Printf("[Hetzner] Using region: %s\n", region)
	fmt.Printf("[Hetzner] Using access key: %s...\n", maskToken(c.accessKey))

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.accessKey, c.secretKey, "")),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: urlStr, SigningRegion: region}, nil
			}),
		),
	)
	if err != nil {
		fmt.Printf("[Hetzner] config error: %v\n", err)
		return nil, err
	}

	s3client := s3.NewFromConfig(cfg)
	resp, err := s3client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		fmt.Printf("[Hetzner] ListBuckets error: %v\n", err)
		if resp != nil {
			fmt.Printf("[Hetzner] ListBuckets response: %+v\n", resp)
		}
		if ae, ok := err.(interface {
			ErrorCode() string
			ErrorMessage() string
		}); ok {
			fmt.Printf("[Hetzner] API ErrorCode: %s, Message: %s\n", ae.ErrorCode(), ae.ErrorMessage())
		}
		return nil, err
	}
	fmt.Printf("[Hetzner] ListBuckets response: %+v\n", resp)
	if len(resp.Buckets) == 0 {
		fmt.Println("[Hetzner] No buckets found.")
	}
	var out []*sharedmodels.ObjectStorageBucket
	for _, b := range resp.Buckets {
		created := time.Time{}
		if b.CreationDate != nil {
			created = *b.CreationDate
		}
		if !created.IsZero() && created.Before(time.Now().AddDate(-1, 0, 0)) {
			fmt.Printf("[Hetzner] WARNING: Bucket %s has suspiciously old timestamp: %v\n", *b.Name, created)
		}
		out = append(out, &sharedmodels.ObjectStorageBucket{
			Name:      *b.Name,
			Region:    region,
			Provider:  sharedmodels.CloudProvider("hetzner"),
			CreatedAt: created,
		})
	}
	return out, nil
}

// DeleteBucket deletes a Hetzner object storage bucket by name (S3-compatible API)
func (c *HetznerObjectStorageClient) DeleteBucket(name string) error {
	fmt.Println("[Hetzner] DeleteBucket called")
	region := "fsn1"
	endpoint := fmt.Sprintf("https://%s.your-objectstorage.com", region)
	parsed, err := url.Parse(endpoint)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return fmt.Errorf("invalid endpoint URL: %s", endpoint)
	}
	urlStr := parsed.String()
	fmt.Printf("[Hetzner] Using endpoint: %s\n", urlStr)
	fmt.Printf("[Hetzner] Using region: %s\n", region)
	fmt.Printf("[Hetzner] Using access key: %s...\n", maskToken(c.accessKey))

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.accessKey, c.secretKey, "")),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{URL: urlStr, SigningRegion: region}, nil
			}),
		),
	)
	if err != nil {
		return fmt.Errorf("[Hetzner] config error: %v", err)
	}

	s3client := s3.NewFromConfig(cfg)
	_, err = s3client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: &name,
	})
	if err != nil {
		fmt.Printf("[Hetzner] DeleteBucket error: %v\n", err)
		if ae, ok := err.(interface {
			ErrorCode() string
			ErrorMessage() string
		}); ok {
			fmt.Printf("[Hetzner] API ErrorCode: %s, Message: %s\n", ae.ErrorCode(), ae.ErrorMessage())
		}
		return err
	}
	fmt.Println("[Hetzner] DeleteBucket succeeded.")
	return nil
}

func maskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:4] + "****" + token[len(token)-4:]
}

// TODO: Implement GetBucket
