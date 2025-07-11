package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"

	"punchbag-cube-testsuite/client/pkg/api"
)

type Formatter struct {
	format string
}

func NewFormatter(format string) *Formatter {
	return &Formatter{
		format: format,
	}
}

// FormatSimulationResult formats simulation results
func (f *Formatter) FormatSimulationResult(result map[string]interface{}) error {
	switch f.format {
	case "json":
		return f.formatJSON(result)
	case "yaml":
		return f.formatYAML(result)
	default:
		return f.formatTable(result)
	}
}

// FormatProviderInfo formats provider information
func (f *Formatter) FormatProviderInfo(info map[string]interface{}) error {
	switch f.format {
	case "json":
		return f.formatJSON(info)
	case "yaml":
		return f.formatYAML(info)
	default:
		return formatProviderInfoTable(info)
	}
}

// FormatClusterList formats cluster list
func (f *Formatter) FormatClusterList(clusters []map[string]interface{}) error {
	switch f.format {
	case "json":
		return f.formatJSON(clusters)
	case "yaml":
		return f.formatYAML(clusters)
	default:
		return f.formatClusterTable(clusters)
	}
}

// FormatProviderOperation formats provider operation results
func (f *Formatter) FormatProviderOperation(result interface{}) error {
	switch f.format {
	case "json":
		return f.formatJSON(result)
	case "yaml":
		return f.formatYAML(result)
	default:
		return f.formatTable(result)
	}
}

func (f *Formatter) formatJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (f *Formatter) formatYAML(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	return encoder.Encode(data)
}

func (f *Formatter) formatTable(data interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Key", "Value"})

	// Convert data to map for table display
	jsonBytes, _ := json.Marshal(data)
	var dataMap map[string]interface{}
	json.Unmarshal(jsonBytes, &dataMap)

	for key, value := range dataMap {
		valueStr := fmt.Sprintf("%v", value)
		table.Append([]string{key, valueStr})
	}

	table.Render()
	return nil
}

func (f *Formatter) formatClusterTable(clusters []map[string]interface{}) error {
	if len(clusters) == 0 {
		fmt.Println("No clusters found")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Provider", "Status", "Location"})

	for _, cluster := range clusters {
		name := getStringValue(cluster, "name")
		provider := getStringValue(cluster, "provider")
		status := getStringValue(cluster, "status")
		location := getStringValue(cluster, "location")

		table.Append([]string{name, provider, status, location})
	}

	table.Render()
	return nil
}

// Legacy printing functions for backward compatibility
func PrintClusters(clusters []*api.Cluster, format string) {
	switch format {
	case "json":
		printJSON(clusters)
	case "yaml":
		printYAML(clusters)
	default:
		printClustersTable(clusters)
	}
}

func PrintAKSClusters(clusters []*api.AKSCluster, format string) {
	switch format {
	case "json":
		printJSON(clusters)
	case "yaml":
		printYAML(clusters)
	default:
		printAKSClustersTable(clusters)
	}
}

func PrintCluster(cluster *api.Cluster, format string) {
	switch format {
	case "json":
		printJSON(cluster)
	case "yaml":
		printYAML(cluster)
	default:
		printClusterTable(cluster)
	}
}

func PrintAKSCluster(cluster *api.AKSCluster, format string) {
	switch format {
	case "json":
		printJSON(cluster)
	case "yaml":
		printYAML(cluster)
	default:
		printAKSClusterTable(cluster)
	}
}

func PrintTestResults(results []*api.TestResult, format string) {
	switch format {
	case "json":
		printJSON(results)
	case "yaml":
		printYAML(results)
	default:
		printTestResultsTable(results)
	}
}

func PrintAKSTestResults(results []*api.AKSTestResult, format string) {
	switch format {
	case "json":
		printJSON(results)
	case "yaml":
		printYAML(results)
	default:
		printAKSTestResultsTable(results)
	}
}

func PrintTestResult(result *api.TestResult, format string) {
	switch format {
	case "json":
		printJSON(result)
	case "yaml":
		printYAML(result)
	default:
		printTestResultTable(result)
	}
}

func PrintAKSTestResult(result *api.AKSTestResult, format string) {
	switch format {
	case "json":
		printJSON(result)
	case "yaml":
		printYAML(result)
	default:
		printAKSTestResultTable(result)
	}
}

// Helper functions
func printJSON(data interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.Encode(data)
}

func printYAML(data interface{}) {
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	encoder.Encode(data)
}

func printClustersTable(clusters []*api.Cluster) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tPROVIDER\tSTATUS\tCREATED")

	for _, cluster := range clusters {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			cluster.ID,
			cluster.Name,
			cluster.CloudProvider,
			cluster.Status,
			cluster.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}

	w.Flush()
}

