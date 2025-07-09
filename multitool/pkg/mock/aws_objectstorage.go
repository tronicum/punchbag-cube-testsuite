package mock

import (
	"errors"

	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

type AwsObjectStorage struct {
	buckets map[string]*sharedmodels.ObjectStorageBucket
}

func NewAwsObjectStorage() *AwsObjectStorage {
	return &AwsObjectStorage{buckets: make(map[string]*sharedmodels.ObjectStorageBucket)}
}

func (a *AwsObjectStorage) CreateBucket(bucket *sharedmodels.ObjectStorageBucket) (*sharedmodels.ObjectStorageBucket, error) {
	if bucket.ID == "" {
		bucket.ID = bucket.Name + "-aws"
	}
	a.buckets[bucket.ID] = bucket
	return bucket, nil
}

func (a *AwsObjectStorage) GetBucket(id string) (*sharedmodels.ObjectStorageBucket, error) {
	b, ok := a.buckets[id]
	if !ok {
		return nil, errors.New("bucket not found")
	}
	return b, nil
}

func (a *AwsObjectStorage) ListBuckets() ([]*sharedmodels.ObjectStorageBucket, error) {
	var out []*sharedmodels.ObjectStorageBucket
	for _, b := range a.buckets {
		out = append(out, b)
	}
	return out, nil
}

func (a *AwsObjectStorage) DeleteBucket(id string) error {
	if _, ok := a.buckets[id]; !ok {
		return errors.New("bucket not found")
	}
	delete(a.buckets, id)
	return nil
}
