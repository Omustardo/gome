package main

import (
	"flag"
	"image/color"
	"math"
	"math/rand"
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
	count        = flag.Int("count", 100, "number of objects to draw")
	windowWidth  = flag.Int("window_width", 1000, "initial window width")
	windowHeight = flag.Int("window_height", 1000, "initial window height")
)

func main() {
	flag.Parse()
	terminate := gome.Initialize("Lots of Cubes", *windowWidth, *windowHeight, "")
	defer terminate()

	genCubes := func(count int) []*model.Model {
		// Try to evenly space the cubes in a grid centered at (0,0)
		scaleMin, scaleMax := float32(50), float32(150)
		countPerRow := int(math.Sqrt(float64(count)))
		var cubes []*model.Model
		for i := 0; i < count; i++ {
			col := &color.NRGBA{util.RandUint8(), util.RandUint8(), util.RandUint8(), 255}
			scale := rand.Float32()*(scaleMax-scaleMin) + scaleMin
			c := &model.Model{
				Mesh: mesh.NewCube(col, gl.Texture{}),
				Entity: entity.Entity{
					Position: mgl32.Vec3{float32(i%countPerRow)*scaleMax - scaleMax*float32(countPerRow)/2.0, float32(i/countPerRow)*scaleMax - scaleMax*float32(countPerRow)/2.0, 0},
					Scale:    mgl32.Vec3{scale, scale, scale},
					Rotation: util.RandQuat(),
				},
			}
			cubes = append(cubes, c)
		}
		return cubes
	}
	cubes := genCubes(*count)

	// target is what the camera is meant to look at and follow. It is not rendered.
	target := &model.Model{
		Mesh:   mesh.NewCube(&color.NRGBA{255, 25, 75, 255}, gl.Texture{}),
		Hidden: true,
	}

	cam := camera.NewRotateCamera(target, 1000)
	rotationPerSecond := mgl32.AnglesToQuat(float32(math.Pi/4)*0.8, float32(math.Pi/4), float32(math.Pi/4)*1.3, mgl32.XYZ)

	ticker := time.NewTicker(time.Second / 60)
	for !view.Window.ShouldClose() {
		fps.Handler.Update()
		glfw.PollEvents() // Reads window events, like keyboard and mouse input.
		// Handler.Update takes current input and stores it. This is necessary to detect things like the start of a keypress.
		keyboard.Handler.Update()
		mouse.Handler.Update()

		ApplyInputs(target)

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

		for _, c := range cubes {
			c.Render()
		}

		// Swaps the buffer that was drawn on to be visible. The visible buffer becomes the one that gets drawn on until it's swapped again.
		view.Window.SwapBuffers()
		<-ticker.C // wait up to 1/60th of a second. This caps framerate to 60 FPS.
	}
}

func ApplyInputs(target *model.Model) {
	var move mgl32.Vec2
	if keyboard.Handler.IsKeyDown(glfw.KeyA) {
		move[0] += -1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyD) {
		move[0] += 1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyW) {
		move[1] += 1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyS) {
		move[1] += -1
	}
	moveSpeed := float32(500)
	move = move.Normalize().Mul(moveSpeed * fps.Handler.DeltaTimeSeconds())
	target.ModifyPosition(move[0], move[1], 0)
}
