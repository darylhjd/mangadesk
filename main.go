package main

import (
	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"
)

var (
	app = tview.NewApplication()
	dex = mangodex.NewDexClient()
)

// Start the program.
func main() {
	app.EnableMouse(true)
	pages := tview.NewPages()

	// Check whether the user is remembered. If they are, then load credentials into the client and refresh token.
	if err := CheckStoredAuth(); err != nil {
		// If error retrieving stored credentials,
		loginPage := LoginPage(pages)
		pages.AddPage(LoginPageID, loginPage, true, true)
		pages.SwitchToPage(LoginPageID)
	} else {
		// If can log in using stored refresh token, then straight away go to logged main page.
		mainPage := LoggedMainPage(pages)
		pages.AddPage(MainPageID, mainPage, true, true)
		pages.SwitchToPage(MainPageID)
	}

	// Run the app.
	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
