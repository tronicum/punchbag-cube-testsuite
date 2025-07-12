package main

import (
	"fmt"
	"strings"
)

// generateEksTerraformBlock generates the Terraform block for an AWS EKS cluster (stub)
func generateEksTerraformBlock(props map[string]interface{}) string {
	name := safeString(props, "name", "example-eks")
	region := safeString(props, "region", "us-west-2")
	nodeCount := safeInt(props, "nodeCount", 2)
	version := safeString(props, "eksVersion", "1.29")
	instanceType := safeString(props, "instanceType", "t3.medium")
	subnetIds := []string{}
	if s, ok := props["subnetIds"].([]interface{}); ok {
		for _, v := range s {
			subnetIds = append(subnetIds, fmt.Sprintf("\"%s\"", v))
		}
	}
	subnetLine := ""
	if len(subnetIds) > 0 {
		subnetLine = fmt.Sprintf("  subnet_ids = [%s]\n", strings.Join(subnetIds, ", "))
	}
	return fmt.Sprintf(`resource "aws_eks_cluster" "example" {
  name     = "%s"
  region   = "%s"
  version  = "%s"
  vpc_config {
    subnet_ids = [%s]
  }
  node_group {
    instance_types = ["%s"]
    desired_capacity = %d
  }
  %s// ...map more fields from JSON as needed
}`,
		name, region, version, strings.Join(subnetIds, ", "), instanceType, nodeCount, subnetLine)
}
