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
	// Create new pages holder.
	pages := tview.NewPages()
	SetInputCaptures(pages) // Set required input captures.

	// Check whether the user is remembered. If they are, then load credentials into the client and refresh token.
	if err := checkAuth(); err != nil {
		// If error retrieving stored credentials.
		ShowLoginPage(pages)
	} else {
		// If can log in using stored refresh token, then straight away go to logged main page.
		ShowMainPage(pages)
	}

	// Run the app. SetRoot also calls SetFocus on the primitive.
	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
