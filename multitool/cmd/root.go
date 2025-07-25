package cmd

import (
	   "fmt"
	   "os"
	   "github.com/spf13/cobra"
	   "text/template"
	   "strings"
	   "runtime"
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
	   awsCmd.Annotations = map[string]string{"group": "Cloud Management Commands"}
	   azureCmd.Annotations = map[string]string{"group": "Cloud Management Commands"}
	   gcpCmd.Annotations = map[string]string{"group": "Cloud Management Commands"}
	   hetznerCmd.Annotations = map[string]string{"group": "Cloud Management Commands"}
	   objectStorageCmd.Annotations = map[string]string{"group": "Cloud ObjectStorage (S3) Commands"}
	rootCmd.PersistentFlags().StringVar(&proxyServer, "server", "", "If set, forward all resource management requests to this cube-server URL (proxy/simulation mode)")

	// Register only the correct top-level commands, matching the new CLI tree structure
rootCmd.AddCommand(awsCmd)           // mt aws ...
rootCmd.AddCommand(azureCmd)         // mt azure ...
rootCmd.AddCommand(gcpCmd)           // mt gcp ...
rootCmd.AddCommand(hetznerCmd)       // mt hetzner ...
rootCmd.AddCommand(dockerCmd)        // mt docker ...
rootCmd.AddCommand(localCmd)         // mt local ...
rootCmd.AddCommand(configCmd)        // mt config ...
rootCmd.AddCommand(testCmd)          // mt test ...
rootCmd.AddCommand(scaffoldCmd)      // mt scaffold ...
rootCmd.AddCommand(objectStorageCmd) // mt objectstorage ...
// k8sctl and k8s-manage are now top-level via their own files
	   // simulate-hetzner-s3 removed per user request
	   rootCmd.SetHelpFunc(printGroupedHelp)
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
			   osys := runtime.GOOS
			   pkgMgr := "Unsupported"
			   if osys == "darwin" {
					   pkgMgr = "Homebrew"
			   } else if osys == "linux" {
					   pkgMgr = "apt or rpm"
			   }
			   fmt.Printf("OS: %s\nPackage Manager: %s\n", osys, pkgMgr)
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
			   var installCmd string
			   if osys == "darwin" {
					   installCmd = fmt.Sprintf("brew install %s", packageName)
			   } else if osys == "linux" {
					   installCmd = fmt.Sprintf("apt install %s or rpm install %s", packageName, packageName)
			   } else {
					   installCmd = "Unsupported"
			   }
			   fmt.Printf("Package: %s\nOS: %s\nCommand: %s\n", packageName, osys, installCmd)

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

var listPackagesCmd = &cobra.Command{
	   Use:   "list-packages",
	 SilenceUsage: true,
	 DisableAutoGenTag: true,
	   Short: "List installed packages based on the detected OS",
	   Run: func(cmd *cobra.Command, cmdArgs []string) {
			   // ...existing code...
	   },
}

// --- Custom help grouping ---
func groupCommandsByAnnotation(cmds []*cobra.Command, annotation string) map[string][]*cobra.Command {
	   groups := make(map[string][]*cobra.Command)
	   for _, c := range cmds {
			   group := c.Annotations[annotation]
			   if group == "" {
					   group = "Other Commands"
			   }
			   groups[group] = append(groups[group], c)
	   }
	   return groups
}

func printGroupedHelp(cmd *cobra.Command, args []string) {
	   cmds := cmd.Commands()
	   groups := groupCommandsByAnnotation(cmds, "group")
	   tmpl := `{{.Long}}

Usage:
  {{.UseLine}}

{{if .HasAvailableSubCommands}}Available Commands:
{{range $group, $cmds := .Groups}}
{{$group}}:
{{range $cmd := $cmds}}  {{$cmd.Name | printf "%-20s"}} {{$cmd.Short}}
{{end}}{{end}}{{end}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}
`
	   data := struct {
			   *cobra.Command
			   Groups map[string][]*cobra.Command
	   }{cmd, groups}
	   t := template.Must(template.New("help").Funcs(template.FuncMap{"trimTrailingWhitespaces": strings.TrimRight}).Parse(tmpl))
	   t.Execute(os.Stdout, data)
}
