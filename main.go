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

	// Set up pages.
	loginPage := LoginPage(pages)

	pages.AddPage(LoginPageID, loginPage, true, true)
	pages.SwitchToPage(LoginPageID)

	// Run the app.
	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
