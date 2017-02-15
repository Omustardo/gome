package gome

import (
	"image/color"
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/asset"
	"github.com/omustardo/gome/input/keyboard"
	"github.com/omustardo/gome/input/mouse"
	"github.com/omustardo/gome/model/mesh"
	"github.com/omustardo/gome/shader"
	"github.com/omustardo/gome/util/fps"
	"github.com/omustardo/gome/view"
)

// Init should be called near the start of the main method.
// It sets up the glfw window, shaders, input, and basic meshes among other things.
func Initialize(windowTitle string, windowWidth, windowHeight int, baseDir string) func() {
	asset.SetBaseDir(baseDir)

	// Initialize gl constants and the glfw window. Note that this must be done before all other gl usage.
	if err := view.Initialize(windowWidth, windowHeight, windowTitle); err != nil {
		log.Fatal(err)
	}

	// Initialize Shaders
	if err := shader.Initialize(); err != nil {
		log.Fatal(err)
	}
	if err := gl.GetError(); err != 0 {
		log.Fatalf("gl error: %v", err)
	}
	shader.Model.SetAmbientLight(&color.NRGBA{60, 60, 60, 0}) // 3D objects don't look 3D in max lighting, so tone it down as a default.
	shader.Model.SetDiffuse(mgl32.Vec3{0, 0, 1}, &color.NRGBA{150, 150, 150, 255})

	// Initialize singletons.
	mouse.Initialize(view.Window)
	keyboard.Initialize(view.Window)
	fps.Initialize()

	// Load standard meshes (cubes, rectangles, etc). These depend on OpenGL buffers, which depend on having an OpenGL
	// context. They must be called sometime after glfw is initialized to work.
	mesh.Initialize()

	return view.Terminate
}
