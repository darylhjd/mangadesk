package pages

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

// downloadChapters : Attempt to download pages
func downloadChapters(pages *tview.Pages, mangaPage *MangaPage, mr *mangodex.MangaResponse, chaps *[]mangodex.ChapterResponse) {
	// Download each chapter.
	// NOTE : Run as a GOROUTINE. Require QueueUpdateDraw
	go func(rows map[int]struct{}) {
		// For each selected chapter to download.
		var errorChaps []string
	ChapterForLoop:
		for r := range rows {
			// Get the corresponding ChapterResponse object.
			chapR := (*chaps)[r-1] // We -1 since the first row is the header.

			// Get chapterName folder name.
			chapterName := generateChapterFolderName(&chapR)

			// Get MangaDex@Home downloader for the chapterName.
			downloader, err := g.Dex.NewMDHomeClient(chapR.Data.ID, "data", chapR.Data.Attributes.Hash, false)
			if err != nil {
				// If error getting downloader, we add this chapterName to the errorPages chapters list.
				errorChaps = append(errorChaps, chapterName)
				continue // Continue on to the next chapterName to download.
			}

			// Create directory to store the downloaded chapters.
			// It is stored in DOWNLOAD_FOLDER/MANGA_NAME/CHAPTER_FOLDER
			chapterFolder := filepath.Join(g.Conf.DownloadDir, mr.Data.Attributes.Title["en"], chapterName)
			if err = os.MkdirAll(chapterFolder, os.ModePerm); err != nil {
				// If error creating folder to store this chapterName, we add this chapterName to errorPages chapters list.
				errorChaps = append(errorChaps, chapterName)
				continue // Continue on to the next chapterName to download.
			}

			// Get each page and save it.
			for pageNum, file := range chapR.Data.Attributes.Data {
				// Get page data.
				image, err := downloader.GetChapterPage(file)
				if err != nil {
					// If error downloading page data.
					continue ChapterForLoop // Continue on to the next chapter.
				}

				// Attempt to write page data into file.
				filename := fmt.Sprintf("%04d%s", pageNum+1, filepath.Ext(file))
				err = ioutil.WriteFile(filepath.Join(chapterFolder, filename), image, os.ModePerm)
				if err != nil {
					// If error storing page data.
					continue ChapterForLoop // Continue on to the next page.
				}
			}

			// Update downloaded column.
			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				dCell := tview.NewTableCell("Y").SetTextColor(g.MangaPageDownloadStatColor)
				mangaPage.ChapterTable.SetCell(r, 2, dCell)
			})
		}

		// Finished downloading all chapters. Now inform user.
		var builder strings.Builder
		builder.WriteString("Last Download Queue finished.\n")
		if len(errorChaps) == 0 {
			builder.WriteString("No errors :)")
		} else {
			builder.WriteString("Errors:\n")
			for _, v := range errorChaps {
				builder.WriteString(v + "\n")
			}
		}
		g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
			OKModal(pages, g.DownloadFinishedModalID, builder.String())
		})
	}(*mangaPage.Selected) // We pass the whole map as a value as we need to clear it later.

	// Clear the stored rows and unmark all chapters
	for k := range *mangaPage.Selected {
		markChapterUnselected(mangaPage.ChapterTable, k)
	}
	mangaPage.Selected = &map[int]struct{}{} // Empty the map
}

// generateChapterFolderName : Generate a folder name for the chapter.
func generateChapterFolderName(cr *mangodex.ChapterResponse) string {
	chapterNum := "unknown"
	if cr.Data.Attributes.Chapter != nil {
		chapterNum = *(cr.Data.Attributes.Chapter)
	}

	// Remove all invalid folder characters from folder name
	chapterName := cr.Data.Attributes.Title
	restrictedChars := []string{"<", ">", ":", "/", "|", "?", "*", "\"", "\\"}
	for s := range restrictedChars {
		chapterName = strings.ReplaceAll(chapterName, restrictedChars[s], "")
	}

	// Use compound name to try to avoid collisions.
	return fmt.Sprintf("%s - %s", chapterNum, chapterName)
}
