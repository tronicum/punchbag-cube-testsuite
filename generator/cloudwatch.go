package main

import (
	"fmt"
	"strings"
)

// generateCloudWatchTerraformBlock generates the Terraform block for AWS CloudWatch Metric Alarm
func generateCloudWatchTerraformBlock(props map[string]interface{}) string {
	name := safeString(props, "name", "example-alarm")
	namespace := safeString(props, "namespace", "AWS/EC2")
	metricName := safeString(props, "metricName", "CPUUtilization")
	comparison := safeString(props, "comparisonOperator", "GreaterThanThreshold")
	threshold := safeInt(props, "threshold", 80)
	period := safeInt(props, "period", 300)
	evaluationPeriods := safeInt(props, "evaluationPeriods", 1)
	statistic := safeString(props, "statistic", "Average")
	alarmActions := []string{}
	if a, ok := props["alarmActions"].([]interface{}); ok {
		for _, v := range a {
			alarmActions = append(alarmActions, fmt.Sprintf("\"%s\"", v))
		}
	}
	alarmActionsLine := ""
	if len(alarmActions) > 0 {
		alarmActionsLine = fmt.Sprintf("  alarm_actions = [%s]\n", strings.Join(alarmActions, ", "))
	}
	return fmt.Sprintf(`resource "aws_cloudwatch_metric_alarm" "example" {
  alarm_name          = "%s"
  namespace           = "%s"
  metric_name         = "%s"
  comparison_operator = "%s"
  threshold           = %d
  period              = %d
  evaluation_periods  = %d
  statistic           = "%s"
%s  // ...map more fields from JSON as needed
}`,
		name, namespace, metricName, comparison, threshold, period, evaluationPeriods, statistic, alarmActionsLine)
}
