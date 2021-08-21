package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tkmcclellan/kocha/pkg/kocha"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kocha",
	Short: "~Welcome to Manga Manager~",
	Long:  `Kocha - A tool for reading manga.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(kocha.Init())
	cobra.CheckErr(rootCmd.Execute())
}
