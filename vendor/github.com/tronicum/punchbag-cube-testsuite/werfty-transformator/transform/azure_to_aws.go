package transform

import (
	"regexp"
	"strings"
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

// ConvertAzureMonitorToAWSCloudWatch converts Azure Monitor resources to AWS CloudWatch
func ConvertAzureMonitorToAWSCloudWatch(tf string) string {
	// Convert Log Analytics Workspace to CloudWatch Log Group
	reLogAnalytics := regexp.MustCompile(`(?s)resource\s+"azurerm_log_analytics_workspace"[^{]*\{.*?\}`)
	logAnalyticsBlocks := reLogAnalytics.FindAllString(tf, -1)

	var awsBlocks []string
	for _, block := range logAnalyticsBlocks {
		name := extractField(block, "name", "example-log-group")
		retention := extractField(block, "retention_in_days", "7")

		awsBlocks = append(awsBlocks, `resource "aws_cloudwatch_log_group" "example" {
  name              = "`+name+`"
  retention_in_days = `+retention+`
}`)
	}

	// Convert Application Insights to CloudWatch Dashboard
	reAppInsights := regexp.MustCompile(`(?s)resource\s+"azurerm_application_insights"[^{]*\{.*?\}`)
	appInsightsBlocks := reAppInsights.FindAllString(tf, -1)

	for _, block := range appInsightsBlocks {
		name := extractField(block, "name", "example-dashboard")
		awsBlocks = append(awsBlocks, `resource "aws_cloudwatch_dashboard" "example" {
  dashboard_name = "`+name+`"
  dashboard_body = jsonencode({
    widgets = [
      {
        type   = "metric"
        properties = {
          metrics = [
            ["AWS/ApplicationELB", "RequestCount"]
          ]
          period = 300
          stat   = "Sum"
          region = "us-east-1"
          title  = "Request Count"
        }
      }
    ]
  })
}`)
	}

	// Remove Azure resources
	tf = reLogAnalytics.ReplaceAllString(tf, "")
	tf = reAppInsights.ReplaceAllString(tf, "")

	// Add AWS resources
	for _, block := range awsBlocks {
		tf += "\n" + block + "\n"
	}

	return tf
}

// ConvertAzureBudgetToAWSBudget converts Azure Budget to AWS Budget
func ConvertAzureBudgetToAWSBudget(tf string) string {
	reBudget := regexp.MustCompile(`(?s)resource\s+"azurerm_consumption_budget_resource_group"[^{]*\{.*?\}`)
	budgetBlocks := reBudget.FindAllString(tf, -1)

	var awsBlocks []string
	for _, block := range budgetBlocks {
		name := extractField(block, "name", "example-budget")
		amount := extractField(block, "amount", "1000")

		awsBlocks = append(awsBlocks, `resource "aws_budgets_budget" "example" {
  name     = "`+name+`"
  budget_type = "COST"
  limit_amount = "`+amount+`"
  limit_unit = "USD"
  time_unit = "MONTHLY"
  time_period_start = "2025-01-01_00:00"

  notification {
    comparison_operator        = "GREATER_THAN"
    threshold                 = 80
    threshold_type            = "PERCENTAGE"
    notification_type         = "ACTUAL"
    subscriber_email_addresses = ["admin@example.com"]
  }
}`)
	}

	// Remove Azure budget resources
	tf = reBudget.ReplaceAllString(tf, "")

	// Add AWS budget resources
	for _, block := range awsBlocks {
		tf += "\n" + block + "\n"
	}

	return tf
}

func extractField(block, field, def string) string {
	re := regexp.MustCompile(field + `\s*=\s*"?([^"\n]+)"?`)
	match := re.FindStringSubmatch(block)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return def
}
