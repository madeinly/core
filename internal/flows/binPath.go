package flows

import (
	"os"
	"path/filepath"

	"github.com/madeinly/core/internal/features/fatal"
)

func BinPath() string {

	binPath, err := os.Executable()

	if err != nil {
		fatal.FatalError(err, "the binary path could not be retrieved")
	}

	return filepath.Dir(binPath)
}
