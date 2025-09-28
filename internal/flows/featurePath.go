package flows

import (
	"os"
	"path"

	"github.com/madeinly/core/internal/features/fatal"
)

// makes sure that the cards folder exist and if not created or return error
func FeaturePath(featurePath string) string {

	binPath := BinPath()

	cardsFolderPath := path.Join(binPath, featurePath)

	_, err := os.Stat(cardsFolderPath)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(cardsFolderPath, 0755)
			if err != nil {
				fatal.FatalError(err, "could not create a feature folder")
				return ""
			}
		} else {
			fatal.FatalError(err, "could not stat for feature folder")
			return ""
		}
	}

	return cardsFolderPath
}
