package view_test

import (
	"log"
	"github.com/omustardo/gome/view"
	"time"
	"github.com/goxjs/glfw"
	"github.com/goxjs/gl"
)

func Example() {
	if err := view.Initialize(1280, 720, "Demo"); err != nil {
		log.Fatal(err)
	}
	defer view.Terminate()

	framerate := time.Second/60
	ticker := time.NewTicker(framerate)
	for !view.Window.ShouldClose() {
		glfw.PollEvents() // Reads window events, like keyboard and mouse input. Necessary to do basic things, like resize the window.

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// TODO: Draw stuff here.

		view.Window.SwapBuffers()
		<-ticker.C // wait up to the framerate cap.
	}
}