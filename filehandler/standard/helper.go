package standard

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func returnBytes(fpath string) ([]byte, error) {
	contents, err := ioutil.ReadFile(filepath.Clean(fpath))
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func isAcceptedImageFile(name string) bool {
	validFiles := []string{".png", ".jpg", ".jpeg", ".gif"}

	for _, ext := range validFiles {
		if strings.HasSuffix(strings.ToLower(name), ext) {
			return true
		}
	}

	return false
}

func isGoFile(name string) bool {
	validFiles := []string{".go"}

	for _, ext := range validFiles {
		if strings.HasSuffix(strings.ToLower(name), ext) {
			return true
		}
	}

	return false
}
