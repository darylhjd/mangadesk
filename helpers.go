package main

/*
This file contain helper functions that would otherwise be too large to fit into main sections of code.
*/

import (
	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	credFile       = "MangaDesk_CredFile"
	DownloadFolder = "downloads"
)

// mainPageTableSelectedFunc : Function for handling when the user press enter on the table.
func mainPageTableSelectedFunc(pages *tview.Pages, table *tview.Table,
	row int, selected *map[int]struct{}, chaps *mangodex.ChapterList) {
	// We add the current selection if the there are no selected rows currently.
	if len(*selected) == 0 {
		(*selected)[row] = struct{}{}
	}

	// Show modal to confirm download.
	ShowModal(pages, DownloadChaptersModalID, "Download selection(s)?", []string{"Yes", "No"},
		func(i int, label string) {
			if label == "Yes" {
				downloadPages(pages, table, selected, chaps)
			}
			pages.RemovePage(DownloadChaptersModalID)
		})
}

// downloadPages : Attempt to download pages
func downloadPages(pages *tview.Pages, table *tview.Table, sRows *map[int]struct{}, chaps *mangodex.ChapterList) {
	// Download each chapter.
	go func(rows map[int]struct{}, list mangodex.ChapterList) {
		// For each chapter.
		for r := range rows {
			// Get the corresponding ChapterResponse object.
			// We -1 since the first row is the header.
			chapR := list.Results[r-1]

			// Create the downloader for the chapter.
			downloader, err := dex.NewMDHomeClient(chapR.Data.ID, "data", chapR.Data.Attributes.Hash, false)
			if err != nil {
				return
			}

			// Get the name of the manga.
			var mangaName string
			for _, relationship := range chapR.Relationships {
				if relationship.Type == "manga" {
					if m, err := dex.ViewManga(relationship.ID); err != nil {
						return
					} else {
						mangaName = m.Data.Attributes.Title["en"]
					}
					break
				}
			}

			// Create folder to store manga.
			chapFolder := filepath.Join(DownloadFolder, mangaName, chapR.Data.Attributes.Chapter)
			if err = os.MkdirAll(chapFolder, os.ModePerm); err != nil {
				return
			}

			// Get each page and save it.
			for _, file := range chapR.Data.Attributes.Data {
				image, err := downloader.GetChapterPage(file)
				if err != nil {
					return
				}
				err = ioutil.WriteFile(filepath.Join(chapFolder, file), image, os.ModePerm)
				if err != nil {
					return
				}
			}
		}
	}(*sRows, *chaps)

	// Clear the stored rows and unmark all chapters
	for r := range *sRows {
		markChapterUnselected(table, r, tcell.ColorBlack, tcell.ColorWhite)
	}
	*sRows = map[int]struct{}{} // Empty the map
}

// markChapterSelected : Mark a chapter as being selected by the user on the main page table.
func markChapterSelected(table *tview.Table, row int, background, text tcell.Color) {
	mangaNameCell := table.GetCell(row, 0)
	chapterCell := table.GetCell(row, 1)

	mangaNameCell.SetBackgroundColor(background).SetTextColor(text)
	chapterCell.SetBackgroundColor(background).SetTextColor(text)
}

// markChapterUnselected : Mark a chapter as being unselected by the user on the main page table.
func markChapterUnselected(table *tview.Table, row int, background, text tcell.Color) {
	mangaNameCell := table.GetCell(row, 0)
	chapterCell := table.GetCell(row, 1)

	mangaNameCell.SetBackgroundColor(background).SetTextColor(text)
	chapterCell.SetBackgroundColor(background).SetTextColor(text)
}

// attemptLoginAndShowMainPage : Attempts to login to MangaDex API and show corresponding main page.
func attemptLoginAndShowMainPage(pages *tview.Pages, form *tview.Form) {
	// Get username and password input.
	u := form.GetFormItemByLabel("Username").(*tview.InputField).GetText()
	p := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	remember := form.GetFormItemByLabel("Remember Me").(*tview.Checkbox).IsChecked()

	// Attempt to login to MangaDex API.
	if err := dex.Login(u, p); err != nil {
		// If there was error during login, we create a Modal to tell the user that the login failed.
		ShowModal(pages, LoginFailureModalID, "Authentication failed\nTry again!", []string{"OK"},
			func(i int, l string) {
				pages.RemovePage(LoginFailureModalID) // Remove the modal once user acknowledge.
			})
		return
	}
	// If successful login.
	// Remember the user's login credentials if user wants it.
	if remember {
		storeLoginDetails()
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
func storeLoginDetails() {
	// Store the refresh tokens into a credential file.
	f, err := os.Create(credFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}()

	// Write refresh token to the file.
	_, err = f.Write([]byte(dex.RefreshToken))
}
