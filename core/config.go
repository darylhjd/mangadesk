package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// configFilePath : The filepath to the configuration file.
var configFilePath = filepath.Join(GetConfDir(), "config.json")

// Defaults for user configuration.
var (
	downloadDir     = "downloads"
	languages       = []string{"en"}
	downloadQuality = "data"
	zipType         = "zip"
)

// UserConfig : This struct contains te user configurable settings.
type UserConfig struct {
	DownloadDir     string   `json:"downloadDir"`
	Languages       []string `json:"languages"`
	DownloadQuality string   `json:"downloadQuality"`
	ForcePort443    bool     `json:"forcePort443"`
	AsZip           bool     `json:"asZip"`
	ZipType         string   `json:"zipType"`
}

// LoadConfiguration : Reads any user configuration settings and will create a default one if it does not exist.
func (m *MangaDesk) LoadConfiguration() error {
	// Make sure the configuration directory exists.
	if err := os.MkdirAll(GetConfDir(), os.ModePerm); err != nil {
		return err
	}

	// Set the current configuration to empty one.
	m.Config = &UserConfig{}

	// Attempt to read user configuration file.
	if confBytes, err := ioutil.ReadFile(configFilePath); err == nil {
		// If no error, attempt unmarshal
		err = json.Unmarshal(confBytes, m.Config)
		if err != nil { // Return error if cannot unmarshal.
			return err
		}
	}
	// Set defaults
	m.Config.SanitiseConfigurations()

	// Save the config file.
	return m.SaveConfiguration()
}

// SaveConfiguration : Save user configuration.
func (m *MangaDesk) SaveConfiguration() error {
	// Format JSON properly for user.
	confBytes, err := json.MarshalIndent(m.Config, "", "\t")
	if err != nil {
		return err
	}

	// Make sure the configuration directory exists. If it already exists, then nothing is done.
	if err = os.MkdirAll(GetConfDir(), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(configFilePath, confBytes, os.ModePerm)
}

// SanitiseConfigurations : Sanitises the configuration to ensure validated fields.
func (c *UserConfig) SanitiseConfigurations() {
	// Download Directory
	if c.DownloadDir == "" {
		c.DownloadDir = downloadDir
	}
	// Expand any environment variables in the path.
	c.DownloadDir = os.ExpandEnv(c.DownloadDir)

	// Languages
	if len(c.Languages) == 0 {
		c.Languages = languages
	}

	// ForcePort443 is false by default.

	// Download Quality
	// Will automatically set to `data` if invalid or no download quality specified.
	if c.DownloadQuality != "data" && c.DownloadQuality != "data-saver" {
		c.DownloadQuality = downloadQuality
	}

	// AsZip is false by default.

	// Set default zip download type. Can be `zip` or `cbz`.
	// Any other invalid entries will default to `zip`.
	if c.ZipType != "zip" && c.ZipType != "cbz" {
		c.ZipType = zipType
	}
}

// GetConfDir : Find the operating system and determine the configuration directory for the application.
func GetConfDir() string {
	// Get the default configuration appDir for the OS.
	configDir, err := os.UserConfigDir()
	if err != nil { // If there is an error, then we use the home appDir.
		configDir, err = os.UserHomeDir()
		if err != nil { // If still fail, then we use `usr` folder in current appDir.
			configDir = "usr"
		}
	}

	// Add the core directory to the path.
	return filepath.Join(configDir, "mangadesk")
}
