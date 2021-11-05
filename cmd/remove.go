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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Remove()
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
