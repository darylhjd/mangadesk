package pages

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
)

// confirmChapterDownloads : Function for handling when the user press enter on the table.
func confirmChapterDownloads(pages *tview.Pages, table *tview.Table,
	selected *map[int]struct{}, row int, mr *mangodex.MangaResponse, chaps *[]mangodex.ChapterResponse) {
	// We add the current selection if the there are no selected rows currently.
	if len(*selected) == 0 {
		(*selected)[row] = struct{}{}
	}

	// Show modal to confirm download.
	ShowModal(pages, g.DownloadChaptersModalID, "Download selection(s)?", []string{"Yes", "No"},
		func(i int, label string) {
			if label == "Yes" {
				// If user confirms to download, then we download the chapters.
				downloadChapters(pages, table, selected, mr, chaps)
			}
			pages.RemovePage(g.DownloadChaptersModalID)
		})
}

// downloadChapters : Attempt to download pages
func downloadChapters(pages *tview.Pages, table *tview.Table, selected *map[int]struct{}, mr *mangodex.MangaResponse, chaps *[]mangodex.ChapterResponse) {
	// Download each chapter.
	// NOTE : Run as a GOROUTINE. Require QueueUpdateDraw
	go func(rows map[int]struct{}) {
		// For each selected chapter to download.
		errorChaps := map[string][]int{}
		for r := range rows {
			// Get the corresponding ChapterResponse object.
			chapR := (*chaps)[r-1] // We -1 since the first row is the header.

			// Get chapter folder name.
			chapter := generateChapterFolderName(&chapR)

			// Get MangaDex@Home downloader for the chapter.
			downloader, err := g.Dex.NewMDHomeClient(chapR.Data.ID, "data", chapR.Data.Attributes.Hash, false)
			if err != nil {
				// If error getting downloader, we add this chapter to the errorPages chapters list.
				errorChaps[chapter] = []int{}
				continue // Continue on to the next chapter to download.
			}

			// Create directory to store the downloaded chapters.
			// It is stored in DOWNLOAD_FOLDER/MANGA_NAME/CHAPTER_FOLDER
			chapterFolder := filepath.Join(g.Conf.DownloadDir, mr.Data.Attributes.Title["en"], chapter)
			if err = os.MkdirAll(chapterFolder, os.ModePerm); err != nil {
				// If error creating folder to store this chapter, we add this chapter to errorPages chapters list.
				errorChaps[chapter] = []int{}
				continue // Continue on to the next chapter to download.
			}

			// Get each page and save it.
			var errorPages []int
			for pageNum, file := range chapR.Data.Attributes.Data {
				// Get page data.
				image, err := downloader.GetChapterPage(file)
				if err != nil {
					// If error downloading page data.
					errorPages = append(errorPages, pageNum+1)
					continue // Continue on to the next page.
				}

				// Attempt to write page data into file.
				filename := fmt.Sprintf("%03d%s", pageNum+1, filepath.Ext(file))
				err = ioutil.WriteFile(filepath.Join(chapterFolder, filename), image, os.ModePerm)
				if err != nil {
					// If error storing page data.
					errorPages = append(errorPages, pageNum+1)
					continue // Continue on to the next page.
				}
			}
			// Tell user that the chapter has been downloaded
			g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				dCell := tview.NewTableCell("Yup!").SetTextColor(tcell.ColorPowderBlue)
				table.SetCell(r, 2, dCell)
			})
		}

		// Finished downloading all chapters. Now inform user.
		var builder strings.Builder
		builder.WriteString("Last Download Queue finished.\n")
		if len(errorChaps) == 0 {
			builder.WriteString("No errors :)")
		} else {
			builder.WriteString("Errors:\n")
			for k, v := range errorChaps {
				builder.WriteString(k + " - ")
				for _, p := range v {
					builder.WriteString(strconv.Itoa(p))
				}
			}
		}
		g.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
			ShowModal(pages, g.DownloadFinishedModalID, builder.String(), []string{"OK"}, func(i int, label string) {
				pages.RemovePage(g.DownloadFinishedModalID)
			})
		})

	}(*selected) // We pass the whole map as a value as we need to clear it.

	// Clear the stored rows and unmark all chapters
	for k := range *selected {
		markChapterUnselected(table, k)
	}
	*selected = map[int]struct{}{} // Empty the map
}

// generateChapterFolderName : Generate a folder name for the chapter.
func generateChapterFolderName(cr *mangodex.ChapterResponse) string {
	chapterNum := "?"
	if cr.Data.Attributes.Chapter != nil {
		chapterNum = *(cr.Data.Attributes.Chapter)
	}
	// Remove all invalid folder characters from folder name
	chapterName := cr.Data.Attributes.Title
	reschar := []string{"<", ">", ":", "/", "|", "?", "*", "\"", "\\"}
	for s := range reschar {
		chapterName = strings.ReplaceAll(chapterName, reschar[s], "")
	}
	// Use compound name to try to avoid collisions.
	return fmt.Sprintf("%s - %s", chapterNum, chapterName)
}
