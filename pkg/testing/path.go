package testing

import (
	"errors"
	"os"
	"strings"
)

// GetFullPath returns full path for a given file regardless whats the current working directory
func GetFullPath(fromRootPath string) (string, error) {
	const rootF = "elastic-search"
	const maxSteps = 10

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	steps := 0
	currentRootPath := wd

	for steps < maxSteps {
		strs := strings.Split(currentRootPath, "/")

		if len(strs) < 2 {
			return "", errors.New("unable to locate the root path")
		}

		if strs[len(strs)-1] == rootF {
			return currentRootPath + fromRootPath, nil
		}

		currentRootPath = strings.Join(strs[:len(strs)-2], "/")
		steps++
	}

	return "", errors.New("unable to locate the root path")
}
