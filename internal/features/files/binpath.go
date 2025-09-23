package files

import (
	"os"
	"path/filepath"
)

func BinPath() (string, error) {

	binPath, err := os.Executable()

	return filepath.Dir(binPath), err
}
