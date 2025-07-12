package generator

import "fmt"

// GenerateS3TerraformBlock generates the Terraform block for an AWS S3 bucket
func GenerateS3TerraformBlock(props map[string]interface{}) string {
	name := SafeString(props, "name", "example-s3-bucket")
	acl := SafeString(props, "acl", "private")
	versioning := SafeBool(props, "versioning", false)
	versioningBlock := ""
	if versioning {
		versioningBlock = "  versioning {\n    enabled = true\n  }\n"
	}
	return fmt.Sprintf("resource \"aws_s3_bucket\" \"example\" {\n  bucket = \"%s\"\n  acl    = \"%s\"\n%s  // ...map more fields from JSON as needed\n}", name, acl, versioningBlock)
}
