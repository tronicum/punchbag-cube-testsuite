package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"gopkg.in/yaml.v2"
)

// Format represents the output format type
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatYAML  Format = "yaml"
)

// Formatter handles output formatting for different data types
type Formatter struct {
	format Format
}

// NewFormatter creates a new formatter with the specified format
func NewFormatter(format Format) *Formatter {
	return &Formatter{format: format}
}

// FormatOutput formats and prints the given data
func (f *Formatter) FormatOutput(data interface{}) error {
	switch f.format {
	case FormatJSON:
		return f.formatJSON(data)
	case FormatYAML:
		return f.formatYAML(data)
	case FormatTable:
		return f.formatTable(data)
	default:
		return fmt.Errorf("unsupported format: %s", f.format)
	}
}

// formatJSON formats data as JSON
func (f *Formatter) formatJSON(data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format JSON: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

// formatYAML formats data as YAML
func (f *Formatter) formatYAML(data interface{}) error {
	output, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to format YAML: %w", err)
	}
	fmt.Println(string(output))
	return nil
}

// formatTable formats data as a table
func (f *Formatter) formatTable(data interface{}) error {
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer writer.Flush()

	switch v := data.(type) {
	case map[string]string:
		fmt.Fprintln(writer, "Key\tValue")
		fmt.Fprintln(writer, "---\t-----")
		for key, value := range v {
			fmt.Fprintf(writer, "%s\t%s\n", key, value)
		}
	case map[string]interface{}:
		fmt.Fprintln(writer, "Key\tValue")
		fmt.Fprintln(writer, "---\t-----")
		for key, value := range v {
			fmt.Fprintf(writer, "%s\t%v\n", key, value)
		}
	case []map[string]interface{}:
		if len(v) == 0 {
			fmt.Println("No data found")
			return nil
		}
		
		// Get headers from first item
		headers := make([]string, 0)
		for key := range v[0] {
			headers = append(headers, key)
		}
		
		// Print headers
		for i, header := range headers {
			if i > 0 {
				fmt.Fprint(writer, "\t")
			}
			fmt.Fprint(writer, header)
		}
		fmt.Fprintln(writer)
		
		// Print separator
		for i := range headers {
			if i > 0 {
				fmt.Fprint(writer, "\t")
			}
			fmt.Fprint(writer, "---")
		}
		fmt.Fprintln(writer)
		
		// Print data rows
		for _, row := range v {
			for i, header := range headers {
				if i > 0 {
					fmt.Fprint(writer, "\t")
				}
				fmt.Fprintf(writer, "%v", row[header])
			}
			fmt.Fprintln(writer)
		}
	default:
		// For complex types like slices of structs, convert to a simplified format
		if clusters, ok := convertToClusterTable(data); ok {
			return f.formatClusterTable(clusters, writer)
		}
		if testResults, ok := convertToTestResultTable(data); ok {
			return f.formatTestResultTable(testResults, writer)
		}
		
		// Fallback to JSON for complex types
		return f.formatJSON(data)
	}
	
	return nil
}

// formatClusterTable formats cluster data in table format
func (f *Formatter) formatClusterTable(clusters []map[string]interface{}, writer *tabwriter.Writer) error {
	if len(clusters) == 0 {
		fmt.Println("No clusters found")
		return nil
	}

	// Headers
	fmt.Fprintln(writer, "ID\tName\tProvider\tStatus\tRegion/Location\tCreated")
	fmt.Fprintln(writer, "---\t----\t--------\t------\t---------------\t-------")

	// Rows
	for _, cluster := range clusters {
		id := truncateString(fmt.Sprintf("%v", cluster["id"]), 15)
		name := truncateString(fmt.Sprintf("%v", cluster["name"]), 20)
		provider := fmt.Sprintf("%v", cluster["provider"])
		status := fmt.Sprintf("%v", cluster["status"])
		
		// Get region or location
		regionLocation := ""
		if region, ok := cluster["region"]; ok && region != "" {
			regionLocation = fmt.Sprintf("%v", region)
		} else if location, ok := cluster["location"]; ok && location != "" {
			regionLocation = fmt.Sprintf("%v", location)
		}
		
		// Format created date
		created := ""
		if createdAt, ok := cluster["created_at"]; ok {
			created = formatTime(fmt.Sprintf("%v", createdAt))
		}

		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n",
			id, name, provider, status, regionLocation, created)
	}

	return nil
}

