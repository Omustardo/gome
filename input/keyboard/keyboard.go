// keyboard handles keyboard interaction with a glfw window.
// Sample usage:
//   TODO
package keyboard

import (
	"fmt"
	"sort"
	"strings"

	"log"

	"github.com/goxjs/glfw"
)

// Handler is the singleton keyboard handler. It should be initialized with keyboard.Initialize(), and then
// all keyboard related input should be obtained though it.
var Handler *handler

// Initialize sets up the keyboard.Handler singleton.
func Initialize(window *glfw.Window) {
	if Handler != nil {
		panic("keyboard.Handler already initialized")
	}
	if window == nil {
		panic("window is nil")
	}
	Handler = &handler{
		state:         make(map[glfw.Key]bool),
		previousState: make(map[glfw.Key]bool),
		keyEventList:  newGLFWKeyEventList(),
	}
	window.SetKeyCallback(Handler.keyEventList.Callback)
}

type glfwKeyEvent struct {
	key      glfw.Key
	scancode int
	action   glfw.Action
	mods     glfw.ModifierKey
}

const eventListCap = 10 // max number of key events between a single call to keyboard.Handler.Update()

type glfwKeyEventList []glfwKeyEvent

func newGLFWKeyEventList() *glfwKeyEventList {
	eventList := glfwKeyEventList(make([]glfwKeyEvent, 0, eventListCap))
	return &eventList
}

// freeze returns the list of key events since it was last called, and limited
// to 'eventListCap' events. It then clears the internal buffer.
func (keyEventList *glfwKeyEventList) freeze() []glfwKeyEvent {
	// The list of key events is double buffered.  This allows the application
	// to process events during a frame without having to worry about new
	// events arriving and growing the list.
	// TODO: Use two buffers rather than assigning and making a new one each time.
	frozen := *keyEventList

	// Note that the warning in Intellij about "assigning to method receiver only propagating to callees and not callers" is wrong.
	// It doesn't realize that we're dereferencing a pointer. See https://play.golang.org/p/VlvvRIrx7v
	*keyEventList = make([]glfwKeyEvent, 0, eventListCap)

	return frozen
}

// Callback is intended to be passed it into glfw.Window's SetKeyCallback method which uses it as an event handler for
// key events. It can also be called directly to simulate key events.
func (keyEventList *glfwKeyEventList) Callback(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	event := glfwKeyEvent{key, scancode, action, mods}
	*keyEventList = append(*keyEventList, event)
}

// handler is the singleton member of the keyboard package. Create it using keyboard.Initialize()
type handler struct {
	// state maps from keys to whether they are pressed.
	state         map[glfw.Key]bool
	previousState map[glfw.Key]bool

	// keyEventList is modified by the glfw window callback, or directly for testing, to keep track of keypresses
	// between calls to handler.Update()
	keyEventList *glfwKeyEventList
}

// process the most recent key events and use them to modify the internal
// handler's view of the keyboard state.
func (h *handler) process(events []glfwKeyEvent) {
	for _, event := range events {
		h.setState(event.key, event.action)
	}
}

func (h *handler) setState(key glfw.Key, action glfw.Action) {
	switch action {
	case glfw.Press:
		h.state[key] = true
		// log.Println("Key:", key, "pressed")
	case glfw.Release:
		h.state[key] = false
		// log.Println("Key:", key, "released")
	}
}

// Update is expected to be called once per frame, or more. It handles any key events since it was last called.
func (h *handler) Update() {
	h.previousState = h.state
	h.state = make(map[glfw.Key]bool)
	for k, pressed := range h.previousState {
		if pressed {
			h.state[k] = true
		}
	}

	// Get a snapshot of key events so incoming ones don't affect the processing.
	// Note that this clears h.keyEventList so it's ready for new events.
	keyEvents := h.keyEventList.freeze()
	h.process(keyEvents)
}

// ====== Helper functions ======

// IsKeyDown returns whether any of the provided keys are currently pressed.
// Example usage:
//  if keyboard.Handler.IsKeyDown(glfw.KeyA, glfw.KeyLeft) {
//     // Move player left
//  }
//  if keyboard.Handler.IsKeyDown(glfw.KeySpace) {
//    // Fire weapon
//  }
func (h *handler) IsKeyDown(keys ...glfw.Key) bool {
	return isKeyDown(h.state, keys...)
}
func (h *handler) WasKeyDown(keys ...glfw.Key) bool {
	return isKeyDown(h.previousState, keys...)
}
func (h *handler) JustPressed(key glfw.Key) bool {
	return h.IsKeyDown(key) && !h.WasKeyDown(key)
}

func isKeyDown(state map[glfw.Key]bool, keys ...glfw.Key) bool {
	if state == nil {
		log.Println("nil keyboard state detected")
		return false
	}
	for _, k := range keys {
		if state[k] {
			return true
		}
	}
	return false
}

func (h *handler) LeftPressed() bool {
	return h.IsKeyDown(glfw.KeyLeft)
}
func (h *handler) RightPressed() bool {
	return h.IsKeyDown(glfw.KeyRight)
}
func (h *handler) UpPressed() bool {
	return h.IsKeyDown(glfw.KeyUp)
}
func (h *handler) DownPressed() bool {
	return h.IsKeyDown(glfw.KeyDown)
}
func (h *handler) SpacePressed() bool {
	return h.IsKeyDown(glfw.KeySpace)
}

func (h *handler) WasLeftPressed() bool {
	return h.WasKeyDown(glfw.KeyLeft)
}
func (h *handler) WasRightPressed() bool {
	return h.WasKeyDown(glfw.KeyRight)
}
func (h *handler) WasUpPressed() bool {
	return h.WasKeyDown(glfw.KeyUp)
}
func (h *handler) WasDownPressed() bool {
	return h.WasKeyDown(glfw.KeyDown)
}
func (h *handler) WasSpacePressed() bool {
	return h.WasKeyDown(glfw.KeySpace)
}

// String prints out all of the currently pressed keys in human readable format.
// TODO: Currently casts the keycode to a character. This works for standard
// letters and numbers, but things like numpad numbers don't work, and obviously
// Shift, Delete, and other longer names couldn't possibly work. Improve this.
func (h *handler) String() string {
	var keys []string
	for key, pressed := range h.state {
		if pressed {
			keys = append(keys, fmt.Sprintf("'%c'", key))
		}
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}
