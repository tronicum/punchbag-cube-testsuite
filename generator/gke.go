package main

import "fmt"

// generateGkeTerraformBlock generates the Terraform block for a GCP GKE cluster (stub)
func generateGkeTerraformBlock(props map[string]interface{}) string {
	name := safeString(props, "name", "example-gke")
	location := safeString(props, "location", "us-central1")
	nodeCount := safeInt(props, "nodeCount", 2)
	// TODO: Map more GKE fields as needed
	return fmt.Sprintf(`resource "google_container_cluster" "example" {
  name     = "%s"
  location = "%s"
  initial_node_count = %d
  // ...map more fields from JSON as needed
}`,
		name, location, nodeCount)
}
