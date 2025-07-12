package main

import "fmt"

// generateS3TerraformBlock generates the Terraform block for an AWS S3 bucket
func generateS3TerraformBlock(props map[string]interface{}) string {
	name := safeString(props, "name", "example-s3-bucket")
	acl := safeString(props, "acl", "private")
	versioning := safeBool(props, "versioning", false)
	versioningBlock := ""
	if versioning {
		versioningBlock = "  versioning {\n    enabled = true\n  }\n"
	}
	return fmt.Sprintf(`resource "aws_s3_bucket" "example" {
  bucket = "%s"
  acl    = "%s"
%s  // ...map more fields from JSON as needed
}`,
		name, acl, versioningBlock)
}
