package settings

import (
	"context"
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
	JWTSalt     string `toml:"jwtSalt"`
	Email       struct {
		User       string `toml:"user"`
		Address    string `toml:"address"`
		Port       string `toml:"port"`
		Password   string `toml:"password"`
		Encryption string `toml:"encryption"`
	} `toml:"email"`
}

var (
	settings     AppSettings
	settingsLock sync.RWMutex
	initOnce     sync.Once
)

// GetSettings ensures that settings are loaded before returning them.
// It will panic if the initial settings load fails.
func GetSettings() AppSettings {
	initOnce.Do(func() {
		if err := setSettings(); err != nil {
			log.Fatalf("FATAL: Failed to load initial settings: %v", err)
		}
	})

	settingsLock.RLock()
	defer settingsLock.RUnlock()
	return settings
}

// GetRawSettings reads the settings.toml file and returns its content as a raw map.
// It also ensures the structured settings are loaded for the application.
func GetRawSettings() (map[string]interface{}, error) {
	initOnce.Do(func() {
		if err := setSettings(); err != nil {
			log.Fatalf("FATAL: Failed to load initial settings: %v", err)
		}
	})

	binPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	settingsPath := filepath.Join(filepath.Dir(binPath), "settings.toml")
	settingsByte, err := os.ReadFile(settingsPath)
	if err != nil {
		return nil, err
	}

	var rawSettings map[string]interface{}
	if err := toml.Unmarshal(settingsByte, &rawSettings); err != nil {
		return nil, err
	}

	return rawSettings, nil
}

// setSettings is the internal function that performs the loading of settings.
func setSettings() error {
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

// WatchSettings continuously monitors the settings.toml file for changes.
func WatchSettings(ctx context.Context) {
	binPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	settingsPath := filepath.Join(filepath.Dir(binPath), "settings.toml")

	lastModTime := time.Time{}
	lastSize := int64(0)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fileInfo, err := os.Stat(settingsPath)
			if err != nil {
				log.Println("Error checking file:", err)
				continue
			}

			modified := fileInfo.ModTime() != lastModTime || fileInfo.Size() != lastSize
			if modified {
				if err := setSettings(); err != nil {
					log.Println("Failed to reload settings:", err)
				} else {
					lastModTime = fileInfo.ModTime()
					lastSize = fileInfo.Size()
				}
			}
		}
	}
}
