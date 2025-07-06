package mock

import (
	"errors"
	"punchbag-cube-testsuite/multitool/pkg/models"
)

type AzureObjectStorage struct {
	buckets map[string]*models.Bucket
}

func NewAzureObjectStorage() *AzureObjectStorage {
	return &AzureObjectStorage{buckets: make(map[string]*models.Bucket)}
}

func (a *AzureObjectStorage) CreateBucket(bucket *models.Bucket) (*models.Bucket, error) {
	if bucket.ID == "" {
		bucket.ID = bucket.Name + "-azure"
	}
	a.buckets[bucket.ID] = bucket
	return bucket, nil
}

func (a *AzureObjectStorage) GetBucket(id string) (*models.Bucket, error) {
	b, ok := a.buckets[id]
	if !ok {
		return nil, errors.New("bucket not found")
	}
	return b, nil
}

func (a *AzureObjectStorage) ListBuckets() ([]*models.Bucket, error) {
	var out []*models.Bucket
	for _, b := range a.buckets {
		out = append(out, b)
	}
	return out, nil
}

func (a *AzureObjectStorage) DeleteBucket(id string) error {
	if _, ok := a.buckets[id]; !ok {
		return errors.New("bucket not found")
	}
	delete(a.buckets, id)
	return nil
}
