package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/tkmcclellan/kocha/internal/models"
	"github.com/tkmcclellan/kocha/pkg/kocha"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove manga",
	Long: `Remove manga from kocha.

This command will remove a manga from kocha and
delete any of the manga's chapters from the user's
file system.
`,
	Run: func(cmd *cobra.Command, args []string) {
		Remove()
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func Remove() {
	manga := kocha.List()

	var options []string
	items := make(map[string]models.Manga)
	for _, m := range manga {
		items[m.ToReadable()] = m
		options = append(options, m.ToReadable())
	}

	questions := []*survey.Question{
		{
			Name: "manga",
			Prompt: &survey.Select{
				Message: "Which manga do you want to delete?",
				Options: options,
			},
		},
		{
			Name:   "confirmation",
			Prompt: &survey.Confirm{Message: "Are you sure you want to delete this manga?"},
		},
	}

	answers := struct {
		Manga        string
		Confirmation bool
	}{}
	err := survey.Ask(questions, &answers)
	if err != nil {
		panic(err)
	}

	if answers.Confirmation {
		manga := items[answers.Manga]
		kocha.Remove(&manga)
	} else {
		return
	}
}
