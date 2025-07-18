package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/tronicum/punchbag-cube-testsuite/multitool/pkg/client"
)

var simHetznerS3Cmd = &cobra.Command{
	Use:   "simulate-hetzner-s3",
	Short: "Start a local Hetzner S3 simulation server",
	Run: func(cmd *cobra.Command, args []string) {
		port := "8081"
		if len(args) > 0 {
			port = args[0]
		}
		mock := client.NewHetznerS3Mock()
		fmt.Printf("Starting Hetzner S3 simulation on http://localhost:%s\n", port)
		http.Handle("/", mock)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start server: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(simHetznerS3Cmd)
}
