package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	gitCommit string
	buildDate string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show build version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mt version: %s\n", version)
		fmt.Printf("commit: %s\n", gitCommit)
		fmt.Printf("built: %s\n", buildDate)
		fmt.Printf("go: %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

// version is set at build time via -ldflags
var version = "dev"
