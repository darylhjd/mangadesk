package service

import (
	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"
	"log"

	"github.com/darylhjd/mangadesk/app/core"
	"github.com/darylhjd/mangadesk/app/ui"
)

// Start : Set up the application.
func Start() {
	// Create new app.
	core.App = &core.MangaDesk{
		Client:     mangodex.NewDexClient(),
		TView:      tview.NewApplication(),
		PageHolder: tview.NewPages(),
	}

	// Show appropriate screen based on restore session result. If user has
	// guestLogin on their config, we don't show the login page.
	if err := core.App.Initialise(); err != nil && !core.App.Config.GuestLogin {
		ui.ShowLoginPage()
	} else {
		ui.ShowMainPage()
	}

	log.Println("Initialised starting screen.")
	ui.SetUniversalHandlers()

	// Run the app.
	log.Println("Running app...")
	if err := core.App.TView.Run(); err != nil {
		log.Println(err)
	}
}

// Shutdown : Shutdown the application.
func Shutdown() {
	core.App.Shutdown()
}