func printClusterTable(cluster *api.Cluster) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintf(w, "ID:\t%s\n", cluster.ID)
	fmt.Fprintf(w, "Name:\t%s\n", cluster.Name)
	fmt.Fprintf(w, "Cloud Provider:\t%s\n", cluster.CloudProvider)
	fmt.Fprintf(w, "Status:\t%s\n", cluster.Status)
	fmt.Fprintf(w, "Created:\t%s\n", cluster.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "Updated:\t%s\n", cluster.UpdatedAt.Format("2006-01-02 15:04:05"))

	if cluster.Config != nil {
		fmt.Fprintf(w, "Configuration:\t\n")
		for key, value := range cluster.Config {
			fmt.Fprintf(w, "  %s:\t%v\n", key, value)
		}
	}

	w.Flush()
}

func printAKSClusterTable(cluster *api.AKSCluster) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintf(w, "ID:\t%s\n", cluster.ID)
	fmt.Fprintf(w, "Name:\t%s\n", cluster.Name)
	fmt.Fprintf(w, "Resource Group:\t%s\n", cluster.ResourceGroup)
	fmt.Fprintf(w, "Location:\t%s\n", cluster.Location)
	fmt.Fprintf(w, "Kubernetes Version:\t%s\n", cluster.KubernetesVersion)
	fmt.Fprintf(w, "Status:\t%s\n", cluster.Status)
	fmt.Fprintf(w, "Node Count:\t%d\n", cluster.NodeCount)
	fmt.Fprintf(w, "Created:\t%s\n", cluster.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(w, "Updated:\t%s\n", cluster.UpdatedAt.Format("2006-01-02 15:04:05"))

	if len(cluster.Tags) > 0 {
		fmt.Fprintf(w, "Tags:\t")
		first := true
		for k, v := range cluster.Tags {
			if !first {
				fmt.Fprintf(w, ", ")
			}
			fmt.Fprintf(w, "%s=%s", k, v)
			first = false
		}
		fmt.Fprintf(w, "\n")
	}

	w.Flush()
}

func printAKSTestResultsTable(results []*api.AKSTestResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tCLUSTER ID\tTEST TYPE\tSTATUS\tDURATION\tSTARTED")

	for _, result := range results {
		duration := ""
		if result.Duration > 0 {
			duration = result.Duration.String()
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			result.ID,
			result.ClusterID,
			result.TestType,
			result.Status,
			duration,
			result.StartedAt.Format("2006-01-02 15:04:05"),
		)
	}

	w.Flush()
}

func printAKSTestResultTable(result *api.AKSTestResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintf(w, "ID:\t%s\n", result.ID)
	fmt.Fprintf(w, "Cluster ID:\t%s\n", result.ClusterID)
	fmt.Fprintf(w, "Test Type:\t%s\n", result.TestType)
	fmt.Fprintf(w, "Status:\t%s\n", result.Status)

	if result.Duration > 0 {
		fmt.Fprintf(w, "Duration:\t%s\n", result.Duration.String())
	}

	fmt.Fprintf(w, "Started:\t%s\n", result.StartedAt.Format("2006-01-02 15:04:05"))

	if result.CompletedAt != nil {
		fmt.Fprintf(w, "Completed:\t%s\n", result.CompletedAt.Format("2006-01-02 15:04:05"))
	}

	if result.ErrorMsg != "" {
		fmt.Fprintf(w, "Error:\t%s\n", result.ErrorMsg)
	}

	if len(result.Details) > 0 {
		fmt.Fprintf(w, "Details:\t\n")
		for key, value := range result.Details {
			fmt.Fprintf(w, "  %s:\t%v\n", key, value)
		}
	}

	w.Flush()
}

