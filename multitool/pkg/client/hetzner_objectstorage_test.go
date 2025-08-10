package client

import (
	"context"
	"testing"

	"github.com/tronicum/punchbag-cube-testsuite/shared/models"
	"github.com/tronicum/punchbag-cube-testsuite/shared/providers/hetzner"
)

func TestHetznerObjectStorage_ListBuckets(t *testing.T) {
	client := hetzner.NewHetznerS3Client()
	buckets, err := client.ListBuckets(context.Background())
	if err != nil {
		// Check if error is due to missing credentials - skip test with warning instead of failing
		if errStr := err.Error(); errStr == "HCLOUD_TOKEN not set. Cannot list Hetzner object storages." {
			t.Skipf("Skipping test due to missing credentials: %v", err)
			return
		}
		t.Fatalf("ListBuckets failed: %v", err)
	}
	if len(buckets) == 0 {
		t.Logf("No buckets found - this may be expected for empty projects")
	}
	for _, b := range buckets {
		if b.Name == "" {
			t.Errorf("Bucket with empty name found: %+v", b)
		}
		if b.Provider != models.CloudProviderHetzner {
			t.Errorf("Bucket provider mismatch: got %v, want %v", b.Provider, models.CloudProviderHetzner)
		}
	}
}
