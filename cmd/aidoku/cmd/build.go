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
		flags := cmd.Flags()

		output, _ := flags.GetString("output")

		web, _ := flags.GetBool("web")
		webTitle, _ := flags.GetString("web-title")
		webDescription, _ := flags.GetString("web-description")
		webIcon, _ := flags.GetString("web-icon")

		webArgs := build.WebTemplateArguments{
			Title:       webTitle,
			Description: webDescription,
			Icon:        webIcon,
		}

		return build.BuildWrapper(args, output, web, webArgs)
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
	buildCmd.Flags().StringP("output", "o", "public", "Output folder")

	buildCmd.Flags().BoolP("web", "w", false, "Generate a landing page for the source list")
	buildCmd.Flags().String("web-title", "An Aidoku source list", "Title of the landing page")
	buildCmd.Flags().String("web-description", "A source list for use with Aidoku.", "Description of the landing page")
	buildCmd.Flags().String("web-icon", "https://aidoku.app/images/favicon-32x32.png", "Icon of the landing page")

	buildCmd.MarkZshCompPositionalArgumentFile(1, "*.aix")
	buildCmd.ValidArgsFunction = func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"aix"}, cobra.ShellCompDirectiveFilterFileExt
	}
}
