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

// The following may be changed depending on user custom configuration

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
	confBytes, err := ioutil.ReadFile(confPath)
	if err != nil { // If error, then file does not exist, so we create the config file.
		// Make sure to `usr` directory exists. If it already exists, then nothing is done.
		if e := os.MkdirAll(UsrDir, os.ModePerm); e != nil {
			return e
		}

		// Set default DownloadDir : "downloads", and format JSON properly for user.
		Conf = UserConfig{DownloadDir: DownloadDir, Languages: Languages}
		return SaveConfiguration(confPath)
	}
	// If no error, then we can load the configuration.
	if err = json.Unmarshal(confBytes, &Conf); err != nil {
		return err
	}

	// Check for defaults
	if Conf.DownloadDir == "" {
		Conf.DownloadDir = DownloadDir
	}
	if len(Conf.Languages) == 0 {
		Conf.Languages = Languages
	}

	// Expand any environment variables in the user provided string
	Conf.DownloadDir = os.ExpandEnv(Conf.DownloadDir)

	// Save the config file.
	return SaveConfiguration(confPath)
}

// SaveConfiguration : Save user configuration.
func SaveConfiguration(path string) error {
	newConf, e := json.MarshalIndent(&Conf, "", "\t")
	if e != nil {
		return e
	}
	return ioutil.WriteFile(path, newConf, os.ModePerm)
}
