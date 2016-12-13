// +build !js

package asset

import (
	"flag"
	"io/ioutil"
	"path/filepath"
)

// TODO: Find a better way to do this.Flags outside of main is bad practice.
// The problem is that loadFile *must* take a relative path to work with http GET requests. For the desktop version,
// relative path should be fine when distributing the game to end users, but it doesn't work for development.
// The standard practice in Golang is to use `go run`, but that generates an executable in a temp dir, with no easy way
// to get a relative path to the assets. It might be possible to do it dynamically with parsing os.Getenv("GOPATH")?
// If going that route, be careful because GOPATH can contain multiple paths, separated by semicolons.
// This is a reasonable solution, as long as it's removed and only relative paths are used when development is done.
var baseDir = flag.String("base_dir", `C:\workspace\Go\src\github.com\omustardo\gome\demos`, "All file paths should be specified relative to this root.")

func loadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(*baseDir, path))
}
