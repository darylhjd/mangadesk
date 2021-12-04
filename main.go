package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rivo/tview"

	"github.com/darylhjd/mangadesk/globals"
	p "github.com/darylhjd/mangadesk/pages"
)

// Start the program.
func main() {
	// Load user configuration.
	if err := globals.LoadUserConfiguration(); err != nil {
		fmt.Println("Unable to read configuration file. Is it formatted correctly?")
		fmt.Println("If in doubt, delete the configuration file to start over!\n\nDetails:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create a new page holder.
	pages := tview.NewPages()

	// Set input captures that are valid for the whole app.
	p.SetUniversalHandlers(pages)

	// Check whether the user is remembered. If they are, then load credentials into the client and refresh token.
	if err := RestoreSession(); err != nil {
		// Prompt login if unable to read stored credentials.
		p.ShowLoginPage(pages)
	} else {
		// If able to log in using stored refresh token, go to logged main page.
		p.ShowMainPage(pages)
	}

	// Run the app. SetRoot also calls SetFocus on the primitive.
	if err := globals.App.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}

// RestoreSession : Check if the user's credentials have been stored before.
// If they are, then read it, and attempt to refresh the token.
// Will return error if any steps fail (no stored credentials, authentication failed).
func RestoreSession() error {
	// Try to read stored credential file.
	content, err := globals.LoadCredentials()
	if err != nil { // If error, then user was not originally logged in.
		fmt.Println("No past session, using Guest account...")
		time.Sleep(time.Millisecond * 750)
		return err
	}

	fmt.Println("Attempting session restore...")
	globals.DexClient.Auth.SetRefreshToken(string(content)) // We set the stored refresh token.

	// Do a refresh of the token to keep it up to date. If the token has already expired, user needs to log in again.
	if err = globals.DexClient.Auth.RefreshSessionToken(); err != nil {
		fmt.Println("Session expired. Please login again.")
		time.Sleep(time.Millisecond * 750)
		return err
	}
	return nil
}
