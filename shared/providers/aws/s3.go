package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	"github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

type S3Client struct {
	client *s3.Client
}


func NewS3Client(ctx context.Context, region string) (*S3Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	return &S3Client{client: s3.NewFromConfig(cfg)}, nil
}

func (c *S3Client) CreateBucket(ctx context.Context, bucket *models.ObjectStorageBucket) error {
	_, err := c.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket.Name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(bucket.Region),
		},
	})
	return err
}

func (c *S3Client) ListBuckets(ctx context.Context) ([]models.ObjectStorageBucket, error) {
	resp, err := c.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	var buckets []models.ObjectStorageBucket
	for _, b := range resp.Buckets {
		buckets = append(buckets, models.ObjectStorageBucket{Name: *b.Name, Provider: models.CloudProviderAWS})
	}
	   return buckets, nil
}

func (c *S3Client) DeleteBucket(ctx context.Context, bucketName string) error {
	   _, err := c.client.DeleteBucket(ctx, &s3.DeleteBucketInput{
			   Bucket: aws.String(bucketName),
	   })
	   return err
}
