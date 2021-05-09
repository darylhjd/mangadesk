package main

import (
	"fmt"
	"github.com/rivo/tview"
	"net/url"
	"strconv"
)

const (
	LoginPageID = "login_page"
	MainPageID  = "main_page"

	LoginFailureModalID = "login_failure_modal"
)

// LoginPage : Page to show login form.
func LoginPage(pages *tview.Pages) *tview.Grid {
	// Create the form
	form := tview.NewForm()
	form.AddInputField("Username", "", 0, nil, nil)
	form.AddPasswordField("Password", "", 0, '*', nil)
	form.AddButton("Login", func() {
		LoginToMangaDex(pages, form) // Try logging in.
	})
	form.AddButton("Guest", func() {
		GuestToMangaDex(pages) // Do not attempt login and just show main page.
	})
	form.SetButtonsAlign(tview.AlignCenter)
	form.SetTitle("Login to MangaDex")
	form.SetBorder(true)

	// Create a new grid for the form. This is to align the form to the centre.
	grid := tview.NewGrid().SetColumns(0, 0, 0).SetRows(0, 0, 0)
	grid.AddItem(form, 1, 1, 1, 1, 0, 0, true)

	return grid
}

// MainPage : Page to show main page. Works for both Guest and Logged account.
func MainPage(pages *tview.Pages) *tview.Grid {
	// Create main page grid.
	grid := tview.NewGrid().SetColumns(-2, -1)
	grid.SetBorder(true)
	grid.SetTitle("Main Page")
	return grid
}

// LoggedMainPage : Convenience wrapper for MainPage but for a logged in user.
func LoggedMainPage(pages *tview.Pages) *tview.Grid {
	mGrid := MainPage(pages)

	// Get list of
}

// GuestMainPage : Convenience wrapper for MainPage but for a guest user.
func GuestMainPage(pages *tview.Pages) *tview.Grid {
	mainGrid := MainPage(pages)

	// Get list of recently updated manga.
	params := url.Values{}
	params.Add("limit", strconv.Itoa(50))
	params.Add("order[createdAt]", "asc")
	chapters, err := dex.ChapterList(params)
	if err != nil {
		panic(err)
	}

	list := tview.NewList()
	list.SetBorder(true)
	for _, c := range chapters.Results {
		list.InsertItem(-1, c.Data.Attributes.Title, fmt.Sprintf("Chapter %s", c.Data.Attributes.Chapter), 0, nil)
	}
	list.SetWrapAround(false)
	list.SetTitle("Recently uploaded chapters")

	mainGrid.AddItem(list, 0, 0, 1, 1, 0, 0, true)
	return mainGrid
}

// ErrorModal : Show a modal to the user if there is an error.
func ErrorModal(pages *tview.Pages, err string, idLabel string) *tview.Modal {
	em := tview.NewModal()
	em.SetText(fmt.Sprintf("Error: %s", err))
	em.AddButtons([]string{"OK"})
	em.SetDoneFunc(func(i int, label string) {
		if label == "OK" {
			pages.RemovePage(idLabel) // Remove the modal once user acknowledge.
		}
	})
	em.SetFocus(0)
	return em
}