// formatTestResultTable formats test result data in table format
func (f *Formatter) formatTestResultTable(results []map[string]interface{}, writer *tabwriter.Writer) error {
	if len(results) == 0 {
		fmt.Println("No test results found")
		return nil
	}

	// Headers
	fmt.Fprintln(writer, "ID\tCluster ID\tTest Type\tStatus\tDuration\tStarted")
	fmt.Fprintln(writer, "---\t----------\t---------\t------\t--------\t-------")

	// Rows
	for _, result := range results {
		id := truncateString(fmt.Sprintf("%v", result["id"]), 15)
		clusterID := truncateString(fmt.Sprintf("%v", result["cluster_id"]), 15)
		testType := fmt.Sprintf("%v", result["test_type"])
		status := fmt.Sprintf("%v", result["status"])
		duration := fmt.Sprintf("%v", result["duration"])
		
		// Format started date
		started := ""
		if startedAt, ok := result["started_at"]; ok {
			started = formatTime(fmt.Sprintf("%v", startedAt))
		}

		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\n",
			id, clusterID, testType, status, duration, started)
	}

	return nil
}

// Helper functions
func convertToClusterTable(data interface{}) ([]map[string]interface{}, bool) {
	// Try to extract cluster data for table formatting
	switch v := data.(type) {
	case []interface{}:
		var clusters []map[string]interface{}
		for _, item := range v {
			if cluster, ok := item.(map[string]interface{}); ok {
				// Check if it looks like a cluster (has id, name, provider fields)
				if _, hasID := cluster["id"]; hasID {
					if _, hasName := cluster["name"]; hasName {
						if _, hasProvider := cluster["provider"]; hasProvider {
							clusters = append(clusters, cluster)
						}
					}
				}
			}
		}
		return clusters, len(clusters) > 0
	}
	return nil, false
}

func convertToTestResultTable(data interface{}) ([]map[string]interface{}, bool) {
	// Try to extract test result data for table formatting
	switch v := data.(type) {
	case []interface{}:
		var results []map[string]interface{}
		for _, item := range v {
			if result, ok := item.(map[string]interface{}); ok {
				// Check if it looks like a test result (has id, cluster_id, test_type fields)
				if _, hasID := result["id"]; hasID {
					if _, hasClusterID := result["cluster_id"]; hasClusterID {
						if _, hasTestType := result["test_type"]; hasTestType {
							results = append(results, result)
						}
					}
				}
			}
		}
		return results, len(results) > 0
	}
	return nil, false
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatTime(timeStr string) string {
	// Simple time formatting - take first 19 characters (YYYY-MM-DDTHH:MM:SS)
	if len(timeStr) >= 19 {
		return timeStr[:19]
	}
	return timeStr
}

// FormatClusters formats cluster data for output
func FormatClusters(clusters interface{}, format Format) error {
	formatter := NewFormatter(format)
	return formatter.FormatOutput(clusters)
}

// FormatTestResults formats test result data for output
func FormatTestResults(results interface{}, format Format) error {
	formatter := NewFormatter(format)
	return formatter.FormatOutput(results)
}

// FormatError formats error messages consistently
func FormatError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

// FormatSuccess formats success messages consistently
func FormatSuccess(message string) {
	fmt.Printf("✓ %s\n", message)
}

// FormatInfo formats info messages consistently
func FormatInfo(message string) {
	fmt.Printf("ℹ %s\n", message)
}

// FormatWarning formats warning messages consistently
func FormatWarning(message string) {
	fmt.Printf("⚠ %s\n", message)
}
