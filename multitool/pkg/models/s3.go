package models

import sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"

type S3Bucket = sharedmodels.ObjectStorageBucket

// type BlobStorage = sharedmodels.ObjectStorageBucket // Uncomment if needed for Azure
// type GCSBucket = sharedmodels.ObjectStorageBucket   // Uncomment if needed for GCP

type S3BucketPolicy = sharedmodels.ObjectStoragePolicy

// type S3BucketPolicyStatement = sharedmodels.ObjectStorageStatement // Uncomment if needed
// type S3BucketVersioning = sharedmodels.ObjectStorageRule           // Uncomment if needed
// type S3BucketLifecycleRule = sharedmodels.ObjectStorageRule        // Uncomment if needed
// type S3LifecycleTransition = sharedmodels.ObjectStorageRule        // Uncomment if needed
