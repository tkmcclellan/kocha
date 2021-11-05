/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"github.com/tkmcclellan/kocha/internal/models"
	"github.com/tkmcclellan/kocha/internal/providers"
	"github.com/tkmcclellan/kocha/internal/util"
	"github.com/tkmcclellan/kocha/pkg/kocha"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		Read()
	},
}

func init() {
	rootCmd.AddCommand(readCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
	completed := make(chan bool)
	var downloadCount int
	for i := low; i <= high; i++ {
		_, err := os.Stat(filepath.Join(util.MangaPath, manga.Dirname(), chapters[i].Dirname()))
		if os.IsNotExist(err) {
			fmt.Println("downloading chapter ", i)
			go provider.DownloadChapter(chapters[i], completed)
			downloadCount++
		}
	}

	if downloadCount > 0 {
		for i := 0; i < downloadCount; i++ {
			<-completed
		}
	}
}

func deleteOutsideRange(manga *models.Manga, chapters []models.Chapter, low int, high int) {
	for i := 0; i < len(chapters); i++ {
		chapterPath := filepath.Join(util.MangaPath, manga.Dirname(), chapters[i].Dirname())
		_, err := os.Stat(chapterPath)
		if err == nil && (i < low || i > high) {
			fmt.Println("deleting chapter ", i)
			err := util.DeleteDir(filepath.Join(manga.Dirname(), chapters[i].Dirname()))
			if err != nil {
				panic(err)
			}
		}
	}
}

func readManga(manga *models.Manga) {
	chapters := manga.Chapters(nil)
	currentChapter := manga.CurrentChapter
	provider, err := providers.FindProvider(manga.Provider)
	if err != nil {
		panic(err)
	}

	high, low := chapterRange(chapters, currentChapter)
	downloadWithinRange(provider, manga, chapters, low, high)
	deleteOutsideRange(manga, chapters, low, high)

	for {
		chapter := chapters[currentChapter]
		filePath := filepath.Join(util.MangaPath, manga.Dirname(), chapter.Dirname(), util.CleanString(chapter.Title)+".pdf")
		err = open.Run(filePath)
		if err != nil {
			panic(err)
		}

		high, low := chapterRange(chapters, currentChapter)
		go downloadWithinRange(provider, manga, chapters, low, high)
		go deleteOutsideRange(manga, chapters, low, high)

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
		err := survey.AskOne(question, &selection, survey.WithValidator(survey.Required))
		if err != nil {
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
