package globals

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// File paths for settings and configuration.
var (
	credFilePath   = filepath.Join(GetConfDir(), "credentials")
	configFilePath = filepath.Join(GetConfDir(), "config.json")
)

// Defaults for user configuration.
var (
	downloadDir     = "downloads"
	languages       = []string{"en"}
	downloadQuality = "data"
	zipType         = "zip"
)

// UserConfig : This struct contains user configurable settings.
type UserConfig struct {
	DownloadDir     string   `json:"downloadDir"`
	Languages       []string `json:"languages"`
	DownloadQuality string   `json:"downloadQuality"`
	ForcePort443    bool     `json:"forcePort443"`
	AsZip           bool     `json:"asZip"`
	ZipType         string   `json:"zipType"`
}

// LoadCredentials : Load saved credentials if there is. Else, return error.
func LoadCredentials() ([]byte, error) {
	return ioutil.ReadFile(credFilePath)
}

// LoadUserConfiguration : Reads any user configuration settings and will create a default one if it does not exist.
func LoadUserConfiguration() error {
	// Set the current configuration to empty one.
	Conf = UserConfig{}

	// Attempt to read user configuration file.
	if confBytes, err := ioutil.ReadFile(configFilePath); err == nil {
		// If no error, attempt unmarshal
		err = json.Unmarshal(confBytes, &Conf)
		if err != nil { // Return error if cannot unmarshal.
			return err
		}
	}
	// Set defaults
	SanitiseConfigurations()

	// Save the config file.
	return SaveConfiguration()
}

// SaveConfiguration : Save user configuration.
func SaveConfiguration() error {
	// Format JSON properly for user.
	confBytes, err := json.MarshalIndent(&Conf, "", "\t")
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
func SanitiseConfigurations() {
	// Download Directory
	if Conf.DownloadDir == "" {
		Conf.DownloadDir = downloadDir
	}
	// Expand any environment variables in the path.
	Conf.DownloadDir = os.ExpandEnv(Conf.DownloadDir)

	// Languages
	if len(Conf.Languages) == 0 {
		Conf.Languages = languages
	}

	// ForcePort443 is false by default.

	// Download Quality
	// Will automatically set to `data` if invalid or no download quality specified.
	if Conf.DownloadQuality != "data" && Conf.DownloadQuality != "data-saver" {
		Conf.DownloadQuality = downloadQuality
	}

	// AsZip is false by default.

	// Set default zip download type. Can be `zip` or `cbz`.
	// Any other invalid entries will default to `zip`.
	if Conf.ZipType != "zip" && Conf.ZipType != "cbz" {
		Conf.ZipType = zipType
	}
}

// GetConfDir : Find the operating system and determine the configuration directory for the application.
func GetConfDir() string {
	appDir := "mangadesk"

	// Get the default configuration appDir for the OS.
	configDir, err := os.UserConfigDir()
	if err != nil { // If there is an error, then we use the home appDir.
		configDir, err = os.UserHomeDir()
		if err != nil { // If still fail, then we use `usr` folder in current appDir.
			configDir = "usr"
		}
	}

	// Add the app directory to the path.
	return filepath.Join(configDir, appDir)
}
