package generator

import "fmt"

// GenerateEksTerraformBlock generates the Terraform block for an EKS cluster from properties
func GenerateEksTerraformBlock(props map[string]interface{}) string {
	name := SafeString(props, "name", "example-eks")
	region := SafeString(props, "region", "us-west-2")
	nodeCount := SafeInt(props, "nodeCount", 3)
	// ...add more fields as needed...
	return fmt.Sprintf(`resource "aws_eks_cluster" "example" {
  name     = "%s"
  region   = "%s"
  // ...map more fields from JSON as needed
}

resource "aws_eks_node_group" "example" {
  cluster_name    = aws_eks_cluster.example.name
  node_group_name = "example"
  node_count      = %d
  // ...map more fields from JSON as needed
}
`, name, region, nodeCount)
}
