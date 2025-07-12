package generator

import "fmt"

// Declarative schema for resource validation
var ResourceSchemas = map[string]map[string][]string{
	"azure": {
		"aks":          {"name", "location", "resourceGroup", "nodeCount"},
		"monitor":      {"name", "resourceGroup", "severity", "criteria"},
		"loganalytics": {"name", "location", "resourceGroup", "sku", "retentionInDays"},
		"appinsights":  {"name", "location", "resourceGroup", "applicationType"},
	},
	"aws": {
		"eks":        {"name", "region", "nodeCount"},
		"s3":         {"name"},
		"cloudwatch": {"name", "namespace", "metricName", "comparisonOperator", "threshold", "period", "evaluationPeriods", "statistic"},
		"budget":     {"name", "amount", "period"},
	},
	"gcp": {
		"gke": {"name", "location", "nodeCount"},
	},
}

// ValidateResourceProperties checks required fields for each resource/provider using the declarative schema
func ValidateResourceProperties(provider, resourceType string, props map[string]interface{}) error {
	schema, ok := ResourceSchemas[provider]
	if !ok {
		return fmt.Errorf("unknown provider: %s", provider)
	}
	required, ok := schema[resourceType]
	if !ok {
		return fmt.Errorf("unknown resource type: %s/%s", provider, resourceType)
	}
	missing := []string{}
	for _, k := range required {
		if _, ok := props[k]; !ok {
			missing = append(missing, k)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required fields for %s/%s: %v", provider, resourceType, missing)
	}
	return nil
}
