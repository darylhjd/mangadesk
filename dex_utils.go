package main

/*
This file contains functions that deal with MangaDex's API.
*/

import (
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
)

const credFile = "credentials"

// LoginToMangaDex : Function to handle logging in and transition to main page.
func LoginToMangaDex(pages *tview.Pages, f *tview.Form) {
	// Get username and password input.
	u := f.GetFormItemByLabel("Username").(*tview.InputField).GetText()
	p := f.GetFormItemByLabel("Password").(*tview.InputField).GetText()
	remember := f.GetFormItemByLabel("Remember Me").(*tview.Checkbox).IsChecked()

	// Attempt to login to MangaDex API.
	if err := dex.Login(u, p); err != nil {
		// If there was error during login, we create a Modal to tell the user that the login failed.
		errorModal := ErrorModal(pages, "Login failed. Try again!", LoginFailureModalID)
		pages.AddPage(LoginFailureModalID, errorModal, true, false)
		pages.ShowPage(LoginFailureModalID) // Show the modal to the user.
		return
	}

	// Remember the user's login credentials if user wants it.
	if remember {
		RememberLoginDetails()
	}
	// Remove the login page as we no longer need it.
	pages.RemovePage(LoginPageID)

	// Then create and switch to main page.
	mainPage := LoggedMainPage(pages)
	pages.AddPage(LoggedMainPageID, mainPage, true, true)
	pages.SwitchToPage(LoggedMainPageID)
}

// LogoutOfMangaDex : Function to handle logging out and transition.
func LogoutOfMangaDex(pages *tview.Pages) {
	// Create a modal to ask if the user is sure they want to logout
	sureModal := tview.NewModal()
	sureModal.SetBorder(true).SetTitle("Logout")
	sureModal.SetText("Are you sure?")
	sureModal.AddButtons([]string{"Yes", "No"})
	sureModal.SetDoneFunc(func(i int, label string) {
		if label == "Yes" {
			err := dex.Logout()
			if err != nil {
				panic(err)
			}
			// We redirect the user to the guest main page and also remove the logged in main page.
			pages.RemovePage(LoggedMainPageID)
			pages.AddPage(GuestMainPageID, GuestMainPage(pages), true, true)
			pages.SwitchToPage(GuestMainPageID)
		}
		// Regardless of what is pressed, we remove the modal from the page.
		pages.RemovePage(LogoutModalID)
	})

	// We show the modal on top of the screen.
	pages.AddPage(LogoutModalID, sureModal, true, true)
	pages.ShowPage(LogoutModalID)
}

// GuestToMangaDex : Do not attempt logging in to API and just show the guest main page.
func GuestToMangaDex(pages *tview.Pages) {
	mainPage := GuestMainPage(pages)

	pages.AddPage(GuestMainPageID, mainPage, true, true)
	pages.SwitchToPage(GuestMainPageID)
}

// CheckStoredAuth : Check if the user's credentials have been stored before.
// If they are, then read it, else return error.
func CheckStoredAuth() error {
	if _, err := os.Stat(credFile); os.IsNotExist(err) {
		return err
	}

	// If the file exists, then we read it.
	content, err := ioutil.ReadFile(credFile)
	if err != nil {
		return err
	}

	// Do a refresh of the token to keep it up to date.
	dex.RefreshToken = string(content)
	return dex.RefreshSessionToken()
}

// RememberLoginDetails : Store the refresh token after logging in.
func RememberLoginDetails() {
	// Store the session and refresh tokens into a config file for now.
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
