package main

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

	LoginFailureModalID   = "login_failure_modal"
	LoginLogoutCfmModalID = "logout_modal"
)

// ShowLoginPage : Page to show login form.
func ShowLoginPage(pages *tview.Pages) {
	// Create the form
	form := tview.NewForm()
	form.AddInputField("Username", "", 0, nil, nil)
	form.AddPasswordField("Password", "", 0, '*', nil)
	form.AddCheckbox("Remember Me", false, nil)
	form.AddButton("Login", func() {
		attemptLoginAndShowMainPage(pages, form)
	})
	form.AddButton("Guest", func() {
		ShowMainPage(pages)
	})
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetTitle("Login to MangaDex")
	form.SetBorder(true)

	// Create a new grid for the form. This is to align the form to the centre.
	grid := tview.NewGrid().SetColumns(0, 0, 0).SetRows(0, 0, 0)
	grid.AddItem(form, 1, 1, 1, 1, 0, 0, true)

	pages.AddPage(LoginPageID, grid, true, false)
	pages.SwitchToPage(LoginPageID)
}

// createMainPage : Creates the basic template for the main page.
// Works for both Guest and Logged account.
func createMainPage(title, user string, c tcell.Color, chaps []mangodex.ChapterResponse) *tview.Grid {
	// Create main page grid.
	grid := tview.NewGrid()
	// 15x15 grid.
	var g []int
	for i := 0; i < 15; i++ {
		g = append(g, -1)
	}
	grid.SetColumns(g...).SetRows(g...)
	grid.SetTitle("Welcome to MangaDex!").
		SetBackgroundColor(tcell.ColorBlack).SetTitleColor(tcell.ColorOrange).SetBorderColor(tcell.ColorDarkGray).
		SetBorder(true)

	// Set username.
	username := tview.NewTextView().SetText(fmt.Sprintf("Logged in as %s", user)).
		SetTextColor(c).SetTextAlign(tview.AlignCenter).
		SetWrap(true).SetWordWrap(true)

	// Add the user info box to the main grid.
	grid.AddItem(username, 0, 0, 1, 15, 0, 0, false).
		AddItem(username, 0, 12, 1, 3, 0, 60+len(user), false)

	// Show chapters.
	list := tview.NewList()
	list.SetTitle(title).SetTitleColor(tcell.ColorBlue).SetBorder(true)
	for _, c := range chaps {
		list.InsertItem(-1, fmt.Sprintf("%s, Chapter %s", c.Data.Attributes.Title, c.Data.Attributes.Chapter),
			"", 0, nil)
	}
	list.SetWrapAround(false)
	// Add the list to the main grid.
	grid.AddItem(list, 1, 0, 14, 15, 0, 0, true).
		AddItem(list, 0, 0, 15, 12, 0, 60+len(user), true)

	return grid
}

// ShowMainPage : Create and show the main page.
func ShowMainPage(pages *tview.Pages) {
	var page *tview.Grid
	if dex.IsLoggedIn() {
		page = createLoggedMainPage()
	} else {
		page = createGuestMainPage()
	}

	pages.AddPage(MainPageID, page, true, false)
	pages.SwitchToPage(MainPageID)
}

// createLoggedMainPage : Convenience wrapper for createMainPage but for a logged in user.
func createLoggedMainPage() *tview.Grid {
	// Get user info.
	u, err := dex.GetLoggedUser()
	if err != nil {
		panic(err)
	}

	// Get chapter responses for logged user's followed manga.
	params := url.Values{}
	params.Add("limit", "50")
	params.Add("locales[]", "en")
	chapters, err := dex.GetUserFollowedMangaChapterFeed(params)
	if err != nil {
		panic(err)
	}

	return createMainPage("Your Manga Feed",
		u.Data.Attributes.Username, tcell.ColorMediumSpringGreen, chapters.Results)
}

// createGuestMainPage : Convenience wrapper for createMainPage but for a guest user.
func createGuestMainPage() *tview.Grid {
	// Get recently uploaded chapters.
	params := url.Values{}
	params.Add("limit", "50")
	params.Add("order[createdAt]", "desc")
	chapters, err := dex.ChapterList(params)
	if err != nil {
		panic(err)
	}

	return createMainPage("Recently Updated Manga",
		"Guest", tcell.ColorSaddleBrown, chapters.Results)
}

// CreateModal : Convenience function to create a modal.
func CreateModal(text string, buttons []string, f func(buttonIndex int, buttonLabel string)) *tview.Modal {
	m := tview.NewModal()
	m.SetText(text).AddButtons(buttons).SetDoneFunc(f).SetFocus(0).SetBorder(true)
	return m
}
