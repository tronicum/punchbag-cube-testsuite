package mock

import (
	"errors"
	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

type HetznerObjectStorage struct {
	buckets map[string]*sharedmodels.ObjectStorageBucket
}

func NewHetznerObjectStorage() *HetznerObjectStorage {
	return &HetznerObjectStorage{buckets: make(map[string]*sharedmodels.ObjectStorageBucket)}
}

func (h *HetznerObjectStorage) CreateBucket(bucket *sharedmodels.ObjectStorageBucket) (*sharedmodels.ObjectStorageBucket, error) {
	if bucket.ID == "" {
		bucket.ID = bucket.Name + "-hetzner"
	}
	h.buckets[bucket.ID] = bucket
	return bucket, nil
}

func (h *HetznerObjectStorage) GetBucket(id string) (*sharedmodels.ObjectStorageBucket, error) {
	b, ok := h.buckets[id]
	if !ok {
		return nil, errors.New("bucket not found")
	}
	return b, nil
}

func (h *HetznerObjectStorage) ListBuckets() ([]*sharedmodels.ObjectStorageBucket, error) {
	var out []*sharedmodels.ObjectStorageBucket
	for _, b := range h.buckets {
		out = append(out, b)
	}
	return out, nil
}

func (h *HetznerObjectStorage) DeleteBucket(id string) error {
	if _, ok := h.buckets[id]; !ok {
		return errors.New("bucket not found")
	}
	delete(h.buckets, id)
	return nil
}
