// +build !js

package asset

import (
	"io/ioutil"
	"path/filepath"
)

func loadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(baseDir, path))
}
