package settings

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// type EmailCredentials struct {
// 	Password string
// 	User     string
// 	Port     string
// 	Imap     string
// 	Smtp     string
// }

// type MadeinlySetting struct {
// 	Version          string
// 	Debug            bool
// 	EmailCredentials EmailCredentials
// }

// type Cors struct {
// 	FrontDomain string
// }

type SettingsModel struct {
	Version float64
	Debug   bool
}

var Settings SettingsModel

func InitSettings() {
	binPath, err := os.Executable()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	settingsPath := filepath.Join(filepath.Dir(binPath), "settings.toml")

	settingsByte, err := os.ReadFile(settingsPath)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	toml.Decode(string(settingsByte), &Settings)

}

func UpdateSettings() {

}
