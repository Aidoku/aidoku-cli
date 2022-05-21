package cmd

import (
	"github.com/Aidoku/aidoku-cli/internal/build"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:           "build <FILES>",
	Short:         "Build a source list from packages",
	Version:       rootCmd.Version,
	Args:          cobra.MinimumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ForceColor {
			color.NoColor = false
		}
		output, _ := cmd.Flags().GetString("output")
		return build.BuildWrapper(args, output)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringP("output", "o", "public", "Output folder")

	buildCmd.MarkZshCompPositionalArgumentFile(1, "*.aix")
	buildCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"aix"}, cobra.ShellCompDirectiveFilterFileExt
	}
}
