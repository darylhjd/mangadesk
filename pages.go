package main

import (
	"fmt"
	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"net/url"
)

const (
	LoginPageID      = "login_page"
	LoggedMainPageID = "logged_main_page"
	GuestMainPageID  = "guest_main_page"

	LoginFailureModalID = "login_failure_modal"
	LogoutModalID       = "logout_modal"
)

// LoginPage : Page to show login form.
func LoginPage(pages *tview.Pages) *tview.Grid {
	// Create the form
	form := tview.NewForm()
	form.AddInputField("Username", "", 0, nil, nil)
	form.AddPasswordField("Password", "", 0, '*', nil)
	form.AddCheckbox("Remember Me", false, nil)
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

// MainPage : Page to show main page. Creates the basic view for the main page.
// Works for both Guest and Logged account.
func MainPage(pages *tview.Pages, user string, chaps []mangodex.ChapterResponse) *tview.Grid {
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

	// Set keyboard triggers.
	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlL {
			LogoutOfMangaDex(pages)
		}
		return event
	})

	// Set username.
	username := tview.NewTextView().SetText(fmt.Sprintf("Logged in as %s", user)).
		SetTextColor(tcell.ColorMediumSpringGreen).SetTextAlign(tview.AlignCenter).
		SetWrap(true).SetWordWrap(true)

	// Add the user info box to the main grid.
	grid.AddItem(username, 0, 0, 1, 15, 0, 0, false).
		AddItem(username, 0, 12, 1, 3, 0, 60+len(user), false)

	// Show chapters.
	list := tview.NewList()
	list.SetTitle("Your Manga Feed!").SetTitleColor(tcell.ColorBlue).SetBorder(true)
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

// LoggedMainPage : Convenience wrapper for MainPage but for a logged in user.
func LoggedMainPage(pages *tview.Pages) *tview.Grid {
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

	return MainPage(pages, u.Data.Attributes.Username, chapters.Results)
}

// GuestMainPage : Convenience wrapper for MainPage but for a guest user.
func GuestMainPage(pages *tview.Pages) *tview.Grid {
	// Get recently uploaded chapters.
	params := url.Values{}
	params.Add("limit", "50")
	params.Add("order[createdAt]", "desc")
	chapters, err := dex.ChapterList(params)
	if err != nil {
		panic(err)
	}

	return MainPage(pages, "Guest", chapters.Results)
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
