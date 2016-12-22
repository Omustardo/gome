package main

import (
	"flag"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"github.com/omustardo/gome/camera"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/input/keyboard"
	"github.com/omustardo/gome/input/mouse"
	"github.com/omustardo/gome/model"
	"github.com/omustardo/gome/model/mesh"
	"github.com/omustardo/gome/shader"
	"github.com/omustardo/gome/util"
	"github.com/omustardo/gome/util/axis"
	"github.com/omustardo/gome/util/fps"
	"github.com/omustardo/gome/view"
)

var (
	count        = flag.Int("count", 100, "number of objects to draw")
	windowWidth  = flag.Int("window_width", 1000, "initial window width")
	windowHeight = flag.Int("window_height", 1000, "initial window height")
)

func init() {
	// log print with .go file and line number.
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()

	// Initialize gl constants and the glfw window. Note that this must be done before all other gl usage.
	if err := view.Initialize(*windowWidth, *windowHeight, "Graphics Demo"); err != nil {
		log.Fatal(err)
	}
	defer view.Terminate()

	// Initialize Shaders
	if err := shader.Initialize(); err != nil {
		log.Fatal(err)
	}
	if err := gl.GetError(); err != 0 {
		log.Fatalf("gl error: %v", err)
	}
	shader.Model.SetAmbientLight(&color.NRGBA{60, 60, 60, 0}) // 3D objects don't look 3D in max lighting, so tone it down.

	// Initialize singletons.
	mouse.Initialize(view.Window)
	keyboard.Initialize(view.Window)
	fps.Initialize()

	// Load standard meshes (cubes, rectangles, etc). These depend on OpenGL buffers, which depend on having an OpenGL
	// context. They must be called sometime after glfw is initialized to work.
	mesh.Initialize()
	axis.Initialize()

	// =========== Done with common initializations. From here on it's specific to this demo.

	genCubes := func(count int) []*model.Model {
		// Try to evenly space the cubes in a grid centered at (0,0)
		scaleMin, scaleMax := float32(50), float32(150)
		countPerRow := int(math.Sqrt(float64(count)))
		var cubes []*model.Model
		for i := 0; i < count; i++ {
			col := &color.NRGBA{util.RandUint8(), util.RandUint8(), util.RandUint8(), 255}
			scale := rand.Float32()*(scaleMax-scaleMin) + scaleMin
			rot := rand.Float32() * 2 * math.Pi
			c := &model.Model{
				Mesh: mesh.NewCube(col, gl.Texture{}),
				Entity: entity.Entity{
					Position: mgl32.Vec3{float32(i%countPerRow)*scaleMax - scaleMax*float32(countPerRow)/2.0, float32(i/countPerRow)*scaleMax - scaleMax*float32(countPerRow)/2.0, 0},
					Scale:    mgl32.Vec3{scale, scale, scale},
					Rotation: mgl32.Vec3{rot, rot, rot},
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

	cam := &camera.TargetCamera{
		Target:       target,
		TargetOffset: mgl32.Vec3{0, 0, 500},
		Up:           mgl32.Vec3{0, 1, 0},
		Near:         0.1,
		Far:          10000,
		FOV:          math.Pi / 2.0,
	}

	rotationPerSecond := float32(math.Pi / 4)

	ticker := time.NewTicker(time.Second / 60)
	for !view.Window.ShouldClose() {
		fps.Handler.Update()
		glfw.PollEvents() // Reads window events, like keyboard and mouse input.
		// Handler.Update takes current input and stores it. This is necessary to detect things like the start of a keypress.
		keyboard.Handler.Update()
		mouse.Handler.Update()

		ApplyInputs(target, cam)

		rotate := func(m *model.Model) {
			m.Rotation[0] += rotationPerSecond * fps.Handler.DeltaTimeSeconds() * 0.8
			m.Rotation[1] += rotationPerSecond * fps.Handler.DeltaTimeSeconds()
			m.Rotation[2] += rotationPerSecond * fps.Handler.DeltaTimeSeconds() * 1.3
		}
		rotate(target)

		cam.Update()

		// Set up Model-View-Projection Matrix and send it to the shader programs.
		mvMatrix := cam.ModelView()
		w, h := view.Window.GetSize()
		pMatrix := cam.ProjectionPerspective(float32(w), float32(h))
		shader.Model.SetMVPMatrix(pMatrix, mvMatrix)

		// Clear screen, then Draw everything
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		axis.DrawXYZAxes()

		for _, c := range cubes {
			c.Render()
		}

		// Swaps the buffer that was drawn on to be visible. The visible buffer becomes the one that gets drawn on until it's swapped again.
		view.Window.SwapBuffers()
		<-ticker.C // wait up to 1/60th of a second. This caps framerate to 60 FPS.
	}
}

func ApplyInputs(target *model.Model, cam camera.Camera) {
	var move mgl32.Vec2
	if keyboard.Handler.IsKeyDown(glfw.KeyA, glfw.KeyLeft) {
		move[0] += -1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyD, glfw.KeyRight) {
		move[0] += 1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyW, glfw.KeyUp) {
		move[1] += 1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyS, glfw.KeyDown) {
		move[1] += -1
	}
	moveSpeed := float32(500)
	move = move.Normalize().Mul(moveSpeed * fps.Handler.DeltaTimeSeconds())
	target.ModifyCenter(move[0], move[1], 0)

	w, h := view.Window.GetSize()
	if mouse.Handler.LeftPressed() {
		move = cam.ScreenToWorldCoord2D(mouse.Handler.Position(), w, h).Sub(target.Center().Vec2())

		move = move.Normalize().Mul(moveSpeed * fps.Handler.DeltaTimeSeconds())
		target.ModifyCenter(move[0], move[1], 0)
	}
}