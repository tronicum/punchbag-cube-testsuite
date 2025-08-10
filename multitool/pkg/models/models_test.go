package models

import "testing"

func TestModelsBasic(t *testing.T) {
	// Just check that type aliasing works
	var _ Cluster
	var _ TestResult
	var _ TestRequest
	var _ NodePool
	var _ ClusterCreateRequest
	var _ LoadTestMetrics
}
