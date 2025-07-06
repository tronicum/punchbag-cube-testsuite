package transform

import (
	"regexp"
)

// ConvertAzureBlobToAWSS3 maps Azure Blob Storage Terraform to AWS S3 Terraform
func ConvertAzureBlobToAWSS3(tf string) string {
	reAccount := regexp.MustCompile(`(?s)resource\s+\"azurerm_storage_account\"[^{]*\{.*?\}`)
	accountBlocks := reAccount.FindAllString(tf, -1)
	var awsBlocks []string
	for _, block := range accountBlocks {
		name := extractField(block, "name", "examplestorageacct")
		region := extractField(block, "location", "us-east-1")
		awsBlocks = append(awsBlocks, `resource "aws_s3_bucket" "example" {
  bucket = "`+name+`"
  region = "`+region+`"
}`)
	}
	// Remove all azurerm_storage_account and azurerm_storage_container blocks

	tf = reAccount.ReplaceAllString(tf, "")
	reContainer := regexp.MustCompile(`(?s)resource\s+\"azurerm_storage_container\"[^{]*\{.*?\}`)

	tf = reContainer.ReplaceAllString(tf, "")
	// Add the new AWS S3 blocks at the end
	if len(awsBlocks) > 0 {
		tf += "\n" + awsBlocks[0] + "\n"
	}
	return tf
}

func extractField(block, field, def string) string {
	re := regexp.MustCompile(field + `\s*=\s*\"([^\"]+)\"`)
	match := re.FindStringSubmatch(block)
	if len(match) > 1 {
		return match[1]
	}
	return def
}
