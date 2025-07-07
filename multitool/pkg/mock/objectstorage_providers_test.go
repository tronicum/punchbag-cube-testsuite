package mock_test

import (
	"testing"

	"github.com/tronicum/punchbag-cube-testsuite/multitool/pkg/mock"
	"github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

func TestStackITObjectStorage(t *testing.T) {
	store := mock.NewStackITObjectStorage()
	bucket := &models.ObjectStorageBucket{Name: "test-bucket", Provider: models.StackIT, Region: "eu01"}
	created, err := store.CreateBucket(bucket)
	if err != nil {
		t.Fatalf("CreateBucket failed: %v", err)
	}
	if created.Name != "test-bucket" || created.Provider != models.StackIT {
		t.Errorf("Unexpected bucket: %+v", created)
	}
	buckets, err := store.ListBuckets()
	if err != nil || len(buckets) != 1 {
		t.Errorf("ListBuckets failed or wrong count: %v, %d", err, len(buckets))
	}
	got, err := store.GetBucket(created.ID)
	if err != nil || got.Name != "test-bucket" {
		t.Errorf("GetBucket failed: %v, %+v", err, got)
	}
	if err := store.DeleteBucket(created.ID); err != nil {
		t.Errorf("DeleteBucket failed: %v", err)
	}
}

func TestHetznerObjectStorage(t *testing.T) {
	store := mock.NewHetznerObjectStorage()
	bucket := &models.ObjectStorageBucket{Name: "hetz-bucket", Provider: models.Hetzner, Region: "fsn1"}
	created, err := store.CreateBucket(bucket)
	if err != nil {
		t.Fatalf("CreateBucket failed: %v", err)
	}
	if created.Name != "hetz-bucket" || created.Provider != models.Hetzner {
		t.Errorf("Unexpected bucket: %+v", created)
	}
	buckets, err := store.ListBuckets()
	if err != nil || len(buckets) != 1 {
		t.Errorf("ListBuckets failed or wrong count: %v, %d", err, len(buckets))
	}
	got, err := store.GetBucket(created.ID)
	if err != nil || got.Name != "hetz-bucket" {
		t.Errorf("GetBucket failed: %v, %+v", err, got)
	}
	if err := store.DeleteBucket(created.ID); err != nil {
		t.Errorf("DeleteBucket failed: %v", err)
	}
}

func TestIonosObjectStorage(t *testing.T) {
	store := mock.NewIonosObjectStorage()
	bucket := &models.ObjectStorageBucket{Name: "ionos-bucket", Provider: models.IONOS, Region: "de"}
	created, err := store.CreateBucket(bucket)
	if err != nil {
		t.Fatalf("CreateBucket failed: %v", err)
	}
	if created.Name != "ionos-bucket" || created.Provider != models.IONOS {
		t.Errorf("Unexpected bucket: %+v", created)
	}
	buckets, err := store.ListBuckets()
	if err != nil || len(buckets) != 1 {
		t.Errorf("ListBuckets failed or wrong count: %v, %d", err, len(buckets))
	}
	got, err := store.GetBucket(created.ID)
	if err != nil || got.Name != "ionos-bucket" {
		t.Errorf("GetBucket failed: %v, %+v", err, got)
	}
	if err := store.DeleteBucket(created.ID); err != nil {
		t.Errorf("DeleteBucket failed: %v", err)
	}
}
