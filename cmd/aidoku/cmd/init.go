package cmd

import (
	"bufio"
	"net/url"
	"os"
	"strings"

	"golang.org/x/exp/slices"
	"golang.org/x/text/language"

	"github.com/Aidoku/aidoku-cli/internal/templates"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func checkPrompt(err error) error {
	if err != nil {
		if err == terminal.InterruptErr {
			color.Red("interrupted")
			os.Exit(1)
		}
	}
	return nil
}

func languageValidator(lang string) bool {
	return slices.Contains([]string{"rust-template", "rust", "as", "c"}, lang)
}

func sourceLanguageValidator(response interface{}) error {
	if response.(string) != "multi" {
		_, err := language.Parse(response.(string))
		return err
	} else {
		return nil
	}
}

var initCommand = &cobra.Command{
	Use:           "init [rust-template|rust|as|c] [DIR]",
	Short:         "Create initial code for an Aidoku source",
	Version:       rootCmd.Version,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		lang, _ := cmd.Flags().GetString("language")
		name, _ := cmd.Flags().GetString("name")
		homepage, _ := cmd.Flags().GetString("homepage")
		nsfw, _ := cmd.Flags().GetInt("nsfw")
		if len(args) < 1 || !languageValidator(args[0]) {
			prompt := &survey.Select{
				Message: "Template",
				Options: []string{"rust-template", "rust", "as", "c"},
			}
			var selection string
			checkPrompt(survey.AskOne(prompt, &selection, survey.WithValidator(survey.Required)))
			args = append(args, selection)
		}
		if len(args) < 2 {
			prompt := &survey.Input{
				Message: "Directory to generate template",
				Default: ".",
			}
			var selection string
			checkPrompt(survey.AskOne(prompt, &selection))
			args = append(args, selection)
		}
		if name == "" {
			prompt := &survey.Input{
				Message: "Source name",
			}
			var result string
			checkPrompt(survey.AskOne(prompt, &result, survey.WithValidator(survey.Required)))
			name = result
		}
		if lang == "" {
			prompt := &survey.Input{
				Message: "Source language",
			}
			var result string
			checkPrompt(survey.AskOne(prompt, &result, survey.WithValidator(sourceLanguageValidator)))
			lang = result
		}
		if homepage == "" {
			prompt := &survey.Input{
				Message: "Source homepage",
			}
			var result string
			checkPrompt(survey.AskOne(prompt, &result, survey.WithValidator(func(response interface{}) error {
				_, err := url.ParseRequestURI(response.(string))
				return err
			})))
			homepage = result
		}
		if nsfw == -1 {
			prompt := &survey.Select{
				Message: "Source NSFW level",
				Options: []string{"None", "Moderate amounts of sex", "Insane Amounts of Sex"},
			}
			var selection int
			checkPrompt(survey.AskOne(prompt, &selection, survey.WithValidator(survey.Required)))
			nsfw = selection
		}

		source := templates.Source{
			Name:     name,
			Homepage: homepage,
			Language: lang,
			Nsfw:     nsfw,
		}

		var err error
		switch args[0] {
		case "rust":
			{
				// Determine if it's a child of a template by checking for a parent package at $output
				if file, err := os.Open(args[1] + "/../../template/Cargo.toml"); err == nil {
					defer file.Close()
					scanner := bufio.NewScanner(file)
					for scanner.Scan() {
						line := scanner.Text()
						if strings.Contains(line, "name = \"") {
							source.TemplateName = strings.TrimSuffix(strings.TrimPrefix(line, "name = \""), "\"")
							break
						}
					}
				}
				err = templates.RustGenerator(args[1], source)
			}
		case "as":
			err = templates.AscGenerator(args[1], source)
		case "c":
			err = templates.ClangGenerator(args[1], source)
		case "rust-template":
			err = templates.RustTemplateGenerator(args[1], source)
		}
		if err != nil {
			color.Red("error: could not generate initial code")
		}
		return err
	},
}

func init() {
	rootCmd.AddCommand(initCommand)
	initCommand.Flags().StringP("name", "n", "", "Source name")
	initCommand.Flags().StringP("language", "l", "", "Source language")
	initCommand.Flags().StringP("homepage", "p", "", "Source homepage")
	initCommand.Flags().Int("nsfw", -1, "Source NSFW level")
}
