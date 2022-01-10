package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/tkmcclellan/kocha/internal/models"
	"github.com/tkmcclellan/kocha/internal/providers"
	"github.com/tkmcclellan/kocha/internal/util"
	"github.com/tkmcclellan/kocha/pkg/kocha"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read manga",
	Long: `Read manga.

This command will attempt to open a PDF of the images in the
selected chapter in the default program for opening PDFs on
the user's system.

If the download mode is anything other than 'All', then chapters
will be deleted and downloaded in the background as the user reads.
Kocha will give the user a reading menu that allows them to open the
next chapter, the previous chapter, [WIP] a specific chapter, or quit reading.

Kocha will keep track of which chapter the user is currently reading
and open that chapter again the next time this command is called.
`,
	Run: func(cmd *cobra.Command, args []string) {
		Read()
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}

func Read() {
	manga := kocha.List()

	var options []string
	items := make(map[string]models.Manga)
	for _, m := range manga {
		items[m.ToReadable()] = m
		options = append(options, m.ToReadable())
	}

	question := &survey.Select{
		Message: "Which manga do you want to read?",
		Options: options,
	}

	var selection string
	err := survey.AskOne(question, &selection, survey.WithValidator(survey.Required))
	if err != nil {
		panic(err)
	}

	selectedManga := items[selection]
	readManga(&selectedManga)
}

const dlrange int = 2

func chapterRange(chapters []models.Chapter, currentChapter int) (int, int) {
	high := currentChapter
	low := currentChapter
	for i := 1; i <= dlrange; i++ {
		if currentChapter+i < len(chapters) {
			high = currentChapter + i
		}
		if currentChapter-i >= 0 {
			low = currentChapter - i
		}
	}

	return high, low
}

func downloadWithinRange(provider providers.Provider, manga *models.Manga, chapters []models.Chapter, low int, high int) {
	for i := low; i <= high; i++ {
		_, err := os.Stat(filepath.Join(util.MangaPath, manga.Dirname(), chapters[i].Dirname()))
		if os.IsNotExist(err) {
			fmt.Println("Downloading chapter", i+1)
			go provider.DownloadChapter(chapters[i])
			// writer := uilive.New()
			// writer.Start()
			// fmt.Fprintln(writer, "Downloading chapter", i+1)
			// go func(writer *uilive.Writer) {
			// provider.DownloadChapter(chapters[i])
			// fmt.Fprintf(writer, "\033[1AChapter %d downloaded        \n\033[1m\033[92m? \033[37mWhat would you like to do?  \033[0m\033[36m[Use arrows to move, type to filter]\033[0m", i)
			// writer.Stop()
			// }(writer)
		}
	}
}

func deleteOutsideRange(manga *models.Manga, chapters []models.Chapter, low int, high int) {
	for i := 0; i < len(chapters); i++ {
		chapterPath := filepath.Join(util.MangaPath, manga.Dirname(), chapters[i].Dirname())
		_, err := os.Stat(chapterPath)
		if err == nil && (i < low || i > high) {
			fmt.Println("Deleting chapter", i)
			err := util.DeleteDir(filepath.Join(manga.Dirname(), chapters[i].Dirname()))
			if err != nil {
				panic(err)
			}
		}
	}
}

func openbrowser(url string) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("x-www-browser", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func readManga(manga *models.Manga) {
	chapters := manga.Chapters(nil)
	currentChapter := manga.CurrentChapter
	provider, err := providers.FindProvider(manga.Provider)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	for {
		chapter := chapters[currentChapter]
		filePath := filepath.Join(util.MangaPath, manga.Dirname(), chapter.Dirname(), util.CleanString(chapter.Title)+".pdf")
		err := openbrowser(filePath)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		high, low := chapterRange(chapters, currentChapter)
		downloadWithinRange(provider, manga, chapters, low, high)
		deleteOutsideRange(manga, chapters, low, high)

		options := []string{}
		if currentChapter+1 < len(chapters) {
			options = append(options, "Next")
		}
		if currentChapter-1 >= 0 {
			options = append(options, "Previous")
		}
		options = append(options, "Quit")

		question := &survey.Select{
			Message: "What would you like to do?",
			Options: options,
		}

		var selection string
		err = survey.AskOne(question, &selection, survey.WithValidator(survey.Required))
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		if selection == "Next" {
			currentChapter += 1
			chapter.Read = true
			chapter.Save()
			manga.CurrentChapter = currentChapter
			manga.Save()
		} else if selection == "Previous" {
			currentChapter -= 1
			manga.CurrentChapter = currentChapter
			manga.Save()
		} else if selection == "Quit" {
			break
		} else {
			panic("Invalid reading menu choice")
		}
	}
}
