package pages

/*
Functions used by the application for downloading of chapters.
*/

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

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

			// Get manga and chapter folder name.
			mangaName, chapterName := generateChapterFolderNames(mr, &chapR)

			// Get MangaDex@Home downloader for the chapterName.
			downloader, err := g.Dex.NewMDHomeClient(chapR.Data.ID, g.Conf.DownloadQuality,
				chapR.Data.Attributes.Hash, g.Conf.ForcePort443)
			if err != nil {
				// If error getting downloader, we add this chapterName to the errorPages chapters list.
				errorChaps = append(errorChaps,
					fmt.Sprintf("%s -> %s", chapterName, err.Error()))
				continue // Continue on to the next chapterName to download.
			}

			// Create directory to store the downloaded chapters.
			// It is stored in DOWNLOAD_FOLDER/MANGA_NAME/CHAPTER_FOLDER
			chapterFolder := filepath.Join(g.Conf.DownloadDir, mangaName, chapterName)
			if err = os.MkdirAll(chapterFolder, os.ModePerm); err != nil {
				// If error creating folder to store this chapterName, we add this chapterName to errorPages chapters list.
				errorChaps = append(errorChaps,
					fmt.Sprintf("%s ->%s", chapterName, err.Error()))
				continue // Continue on to the next chapterName to download.
			}

			// Get each page and save it.
			// Note that the moment one page fails to download, the whole chapter is skipped.
			pageFiles := chapR.Data.Attributes.Data
			if g.Conf.DownloadQuality == "data-saver" {
				pageFiles = chapR.Data.Attributes.DataSaver
			}
			for pageNum, file := range pageFiles {
				// Get page data.
				image, err := downloader.GetChapterPage(file)
				if err != nil {
					// If error downloading page data.
					errorChaps = append(errorChaps,
						fmt.Sprintf("%s -> %s", chapterName, err.Error()))
					continue ChapterForLoop // Continue on to the next chapter.
				}

				// Attempt to write page data into file.
				filename := fmt.Sprintf("%04d%s", pageNum+1, filepath.Ext(file))
				err = ioutil.WriteFile(filepath.Join(chapterFolder, filename), image, os.ModePerm)
				if err != nil {
					// If error storing page data.
					errorChaps = append(errorChaps,
						fmt.Sprintf("%s -> %s", chapterName, err.Error()))
					continue ChapterForLoop // Continue on to the next page.
				}
			}

			// Check whether to also save as zip file type
			if g.Conf.AsZip {
				err := saveAsZipFolder(chapterFolder)
				if err != nil {
					// If error saving zip folder
					errorChaps = append(errorChaps,
						fmt.Sprintf("%s -> %s", chapterName, err.Error()))
					continue // Continue on to the next chapterName to download.
				}
				// Remove unzipped folder. Ignore any errors.
				_ = os.RemoveAll(chapterFolder)
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

// saveAsZipFolder : This function creates a zip folder to store a chapter download.
func saveAsZipFolder(chapterFolder string) error {
	zipFile, err := os.Create(fmt.Sprintf("%s.%s", chapterFolder, g.Conf.ZipType))
	if err != nil {
		return err
	}
	defer func() {
		_ = zipFile.Close()
	}()

	w := zip.NewWriter(zipFile)
	defer func() {
		_ = w.Close()
	}()

	return filepath.WalkDir(chapterFolder, func(path string, d fs.DirEntry, err error) error {
		// Stop walking immediately if encounter error
		if err != nil {
			return err
		}
		// Skip if a DirEntry is a folder. By right, this shouldn't happen since any downloads will
		// just contain PNGs or JPEGs but it's here just in case.
		if d.IsDir() {
			return nil
		}

		// Open the original image file.
		fileOriginal, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() {
			_ = fileOriginal.Close()
		}()

		// Create designated file in zip folder for current image.
		// Use custom header to set modified timing.
		// This fixes zip parsing issues in certain situations.
		fh := zip.FileHeader{
			Name:     d.Name(),
			Modified: time.Now(),
			Method:   zip.Deflate, // Consistent with w.Create() source code.
		}
		fileZip, err := w.CreateHeader(&fh)
		if err != nil {
			return err
		}

		// Copy the original file into its designated file in the zip archive.
		_, err = io.Copy(fileZip, fileOriginal)
		if err != nil {
			return err
		}
		return nil
	})
}

// generateChapterFolderNames : Generate folder names for the chapter and manga.
// Returns the name for the manga folder, then the name for the chapter folder.
func generateChapterFolderNames(mr *mangodex.MangaResponse, cr *mangodex.ChapterResponse) (string, string) {
	chapterNum := "-"
	if cr.Data.Attributes.Chapter != nil {
		chapterNum = *(cr.Data.Attributes.Chapter)
	}

	mangaName := mr.Data.Attributes.Title["en"]
	// Use compound name to try to avoid collisions.
	generatedName := fmt.Sprintf("Chapter%s_[%s-%s]_%s_%s",
		chapterNum, strings.ToUpper(cr.Data.Attributes.TranslatedLanguage), g.Conf.DownloadQuality,
		cr.Data.Attributes.Title, strings.SplitN(cr.Data.ID, "-", 2)[0])

	// Remove all invalid folder characters from folder name
	restrictedChars := []string{"<", ">", ":", "/", "|", "?", "*", "\"", "\\"}
	for s := range restrictedChars {
		mangaName = strings.ReplaceAll(mangaName, restrictedChars[s], "")
		generatedName = strings.ReplaceAll(generatedName, restrictedChars[s], "")
	}
	return mangaName, generatedName
}
