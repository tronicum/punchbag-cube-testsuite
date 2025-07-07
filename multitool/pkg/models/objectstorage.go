package models

import sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"

// Bucket is now an alias for the canonical ObjectStorageBucket in shared models.
type Bucket = sharedmodels.ObjectStorageBucket

// ObjectStorage defines generic CRUD operations for object storage buckets.
type ObjectStorage interface {
	CreateBucket(bucket *Bucket) (*Bucket, error)
	GetBucket(id string) (*Bucket, error)
	ListBuckets() ([]*Bucket, error)
	DeleteBucket(id string) error
}