func formatProviderValidationTable(data map[string]interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Property", "Value"})
	table.SetBorder(false)

	if provider, ok := data["provider"].(string); ok {
		table.Append([]string{"Provider", provider})
	}
	if status, ok := data["status"].(string); ok {
		table.Append([]string{"Status", status})
	}
	if timestamp, ok := data["timestamp"].(string); ok {
		table.Append([]string{"Timestamp", timestamp})
	}

	// Add regions/locations
	if regions, ok := data["regions"].([]interface{}); ok {
		regionStr := ""
		for i, region := range regions {
			if i > 0 {
				regionStr += ", "
			}
			regionStr += fmt.Sprintf("%v", region)
		}
		table.Append([]string{"Regions", regionStr})
	}
	if locations, ok := data["locations"].([]interface{}); ok {
		locationStr := ""
		for i, location := range locations {
			if i > 0 {
				locationStr += ", "
			}
			locationStr += fmt.Sprintf("%v", location)
		}
		table.Append([]string{"Locations", locationStr})
	}

	table.Render()
	return nil
}

func formatProviderInfoTable(data map[string]interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Property", "Value"})
	table.SetBorder(false)

	if provider, ok := data["provider"].(string); ok {
		table.Append([]string{"Provider", provider})
	}
	if name, ok := data["name"].(string); ok {
		table.Append([]string{"Name", name})
	}
	if description, ok := data["description"].(string); ok {
		table.Append([]string{"Description", description})
	}
	if documentation, ok := data["documentation"].(string); ok {
		table.Append([]string{"Documentation", documentation})
	}
	if pricingModel, ok := data["pricing_model"].(string); ok {
		table.Append([]string{"Pricing Model", pricingModel})
	}

	// Add supported features
	if features, ok := data["supported_features"].([]interface{}); ok {
		featureStr := ""
		for i, feature := range features {
			if i > 0 {
				featureStr += ", "
			}
			featureStr += fmt.Sprintf("%v", feature)
		}
		table.Append([]string{"Features", featureStr})
	}

	table.Render()
	return nil
}
		for key, value := range result.Details {
			fmt.Fprintf(w, "  %s:\t%v\n", key, value)
		}
	}

	w.Flush()
}

func formatProviderValidationTable(data map[string]interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Property", "Value"})
	table.SetBorder(false)

	if provider, ok := data["provider"].(string); ok {
		table.Append([]string{"Provider", provider})
	}
	if status, ok := data["status"].(string); ok {
		table.Append([]string{"Status", status})
	}
	if timestamp, ok := data["timestamp"].(string); ok {
		table.Append([]string{"Timestamp", timestamp})
	}

	// Add regions/locations
	if regions, ok := data["regions"].([]interface{}); ok {
		regionStr := ""
		for i, region := range regions {
			if i > 0 {
				regionStr += ", "
			}
			regionStr += fmt.Sprintf("%v", region)
		}
		table.Append([]string{"Regions", regionStr})
	}
	if locations, ok := data["locations"].([]interface{}); ok {
		locationStr := ""
		for i, location := range locations {
			if i > 0 {
				locationStr += ", "
			}
			locationStr += fmt.Sprintf("%v", location)
		}
		table.Append([]string{"Locations", locationStr})
	}

	table.Render()
	return nil
}

func formatProviderInfoTable(data map[string]interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Property", "Value"})
	table.SetBorder(false)

	if provider, ok := data["provider"].(string); ok {
		table.Append([]string{"Provider", provider})
	}
	if name, ok := data["name"].(string); ok {
		table.Append([]string{"Name", name})
	}
	if description, ok := data["description"].(string); ok {
		table.Append([]string{"Description", description})
	}
	if documentation, ok := data["documentation"].(string); ok {
		table.Append([]string{"Documentation", documentation})
	}
	if pricingModel, ok := data["pricing_model"].(string); ok {
		table.Append([]string{"Pricing Model", pricingModel})
	}

	// Add supported features
	if features, ok := data["supported_features"].([]interface{}); ok {
		featureStr := ""
		for i, feature := range features {
			if i > 0 {
				featureStr += ", "
			}
			featureStr += fmt.Sprintf("%v", feature)
		}
		table.Append([]string{"Features", featureStr})
	}

	table.Render()
	return nil
}
