package mock

import (
	"errors"
	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

type IonosObjectStorage struct {
	buckets map[string]*sharedmodels.ObjectStorageBucket
}

func NewIonosObjectStorage() *IonosObjectStorage {
	return &IonosObjectStorage{buckets: make(map[string]*sharedmodels.ObjectStorageBucket)}
}

func (i *IonosObjectStorage) CreateBucket(bucket *sharedmodels.ObjectStorageBucket) (*sharedmodels.ObjectStorageBucket, error) {
	if bucket.ID == "" {
		bucket.ID = bucket.Name + "-ionos"
	}
	i.buckets[bucket.ID] = bucket
	return bucket, nil
}

func (i *IonosObjectStorage) GetBucket(id string) (*sharedmodels.ObjectStorageBucket, error) {
	b, ok := i.buckets[id]
	if !ok {
		return nil, errors.New("bucket not found")
	}
	return b, nil
}

func (i *IonosObjectStorage) ListBuckets() ([]*sharedmodels.ObjectStorageBucket, error) {
	var out []*sharedmodels.ObjectStorageBucket
	for _, b := range i.buckets {
		out = append(out, b)
	}
	return out, nil
}

func (i *IonosObjectStorage) DeleteBucket(id string) error {
	if _, ok := i.buckets[id]; !ok {
		return errors.New("bucket not found")
	}
	delete(i.buckets, id)
	return nil
}
