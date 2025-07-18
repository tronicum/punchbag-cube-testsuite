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
type HetznerS3Mock struct {
	buckets map[string]map[string]S3Object // bucket -> key -> object
	mu      sync.RWMutex
}

// NewHetznerS3Mock creates a new in-memory Hetzner S3 mock.
func NewHetznerS3Mock() *HetznerS3Mock {
	return &HetznerS3Mock{
		buckets: make(map[string]map[string]S3Object),
	}
}

// CreateBucket creates a new bucket.
func (m *HetznerS3Mock) CreateBucket(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.buckets[name]; exists {
		return fmt.Errorf("bucket already exists")
	}
	m.buckets[name] = make(map[string]S3Object)
	return nil
}

// PutObject puts an object into a bucket.
func (m *HetznerS3Mock) PutObject(bucket, key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.buckets[bucket]
	if !ok {
		return fmt.Errorf("bucket not found")
	}
	b[key] = S3Object{Key: key, Value: value}
	return nil
}

// GetObject retrieves an object from a bucket.
func (m *HetznerS3Mock) GetObject(bucket, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.buckets[bucket]
	if !ok {
		return nil, fmt.Errorf("bucket not found")
	}
	obj, ok := b[key]
	if !ok {
		return nil, fmt.Errorf("object not found")
	}
	return obj.Value, nil
}

// ListBuckets lists all buckets.
func (m *HetznerS3Mock) ListBuckets() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []string
	for k := range m.buckets {
		out = append(out, k)
	}
	return out
}

// ListObjects lists all object keys in a bucket.
func (m *HetznerS3Mock) ListObjects(bucket string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	b, ok := m.buckets[bucket]
	if !ok {
		return nil, fmt.Errorf("bucket not found")
	}
	var out []string
	for k := range b {
		out = append(out, k)
	}
	return out, nil
}

// DeleteObject deletes an object from a bucket.
func (m *HetznerS3Mock) DeleteObject(bucket, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.buckets[bucket]
	if !ok {
		return fmt.Errorf("bucket not found")
	}
	if _, ok := b[key]; !ok {
		return fmt.Errorf("object not found")
	}
	delete(b, key)
	return nil
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
