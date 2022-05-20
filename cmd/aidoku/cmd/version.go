package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func NewVersionCmd(version, commit, date, builtBy string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(FormatVersion(version, commit, date, builtBy))
		},
	}
}

func FormatVersion(version, commit, date, builtBy string) string {
	var output strings.Builder
	output.WriteString("aidoku-cli version ")
	output.WriteString(strings.TrimPrefix(version, "v"))
	if commit != "" {
		output.WriteString(", commit ")
		output.WriteString(commit)
	}
	if date != "" {
		output.WriteString(", built at ")
		output.WriteString(date)
	}
	if builtBy != "" {
		output.WriteString(", built by ")
		output.WriteString(builtBy)
	}
	return output.String()
}
