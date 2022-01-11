package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/fatih/color"
	"github.com/tkmcclellan/kocha/pkg/kocha"
)

var cfgFile string

var ascii string = `
   __   __)           
  (, ) /       /)     
    /(   ____ (/   _  
 ) /  \_(_)(__/ )_(_(_
(_/                   
                      
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kocha",
	Short: "~Welcome to Kocha~",
	Long:  `Kocha ~ Black Tea.`,
	Run: func(cmd *cobra.Command, args []string) {
		color.Yellow(ascii)

		for {
			question := &survey.Select{
				Message: "What would you like to do?",
				Options: []string{"Add", "Read", "List", "Remove", "Quit"},
			}

			var selection string
			err := survey.AskOne(question, &selection, survey.WithValidator(survey.Required))
			if err != nil {
				panic(err)
			}

			if selection == "Quit" {
				break
			} else if selection == "Add" {
				Add("", "", "", "")
			} else if selection == "List" {
				List()
			} else if selection == "Remove" {
				Remove()
			} else if selection == "Read" {
				Read()
			} else {
				panic("Invalid command")
			}
		}
		color.Yellow("Thank you for using Kocha!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(kocha.Init())
	cobra.CheckErr(rootCmd.Execute())
}
