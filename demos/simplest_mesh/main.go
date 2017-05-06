package main

import (
	"flag"
	"image/color"
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"github.com/omustardo/gome"
	"github.com/omustardo/gome/camera"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/input/keyboard"
	"github.com/omustardo/gome/input/mouse"
	"github.com/omustardo/gome/model"
	"github.com/omustardo/gome/model/mesh"
	"github.com/omustardo/gome/shader"
	"github.com/omustardo/gome/util"
	"github.com/omustardo/gome/util/fps"
	"github.com/omustardo/gome/view"
)

var (
	windowWidth  = flag.Int("window_width", 1000, "initial window width")
	windowHeight = flag.Int("window_height", 1000, "initial window height")
)

func main() {
	flag.Parse()
	terminate := gome.Initialize("Generated Cube", *windowWidth, *windowHeight, "")
	defer terminate()

	target := &model.Model{
		Mesh: mesh.NewCube(&color.NRGBA{255, 25, 75, 255}, gl.Texture{}),
		Entity: entity.Entity{
			Position: mgl32.Vec3{},
			Scale:    mgl32.Vec3{100, 100, 100},
			Rotation: mgl32.QuatIdent(),
		},
	}

	cam := camera.NewOrbitCamera(target, 500)

	rotationPerSecond := mgl32.AnglesToQuat(float32(math.Pi/4)*0.8, float32(math.Pi/4), float32(math.Pi/4)*1.3, mgl32.XYZ)

	ticker := time.NewTicker(time.Second / 60)
	for !view.Window.ShouldClose() {
		fps.Handler.Update()
		glfw.PollEvents() // Reads window events, like keyboard and mouse input.
		// Handler.Update takes current input and stores it. This is necessary to detect things like the start of a keypress.
		keyboard.Handler.Update()
		mouse.Handler.Update()

		target.ModifyRotationLocalQ(util.ScaleQuatRotation(rotationPerSecond, fps.Handler.DeltaTimeSeconds()))

		cam.Update(fps.Handler.DeltaTime())

		// Set up Model-View-Projection Matrix and send it to the shader programs.
		mvMatrix := cam.ModelView()
		w, h := view.Window.GetSize()
		pMatrix := cam.ProjectionPerspective(float32(w), float32(h))
		shader.Model.SetMVPMatrix(pMatrix, mvMatrix)

		// Clear screen, then Draw everything
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		model.RenderXYZAxes()

		target.Render()

		// Swaps the buffer that was drawn on to be visible. The visible buffer becomes the one that gets drawn on until it's swapped again.
		view.Window.SwapBuffers()
		<-ticker.C // wait up to 1/60th of a second. This caps framerate to 60 FPS.
	}
}
