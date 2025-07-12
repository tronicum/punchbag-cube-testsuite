package main

import "fmt"

// generateAwsBudgetTerraformBlock generates the Terraform block for AWS Budgets
func generateAwsBudgetTerraformBlock(props map[string]interface{}) string {
	name := safeString(props, "name", "example-budget")
	amount := safeInt(props, "amount", 1000)
	period := safeString(props, "period", "MONTHLY")
	return fmt.Sprintf(`resource "aws_budgets_budget" "example" {
  name              = "%s"
  budget_type       = "COST"
  limit_amount      = "%d"
  limit_unit        = "USD"
  time_period_start = "2023-01-01_00:00"
  time_unit         = "%s"
  // ...map more fields from JSON as needed
}`,
		name, amount, period)
}
