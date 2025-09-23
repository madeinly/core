package safetyControl

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/madeinly/core/internal/features/fatal"
	"golang.org/x/sys/unix"
)

// --- permission constants -----------------------------------------------
const (
	// root-level db.sqlite
	RootSQLitePerm = unix.R_OK | unix.W_OK // 0660 (rw-rw----)

	// every other *.sqlite / *.sqlite-* file
	OtherSQLitePerm = unix.R_OK // 0440 (r--r-----)

	// everything else (regular files & directories)
	DefaultDirPerm  = unix.R_OK | unix.W_OK | unix.X_OK
	DefaultFilePerm = unix.R_OK | unix.W_OK
)

// FilesIntegrity walks the directory that contains the executable
// and returns the paths of files/directories that the current user
// *cannot* read and write (and, for directories, traverse).
// The slice is empty when everything passes the checks.
func FilesIntegrity() ([]string, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	root := filepath.Dir(exe)

	var bad []string
	err = filepath.WalkDir(root, func(p string, d os.DirEntry, _ error) error {
		var mode uint32

		switch {
		case d.IsDir():
			mode = DefaultDirPerm
		case p == filepath.Join(root, "db.sqlite"):
			mode = RootSQLitePerm
		case isSQLiteFile(p):
			mode = OtherSQLitePerm
		default:
			mode = DefaultFilePerm
		}

		if unix.Access(p, mode) != nil {
			bad = append(bad, p)
			if d.IsDir() {
				return filepath.SkipDir
			}
		}
		return nil
	})

	return bad, err
}

func isSQLiteFile(name string) bool {
	base := filepath.Base(name)
	return strings.HasSuffix(base, ".sqlite") || strings.Contains(base, ".sqlite-")
}

func RootPath() string {
	binPath, err := os.Executable()
	if err != nil {
		fatal.FatalError(err, "Could not determine the application's root path")
	}
	return path.Dir(binPath)
}
