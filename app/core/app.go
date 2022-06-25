package core

import (
	"fmt"
	"log"
	"os"

	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"
)

// App : Global App variable.
var (
	App        *MangaDesk
	AppVersion = "MangaDesk v0.7.8"
)

// MangaDesk : The client for this application.
type MangaDesk struct {
	Client *mangodex.DexClient

	TView      *tview.Application
	PageHolder *tview.Pages

	Config  *UserConfig
	LogFile *os.File
}

// Initialise : Initialise the app. Return error if unable to restore previous session.
func (m *MangaDesk) Initialise() error {
	// Set up logging.
	if err := m.setUpLogging(); err != nil {
		fmt.Println("Unable to set up logging...")
		fmt.Println("Application will not continue.")
		os.Exit(1)
	}

	// Load user configuration.
	if err := m.loadConfiguration(); err != nil {
		log.Println("Unable to read configuration file. Is it formatted correctly?")
		log.Println("If in doubt, delete the configuration file to start over!\n\nDetails:")
		log.Println(err.Error())
		os.Exit(1)
	}

	// Set the page holder as the application root and focus on it.
	m.TView.SetRoot(m.PageHolder, true).SetFocus(m.PageHolder)

	// Try to restore the last session so the user does not need to log in again.
	return m.restoreSession()
}

// Shutdown : Stop all services such as logging and let the application shut down gracefully.
func (m *MangaDesk) Shutdown() {
	// Stop all necessary services, such as logging.

	// Sync the screen to make sure that the terminal screen is not corrupted.
	App.TView.Sync()
	App.TView.Stop()

	// Stop the logging
	if err := m.stopLogging(); err != nil {
		fmt.Println("Error while closing log file!")
	}
	fmt.Println("Application shutdown.")
}
