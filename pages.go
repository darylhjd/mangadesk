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
	"strings"
)

const (
	LoginPageID  = "login_page"
	MainPageID   = "main_page"
	MangaPageID  = "manga_page"
	HelpPageID   = "help_page"
	SearchPageID = "search_page"

	LoginLogoutFailureModalID   = "login_failure_modal"
	LoginLogoutCfmModalID       = "logout_modal"
	StoreCredentialErrorModalID = "store_cred_error_modal"
	DownloadChaptersModalID     = "download_chapters_modal"
	DownloadErrorModalID        = "download_error_modal"
	GenericAPIErrorModalID      = "api_error_modal"
)

// ShowLoginPage : Show login page.
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
		AddItem(form, 1, 1, 1, 1, 32, 70, true)

	pages.AddPage(LoginPageID, grid, true, false)
	app.SetFocus(grid)
	pages.SwitchToPage(LoginPageID)
}

// ShowMainPage : Show the main page. Can be for logged user or guest user.
func ShowMainPage(pages *tview.Pages) {
	// Create the base main grid.
	// 15x15 grid.
	var g []int
	for i := 0; i < 15; i++ { // This is to create 15 grids.
		g = append(g, -1)
	}
	grid := tview.NewGrid().SetColumns(g...).SetRows(g...)
	// Set grid attributes.
	grid.SetTitleColor(tcell.ColorOrange).
		SetBorderColor(tcell.ColorLightGrey).
		SetBorder(true)

	// Create the base main table.
	table := tview.NewTable()
	table.SetSelectable(true, false). // Sets only the rows to be selectable
						SetSeparator('|').
						SetBordersColor(tcell.ColorGrey).
						SetTitleColor(tcell.ColorLightSkyBlue).
						SetBorder(true)

	// Add the table to the grid.
	grid.AddItem(table, 0, 0, 15, 15, 0, 0, true)

	if dex.IsLoggedIn() {
		setUpLoggedMainPage(pages, grid, table)
	} else {
		params := url.Values{}
		params.Add("limit", "75")
		grid.SetTitle("Welcome to MangaDex, [red]Guest!")
		table.SetTitle("Recently updated manga")
		setUpGuestMainPage(pages, table, params)
	}

	pages.AddPage(MainPageID, grid, true, false)
	app.SetFocus(grid)
	pages.SwitchToPage(MainPageID)
}

// setUpLoggedMainPage : Set up the main page for a logged user.
func setUpLoggedMainPage(pages *tview.Pages, grid *tview.Grid, table *tview.Table) {
	// For logged users, we fill the table with their followed manga.
	// Get user information
	username := "?"
	if u, err := dex.GetLoggedUser(); err == nil {
		username = u.Data.Attributes.Username
	}
	grid.SetTitle(fmt.Sprintf("Welcome to MangaDex, [lightgreen]%s!", username))
	table.SetTitle("Your followed manga")

	titleColor := tcell.ColorLightGoldenrodYellow
	statusColor := tcell.ColorSaddleBrown

	// Set up table
	mangaTitleHeader := tview.NewTableCell("Manga").
		SetAlign(tview.AlignCenter).
		SetTextColor(titleColor).
		SetSelectable(false)
	statusHeader := tview.NewTableCell("Status").
		SetAlign(tview.AlignCenter).
		SetTextColor(statusColor).
		SetSelectable(false)
	table.SetCell(0, 0, mangaTitleHeader).
		SetCell(0, 1, statusHeader).
		SetFixed(1, 0)

	go func() {
		// Get the manga list.
		mangaList, err := dex.GetUserFollowedMangaList(50, 0)
		if err != nil {
			app.QueueUpdateDraw(func() {
				ShowModal(pages, GenericAPIErrorModalID, "Error getting followed manga.", []string{"OK"},
					func(i int, label string) {
						pages.RemovePage(GenericAPIErrorModalID)
					})
			})
			return
		}

		// Set up the selected function.
		table.SetSelectedFunc(func(row, column int) {
			ShowMangaPage(pages, &(mangaList.Results[row-1]))
		})

		// Add the entries into the table.
		for i, mr := range mangaList.Results {
			// Create the manga title cell
			mtCell := tview.NewTableCell(fmt.Sprintf("%-50s", mr.Data.Attributes.Title["en"])).
				SetMaxWidth(50)
			mtCell.Color = titleColor

			// Create status cell
			status := "-"
			if mr.Data.Attributes.Status != nil {
				status = strings.Title(*mr.Data.Attributes.Status)
			}
			sCell := tview.NewTableCell(fmt.Sprintf("%-15s", status)).
				SetMaxWidth(15)
			sCell.Color = statusColor

			app.QueueUpdateDraw(func() {
				table.SetCell(i+1, 0, mtCell).
					SetCell(i+1, 1, sCell)
			})
		}
	}()
}

