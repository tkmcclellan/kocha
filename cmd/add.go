package cmd

import (
	"github.com/spf13/cobra"

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
		provider, _ := cmd.Flags().GetString("provider")
		dlmode, _ := cmd.Flags().GetString("dlmode")

		kocha.Init()
		kocha.Add(provider, dlmode)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().String("provider", "mangakakolot", "[mangakakolot]")
	addCmd.Flags().String("dlmode", "dynamic", "[dynamic, all, none]")
}
