package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Verbose bool
	Version = "develop"
)

var rootCmd = &cobra.Command{
	Use:     "aidoku",
	Short:   "Aidoku development toolkit",
	Version: Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}
