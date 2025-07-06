package generator

import "fmt"

// GenerateGkeTerraformBlock generates the Terraform block for a GCP GKE cluster (stub)
func GenerateGkeTerraformBlock(props map[string]interface{}) string {
	name := SafeString(props, "name", "example-gke")
	location := SafeString(props, "location", "us-central1")
	nodeCount := SafeInt(props, "nodeCount", 2)
	return fmt.Sprintf("resource \"google_container_cluster\" \"example\" {\n  name     = \"%s\"\n  location = \"%s\"\n  initial_node_count = %d\n  // ...map more fields from JSON as needed\n}", name, location, nodeCount)
}
