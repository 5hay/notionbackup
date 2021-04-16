package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kjk/notionapi"
)

const version = "1.0.0"

type App struct {
	client     *notionapi.Client
	pageID     string
	exportType string
	exportDir  string
}

func main() {
	authToken := os.Getenv("NOTION_TOKEN")
	pageID := os.Getenv("NOTION_PAGEID")
	if authToken == "" || pageID == "" {
		log.Fatalln("You have to set the env vars NOTION_TOKEN and NOTION_PAGEID.")
	}

	app := &App{
		client: &notionapi.Client{
			AuthToken: authToken,
		},
		pageID:     pageID,
		exportType: os.Getenv("NOTION_EXPORTTYPE"),
	}

	var err error
	app.exportDir = os.Getenv("NOTION_EXPORTDIR")
	if app.exportDir == "" {
		app.exportDir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("notionbackup (v%s) | Starting the export process ...\n", version)

	startTime := time.Now()
	exportURL, err := app.exportPageURL(true)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Notion export successful. Starting to download the exported .zip file now ...\n")

	bytesWritten := app.saveToFile(exportURL)

	log.Printf("Notion export (page id: %s) took %s, %d bytes written.\n", app.pageID, time.Since(startTime).String(), bytesWritten)
}

func (app *App) exportPageURL(recursive bool) (string, error) {
	if app.exportType == "" {
		app.exportType = "markdown"
	}

	// support full url, like https://www.notion.so/username/PageName-abcdefghia1f4505762g63874a1e97yz
	if strings.HasPrefix(app.pageID, "https://") {
		app.pageID = notionapi.ExtractNoDashIDFromNotionURL(app.pageID)
	}

	return app.client.RequestPageExportURL(app.pageID, app.exportType, recursive)
}

func (app *App) saveToFile(exportURL string) int64 {
	fileName := fmt.Sprintf("%s_%s.zip", time.Now().Format("20060102150405"), app.pageID)

	if err := os.MkdirAll(app.exportDir, 0755); err != nil {
		log.Fatal(err)
	}

	sep := string(os.PathSeparator)
	if strings.HasSuffix(app.exportDir, sep) {
		sep = ""
	}

	path := strings.Join([]string{app.exportDir, fileName}, sep)

	file, err := os.Create(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	resp, err := http.Get(exportURL)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	bytesWritten, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return bytesWritten
}
