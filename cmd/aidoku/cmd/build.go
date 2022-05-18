package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/beerpiss/aidoku-cli/internal/build"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:           "build <FILES>",
	Short:         "Build a source list from packages",
	Args:          cobra.MinimumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		os.RemoveAll(output)
		var fileList []string
		for _, arg := range args {
			files, err := filepath.Glob(arg)
			if err != nil {
				color.Red("error: invalid glob pattern %s", arg)
				continue
			}
			fileList = append(fileList, files...)
		}
		if len(fileList) == 0 {
			return errors.New("no files given")
		}
		os.MkdirAll(output, os.FileMode(0644))
		os.MkdirAll(output+"/icons", os.FileMode(0644))
		os.MkdirAll(output+"/sources", os.FileMode(0644))
		build.BuildSource(fileList, output)
		return nil
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
