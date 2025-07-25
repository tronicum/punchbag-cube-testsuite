package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check environment and dependencies for multitool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Environment diagnostics:")
		fmt.Printf("  OS: %s\n", runtime.GOOS)
		fmt.Printf("  Arch: %s\n", runtime.GOARCH)
		fmt.Printf("  Go version: %s\n", runtime.Version())
		check("git", "--version")
		check("docker", "version")
		check("kubectl", "version", "--client=true")
	},
}

func check(bin string, args ...string) {
	fmt.Printf("  %s: ", bin)
	cmd := exec.Command(bin, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("not found or error: %v\n", err)
	}
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
