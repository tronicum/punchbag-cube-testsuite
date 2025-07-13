// DeleteBucket deletes a bucket in Hetzner S3-compatible object storage.



package hetzner


import (
	"context"
	"os"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

type HetznerS3Client struct {
	accessKey string
	secretKey string
	region    string
}

// NewHetznerS3Client creates a Hetzner S3-compatible client using environment variables or defaults.
// Looks for HETZNER_S3_ACCESS_KEY, HETZNER_S3_SECRET_KEY, HETZNER_S3_REGION (default: fsn1)
func NewHetznerS3Client() *HetznerS3Client {
	accessKey := os.Getenv("HETZNER_S3_ACCESS_KEY")
	secretKey := os.Getenv("HETZNER_S3_SECRET_KEY")
	region := os.Getenv("HETZNER_S3_REGION")
	if region == "" {
		region = "fsn1"
	}
	return &HetznerS3Client{accessKey: accessKey, secretKey: secretKey, region: region}
}

// NewHetznerS3ClientWithKeys creates a Hetzner S3-compatible client with explicit credentials.
func NewHetznerS3ClientWithKeys(accessKey, secretKey, region string) *HetznerS3Client {
	if region == "" {
		region = "fsn1"
	}
	return &HetznerS3Client{accessKey: accessKey, secretKey: secretKey, region: region}
}

func (c *HetznerS3Client) CreateBucket(ctx context.Context, bucket *models.ObjectStorageBucket) error {
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
	_, err = s3Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucket.Name,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *HetznerS3Client) ListBuckets(ctx context.Context) ([]models.ObjectStorageBucket, error) {
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
		return nil, err
	}
	s3Client := s3.NewFromConfig(cfg)
	resp, err := s3Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	buckets := make([]models.ObjectStorageBucket, 0, len(resp.Buckets))
	for _, b := range resp.Buckets {
		buckets = append(buckets, models.ObjectStorageBucket{
			Name:     aws.ToString(b.Name),
			Provider: models.CloudProviderHetzner,
			Region:   c.region,
		})
	}
	return buckets, nil
}

// DeleteBucket deletes a bucket in Hetzner S3-compatible object storage.
func (c *HetznerS3Client) DeleteBucket(ctx context.Context, bucketName string) error {
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


