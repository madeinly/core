package settings

import (
	"fmt"
	"os"
	"path/filepath"
)

type EmailCredentials struct {
	Password string
	User     string
	Port     string
	Imap     string
	Smtp     string
}

type MadeinlySetting struct {
	Version          string
	Debug            bool
	EmailCredentials EmailCredentials
}

func InitSettings() {
	binPath, err := os.Executable()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	settingsPath := filepath.Join(binPath, "settings.toml")

	settingsFile, err := os.ReadFile(settingsPath)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(settingsFile))

}

func UpdateSettings() {

}
