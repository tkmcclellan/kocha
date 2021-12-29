package cmd

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/tkmcclellan/kocha/pkg/kocha"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List manga",
	Long: `List tracked manga.

Lists manga tracked by kocha in a table format.
`,
	Run: func(cmd *cobra.Command, args []string) {
		List()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func List() {
	manga := kocha.List()
	data := [][]string{}
	for _, m := range manga {
		data = append(data, m.ToRow())
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader([]string{"ID", "Title", "Authors", "Uri", "Updated At", "Download Mode"})

	table.AppendBulk(data)
	table.Render()
}
