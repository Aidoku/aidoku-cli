package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Verbose    bool
	ForceColor bool

	version = "develop"
	commit  string
	date    string
	builtBy string
)

var rootCmd = &cobra.Command{
	Use:   "aidoku",
	Short: "Aidoku development toolkit",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVar(&ForceColor, "force-color", false, "always output with color")
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	formattedVersion := FormatVersion(version, commit, date, builtBy)
	rootCmd.SetVersionTemplate(formattedVersion)
	rootCmd.Version = formattedVersion

	rootCmd.AddCommand(NewVersionCmd(version, commit, date, builtBy))
}