// setUpGuestMainPage : Set up the main page for a guest user.
func setUpGuestMainPage(pages *tview.Pages, table *tview.Table, params url.Values) {
	// For guest users, we fill the table with recently updated manga.
	titleColor := tcell.ColorOrange
	descColor := tcell.ColorLightGrey
	tagColor := tcell.ColorLightSteelBlue

	// Set up the table.
	mangaTitleHeader := tview.NewTableCell("Manga").
		SetAlign(tview.AlignCenter).
		SetTextColor(titleColor).
		SetSelectable(false)
	descHeader := tview.NewTableCell("Description").
		SetAlign(tview.AlignCenter).
		SetTextColor(descColor).
		SetSelectable(false)
	tagHeader := tview.NewTableCell("Tags").
		SetAlign(tview.AlignCenter).
		SetTextColor(tagColor).
		SetSelectable(false)
	table.SetCell(0, 0, mangaTitleHeader).
		SetCell(0, 1, descHeader).
		SetCell(0, 2, tagHeader).
		SetFixed(1, 0)

	go func() {
		// Get recently updated manga.
		mangaList, err := dex.MangaList(params)
		if err != nil {
			app.QueueUpdateDraw(func() {
				ShowModal(pages, GenericAPIErrorModalID, "Error loading manga list.", []string{"OK"},
					func(i int, label string) {
						pages.RemovePage(GenericAPIErrorModalID)
					})
			})
			return
		}

		// Set up the selected function.
		table.SetSelectedFunc(func(row, column int) {
			ShowMangaPage(pages, &(mangaList.Results[row-1]))
		})

		// Add the entries into the table
		for i, mr := range mangaList.Results {
			// Create the manga title cell
			mtCell := tview.NewTableCell(fmt.Sprintf("%-40s", mr.Data.Attributes.Title["en"])).
				SetMaxWidth(40)
			mtCell.Color = titleColor

			// Create the description cell
			desc := tview.Escape(fmt.Sprintf("%-70s", mr.Data.Attributes.Description["en"]))
			descCell := tview.NewTableCell(desc).
				SetMaxWidth(70)
			descCell.Color = descColor

			// Create the tag cell
			tags := make([]string, len(mr.Data.Attributes.Tags))
			for i, tag := range mr.Data.Attributes.Tags {
				tags[i] = tag.Attributes.Name["en"]
			}
			tagCell := tview.NewTableCell(strings.Join(tags, ", "))
			tagCell.Color = tagColor

			app.QueueUpdateDraw(func() {
				table.SetCell(i+1, 0, mtCell).
					SetCell(i+1, 1, descCell).
					SetCell(i+1, 2, tagCell)
			})
		}
	}()
}

// ShowMangaPage : Show the manga page.
func ShowMangaPage(pages *tview.Pages, mr *mangodex.MangaResponse) {
	// Create the base main grid.
	// 15x15 grid.
	var g []int
	for i := 0; i < 15; i++ { // This is to create 15 grids.
		g = append(g, -1)
	}
	grid := tview.NewGrid().SetColumns(g...).SetRows(g...)
	// Set grid attributes
	grid.SetTitleColor(tcell.ColorOrange).
		SetBorderColor(tcell.ColorLightGrey).
		SetTitle("Manga Information").
		SetBorder(true)

	// Set input handlers for this page
	setMangaPageHandlers(pages, grid)

	// Use a textview for basic information of the manga.
	info := tview.NewTextView()
	go func() {
		setMangaInfo(info, mr)
	}()
	// Set textview attributes
	info.SetWrap(true).SetWordWrap(true).
		SetBorderColor(tcell.ColorLightGrey).
		SetTitleColor(tcell.ColorLightSkyBlue).
		SetTitle("About").
		SetBorder(true)

	// Use a table to show the chapters for the manga.
	table := tview.NewTable()
	// Set chapter headers
	numHeader := tview.NewTableCell("Chap").
		SetTextColor(tcell.ColorLightYellow).
		SetSelectable(false)
	titleHeader := tview.NewTableCell("Name").
		SetTextColor(tcell.ColorLightSkyBlue).
		SetSelectable(false)
	table.SetCell(0, 0, numHeader).
		SetCell(0, 1, titleHeader).
		SetFixed(1, 0)
	go func() {
		setMangaChaptersTable(pages, table, mr)
	}()
	// Set table attributes
	table.SetSelectable(true, false).
		SetSeparator('|').
		SetBordersColor(tcell.ColorGrey).
		SetTitle("Read").
		SetTitleColor(tcell.ColorLightSkyBlue).
		SetBorder(true)

	grid.AddItem(info, 0, 0, 5, 15, 0, 0, false).
		AddItem(table, 5, 0, 10, 15, 0, 0, true).
		AddItem(info, 0, 0, 15, 5, 0, 80, false).
		AddItem(table, 0, 5, 15, 10, 0, 80, true)

	pages.AddPage(MangaPageID, grid, true, false)
	app.SetFocus(grid)
	pages.SwitchToPage(MangaPageID)
}

