package models

type S3Bucket struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Region      string `json:"region"`
	CreatedAt   string `json:"created_at"`
	StorageClass string `json:"storage_class"`
}
