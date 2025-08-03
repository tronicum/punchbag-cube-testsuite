package aws

import (
	"context"
	"os"
	"testing"

	"github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

func TestS3BucketPolicyVersioningLifecycle(t *testing.T) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		t.Skip("AWS_REGION not set")
	}
	client, err := NewS3Client(context.Background(), region)
	if err != nil {
		t.Fatalf("failed to create S3 client: %v", err)
	}
	bucketName := "test-bucket-policy-versioning-lifecycle-123456"
	bucket := &models.ObjectStorageBucket{
		Name:     bucketName,
		Provider: models.CloudProviderAWS,
		Region:   region,
	}
	// Create bucket
	err = client.CreateBucket(context.Background(), bucket)
	if err != nil {
		t.Fatalf("failed to create bucket: %v", err)
	}
	// Set versioning
	err = client.SetBucketVersioning(context.Background(), bucketName, true)
	if err != nil {
		t.Errorf("failed to enable versioning: %v", err)
	}
	// Set policy
	policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:GetObject","Resource":"arn:aws:s3:::` + bucketName + `/*"}]}`
	err = client.SetBucketPolicy(context.Background(), bucketName, policy)
	if err != nil {
		t.Errorf("failed to set bucket policy: %v", err)
	}
	// Set lifecycle
	lifecycle := `{"Rules":[{"ID":"expire-objects","Status":"Enabled","Expiration":{"Days":1}}]}`
	err = client.SetBucketLifecycle(context.Background(), bucketName, lifecycle)
	if err != nil {
		t.Errorf("failed to set lifecycle: %v", err)
	}
	// Cleanup
	err = client.DeleteBucket(context.Background(), bucketName)
	if err != nil {
		t.Errorf("failed to delete bucket: %v", err)
	}
}
