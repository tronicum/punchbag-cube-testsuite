package models

type S3Bucket struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Region       string `json:"region"`
	Provider     string `json:"provider"`
	CreatedAt    string `json:"created_at"`
	StorageClass string `json:"storage_class"`
}

type BlobStorage struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Region    string `json:"region"`
	Provider  string `json:"provider"`
	CreatedAt string `json:"created_at"`
	Tier      string `json:"tier"`
}

type GCSBucket struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Region       string `json:"region"`
	Provider     string `json:"provider"`
	CreatedAt    string `json:"created_at"`
	StorageClass string `json:"storage_class"`
}
