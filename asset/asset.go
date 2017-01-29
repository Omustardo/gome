// Package asset provides tools to manage loading, use, and unloading of assets, such as images and audio.
package asset

import "fmt"

// TODO: Find a better way to do this.
// The problem is that the js version of loadFile must take a relative path to work with http GET requests.
// For the desktop version, relative path should be fine when distributing the game to end users, but it doesn't work
// for development. The standard practice in Golang is to use `go run`, but that generates an executable in a temp dir,
// with no easy way to get a relative path to the assets. It might be possible to do it dynamically with parsing
// os.Getenv("GOPATH")? If going that route, be careful because GOPATH can contain multiple paths, separated by semicolons.
// Explicitly setting baseDir is a reasonable solution, as long as it's removed and only relative paths are used when development is done.
var baseDir string // = flag.String("base_dir", `C:\workspace\Go\src\github.com\omustardo\gome\demos`, "All file paths should be specified relative to this root.")

func Initialize(baseDirectory string) {
	fmt.Println("Setting up asset path...")
	baseDir = baseDirectory
}
