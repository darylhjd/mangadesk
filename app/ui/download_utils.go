package ui

import (
	"archive/zip"
	"fmt"
	"github.com/darylhjd/mangadesk/app/core"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"
)

func (p *MangaPage) downloadChapters() {
	// Get a reference to the current selection
	selection := p.Selected
	// Reset.
	p.Selected = map[int]struct{}{}

	// Clear selection highlight
	for index := range selection {
		markChapterUnselected(p.Table, index)
	}

	// Download the selected chapters.
	for index := range selection {
		// Get the reference to the chapter.
		chapter := p.Table.GetCell(index, 0).GetReference().(*mangodex.Chapter)

		downloadClient, err := core.App.Client.AtHome.NewMDHomeClient(
			chapter, core.App.Config.DownloadQuality, core.App.Config.ForcePort443)
	}
}

// downloadChapters : Attempt to download Holder
func downloadChapters(pages *tview.Pages, mangaPage *MangaPage, m *mangodex.Manga, chaps *[]mangodex.Chapter) {
	// Download each chapter.
	// NOTE : Initialise as a GOROUTINE. Require QueueUpdateDraw
	go func(rows map[int]struct{}) {
		// For each selected chapter to download.
		var errorChaps []string
	ChapterForLoop:
		for r := range rows {
			// Get the corresponding Chapter object.
			chap := (*chaps)[r-1] // We -1 since the first row is the header.

			// Get manga and chapter folder name.
			mangaName, chapterName := generateChapterFolderNames(m, &chap)

			// Get MangaDex@Home downloader for the chapterName.
			downloader, err := core.DexClient.NewMDHomeClient(chap.ID, core.Conf.DownloadQuality,
				chap.Attributes.Hash, core.Conf.ForcePort443)
			if err != nil {
				// If error getting downloader, we add this chapterName to the errorPages chapters list.
				errorChaps = append(errorChaps,
					fmt.Sprintf("%s -> %s", chapterName, err.Error()))
				continue // Continue on to the next chapterName to download.
			}

			// Create directory to store the downloaded chapters.
			// It is stored in DOWNLOAD_FOLDER/MANGA_NAME/CHAPTER_FOLDER
			chapterFolder := filepath.Join(core.Conf.DownloadDir, mangaName, chapterName)
			if err = os.MkdirAll(chapterFolder, os.ModePerm); err != nil {
				// If error creating folder to store this chapterName, we add this chapterName to errorPages chapters list.
				errorChaps = append(errorChaps,
					fmt.Sprintf("%s ->%s", chapterName, err.Error()))
				continue // Continue on to the next chapterName to download.
			}

			// Get each page and save it.
			// Note that the moment one page fails to download, the whole chapter is skipped.
			pageFiles := chap.Attributes.Data
			if core.Conf.DownloadQuality == "data-saver" {
				pageFiles = chap.Attributes.DataSaver
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
			if core.Conf.AsZip {
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
			core.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
				dCell := tview.NewTableCell("Y").SetTextColor(MangaPageDownloadStatColor)
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
		core.App.QueueUpdateDraw(func() { // GOROUTINE : Require QueueUpdateDraw
			okModal(pages, DownloadFinishedModalID, builder.String())
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
	zipFile, err := os.Create(fmt.Sprintf("%s.%s", chapterFolder, core.Conf.ZipType))
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

// getDownloadFolder : Get the download folder for a manga's chapter.
func (p *MangaPage) getDownloadFolder(chapter *mangodex.Chapter) string {
	mangaName := p.Manga.GetTitle("en")
	chapterName := fmt.Sprintf("Chapter%s [%s-%s] %s_%s",
		chapter.GetChapterNum(), strings.ToUpper(chapter.Attributes.TranslatedLanguage), core.App.Config.DownloadQuality,
		chapter.GetTitle(), strings.SplitN(chapter.ID, "-", 2)[0])

	// Remove invalid characters from the folder name
	restricted := []string{"<", ">", ":", "/", "|", "?", "*", "\"", "\\", "."}
	for _, c := range restricted {
		mangaName = strings.ReplaceAll(mangaName, c, "")
		chapterName = strings.ReplaceAll(chapterName, c, "")
	}

	folder := filepath.Join(core.App.Config.DownloadDir, mangaName, chapterName)
	// If the user wants to download as a zip, then we check for the presence of the zip folder.
	if core.App.Config.AsZip {
		folder = fmt.Sprintf("%s.%s", folder, core.App.Config.ZipType)
	}
	return folder
}
