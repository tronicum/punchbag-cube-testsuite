// This file is now deprecated. All model types have been moved to shared/models/models.go for true sharing.
// Please use only sharedmodels from github.com/tronicum/punchbag-cube-testsuite/shared/models

package models

import (
	sharedmodels "github.com/tronicum/punchbag-cube-testsuite/shared/models"
)

// Deprecated: Use sharedmodels.Cluster instead.
type Cluster = sharedmodels.Cluster

// Deprecated: Use sharedmodels.TestResult instead.
type TestResult = sharedmodels.TestResult

// Deprecated: Use sharedmodels.TestRequest instead.
type TestRequest = sharedmodels.TestRequest

// Deprecated: Use sharedmodels.NodePool instead.
type NodePool = sharedmodels.NodePool

// Deprecated: Use sharedmodels.ClusterCreateRequest instead.
type ClusterCreateRequest = sharedmodels.ClusterCreateRequest

// Deprecated: Use sharedmodels.LoadTestMetrics instead.
type LoadTestMetrics = sharedmodels.LoadTestMetrics
