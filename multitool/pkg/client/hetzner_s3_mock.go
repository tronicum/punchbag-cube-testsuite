package client

import (
	"fmt"
	"net/http"
	"sync"
)

// S3Object represents a simple S3 object (key, value).
type S3Object struct {
	Key   string
	Value []byte
}

// HetznerS3Mock simulates a minimal Hetzner S3 API in-memory.
// Deprecated: Use shared/providers/hetzner/s3_mock.go instead.
// This file is now a stub for legacy references.
	for k := range b {
// Deprecated: Use shared/providers/hetzner/s3_mock.go instead.
// This file is intentionally left blank. All Hetzner S3 mock logic has been migrated to shared/providers/hetzner/s3_mock.go.
}

// DeleteBucket deletes a bucket and all its objects.
func (m *HetznerS3Mock) DeleteBucket(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.buckets[name]; !ok {
		return fmt.Errorf("bucket not found")
	}
	delete(m.buckets, name)
	return nil
}

// ServeHTTP implements http.Handler for S3-like API endpoints (minimal, for simulation mode).
func (m *HetznerS3Mock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement minimal S3 REST API simulation (bucket/object CRUD via HTTP)
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Hetzner S3 simulation not yet implemented"))
}
