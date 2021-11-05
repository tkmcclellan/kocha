package cmd

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/tkmcclellan/kocha/internal/models"
	"github.com/tkmcclellan/kocha/pkg/kocha"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		providerType, _ := cmd.Flags().GetString("provider")
		dlmode, _ := cmd.Flags().GetString("dlmode")
		name := strings.Join(args, " ")

		Add(providerType, dlmode, name)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().String("provider", "", "[mangakakalot]")
	addCmd.Flags().String("dlmode", "", "[dynamic, all, none]")
}

func Add(providerType string, dlmode string, name string) {
	if len(name) == 0 {
		question := &survey.Input{
			Message: "What manga do you want to add?",
		}

		err := survey.AskOne(question, &name, survey.WithValidator(survey.Required))
		if err != nil {
			panic(err)
		}
	}

	if len(providerType) == 0 {
		question := &survey.Select{
			Message: "Provider:",
			Options: []string{"Mangakakalot"},
			Default: "Mangakakalot",
		}

		var selection string
		err := survey.AskOne(question, &selection, survey.WithValidator(survey.Required))
		if err != nil {
			panic(err)
		}

		providerType = strings.ToLower(selection)
	}

	if len(dlmode) == 0 {
		question := &survey.Select{
			Message: "Download Mode:",
			Options: []string{
				"Dynamic: download as you read",
				"All: download all chapters now",
				"None: don't download any chapters"},
			Default: "Dynamic: download as you read",
		}

		var selection string
		err := survey.AskOne(question, &selection, survey.WithValidator(survey.Required))
		if err != nil {
			panic(err)
		}

		if strings.Contains(selection, "All") {
			dlmode = "all"
		} else if strings.Contains(selection, "Dynamic") {
			dlmode = "dynamic"
		} else if strings.Contains(selection, "None") {
			dlmode = "none"
		} else {
			panic("invalid download mode")
		}
	}

	var page uint64 = 1
	searchResults := kocha.Search(name, providerType, page)
	for {
		var options []string
		items := make(map[string]models.Manga)
		if page > 1 {
			options = append(options, "Previous")
		}
		for _, manga := range searchResults.Manga {
			items[manga.ToReadable()] = manga
			options = append(options, manga.ToReadable())
		}
		if page < searchResults.TotalPages {
			options = append(options, "Next")
		}

		question := &survey.Select{
			Message: "Choose a manga:",
			Options: options,
		}

		var selection string
		err := survey.AskOne(question, &selection)
		if err != nil {
			panic(err)
		}

		if selection == "Next" {
			page += 1
			continue
		}
		if selection == "Previous" {
			page -= 1
			continue
		}
		manga := items[selection]
		kocha.Add(&manga, dlmode)
		return
	}
}
