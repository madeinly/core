package flows

import (
	"os"
	"path/filepath"

	"github.com/madeinly/core/internal/features/fatal"
	"github.com/madeinly/core/internal/features/random"
)

// TempFile creates a new temp file inside BinPath()/.tmp and returns it.
// The caller must close it.  On any failure the process terminates via
// fatal.FatalError.
func TempFile() *os.File {
	name, err := random.RandomString(15, random.Upper|random.Lower|random.Digits)
	if err != nil {
		fatal.FatalError(err, "could not generate random string")
	}

	tmpDir := filepath.Join(BinPath(), ".tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		fatal.FatalError(err, "could not create .tmp folder")
	}

	f, err := os.OpenFile(filepath.Join(tmpDir, name),
		os.O_CREATE|os.O_EXCL|os.O_RDWR, 0600)
	if err != nil {
		fatal.FatalError(err, "could not open temp file")
	}
	return f
}
