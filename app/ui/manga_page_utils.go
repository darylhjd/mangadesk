package ui

import (
	"archive/zip"
	"fmt"
	"github.com/darylhjd/mangadesk/app/core"
	"github.com/darylhjd/mangadesk/app/ui/utils"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"
)

const (
	maxRetries = 5
)

// downloadChapters : Download current chapters specified by the user.
func (p *MangaPage) downloadChapters(selection map[int]struct{}, attemptNo int) {
	// Download the selected chapters.
	errored := map[int]struct{}{}
	for index := range selection {
		// Get the reference to the chapter.
		var (
			chapter *mangodex.Chapter
			ok      bool
		)
		if chapter, ok = p.Table.GetCell(index, 0).GetReference().(*mangodex.Chapter); !ok {
			return
		}

		// Save the current chapter.
		err := p.saveChapter(chapter)
		if err != nil {
			// If there was an error saving current chapter, we skip and continue trying next chapters.
			msg := fmt.Sprintf("Error saving %s - Chapter: %s, %s - %s",
				p.Manga.GetTitle("en"), chapter.GetChapterNum(), chapter.GetTitle(), err.Error())
			log.Println(msg)
			errored[index] = struct{}{}
			continue
		}

		core.App.TView.QueueUpdateDraw(func() {
			downloadCell := tview.NewTableCell(readStatus).SetTextColor(utils.MangaPageDownloadStatColor)
			p.Table.SetCell(index, 2, downloadCell)
		})
	}

	var (
		msg     strings.Builder
		modal   *tview.Modal
		modalID string
	)
	// Use unique ID for this particular download.
	modalID = fmt.Sprintf("%s - %s - %v", utils.DownloadFinishedModalID, p.Manga.GetTitle("en"), selection)

	msg.WriteString("Last Download Queue finished.\n")
	msg.WriteString(fmt.Sprintf("Manga: %s\n", p.Manga.GetTitle("en")))
	if len(errored) != 0 {
		// If there were errors, we ask the user whether we want to retry,
		// but we do not retry after a certain amount of re-attempts.
		msg.WriteString("We encountered some errors! Check the log for more details.")
		if attemptNo < maxRetries {
			msg.WriteString("\nRetry failed downloads?")
			modal = confirmModal(modalID, msg.String(), "Retry", func() {
				go p.downloadChapters(errored, attemptNo+1)
			})
		} else {
			msg.WriteString("\nMaximum retries reached.")
			modal = okModal(modalID, msg.String())
		}
	} else {
		msg.WriteString("No errors :>")
		modal = okModal(modalID, msg.String())
	}
	core.App.TView.QueueUpdateDraw(func() {
		ShowModal(modalID, modal)
	})
}

// saveChapter : Save a chapter.
func (p *MangaPage) saveChapter(chapter *mangodex.Chapter) error {
	downloader, err := core.App.Client.AtHome.NewMDHomeClient(
		chapter, core.App.Config.DownloadQuality, core.App.Config.ForcePort443)
	if err != nil {
		return err
	}

	// Create directory to store the current chapter.
	downloadFolder := p.getDownloadFolder(chapter)
	if err = os.MkdirAll(downloadFolder, os.ModePerm); err != nil {
		return err
	}

	// Get the pages to download
	pages := chapter.Attributes.Data
	if core.App.Config.DownloadQuality == "data-saver" {
		pages = chapter.Attributes.DataSaver
	}

	// Save each page.
	for num, page := range pages {
		// Get image data.
		image, err := downloader.GetChapterPage(page)
		if err != nil {
			return err
		}

		filename := fmt.Sprintf("%04d%s", num+1, filepath.Ext(page))
		filePath := filepath.Join(downloadFolder, filename)
		// Save image
		if err = ioutil.WriteFile(filePath, image, os.ModePerm); err != nil {
			return err
		}
	}

	// If user wants to save the downloads as a zip, then do so.
	if core.App.Config.AsZip {
		if err = p.saveAsZipFolder(downloadFolder); err != nil {
			return err
		}
	}
	return nil
}

