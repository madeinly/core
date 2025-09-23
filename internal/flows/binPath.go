package flows

import (
	"github.com/madeinly/core/internal/features/fatal"
	"github.com/madeinly/core/internal/features/files"
)

func BinPath() string {

	binPath, err := files.BinPath()

	if err != nil {
		fatal.FatalError(err, "the binary path could not be retrieved")
	}

	return binPath
}
