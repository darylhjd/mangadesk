package main

/*
This file contains functions that deal with MangaDex's API.
*/

import (
	"github.com/rivo/tview"
)

// LoginToMangaDex : Function to handle logging in and transition to main page.
func LoginToMangaDex(pages *tview.Pages, f *tview.Form) {
	u := f.GetFormItemByLabel("Username").(*tview.InputField).GetText()
	p := f.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	err := dex.Login(u, p)
	if err != nil {
		// If there was error during login, we create a Modal to tell the user that the login failed.
		errorModal := ErrorModal(pages, "Login failed. Try again!", LoginFailureModalID)

		pages.AddPage(LoginFailureModalID, errorModal, true, false)
		pages.ShowPage(LoginFailureModalID) // Show the modal to the user.
	} else {
		// Remove the login page as we no longer need it.
		pages.RemovePage(LoginPageID)

		// Then create the main page.
		mainPage := GuestMainPage(pages)

		pages.AddPage(MainPageID, mainPage, true, true)
		pages.SwitchToPage(MainPageID)
	}
}

// GuestToMangaDex : Do not attempt logging in to API and just show the guest main page.
func GuestToMangaDex(pages *tview.Pages) {
	mainPage := GuestMainPage(pages)

	pages.AddPage(MainPageID, mainPage, true, true)
	pages.SwitchToPage(MainPageID)
}
