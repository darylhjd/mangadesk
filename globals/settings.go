package globals

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	UsrDir         = "usr"
	CredFileName   = "cred"
	ConfigFileName = "usr_config.json"
)

// The following are defaults for user configuration.

var (
	DownloadDir = "downloads"
	Languages   = []string{"en"}
)

// UserConfig : This struct contains information for user configurable settings.
type UserConfig struct {
	DownloadDir string   `json:"downloadDir"`
	Languages   []string `json:"languages"`
}

// LoadUserConfiguration : Reads any user configuration settings and will create a default one if it does not exist.
func LoadUserConfiguration() error {
	// Path to user configuration file.
	confPath := filepath.Join(UsrDir, ConfigFileName)

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

	// Make sure to `usr` directory exists. If it already exists, then nothing is done.
	if err = os.MkdirAll(UsrDir, os.ModePerm); err != nil {
		return err
	}
	return ioutil.WriteFile(path, confBytes, os.ModePerm)
}

// SetDefaultConfigurations : Sets default configurations.
func SetDefaultConfigurations() {
	if Conf.DownloadDir == "" {
		Conf.DownloadDir = DownloadDir
	}
	if len(Conf.Languages) == 0 {
		Conf.Languages = Languages
	}
}
