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
	Use:   "mt",
	Short: "A CLI tool for cloud management and system operations",
	Long:  "mt is a CLI tool designed for managing cloud resources, installing packages, and handling Docker registries.",
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

// Register only the correct top-level commands, matching the new CLI tree structure
rootCmd.AddCommand(awsCmd)      // mt aws ...
rootCmd.AddCommand(azureCmd)    // mt azure ...
rootCmd.AddCommand(gcpCmd)      // mt gcp ...
rootCmd.AddCommand(hetznerCmd)  // mt hetzner ...
rootCmd.AddCommand(dockerCmd)   // mt docker ...
rootCmd.AddCommand(k8sCmd)      // mt k8s ...
rootCmd.AddCommand(localCmd)    // mt local ...
rootCmd.AddCommand(configCmd)   // mt config ...
rootCmd.AddCommand(testCmd)     // mt test ...
rootCmd.AddCommand(scaffoldCmd) // mt scaffold ...
}

// Add basic cloud management functionality
var manageCloudCmd = &cobra.Command{
	Use:   "manage-cloud",
	Short: "Manage cloud resources",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Managing cloud resources...")
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
		osys := runtime.GOOS
		data := map[string]string{
			"Package": packageName,
			"OS":      osys,
			"Command": "Unsupported",
		}
		var installCmd string
		if osys == "darwin" {
			installCmd = fmt.Sprintf("brew install %s", packageName)
			data["Command"] = installCmd
		} else if osys == "linux" {
			installCmd = fmt.Sprintf("apt install %s or rpm install %s", packageName, packageName)
			data["Command"] = installCmd
		}
		formatOutput(data, "table") // Replace "table" with "json" or "yaml" as needed

		relink, _ := cmd.Flags().GetBool("relink")
		if relink {
			mtPath, err := os.Executable()
			if err != nil {
				fmt.Println("Could not determine mt binary path:", err)
				return
			}
			symlinkPath := "/usr/local/bin/mt"
			_ = os.Remove(symlinkPath) // Remove if exists
			err = os.Symlink(mtPath, symlinkPath)
			if err != nil {
				fmt.Printf("Failed to create symlink at %s: %v\n", symlinkPath, err)
			} else {
				fmt.Printf("Symlinked mt binary to %s\n", symlinkPath)
			}
		}
	},
}

// ...existing code...

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
