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
		fmt.Println(err.Error())
		os.Exit(1) // Exit program on error.
	}

	// Create new pages holder.
	pages := tview.NewPages()

	// Set required input captures that are valid for the whole app.
	p.SetUniversalInputCaptures(pages)

	// Check whether the user is remembered. If they are, then load credentials into the client and refresh token.
	if err := checkAuth(); err != nil {
		// If error retrieving stored credentials.
		p.ShowLoginPage(pages)
	} else {
		// If can log in using stored refresh token, then straight away go to logged main page.
		p.ShowMainPage(pages)
	}

	// Run the app. SetRoot also calls SetFocus on the primitive.
	if err := g.App.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}

// checkAuth : Check if the user's credentials have been stored before.
// If they are, then read it, and attempt to refresh the token.
// Will return error if any steps fail (authentication failed).
func checkAuth() error {
	// Location of the credentials file.
	credFilePath := filepath.Join(g.UsrDir, g.CredFileName)
	if _, err := os.Stat(credFilePath); os.IsNotExist(err) {
		return err
	}

	// If the file exists, then we read it.
	content, err := ioutil.ReadFile(g.CredFileName)
	if err != nil {
		return err
	}

	g.Dex.RefreshToken = string(content) // We set the stored refresh token.

	// Do a refresh of the token to keep it up to date.
	return g.Dex.RefreshSessionToken()
}
