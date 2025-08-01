package settings

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

type AppSettings struct {
	Version     string `toml:"version"`
	Debug       bool   `toml:"debug"`
	FrontDomain string `toml:"frontDomain"`
	Address     string `toml:"address"`
	Port        string `toml:"port"`
}

var (
	settings     AppSettings
	settingsLock sync.RWMutex
)

func GetSettings() AppSettings {
	settingsLock.RLock()
	defer settingsLock.RUnlock()
	return settings
}

func SetSettings() error {
	binPath, err := os.Executable()
	if err != nil {
		return err
	}

	settingsPath := filepath.Join(filepath.Dir(binPath), "settings.toml")
	settingsByte, err := os.ReadFile(settingsPath)
	if err != nil {
		return err
	}

	// Create new temporary struct
	var newSettings AppSettings

	// Decode into temporary struct
	if err := toml.Unmarshal(settingsByte, &newSettings); err != nil {
		return err
	}

	// Lock and update
	settingsLock.Lock()
	settings = newSettings
	settingsLock.Unlock()

	return nil
}

func WatchSettings() {
	binPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	settingsPath := filepath.Join(filepath.Dir(binPath), "settings.toml")

	lastModTime := time.Time{}
	lastSize := int64(0)

	for {
		fileInfo, err := os.Stat(settingsPath)
		if err != nil {
			log.Println("Error checking file:", err)
			time.Sleep(1 * time.Second)
			continue
		}

		modified := fileInfo.ModTime() != lastModTime || fileInfo.Size() != lastSize
		if modified {
			if err := SetSettings(); err != nil {
				log.Println("Failed to reload settings:", err)
			} else {
				lastModTime = fileInfo.ModTime()
				lastSize = fileInfo.Size()
			}
		}

		time.Sleep(1 * time.Second)
	}
}
