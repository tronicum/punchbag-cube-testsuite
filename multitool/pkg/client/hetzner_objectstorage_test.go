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
		t.Fatalf("ListBuckets failed: %v", err)
	}
	if len(buckets) == 0 {
		t.Errorf("No buckets found, expected at least one (or test with empty project)")
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