// ShowHelpPage : Show the help page to the user.
func ShowHelpPage(pages *tview.Pages) {
	helpText := "Keyboard Mappings\n" +
		"-----------------------------\n\n" +
		"Universal\n" +
		fmt.Sprintf("%-15s:%15s\n", "Ctrl + L", "Login/Logout") +
		fmt.Sprintf("%-15s:%15s\n", "Ctrl + H", "Show Help") +
		fmt.Sprintf("%-15s:%15s\n\n", "Ctrl + S", "Search") +
		"Manga Page\n" +
		fmt.Sprintf("%-15s:%15s\n", "Ctrl + E", "Select mult.") +
		fmt.Sprintf("%-15s:%15s\n\n", "Enter", "Start download") +
		"Works for most pages\n" +
		fmt.Sprintf("%-15s:%15s\n", "Esc", "Go back")

	help := tview.NewTextView()
	help.SetText(helpText).
		SetTextAlign(tview.AlignCenter).
		SetBorderColor(tcell.ColorLightGrey).
		SetBorder(true)

	// Create a new grid for the text view so we can align it to the center.
	grid := tview.NewGrid().SetColumns(0, 0, 0, 0).SetRows(0, 0, 0, 0).
		AddItem(help, 0, 0, 4, 4, 0, 0, false).
		AddItem(help, 1, 1, 2, 2, 30, 100, false)
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			pages.RemovePage(HelpPageID)
		}
		return event
	})

	pages.AddPage(HelpPageID, grid, true, false)
	app.SetFocus(grid)
	pages.SwitchToPage(HelpPageID)
}

func ShowSearchPage(pages *tview.Pages) {
	// Create the base main grid.
	// 15x15 grid.
	var g []int
	for i := 0; i < 15; i++ { // This is to create 15 grids.
		g = append(g, -1)
	}
	grid := tview.NewGrid().SetColumns(g...).SetRows(g...)
	// Set grid attributes
	grid.SetTitleColor(tcell.ColorOrange).
		SetTitle("Search for Manga").
		SetBorderColor(tcell.ColorLightGrey).
		SetBorder(true)

	// Set input handlers for this case
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			pages.RemovePage(SearchPageID)
		}
		return event
	})

	// Create table to show manga list.
	table := tview.NewTable()
	table.SetSelectable(true, false). // Sets only the rows to be selectable
						SetSeparator('|').
						SetBordersColor(tcell.ColorGrey).
						SetTitleColor(tcell.ColorLightSkyBlue).
						SetTitle("Search Results. [grey]Press Tab to go back to search bar.").
						SetBorder(true)

	// Create a form for the searching
	search := tview.NewForm()
	// Set form attributes
	search.SetButtonsAlign(tview.AlignLeft).
		SetLabelColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorDarkSlateGrey).
		SetFieldTextColor(tcell.ColorFloralWhite).
		SetButtonBackgroundColor(tcell.ColorDodgerBlue)
	// Add form fields
	search.AddInputField("Search Manga:", "", 0, nil, nil).
		AddButton("Search", func() {
			searchTerm := search.GetFormItemByLabel("Search Manga:").(*tview.InputField).GetText()

			params := url.Values{}
			params.Add("limit", "75")
			params.Add("title", searchTerm)
			setUpGuestMainPage(pages, table, params)
			app.SetFocus(table)
		})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyTab:
			app.SetFocus(search)
		}
		return event
	})

	grid.AddItem(search, 0, 0, 3, 15, 0, 0, true).
		AddItem(table, 3, 0, 12, 15, 0, 0, false)

	pages.AddPage(SearchPageID, grid, true, false)
	app.SetFocus(grid)
	pages.ShowPage(SearchPageID)
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
