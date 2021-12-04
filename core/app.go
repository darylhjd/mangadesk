package core

import (
	"fmt"
	"log"
	"os"

	"github.com/darylhjd/mangadesk/pages"
	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"
)

// App : Global App variable.
var App *MangaDesk

// MangaDesk : The client for this application.
type MangaDesk struct {
	Client *mangodex.DexClient

	ViewApp    *tview.Application
	PageHolder *tview.Pages

	Config  *UserConfig
	LogFile *os.File
}

// Initialise : Set up the application.
func Initialise() {
	// Create new app
	App = &MangaDesk{
		Client:     mangodex.NewDexClient(),
		ViewApp:    tview.NewApplication(),
		PageHolder: tview.NewPages(),
	}

	// Set up logging.
	if err := App.SetUpLogging(); err != nil {
		fmt.Println("Unable to set up logging...")
		fmt.Println("Application will not continue.")
		os.Exit(1)
	}

	// Load user configuration.
	if err := App.LoadConfiguration(); err != nil {
		log.Println("Unable to read configuration file. Is it formatted correctly?")
		log.Println("If in doubt, delete the configuration file to start over!\n\nDetails:")
		log.Println(err.Error())
		os.Exit(1)
	}

	// Set input captures that are valid for the whole core.
	pages.SetUniversalHandlers()

	// Try to restore the last session so the user does not need to log in again.
	if err := App.RestoreSession(); err != nil {
		App.ShowLoginPage()
	} else {
		App.ShowMainPage()
	}

	// Set the page holder as the application root and focus on it.
	App.ViewApp.SetRoot(App.PageHolder, true).SetFocus(App.PageHolder)
}

// Run : Run the app.
func (m *MangaDesk) Run() error {
	return m.ViewApp.Run()
}

// Shutdown : Stop all services such as logging and let the application shut down gracefully.
func (m *MangaDesk) Shutdown() {
	// Stop all necessary services, such as logging.

	// Sync the screen to make sure that the terminal screen is not corrupted.
	App.ViewApp.Sync()

	// Stop the logging
	if err := m.StopLogging(); err != nil {
		fmt.Println("Error while closing log file!")
	}
	fmt.Println("ViewApp stopped.")
}
