// +build !js

package asset

import (
	"io/ioutil"
	"path/filepath"
)

// TODO: Make this a flag, or find a better way to do it. Maybe by using os.Getenv("GOPATH")? If going that route, be careful because GOPATH can contain multiple paths, separated by semicolons.
// The issue is that `go run` builds to a temporary directory, so we can't use relative file paths like "assets\sample.png"
// `go build` and then running the exe works, but is slow and much less convenient.
const assetDir = `C:\workspace\Go\src\github.com\omustardo\gome\demos`

func loadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(assetDir, path))
}
