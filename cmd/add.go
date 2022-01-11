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
	Short: "Add manga",
	Long: `Add manga to kocha.

Adding manga to kocha acts the same as adding an e-book or
audiobook to your virtual library. Once you add a manga to kocha, you can
perform actions on it like editing saved information, removing it from your
library, downloading certain chapters of your manga, etc.
`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		dlmode, _ := cmd.Flags().GetString("dlmode")
		provider, _ := cmd.Flags().GetString("provider")
		name := strings.Join(args, " ")

		Add(url, dlmode, provider, name)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().String("provider", "", "[mangakakalot]")
	addCmd.Flags().StringP("dlmode", "d", "", "[dynamic, all, none]")
	addCmd.Flags().StringP("url", "u", "", "add by url")
}

func addMode() (selection string) {
	addMode := &survey.Select{
		Message: "How would you like to add a manga?",
		Options: []string{"Search", "Url"},
	}

	err := survey.AskOne(addMode, &selection, survey.WithValidator(survey.Required))
	if err != nil {
		panic(err)
	}
	return
}

// This exact function was used in both the AddFromQuery
// and AddFromUrl functions so it was broken out for reuse.
func parseDlmode(dlmodeFlag string) (dlmode string) {
	if len(dlmodeFlag) == 0 {
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
	} else {
		dlmode = dlmodeFlag
	}
	return
}

// Parent Add function. This is necessary to support the direct
// invocation of the add function as well as supporting the
// pathway through the interactive interface.
func Add(url string, dlmode string, provider string, name string) {
	if len(url) > 0 {
		AddFromUrl(dlmode, url)
		return
	}

	if len(name) > 0 {
		AddFromQuery(provider, dlmode, name)
		return
	}

	var mode = addMode()

	switch {
	case strings.Contains(mode, "Url"):
		AddFromUrl(dlmode, url)
	case strings.Contains(mode, "Search"):
		AddFromQuery(provider, dlmode, name)
	}
}

// The AddFromUrl function does not require a name or provider
// parameter so it was broken out to be more explicit.
func AddFromUrl(dlmode string, url string) {
	if len(url) == 0 {
		question := &survey.Input{
			Message: "Enter the url",
		}

		err := survey.AskOne(question, &url, survey.WithValidator(survey.Required))
		if err != nil {
			panic(err)
		}
	}

	dlmode = parseDlmode(dlmode)

	kocha.AddFromUrl(url, dlmode)
}

// Ask the user for a provider, download mode, and name for a manga
// and then add it to the database.
func AddFromQuery(providerType string, dlmode string, name string) {
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

	dlmode = parseDlmode(dlmode)

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
