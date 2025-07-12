package mock

import (
	"errors"

	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

type StackITObjectStorage struct {
	buckets map[string]*sharedmodels.ObjectStorageBucket
}

func NewStackITObjectStorage() *StackITObjectStorage {
	return &StackITObjectStorage{buckets: make(map[string]*sharedmodels.ObjectStorageBucket)}
}

func (s *StackITObjectStorage) CreateBucket(bucket *sharedmodels.ObjectStorageBucket) (*sharedmodels.ObjectStorageBucket, error) {
	if bucket.ID == "" {
		bucket.ID = bucket.Name + "-stackit"
	}
	s.buckets[bucket.ID] = bucket
	return bucket, nil
}

func (s *StackITObjectStorage) GetBucket(id string) (*sharedmodels.ObjectStorageBucket, error) {
	b, ok := s.buckets[id]
	if !ok {
		return nil, errors.New("bucket not found")
	}
	return b, nil
}

func (s *StackITObjectStorage) ListBuckets() ([]*sharedmodels.ObjectStorageBucket, error) {
	var out []*sharedmodels.ObjectStorageBucket
	for _, b := range s.buckets {
		out = append(out, b)
	}
	return out, nil
}

func (s *StackITObjectStorage) DeleteBucket(id string) error {
	if _, ok := s.buckets[id]; !ok {
		return errors.New("bucket not found")
	}
	delete(s.buckets, id)
	return nil
}
