// view sets up and handles a glfw window.
package view

import (
	"fmt"

	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
)

// Window is the singleton glfw window. It should be initialized with view.Initialize(), and then
// all window related logic can access it directly.
var Window *glfw.Window

// Initialize sets up the singleton view.Window.
// Sample Usage:
//   if err := view.Initialize(*windowWidth, *windowHeight, "Window Demo"); err != nil {
//     log.Fatal(err)
//   }
//   defer view.Terminate()
//
func Initialize(width, height int, windowName string) error {
	if Window != nil {
		panic("view.Window already initialized")
	}
	fmt.Println("Creating Window...")

	err := glfw.Init(gl.ContextWatcher)
	if err != nil {
		return fmt.Errorf("unable to initialize glfw: %v", err)
	}
	glfw.WindowHint(glfw.Samples, 16) // Anti-aliasing.

	// Window hints to require OpenGL 3.2 or above, and to disable deprecated functions. https://open.gl/context#GLFW
	// These hints are not supported since we're using goxjs/glfw rather than the regular glfw, but should be used in a
	// standard desktop glfw project. TODO: Add support for these in goxjs/glfw/hint_glfw.go or consider using a conditional build rule.
	//glfw.WindowHint(glfw.ContextVersionMajor, 3)
	//glfw.WindowHint(glfw.ContextVersionMinor, 2)
	//glfw.WindowHint(glfw.OpenGLProfile, glfw.OPENGL_CORE_PROFILE)
	//glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)

	// Note CreateWindow ignores input size for WebGL/HTML canvas - it expands to fill browser window. This still matters for desktop.
	window, err := glfw.CreateWindow(width, height, windowName, nil, nil)
	if err != nil {
		return fmt.Errorf("unable to create glfw window: %v", err)
	}

	window.MakeContextCurrent()
	fmt.Printf("OpenGL: %s %s %s; %v samples.\n", gl.GetString(gl.VENDOR), gl.GetString(gl.RENDERER), gl.GetString(gl.VERSION), gl.GetInteger(gl.SAMPLES))
	fmt.Printf("GLSL: %s.\n", gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	glfw.SwapInterval(1) // Vsync.

	gl.ClearColor(0, 0, 0, 1) // Background Color
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// If a triangle is "facing away" from the camera, then don't draw it. https://www.opengl.org/wiki/Face_Culling
	// NOTE: If triangles appear to be missing, this is probably the cause. The order that vertices are listed matters.
	gl.FrontFace(gl.CCW) // This should be the default, but set it to be safe. https://www.opengl.org/sdk/docs/man/html/glFrontFace.xhtml
	gl.Enable(gl.CULL_FACE)
	gl.CullFace(gl.BACK)

	// Accept fragment if it's closer to the camera than the former one. TODO: Consider LEQUAL so draw order takes precedence for equal depth.
	// This makes
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// Set up a callback for when the window is resized. Call it once to properly initialize.
	framebufferSizeCallback := func(w *glfw.Window, framebufferSizeX, framebufferSizeY int) {
		gl.Viewport(0, 0, framebufferSizeX, framebufferSizeY)
	}
	framebufferSizeX, framebufferSizeY := window.GetFramebufferSize()
	framebufferSizeCallback(window, framebufferSizeX, framebufferSizeY)
	window.SetFramebufferSizeCallback(framebufferSizeCallback)

	Window = window
	if err := gl.GetError(); err != 0 {
		return fmt.Errorf("gl error: %v", err)
	}
	return nil
}

func Terminate() {
	glfw.Terminate()
}
