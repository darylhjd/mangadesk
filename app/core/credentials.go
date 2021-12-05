package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// credFilePath : The filepath to the credentials file.
var credFilePath = filepath.Join(getConfDir(), "credentials")

// StoreCredentials : Store the refresh token.
func (m *MangaDesk) StoreCredentials() error {
	if err := ioutil.WriteFile(credFilePath, []byte(m.Client.Auth.GetRefreshToken()), os.ModePerm); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// DeleteCredentials : Delete saved credentials from the system.
func (m *MangaDesk) DeleteCredentials() {
	if err := os.Remove(credFilePath); err != nil {
		log.Println(err)
	}
}

// loadCredentials : Load saved credentials into the client if there is.
func (m *MangaDesk) loadCredentials() error {
	content, err := ioutil.ReadFile(credFilePath)
	if err != nil {
		return err
	}
	m.Client.Auth.SetRefreshToken(string(content))
	return nil
}

// restoreSession : Check if the user's credentials have been stored before.
// If they are, then read it, and attempt to refresh the token.
// Will return error if any steps fail (no stored credentials, authentication failed).
func (m *MangaDesk) restoreSession() error {
	// Try to read stored credential file.
	if err := m.loadCredentials(); err != nil { // If error, then user was not originally logged in.
		fmt.Println("No past session, using Guest account...")
		time.Sleep(time.Millisecond * 750)
		return err
	}

	fmt.Println("Attempting session restore...")
	// Do a refresh of the token to keep it up to date. If the token has already expired, user needs to log in again.
	if err := m.Client.Auth.RefreshSessionToken(); err != nil {
		fmt.Println("Session expired. Please login again.")
		time.Sleep(time.Millisecond * 750)
		return err
	}
	return nil
}
