package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var listExamplesCmd = &cobra.Command{
	Use:   "list-examples",
	Short: "List all available example YAML/JSON configs",
	Run: func(cmd *cobra.Command, args []string) {
		dir := "examples"
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read examples/: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Available examples:")
		for _, f := range files {
			if !f.IsDir() && (filepath.Ext(f.Name()) == ".yaml" || filepath.Ext(f.Name()) == ".json") {
				fmt.Println("-", f.Name())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listExamplesCmd)
}
