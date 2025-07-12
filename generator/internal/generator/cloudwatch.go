package generator

import (
	"fmt"
)

// GenerateCloudWatchTerraformBlock generates the Terraform block for AWS CloudWatch Metric Alarm
func GenerateCloudWatchTerraformBlock(props map[string]interface{}) string {
	name := SafeString(props, "name", "example-alarm")
	namespace := SafeString(props, "namespace", "AWS/EC2")
	metricName := SafeString(props, "metricName", "CPUUtilization")
	comparison := SafeString(props, "comparisonOperator", "GreaterThanThreshold")
	threshold := SafeInt(props, "threshold", 80)
	period := SafeInt(props, "period", 300)
	evaluationPeriods := SafeInt(props, "evaluationPeriods", 1)
	statistic := SafeString(props, "statistic", "Average")
	return fmt.Sprintf("resource \"aws_cloudwatch_metric_alarm\" \"example\" {\n  alarm_name = \"%s\"\n  namespace = \"%s\"\n  metric_name = \"%s\"\n  comparison_operator = \"%s\"\n  threshold = %d\n  period = %d\n  evaluation_periods = %d\n  statistic = \"%s\"\n  // ...map more fields from JSON as needed\n}", name, namespace, metricName, comparison, threshold, period, evaluationPeriods, statistic)
}