// saveAsZipFolder : This function creates a zip folder to store a chapter download.
func (p *MangaPage) saveAsZipFolder(chapterFolder string) error {
	// Create a temporary zip folder to store the zip files. This is because the current images
	// are also stored in their own zip directory as returned from getDownloadFolder.
	tempZip := fmt.Sprintf("%s.%s", chapterFolder, "temp")

	var (
		zipFile *os.File
		err     error
	)

	// Create necessary writers
	if zipFile, err = os.Create(tempZip); err != nil {
		return err
	}
	w := zip.NewWriter(zipFile)

	// Saving the actual files.
	if err = filepath.WalkDir(chapterFolder, func(path string, d fs.DirEntry, err error) error {
		// Stop walking immediately if encounter error
		if err != nil {
			return err
		}
		// Skip if a DirEntry is a folder. By right, this shouldn't happen since any downloads will
		// just contain PNGs or JPEGs, but it's here just in case.
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
		// Fixes zip parsing issues in certain situations.
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
	}); err != nil {
		return err
	}

	// Close the files.
	if err = w.Close(); err != nil {
		return err
	}
	if err = zipFile.Close(); err != nil {
		return err
	}

	// We remove the current unzipped image folder, and rename the temp zip to the real zip.
	if err = os.RemoveAll(chapterFolder); err != nil {
		return err
	}
	if err = os.Rename(tempZip, chapterFolder); err != nil {
		return err
	}
	return err
}

// getDownloadFolder : Get the download folder for a manga's chapter.
func (p *MangaPage) getDownloadFolder(chapter *mangodex.Chapter) string {
	mangaName := p.Manga.GetTitle("en")
	chapterName := fmt.Sprintf("Chapter %s [%s-%s] %s - %s",
		chapter.GetChapterNum(), chapter.Attributes.TranslatedLanguage, core.App.Config.DownloadQuality,
		chapter.GetTitle(), strings.SplitN(chapter.ID, "-", 2)[0])

	// Remove invalid characters from the folder name
	restricted := []string{"<", ">", ":", "/", "|", "?", "*", "\"", "\\", "."}
	for _, c := range restricted {
		mangaName = strings.ReplaceAll(mangaName, c, "-")
		chapterName = strings.ReplaceAll(chapterName, c, "-")
	}

	folder := filepath.Join(core.App.Config.DownloadDir, mangaName, chapterName)
	// If the user wants to download as a zip, then we check for the presence of the zip folder.
	if core.App.Config.AsZip {
		folder = fmt.Sprintf("%s.%s", folder, core.App.Config.ZipType)
	}
	return folder
}

// toggleReadMarkers : Toggle read status for selected chapters.
func (p *MangaPage) toggleReadMarkers(selection map[int]struct{}) {
	// Check if the user is logged in. If they are not, we tell them that they cannot toggle without logging in.
	if !core.App.Client.Auth.IsLoggedIn() {
		log.Printf("Attmpted toggling read marker while not logged in. Informing user...")
		core.App.TView.QueueUpdateDraw(func() {
			modal := okModal(utils.ToggleReadChapterModalErrorID, "You need to log in to toggle read status!")
			ShowModal(utils.ToggleReadChapterModalErrorID, modal)
		})
		return
	}

	// For each selection, we separate into make-read, make-unread bins.
	var read, unRead []string
	for row := range selection {
		var (
			chapter *mangodex.Chapter
			ok      bool
		)
		// Get the chapter for this row.
		if chapter, ok = p.Table.GetCell(row, 0).GetReference().(*mangodex.Chapter); !ok {
			return
		}

		// Get the read/unread status, and split accordingly.
		statusCell := p.Table.GetCell(row, 4)
		if statusCell.Text == readStatus { // If it was originally read, we toggle to unread.
			unRead = append(unRead, chapter.ID)
			statusCell.SetText("")
		} else {
			read = append(read, chapter.ID)
			statusCell.SetText(readStatus)
		}
	}

	// TODO : Wait for API fix.
	// Send the request.
	if _, err := core.App.Client.Chapter.SetReadUnreadMangaChapters(p.Manga.ID, read, unRead); err != nil {
		// Error sending request, tell the user.
		log.Printf("Unable to update read markers: %s\n", err.Error())
		core.App.TView.QueueUpdateDraw(func() {
			modal := okModal(utils.ToggleReadChapterModalErrorID,
				"Error updating read markers.\n\nCheck log for details.")
			ShowModal(utils.ToggleReadChapterModalErrorID, modal)
		})
	}
}
