package models

// Bucket is a generic abstraction for S3-like storage buckets across providers.
type Bucket struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Region       string `json:"region"`
	Provider     string `json:"provider"`
	CreatedAt    string `json:"created_at"`
	StorageClass string `json:"storage_class,omitempty"`
	Tier         string `json:"tier,omitempty"`
}

// ObjectStorage defines generic CRUD operations for object storage buckets.
type ObjectStorage interface {
	CreateBucket(bucket *Bucket) (*Bucket, error)
	GetBucket(id string) (*Bucket, error)
	ListBuckets() ([]*Bucket, error)
	DeleteBucket(id string) error
}
