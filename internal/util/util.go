package util

import (
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/signintech/gopdf"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var BasePath string
var MangaPath string

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	BasePath = path.Join(homedir, ".kocha")
	MangaPath = path.Join(BasePath, "manga")

	_, err = os.Stat(BasePath)
	if os.IsNotExist(err) {
		os.Mkdir(BasePath, 0775)
	}

	dbpath := path.Join(BasePath, "manga.db")
	db, err = gorm.Open(sqlite.Open(dbpath), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic(err)
	}
}

func Database() (*gorm.DB, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	return db.WithContext(ctx), cancel
}

func DownloadImage(dirpath string, title string, src string, completed chan bool) {
	_, filename := path.Split(src)
	filenameSplit := strings.Split(filename, ".")
	extension := filenameSplit[len(filenameSplit)-1]
	imagePath := path.Join(dirpath, title+"."+extension)

	res, err := http.Get(src)
	if err != nil {
		fmt.Println(err)
		if completed != nil {
			completed <- false
		}
		return
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		if completed != nil {
			completed <- false
		}
		return
	}

	Write(imagePath, bytes)

	if completed != nil {
		completed <- true
	}
}

func ChapterToPdf(chapterPath string) error {
	relativePath := path.Join(MangaPath, chapterPath)
	files, err := os.ReadDir(relativePath)
	if err != nil {
		return err
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	for _, imgFile := range files {
		if filepath.Ext(imgFile.Name()) == ".jpg" || filepath.Ext(imgFile.Name()) == ".jpeg" || filepath.Ext(imgFile.Name()) == ".png" {
			imgPath := filepath.Join(relativePath, imgFile.Name())
			if reader, err := os.Open(imgPath); err == nil {
				defer reader.Close()
				im, _, err := image.DecodeConfig(reader)
				if err != nil {
					return err
				}
				// fmt.Printf("%s %d %d\n", imgFile.Name(), im.Width, im.Height)
				imgRect := gopdf.Rect{W: float64(im.Width), H: float64(im.Height)}
				pdf.AddPageWithOption(gopdf.PageOption{PageSize: &imgRect})
				pdf.Image(imgPath, 0, 0, &imgRect)
			} else {
				return err
			}
		}
	}

	_, chapterName := path.Split(chapterPath)
	pdfPath := filepath.Join(relativePath, chapterName+".pdf")

	pdf.WritePdf(pdfPath)

	return nil
}

func Write(relativePath string, data []byte) error {
	relativePath = path.Join(MangaPath, relativePath)
	parentDir, _ := path.Split(relativePath)
	err := os.MkdirAll(parentDir, 0775)
	if err != nil {
		return err
	}
	return os.WriteFile(relativePath, data, 0775)
}

func DeleteDir(dirpath string) error {
	dirPath := path.Join(MangaPath, dirpath)
	return os.RemoveAll(dirPath)
}

func CleanString(text string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(text, "")
	return strcase.ToSnake(processedString)
}

func IntRangeContains(intRange []int, element int) bool {
	for _, v := range intRange {
		if v == element {
			return true
		}
	}

	return false
}
