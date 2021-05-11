package main

/*
This file contain helper functions that would otherwise be too large to fit into main sections of code.
*/

import (
	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
)

const credFile = "MangaDesk_CredFile"

// downloadPages : Attempt to download pages
func downloadPages(pages *tview.Pages, table *tview.Table, chaps *mangodex.ChapterList) {
	// Get selected rows and columns (we can ignore columns since we only set the table to select rows)
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
		errorModal := CreateModal("Authentication failed\nTry again!", []string{"OK"}, func(i int, l string) {
			pages.RemovePage(LoginFailureModalID) // Remove the modal once user acknowledge.
		})
		pages.AddPage(LoginFailureModalID, errorModal, true, false)
		pages.ShowPage(LoginFailureModalID) // Show the modal to the user.
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
