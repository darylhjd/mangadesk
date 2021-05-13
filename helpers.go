package main

/*
This file contain helper functions that would otherwise be too large to fit into main sections of code.
*/

import (
	"fmt"
	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	credFile       = "MangaDesk_CredFile"
	DownloadFolder = "downloads"
)

// confirmChapterDownloads : Function for handling when the user press enter on the table.
func confirmChapterDownloads(pages *tview.Pages, table *tview.Table,
	selected *map[int]struct{}, row int, mr *mangodex.MangaResponse, chaps *mangodex.ChapterList) {
	// We add the current selection if the there are no selected rows currently.
	if len(*selected) == 0 {
		(*selected)[row] = struct{}{}
	}

	// Show modal to confirm download.
	ShowModal(pages, DownloadChaptersModalID, "Download selection(s)?", []string{"Yes", "No"},
		func(i int, label string) {
			if label == "Yes" {
				downloadChapters(pages, table, selected, mr, chaps)
			}
			pages.RemovePage(DownloadChaptersModalID)
		})
}

// downloadChapters : Attempt to download pages
func downloadChapters(pages *tview.Pages, table *tview.Table, selected *map[int]struct{},
	mr *mangodex.MangaResponse, chaps *mangodex.ChapterList) {
	// Download each chapter.
	go func(rows map[int]struct{}, list mangodex.ChapterList) {
		// For each chapter.
		for r := range rows {
			// Get the corresponding ChapterResponse object.
			// We -1 since the first row is the header.
			chapR := list.Results[r-1]

			// Create folder to store manga.
			chapter := "-"
			if chapR.Data.Attributes.Chapter != nil {
				chapter = *(chapR.Data.Attributes.Chapter)
			}

			// Create the downloader for the chapter.
			downloader, err := dex.NewMDHomeClient(chapR.Data.ID, "data", chapR.Data.Attributes.Hash, false)
			if err != nil {
				app.QueueUpdateDraw(func() {
					ShowModal(pages, DownloadErrorModalID,
						fmt.Sprintf("Could not get server to download Chapter %s", chapter),
						[]string{"OK"}, func(i int, label string) {
							pages.RemovePage(DownloadErrorModalID)
						})
				})
				continue
			}

			chapFolder := filepath.Join(DownloadFolder, mr.Data.Attributes.Title["en"], chapter)
			if err = os.MkdirAll(chapFolder, os.ModePerm); err != nil {
				return
			}

			// Get each page and save it.
			var errored []int
			for pageNum, file := range chapR.Data.Attributes.Data {
				image, err := downloader.GetChapterPage(file)
				if err != nil {
					errored = append(errored, pageNum+1)
					continue
				}
				err = ioutil.WriteFile(filepath.Join(chapFolder, file), image, os.ModePerm)
				if err != nil {
					errored = append(errored, pageNum+1)
					continue
				}
			}

			if len(errored) != 0 {
				app.QueueUpdateDraw(func() {
					ShowModal(pages, DownloadErrorModalID,
						fmt.Sprintf("Errors downloading pages: %v", errored), []string{"OK"},
						func(i int, label string) {
							pages.RemovePage(DownloadErrorModalID)
						})
				})
			}
		}
	}(*selected, *chaps)

	// Clear the stored rows and unmark all chapters
	for k := range *selected {
		markChapterUnselected(table, k)
	}
	*selected = map[int]struct{}{} // Empty the map
}

// markChapterSelected : Mark a chapter as being selected by the user on the main page table.
func markChapterSelected(table *tview.Table, row int) {
	chapterCell := table.GetCell(row, 0)
	chapterCell.SetBackgroundColor(tcell.ColorLimeGreen).SetTextColor(tcell.ColorBlack)

	titleCell := table.GetCell(row, 1)
	titleCell.SetBackgroundColor(tcell.ColorLimeGreen).SetTextColor(tcell.ColorBlack)
}

// markChapterUnselected : Mark a chapter as being unselected by the user on the main page table.
func markChapterUnselected(table *tview.Table, row int) {
	cCell := table.GetCell(row, 0)
	cCell.SetTextColor(tcell.ColorLightYellow).SetBackgroundColor(tcell.ColorBlack)

	tCell := table.GetCell(row, 1)
	tCell.SetTextColor(tcell.ColorLightSkyBlue).SetBackgroundColor(tcell.ColorBlack)
}

