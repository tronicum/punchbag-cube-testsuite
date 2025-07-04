package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"punchbag-cube-testsuite/client/pkg/api"

	"gopkg.in/yaml.v3"
)

// PrintClusters prints a list of clusters in the specified format
func PrintClusters(clusters []*api.AKSCluster, format string) {
	switch format {
	case "json":
		printJSON(clusters)
	case "yaml":
		printYAML(clusters)
	default:
		printClustersTable(clusters)
	}
}

// PrintCluster prints a single cluster in the specified format
func PrintCluster(cluster *api.AKSCluster, format string) {
	switch format {
	case "json":
		printJSON(cluster)
	case "yaml":
		printYAML(cluster)
	default:
		printClusterTable(cluster)
	}
}

// PrintTestResults prints a list of test results in the specified format
func PrintTestResults(results []*api.AKSTestResult, format string) {
	switch format {
	case "json":
		printJSON(results)
	case "yaml":
		printYAML(results)
	default:
		printTestResultsTable(results)
	}
}

// PrintTestResult prints a single test result in the specified format
func PrintTestResult(result *api.AKSTestResult, format string) {
	switch format {
	case "json":
		printJSON(result)
	case "yaml":
		printYAML(result)
	default:
		printTestResultTable(result)
	}
}

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

func printClustersTable(clusters []*api.AKSCluster) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tRESOURCE GROUP\tLOCATION\tSTATUS\tNODES\tCREATED")
	
	for _, cluster := range clusters {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%d\t%s\n",
			cluster.ID,
			cluster.Name,
			cluster.ResourceGroup,
			cluster.Location,
			cluster.Status,
			cluster.NodeCount,
			cluster.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}
	
	w.Flush()
}

func printClusterTable(cluster *api.AKSCluster) {
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

func printTestResultsTable(results []*api.AKSTestResult) {
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

func printTestResultTable(result *api.AKSTestResult) {
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
	
	w.Flush()
	
	if len(result.Details) > 0 {
		fmt.Println("\nTest Details:")
		printJSON(result.Details)
	}
}
