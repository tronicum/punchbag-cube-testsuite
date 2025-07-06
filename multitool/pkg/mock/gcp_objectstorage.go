package mock

import (
	"errors"
	"punchbag-cube-testsuite/multitool/pkg/models"
)

type GcpObjectStorage struct {
	buckets map[string]*models.Bucket
}

func NewGcpObjectStorage() *GcpObjectStorage {
	return &GcpObjectStorage{buckets: make(map[string]*models.Bucket)}
}

func (g *GcpObjectStorage) CreateBucket(bucket *models.Bucket) (*models.Bucket, error) {
	if bucket.ID == "" {
		bucket.ID = bucket.Name + "-gcp"
	}
	g.buckets[bucket.ID] = bucket
	return bucket, nil
}

func (g *GcpObjectStorage) GetBucket(id string) (*models.Bucket, error) {
	b, ok := g.buckets[id]
	if !ok {
		return nil, errors.New("bucket not found")
	}
	return b, nil
}

func (g *GcpObjectStorage) ListBuckets() ([]*models.Bucket, error) {
	var out []*models.Bucket
	for _, b := range g.buckets {
		out = append(out, b)
	}
	return out, nil
}

func (g *GcpObjectStorage) DeleteBucket(id string) error {
	if _, ok := g.buckets[id]; !ok {
		return errors.New("bucket not found")
	}
	delete(g.buckets, id)
	return nil
}
