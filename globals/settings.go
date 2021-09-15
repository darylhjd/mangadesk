package globals

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

const (
	CredFileName   = "credentials"
	ConfigFileName = "config.json"
)

// The following are defaults for user configuration.

var (
	DownloadDir     = "downloads"
	Languages       = []string{"en"}
	DownloadQuality = "data"
	ZipType         = "zip"
)

// UserConfig : This struct contains information for user configurable settings.
type UserConfig struct {
	DownloadDir     string   `json:"downloadDir"`
	Languages       []string `json:"languages"`
	DownloadQuality string   `json:"downloadQuality"`
	ForcePort443    bool     `json:"forcePort443"`
	AsZip           bool     `json:"asZip"`
	ZipType         string   `json:"zipType"`
}

// LoadUserConfiguration : Reads any user configuration settings and will create a default one if it does not exist.
func LoadUserConfiguration() error {
	// Path to user configuration file.
	confPath := filepath.Join(GetConfDir(), ConfigFileName)

	// Attempt to read user configuration file.
	if confBytes, err := ioutil.ReadFile(confPath); err != nil { // If error, assume file does not exist.
		// Set defaults and save new configuration.
		SetDefaultConfigurations()
		return SaveConfiguration(confPath)
	} else if err = json.Unmarshal(confBytes, &Conf); err != nil { // If no error reading, then unmarshal.
		return err
	}

	// Check for defaults
	SetDefaultConfigurations()
	// Expand any environment variables in the user provided string
	Conf.DownloadDir = os.ExpandEnv(Conf.DownloadDir)

	// Save the config file.
	return SaveConfiguration(confPath)
}

// SaveConfiguration : Save user configuration.
func SaveConfiguration(path string) error {
	// Format JSON properly for user.
	confBytes, err := json.MarshalIndent(&Conf, "", "\t")
	if err != nil {
		return err
	}

	// Make sure the configuration directory exists. If it already exists, then nothing is done.
	if err = os.MkdirAll(GetConfDir(), os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(path, confBytes, os.ModePerm)
}

// SetDefaultConfigurations : Sets default configurations.
func SetDefaultConfigurations() {
	// Set default download directory if not set.
	if Conf.DownloadDir == "" {
		Conf.DownloadDir = DownloadDir
	}

	// Set default language if not set.
	if len(Conf.Languages) == 0 {
		Conf.Languages = Languages
	}

	// ForcePort443 is false by default.

	// Set default download quality.
	// Will automatically set to `data` if invalid or no download quality specified.
	if Conf.DownloadQuality != "data" && Conf.DownloadQuality != "data-saver" {
		Conf.DownloadQuality = DownloadQuality
	}

	// AsZip is false by default.

	// Set default zip download type. Can be `zip` or `cbz`.
	// Any other invalid entries will default to `zip`.
	if Conf.ZipType != "zip" && Conf.ZipType != "cbz" {
		Conf.ZipType = ZipType
	}
}

// GetConfDir : Find the operating system and determine the configuration directory for the application.
func GetConfDir() string {
	directory := "mangadesk"

	// initialise empty variable here so can be modified below
	UsrDir := ""

	// Looks up XDG_CONFIG_HOME in the environment to check for system set config directory.
	// Supposedly, Linux, BSD, and apparently MacOS uses this variable.
	if configDir, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		UsrDir = filepath.Join(configDir, directory)
	} else {
		switch runtime.GOOS {
		case "linux", "freebsd":
			UsrDir = filepath.Join(os.Getenv("HOME"), ".config", directory)
		case "darwin":
			UsrDir = filepath.Join(os.Getenv("HOME"), "Library", "Preferences", directory)
		case "windows":
			// LOCALAPPDATA always available on Windows environments I believe
			UsrDir = filepath.Join(os.Getenv("LOCALAPPDATA"), directory)
		default:
			// Save in `usr` folder in current directory which the application is run from.
			UsrDir = "usr"
		}
	}
	return UsrDir
}
