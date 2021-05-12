package main

/*
This file contains functions to generate pages for the application.
*/

import (
	"fmt"
	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"net/url"
)

const (
	LoginPageID = "login_page"
	MainPageID  = "main_page"

	LoginFailureModalID     = "login_failure_modal"
	LoginLogoutCfmModalID   = "logout_modal"
	DownloadChaptersModalID = "download_chapters_modal"
)

// ShowLoginPage : Page to show login form.
func ShowLoginPage(pages *tview.Pages) {
	// Create the form
	form := tview.NewForm()
	// Set form attributes.
	form.SetButtonsAlign(tview.AlignCenter).
		SetLabelColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorDarkSlateGrey).
		SetFieldTextColor(tcell.ColorFloralWhite).
		SetButtonBackgroundColor(tcell.ColorDodgerBlue).
		SetTitle("Login to MangaDex").
		SetTitleColor(tcell.ColorOrange).
		SetBorder(true).
		SetBorderColor(tcell.ColorGrey)
	// Add form fields.
	form.AddInputField("Username", "", 0, nil, nil).
		AddPasswordField("Password", "", 0, '*', nil).
		AddCheckbox("Remember Me", false, nil).
		AddButton("Login", func() {
			attemptLoginAndShowMainPage(pages, form)
		}).
		AddButton("Guest", func() {
			ShowMainPage(pages)
		})

	// Create a new grid for the form. This is to align the form to the centre.
	grid := tview.NewGrid().SetColumns(0, 0, 0).SetRows(0, 0, 0)
	grid.AddItem(form, 0, 0, 3, 3, 0, 0, true).
		AddItem(form, 1, 1, 1, 1, 20, 50, true)

	pages.AddPage(LoginPageID, grid, true, false)
	app.SetFocus(grid)
	pages.SwitchToPage(LoginPageID)
}

// ShowMainPage : Creates the basic template for the main page.
// Works for both Guest and Logged account.
func ShowMainPage(pages *tview.Pages) {
	// Whether to show guest or logged main page.
	var (
		chaps      *mangodex.ChapterList
		mainTitle  string
		tableTitle string
		params     = url.Values{}
	)
	if dex.IsLoggedIn() {
		u, err := dex.GetLoggedUser()
		if err != nil {
			panic(err)
		}
		mainTitle = fmt.Sprintf("Welcome to MangaDex, [green]%s!", u.Data.Attributes.Username)
		tableTitle = "Feed"
		params.Add("limit", "50")
		params.Add("locales[]", "en")
		chaps, err = dex.GetUserFollowedMangaChapterFeed(params)
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		mainTitle = "Welcome to MangaDex, [red]Guest!"
		tableTitle = "Recently Uploaded Chapters"
		params.Add("limit", "50")
		params.Add("order[createdAt]", "desc")
		chaps, err = dex.ChapterList(params)
		if err != nil {
			panic(err)
		}
	}

	// Create main page grid.
	// 15x15 grid.
	var g []int
	for i := 0; i < 15; i++ { // This is to create 15 grids.
		g = append(g, -1)
	}
	grid := tview.NewGrid().SetColumns(g...).SetRows(g...)
	// Set grid attributes.
	grid.SetTitle(mainTitle).
		SetBackgroundColor(tcell.ColorBlack).
		SetTitleColor(tcell.ColorOrange).
		SetBorderColor(tcell.ColorLightGrey).
		SetBorder(true)

	// Use a table to show chapter information.
	// len(chaps)x2 table.
	table := tview.NewTable()
	// Set the table header
	mangaHeader := tview.NewTableCell("MANGA").
		SetAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorPapayaWhip).
		SetSelectable(false)
	chapterHeader := tview.NewTableCell("CHAPTER").
		SetAlign(tview.AlignLeft).
		SetTextColor(tcell.ColorPapayaWhip).
		SetSelectable(false)
	table.SetCell(0, 0, mangaHeader). // This sets the first row as fixed as it is used as the header.
						SetCell(0, 1, chapterHeader).
						SetFixed(1, 0)
	// Set table attributes
	table.SetSelectable(true, false).
		SetSelectedStyle(tcell.Style{}.Background(tcell.ColorPeachPuff).Foreground(tcell.ColorBlack)).
		SetSeparator('|').
		SetBordersColor(tcell.ColorGrey).
		SetTitle(tableTitle).
		SetTitleColor(tcell.ColorLightSkyBlue).
		SetBorder(true)

	// Set custom input handlers for the table.
	var selected = map[int]struct{}{}
	setMainPageTableInputCaptures(table, &selected)
	table.SetSelectedFunc(func(row, col int) {
		mainPageTableSelectedFunc(pages, table, row, &selected, chaps)
	})

	// Add the rows of data into the table.
	for index, chapR := range chaps.Results {
		titleCell := tview.NewTableCell(fmt.Sprintf("%-50s", chapR.Data.Attributes.Title)).
			SetMaxWidth(50)
		table.SetCell(index+1, 0, titleCell)
		table.SetCellSimple(index+1, 1, chapR.Data.Attributes.Chapter)
	}

	// Add the list to the main grid.
	grid.AddItem(table, 0, 0, 15, 15, 0, 0, true)

	pages.AddPage(MainPageID, grid, true, false)
	app.SetFocus(grid)
	pages.SwitchToPage(MainPageID)
}

// ShowModal : Convenience function to create a modal.
func ShowModal(pages *tview.Pages, label, text string, buttons []string, f func(buttonIndex int, buttonLabel string)) {
	m := tview.NewModal()
	// Set modal attributes
	m.SetText(text).
		AddButtons(buttons).
		SetDoneFunc(f).
		SetFocus(0).
		SetBackgroundColor(tcell.ColorDarkSlateGrey)

	pages.AddPage(label, m, true, false)
	pages.ShowPage(label)
}
