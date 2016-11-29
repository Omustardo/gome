// +build js

package util

import "log"

// SaveScreenshot reads pixel data from OpenGL buffers, so it must be run in the same main thread as the rest
// of OpenGL.
// TODO: write to file in a goroutine and return a (chan err), or just ignore slow errors. Handling errors that can be caught immediately is fine. Blocking while writing to file adds way too much delay.
func SaveScreenshot(_, _ int, _ string) error {
	log.Println("screenshot function unimplemented")
	return nil
}
