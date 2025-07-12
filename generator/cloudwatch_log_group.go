package main

import "fmt"

// generateCloudWatchLogGroupTerraformBlock generates the Terraform block for AWS CloudWatch Log Group
func generateCloudWatchLogGroupTerraformBlock(props map[string]interface{}) string {
	name := safeString(props, "name", "example-log-group")
	retention := safeInt(props, "retentionInDays", 14)
	return fmt.Sprintf(`resource "aws_cloudwatch_log_group" "example" {
  name              = "%s"
  retention_in_days = %d
  // ...map more fields from JSON as needed
}`,
		name, retention)
}
