package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rivo/tview"

	g "github.com/darylhjd/mangadesk/globals"
	p "github.com/darylhjd/mangadesk/pages"
)

// Start the program.
func main() {
	// Load user configuration.
	if err := g.LoadUserConfiguration(); err != nil {
		fmt.Println("Unable to read configuration file. Is it set correctly?")
		fmt.Println("If in doubt, delete the configuration file to start over!\n\nDetails:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create new pages holder.
	pages := tview.NewPages()

	// Set required input captures that are valid for the whole app.
	p.SetUniversalHandlers(pages)

	// Check whether the user is remembered. If they are, then load credentials into the client and refresh token.
	if err := CheckAuth(); err != nil {
		// Prompt login if unable to read stored credentials.
		p.ShowLoginPage(pages)
	} else {
		// If able to log in using stored refresh token, go to logged main page.
		p.ShowMainPage(pages)
	}

	// Run the app. SetRoot also calls SetFocus on the primitive.
	if err := g.App.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}

// CheckAuth : Check if the user's credentials have been stored before.
// If they are, then read it, and attempt to refresh the token.
// Will return error if any steps fail (no stored credentials, authentication failed).
func CheckAuth() error {
	// Try to read stored credential file.
	content, err := ioutil.ReadFile(filepath.Join(g.UsrDir, g.CredFileName))
	if err != nil { // If error, then file does not exist.
		return err
	}

	fmt.Println("Welcome back!")
	fmt.Println("Restoring session...")
	g.Dex.RefreshToken = string(content) // We set the stored refresh token.

	// Do a refresh of the token to keep it up to date.
	return g.Dex.RefreshSessionToken()
}