// setMangaInfo : Populate the manga page about section,
func setMangaInfo(info *tview.TextView, mr *mangodex.MangaResponse) {
	// Get author and artist information
	authorId := ""
	for _, r := range mr.Relationships {
		if r.Type == "author" {
			authorId = r.ID
			break
		}
	}
	var author string
	if authorId == "" {
		author = "-"
	} else {
		a, err := dex.GetAuthor(authorId)
		if err != nil {
			author = "-"
		} else {
			author = a.Data.Attributes.Name
		}
	}

	status := "-"
	if mr.Data.Attributes.Status != nil {
		status = strings.Title(*mr.Data.Attributes.Status)
	}

	infoText := fmt.Sprintf("Title: %s\n\nAuthor: %s\nStatus: %s\n\nDescription:\n%s",
		mr.Data.Attributes.Title["en"], author, status,
		strings.SplitN(tview.Escape(mr.Data.Attributes.Description["en"]), "\n", 2)[0])
	info.SetText(infoText)
	app.Draw()
}

// setMangaChaptersTable : Populate the manga page chapter table.
func setMangaChaptersTable(pages *tview.Pages, table *tview.Table, mr *mangodex.MangaResponse) {
	// Get chapter feed for this manga.
	params := url.Values{}
	params.Set("limit", "500")
	params.Set("locales[]", "en")
	params.Set("order[chapter]", "desc")
	cl, err := dex.MangaFeed(mr.Data.ID, params)
	if err != nil {
		ShowModal(pages, GenericAPIErrorModalID, "Error getting manga feed", []string{"OK"},
			func(i int, label string) {
				pages.RemovePage(GenericAPIErrorModalID)
			})
		return
	}

	// Set input handlers for the table
	selected := map[int]struct{}{}
	setMangaPageTableHandlers(pages, table, &selected, mr, cl)

	for i, cr := range cl.Results {
		app.QueueUpdateDraw(func() {
			// Chapter cell
			c := "-"
			if cr.Data.Attributes.Chapter != nil {
				c = *cr.Data.Attributes.Chapter
			}
			cCell := tview.NewTableCell(fmt.Sprintf("%-5s", c)).SetMaxWidth(5).
				SetTextColor(tcell.ColorLightYellow)

			// Title cell
			tCell := tview.NewTableCell(cr.Data.Attributes.Title).
				SetTextColor(tcell.ColorLightSkyBlue)

			table.SetCell(i+1, 0, cCell)
			table.SetCell(i+1, 1, tCell)
		})
	}
}

/*
Confirmed functions
*/

// attemptLoginAndShowMainPage : Attempts to login to MangaDex API and show corresponding main page.
func attemptLoginAndShowMainPage(pages *tview.Pages, form *tview.Form) {
	// Get username and password input.
	u := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
	p := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	remember := form.GetFormItemByLabel("Remember Me").(*tview.Checkbox).IsChecked()

	// Attempt to login to MangaDex API.
	if err := dex.Login(u, p); err != nil {
		// If there was error during login, we create a Modal to tell the user that the login failed.
		ShowModal(pages, LoginLogoutFailureModalID, "Authentication failed\nTry again!", []string{"OK"},
			func(i int, l string) {
				pages.RemovePage(LoginLogoutFailureModalID) // Remove the modal once user acknowledge.
			})
		return
	}
	// If successful login.
	// Remember the user's login credentials if user wants it.
	if remember {
		storeLoginDetails(pages)
	}
	// Remove the login page as we no longer need it.
	pages.RemovePage(LoginPageID)

	// Then create and switch to main page.
	ShowMainPage(pages)
}

// checkAuth : Check if the user's credentials have been stored before.
// If they are, then read it, and attempt to refresh the token.
// Will return error if any steps fail (authentication failed).
func checkAuth() error {
	if _, err := os.Stat(credFile); os.IsNotExist(err) {
		return err
	}

	// If the file exists, then we read it.
	content, err := ioutil.ReadFile(credFile)
	if err != nil {
		return err
	}

	// Do a refresh of the token to keep it up to date.
	dex.RefreshToken = string(content) // We set the stored refresh token.
	return dex.RefreshSessionToken()
}

// storeLoginDetails : Store the refresh token after logging in.
func storeLoginDetails(pages *tview.Pages) {
	// Store the refresh tokens into a credential file.
	f, err := os.Create(credFile)
	if err != nil {
		ShowModal(pages, StoreCredentialErrorModalID, "Error storing credentials.", []string{"OK"},
			func(i int, l string) {
				pages.RemovePage(StoreCredentialErrorModalID)
			})
	}
	defer func() {
		_ = f.Close()
	}()

	// Write refresh token to the file.
	_, err = f.Write([]byte(dex.RefreshToken))
}
