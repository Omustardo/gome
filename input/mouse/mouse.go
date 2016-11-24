// mouse handles mouse interaction with a glfw window.
// Sample usage:
//   TODO
package mouse

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/glfw"
)

// Handler is the singleton mouse handler. It should be initialized with mouse.Initialize(), and then
// all mouse related input should be obtained though it.
var Handler *handler

// Initialize sets up the mouse.Handler singleton.
func Initialize(window *glfw.Window) {
	if Handler != nil {
		panic("mouse.Handler already initialized")
	}
	if window == nil {
		panic("window is nil")
	}
	Handler = &handler{
		buttonsBuffer:   make(map[glfw.MouseButton]bool),
		buttons:         make(map[glfw.MouseButton]bool),
		previousButtons: make(map[glfw.MouseButton]bool),
	}
	window.SetMouseButtonCallback(Handler.mouseButtonCallback)
	window.SetCursorPosCallback(Handler.cursorPosCallback)
	window.SetScrollCallback(Handler.scrollCallback)
}

// handler is the singleton member of the keyboard package. Create it using keyboard.Initialize()
type handler struct {
	// Each variable has three versions to hold three update calls worth of data.
	// For example, the buttons buffers hold mouse button presses.
	// stateBuffer is what's changed by the glfw.Window's callback whenever a key is actually pressed. Nothing should read or modify this except the window callback.
	// When Update is called, the buffer moves to the current state. All game logic reads from this in order to prevent the window callbacks from changing data out from under them.
	// When Update is called, the current state is moved to the previous state. We need to keep track of the previous data to do comparisons, like to see when a button was just pressed.

	// State maps from buttons to whether they are pressed.
	buttonsBuffer, buttons, previousButtons map[glfw.MouseButton]bool

	// position is the screen coordinate where the mouse pointer is.
	positionBuffer, position, previousPosition mgl32.Vec2

	// Scroll holds how much scrolling has occurred since the start of the program.
	// PreviousScroll is how much scrolling occurred since the start of the program, ignoring anything more recent than the last call to mouse.Update()
	// To determine changes, subtract the two.
	// The Y value is the standard forward/back, while the left/right scrolling available on some mice is in the X value.
	// While glfw says the value is a float, I've only seen it as integers. One "tick" is +1 or -1 depending on direction.
	// Positive for Forward/Left. Negative for Back/Right. 0 by default.
	scrollBuffer, scroll, previousScroll mgl32.Vec2
}

// mouseButtonCallback is a function for glfw to call when a button event occurs.
func (h *handler) mouseButtonCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
	// Note that this will overwrite any unhandled actions.
	// For example, if you press the left mouse button and then release it without calling
	// handler.Update() in between, it will appear as if no action was taken.
	// I think this is fine, since you shouldn't be able to click and release within 1/60th of a second.
	// If there are noticeable missing key presses, then this is almost certainly the problem.
	setState(h.buttonsBuffer, button, action)
}

// cursorPosCallback is a function for glfw to call when a button event occurs.
func (h *handler) cursorPosCallback(_ *glfw.Window, xpos, ypos float64) {
	// log.Println("got cursor pos event:", xpos, ypos)
	h.positionBuffer[0] = float32(xpos)
	h.positionBuffer[1] = float32(ypos)
}

// scrollCallback is a function for glfw to call when a scroll wheel event occurs.
func (h *handler) scrollCallback(_ *glfw.Window, xoff, yoff float64) {
	// log.Println("got scroll event:", xoff, yoff)
	h.scrollBuffer[0] += float32(xoff)
	h.scrollBuffer[1] += float32(yoff)
}

func setState(state map[glfw.MouseButton]bool, button glfw.MouseButton, action glfw.Action) {
	if state == nil {
		log.Println("found nil mouse button state")
		return
	}
	switch action {
	case glfw.Press:
		state[button] = true
	case glfw.Release:
		state[button] = false
	}
}

// Update is expected to be called once per frame, or more. It handles any mouse events since it was last called.
func (h *handler) Update() {
	h.previousButtons = h.buttons
	h.buttons = h.buttonsBuffer
	h.buttonsBuffer = make(map[glfw.MouseButton]bool)
	for button, pressed := range h.buttons {
		if pressed {
			h.buttonsBuffer[button] = true
		}
	}

	h.previousPosition = h.position
	h.position = h.positionBuffer

	h.previousScroll = h.scroll
	h.scroll = h.scrollBuffer
}

func (h *handler) LeftPressed() bool {
	return h.buttons[glfw.MouseButtonLeft]
}
func (h *handler) RightPressed() bool {
	return h.buttons[glfw.MouseButtonRight]
}
func (h *handler) WasLeftPressed() bool {
	return h.previousButtons[glfw.MouseButtonLeft]
}
func (h *handler) WasRightPressed() bool {
	return h.previousButtons[glfw.MouseButtonRight]
}

// Position returns the screen coordinate where the mouse pointer is.
// (0,0) is the top left of the drawable region (i.e. not including the title bar in a desktop environment).
// Down and right are positive. Up and left are negative.
func (h *handler) Position() mgl32.Vec2 {
	return h.position
}
func (h *handler) PreviousPosition() mgl32.Vec2 {
	return h.previousPosition
}

// Scroll returns the amount of scrolling done in the previous Update() call.
// The Y value is the standard forward/back, while the left/right scrolling available on some mice is in the X value.
// Positive for Forward/Left. Negative for Back/Right. 0 by default.
func (h *handler) Scroll() mgl32.Vec2 {
	return h.scroll.Sub(h.previousScroll)
}

// ScrollTotal returns the amount of scrolling done since the start of the mouse handler.
func (h *handler) ScrollTotal() mgl32.Vec2 {
	return h.scroll
}
