package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"text/tabwriter"

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

// Rename CLI tool to multitool
var rootCmd = &cobra.Command{
	Use:   "multitool",
	Short: "A CLI tool for cloud management and system operations",
	Long:  "Multitool is a CLI tool designed for managing cloud resources, installing packages, and handling Docker registries.",
}

// Execute the root command
func Execute() {
	fmt.Println("Executing root command...")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing root command:", err)
	}
}

var proxyServer string

// Initialize commands
func init() {
	rootCmd.PersistentFlags().StringVar(&proxyServer, "server", "", "If set, forward all resource management requests to this cube-server URL (proxy/simulation mode)")
	rootCmd.AddCommand(manageCloudCmd)
	rootCmd.AddCommand(osDetectCmd)
	rootCmd.AddCommand(packageInstallCmd)
	rootCmd.AddCommand(dockerRegistryCmd)
	rootCmd.AddCommand(listPackagesCmd)

	// Add enhanced cluster and test commands
	rootCmd.AddCommand(clusterCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(configCmd)

	// Legacy k8s commands (deprecated)
	k8sCmd.AddCommand(k8sGetCmd)
	k8sCmd.AddCommand(k8sCreateCmd)
	k8sCmd.AddCommand(k8sDeleteCmd)
	rootCmd.AddCommand(k8sCmd)
}

// Add basic cloud management functionality
var manageCloudCmd = &cobra.Command{
	Use:   "manage-cloud",
	Short: "Manage cloud resources",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Managing cloud resources...")
	},
}

// Add Kubernetes cluster management commands
var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Manage Kubernetes clusters",
}

var k8sGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Kubernetes cluster information",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please specify the cloud provider: aws or azure")
			return
		}
		provider := args[0]
		fmt.Printf("Fetching Kubernetes cluster information for provider: %s\n", provider)
		// Add logic to fetch cluster information
	},
}

var k8sCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please specify the cloud provider: aws or azure")
			return
		}
		provider := args[0]
		fmt.Printf("Creating Kubernetes cluster for provider: %s\n", provider)
		// Add logic to create cluster
	},
}

var k8sDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please specify the cloud provider: aws or azure")
			return
		}
		provider := args[0]
		fmt.Printf("Deleting Kubernetes cluster for provider: %s\n", provider)
		// Add logic to delete cluster
	},
}

var osDetectCmd = &cobra.Command{
	Use:   "os-detect",
	Short: "Detect the operating system and package manager",
	Run: func(cmd *cobra.Command, args []string) {
		os := runtime.GOOS
		data := map[string]string{
			"OS":              os,
			"Package Manager": "Unsupported",
		}
		if os == "darwin" {
			data["Package Manager"] = "Homebrew"
		} else if os == "linux" {
			data["Package Manager"] = "apt or rpm"
		}
		formatOutput(data, "table") // Replace "table" with "json" or "yaml" as needed
	},
}

var packageInstallCmd = &cobra.Command{
	Use:   "install-package",
	Short: "Install a package based on the detected OS",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify a package to install.")
			return
		}
		packageName := args[0]
		os := runtime.GOOS
		data := map[string]string{
			"Package": packageName,
			"OS":      os,
			"Command": "Unsupported",
		}
		if os == "darwin" {
			data["Command"] = fmt.Sprintf("brew install %s", packageName)
		} else if os == "linux" {
			data["Command"] = fmt.Sprintf("apt install %s or rpm install %s", packageName, packageName)
		}
		formatOutput(data, "table") // Replace "table" with "json" or "yaml" as needed
	},
}

var dockerRegistryCmd = &cobra.Command{
	Use:   "docker-registry",
	Short: "Manage Docker registries",
	Run: func(cmd *cobra.Command, args []string) {
		data := map[string]string{
			"Command": "Docker registry management commands will be implemented here.",
		}
		formatOutput(data, "table") // Replace "table" with "json" or "yaml" as needed
	},
}

var dockerRegistryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Docker registries",
	Run: func(cmd *cobra.Command, args []string) {
		data := map[string]string{
			"Registry": "Docker Hub",
			"Status":   "Logged In",
		}
		formatOutput(data, "table") // Replace "table" with "json" or "yaml" as needed
	},
}

var dockerRegistryLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to a Docker registry",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			fmt.Println("Usage: docker-registry login <registry> <username>")
			return
		}
		registry := args[0]
		username := args[1]
		fmt.Printf("Logging in to %s as %s\n", registry, username)
		// Add logic for Docker login
	},
}

var dockerRegistryLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from a Docker registry",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Usage: docker-registry logout <registry>")
			return
		}
		registry := args[0]
		fmt.Printf("Logging out from %s\n", registry)
		// Add logic for Docker logout
	},
}

var listPackagesCmd = &cobra.Command{
	Use:   "list-packages",
	Short: "List installed packages based on the detected OS",
	Run: func(cmd *cobra.Command, cmdArgs []string) {
		os := runtime.GOOS
		var command string
		var args []string
		if os == "darwin" {
			command = "brew"
			args = []string{"list"}
		} else if os == "linux" {
			if _, err := exec.LookPath("apt"); err == nil {
				command = "apt"
				args = []string{"list", "--installed"}
			} else if _, err := exec.LookPath("rpm"); err == nil {
				command = "rpm"
				args = []string{"-qa"}
			} else if _, err := exec.LookPath("pacman"); err == nil {
				command = "pacman"
				args = []string{"-Q"}
			} else {
				fmt.Println("No supported package manager found.")
				return
			}
		} else if os == "windows" {
			if _, err := exec.LookPath("choco"); err == nil {
				command = "choco"
				args = []string{"list", "-lo"}
			} else if _, err := exec.LookPath("winget"); err == nil {
				command = "winget"
				args = []string{"list"}
			} else {
				fmt.Println("No supported package manager found.")
				return
			}
		} else {
			fmt.Println("Unsupported OS")
			return
		}
		output, err := exec.Command(command, args...).Output()
		if err != nil {
			fmt.Println("Error executing command:", err)
			return
		}
		data := map[string]string{
			"Command": fmt.Sprintf("%s %s", command, args),
			"Output":  string(output),
		}
		formatOutput(data, "table") // Replace "table" with "json" or "yaml" as needed
	},
}

// Add a generate-terraform subcommand for multicloud codegen
var generateTerraformCmd = &cobra.Command{
	Use:   "generate-terraform --input <input.json> --output <output.tf> --provider <provider>",
	Short: "Generate Terraform code for a resource (multicloud)",
	Long: `Generate Terraform code for a resource using the multicloud generator.

Examples:
  multitool generate-terraform --input test_aks.json --output test_aks.tf --provider azure
  multitool generate-terraform --input test_eks.json --output test_eks.tf --provider aws
  multitool generate-terraform --input test_gke.json --output test_gke.tf --provider gcp`,
	Run: func(cmd *cobra.Command, args []string) {
		input, _ := cmd.Flags().GetString("input")
		output, _ := cmd.Flags().GetString("output")
		provider, _ := cmd.Flags().GetString("provider")
		if input == "" || output == "" || provider == "" {
			fmt.Println("--input, --output, and --provider are required")
			os.Exit(1)
		}
		// Call the generator binary with correct args
		cmdline := []string{"run", "generator/main.go", "--generate-terraform", "--input", input, "--output", output, "--provider", provider}
		c := exec.Command("go", cmdline...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			fmt.Printf("Terraform generation failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Terraform code written to %s\n", output)
	},
}

func init() {
	dockerRegistryCmd.AddCommand(dockerRegistryListCmd)
	dockerRegistryCmd.AddCommand(dockerRegistryLoginCmd)
	dockerRegistryCmd.AddCommand(dockerRegistryLogoutCmd)
	generateTerraformCmd.Flags().String("input", "", "Input JSON file")
	generateTerraformCmd.Flags().String("output", "", "Output Terraform file")
	generateTerraformCmd.Flags().String("provider", "", "Cloud provider: azure|aws|gcp")
	rootCmd.AddCommand(generateTerraformCmd)
}

func formatOutput(data interface{}, format string) {
	switch format {
	case "json":
		output, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			fmt.Println("Error formatting JSON:", err)
			return
		}
		fmt.Println(string(output))
	case "yaml":
		output, err := yaml.Marshal(data)
		if err != nil {
			fmt.Println("Error formatting YAML:", err)
			return
		}
		fmt.Println(string(output))
	case "table":
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(writer, "Key\tValue")
		for k, v := range data.(map[string]string) {
			fmt.Fprintf(writer, "%s\t%s\n", k, v)
		}
		writer.Flush()
	default:
		fmt.Println("Unsupported format. Use 'json', 'yaml', or 'table'.")
	}
}
