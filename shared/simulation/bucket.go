package simulation

import (
	   "encoding/json"
	   "fmt"
	   "os"
	   "sync"
)

// BucketStore handles bucket simulation state and persistence
// Used by SimulationService for all bucket operations

type BucketStore struct {
	mu         sync.Mutex
	buckets    map[string]map[string]interface{} // provider -> bucketName -> bucketInfo
	persistPath string
}

func NewBucketStore(persistPath string) *BucketStore {
	bs := &BucketStore{
		buckets: make(map[string]map[string]interface{}),
		persistPath: persistPath,
	}
	bs.Load()
	return bs
}

func (bs *BucketStore) Load() {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	f, err := os.Open(bs.persistPath)
	if err != nil {
		return // no file yet
	}
	defer f.Close()
	_ = json.NewDecoder(f).Decode(&bs.buckets)
}

func (bs *BucketStore) Save() {
	   f, err := os.Create(bs.persistPath)
	   if err != nil {
			   fmt.Fprintf(os.Stderr, "[BUCKET DEBUG] Failed to write persist file: %s: %v\n", bs.persistPath, err)
			   return
	   }
	   defer f.Close()
	   if err := json.NewEncoder(f).Encode(bs.buckets); err != nil {
			   fmt.Fprintf(os.Stderr, "[BUCKET DEBUG] Failed to encode buckets to persist file: %s: %v\n", bs.persistPath, err)
	   } else {
			   fmt.Fprintf(os.Stderr, "[BUCKET DEBUG] Persisted buckets to file: %s\n", bs.persistPath)
	   }
}

func (bs *BucketStore) Create(provider, name, region string) map[string]interface{} {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if bs.buckets[provider] == nil {
		bs.buckets[provider] = make(map[string]interface{})
	}
	bucket := map[string]interface{}{
		"bucket":   name,
		"provider": provider,
		"region":   region,
		"status":   "created",
	}
	bs.buckets[provider][name] = bucket
	bs.Save()
	return bucket
}

func (bs *BucketStore) Delete(provider, name string) (bool, map[string]interface{}) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	if bs.buckets[provider] == nil {
		return false, map[string]interface{}{"error": "provider not found"}
	}
	if _, ok := bs.buckets[provider][name]; !ok {
		return false, map[string]interface{}{"error": "bucket not found"}
	}
	delete(bs.buckets[provider], name)
	bs.Save()
	return true, map[string]interface{}{"bucket": name, "status": "deleted"}
}

func (bs *BucketStore) List(provider string) []map[string]interface{} {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	buckets := make([]map[string]interface{}, 0)
	for _, b := range bs.buckets[provider] {
		if bucket, ok := b.(map[string]interface{}); ok {
			buckets = append(buckets, bucket)
		}
	}
	return buckets
}
