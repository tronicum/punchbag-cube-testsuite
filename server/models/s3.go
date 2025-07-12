package models

type S3BucketPolicy struct {
	Version   string                    `json:"version"`
	Statement []S3BucketPolicyStatement `json:"statement"`
}

type S3BucketPolicyStatement struct {
	Effect    string                 `json:"effect"`
	Principal map[string]interface{} `json:"principal"`
	Action    []string               `json:"action"`
	Resource  []string               `json:"resource"`
	Condition map[string]interface{} `json:"condition,omitempty"`
}

type S3BucketVersioning struct {
	Enabled bool `json:"enabled"`
}

type S3BucketLifecycleRule struct {
	ID                          string                  `json:"id"`
	Prefix                      string                  `json:"prefix,omitempty"`
	Status                      string                  `json:"status"`
	Transitions                 []S3LifecycleTransition `json:"transitions,omitempty"`
	ExpirationDays              int                     `json:"expiration_days,omitempty"`
	NoncurrentVersionExpiration int                     `json:"noncurrent_version_expiration,omitempty"`
}

type S3LifecycleTransition struct {
	Days         int    `json:"days"`
	StorageClass string `json:"storage_class"`
}

type S3Bucket struct {
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	Region       string                  `json:"region"`
	Provider     string                  `json:"provider"`
	CreatedAt    string                  `json:"created_at"`
	StorageClass string                  `json:"storage_class"`
	Policy       *S3BucketPolicy         `json:"policy,omitempty"`
	Versioning   *S3BucketVersioning     `json:"versioning,omitempty"`
	Lifecycle    []S3BucketLifecycleRule `json:"lifecycle,omitempty"`
}
