package mock

import (
	"errors"

	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

type GcpObjectStorage struct {
	buckets map[string]*sharedmodels.ObjectStorageBucket
}

func NewGcpObjectStorage() *GcpObjectStorage {
	return &GcpObjectStorage{buckets: make(map[string]*sharedmodels.ObjectStorageBucket)}
}

func (g *GcpObjectStorage) CreateBucket(bucket *sharedmodels.ObjectStorageBucket) (*sharedmodels.ObjectStorageBucket, error) {
	if bucket.ID == "" {
		bucket.ID = bucket.Name + "-gcp"
	}
	g.buckets[bucket.ID] = bucket
	return bucket, nil
}

func (g *GcpObjectStorage) GetBucket(id string) (*sharedmodels.ObjectStorageBucket, error) {
	b, ok := g.buckets[id]
	if !ok {
		return nil, errors.New("bucket not found")
	}
	return b, nil
}

func (g *GcpObjectStorage) ListBuckets() ([]*sharedmodels.ObjectStorageBucket, error) {
	var out []*sharedmodels.ObjectStorageBucket
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
