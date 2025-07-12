package generator

import "fmt"

// GenerateAwsBudgetTerraformBlock generates the Terraform block for AWS Budgets
func GenerateAwsBudgetTerraformBlock(props map[string]interface{}) string {
	name := SafeString(props, "name", "example-budget")
	amount := SafeInt(props, "amount", 1000)
	period := SafeString(props, "period", "MONTHLY")
	return fmt.Sprintf("resource \"aws_budgets_budget\" \"example\" {\n  name = \"%s\"\n  budget_type = \"COST\"\n  limit_amount = \"%d\"\n  limit_unit = \"USD\"\n  time_period_start = \"2023-01-01_00:00\"\n  time_unit = \"%s\"\n  // ...map more fields from JSON as needed\n}", name, amount, period)
}
