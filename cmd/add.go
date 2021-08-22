package cmd

import (
	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
	util "github.com/tkmcclellan/kocha/internal"
	"github.com/tkmcclellan/kocha/internal/providers"
	"github.com/tkmcclellan/kocha/pkg/kocha"

	"bytes"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func getc(f *os.File) (byte, error) {
	b := make([]byte, 1)
	_, err := f.Read(b)
	return b[0], err
}

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
		kocha.Init()

		provider_type, _ := cmd.Flags().GetString("provider")
		// dlmode, _ := cmd.Flags().GetString("dlmode")
		name := strings.Join(args, " ")

		provider, err := providers.FindProvider(provider_type)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		search_results, err := provider.Search(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		writer := uilive.New()
		writer.Start()
		defer writer.Stop()

		ch := make(chan string)
		go util.MonitorStdin(ch)
		defer close(ch)
		defer util.Cleanup()

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go util.MonitorSigint(c)

		for {
			var buffer bytes.Buffer
			for i, m := range search_results.Manga {
				buffer.WriteString(fmt.Sprintf("id: %d\ttitle: %s\tauthors: %s\n%s\n", i, m.Title, strings.Join(m.Authors, ","), strings.Repeat("~", 100)))
			}
			fmt.Fprintf(writer, buffer.String())
			stdin, _ := <-ch
			if stdin == "q" {
				break
			} else {
				buffer.Reset()
				buffer.WriteString("here")
				fmt.Fprintf(writer, buffer.String())
			}
			time.Sleep(time.Millisecond * 100)
		}

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().String("provider", "mangakakalot", "[mangakakalot]")
	addCmd.Flags().String("dlmode", "dynamic", "[dynamic, all, none]")
}
