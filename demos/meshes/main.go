// Demo of loading and displaying various meshes.
package main

import (
	"flag"
	"log"
	"math"
	"os"
	"time"

	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"github.com/omustardo/gome/asset"
	"github.com/omustardo/gome/camera"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/input/keyboard"
	"github.com/omustardo/gome/input/mouse"
	"github.com/omustardo/gome/model"
	"github.com/omustardo/gome/model/mesh"
	"github.com/omustardo/gome/shader"
	"github.com/omustardo/gome/util/axis"
	"github.com/omustardo/gome/util/fps"
	"github.com/omustardo/gome/view"
)

var (
	windowWidth  = flag.Int("window_width", 1000, "initial window width")
	windowHeight = flag.Int("window_height", 1000, "initial window height")

	frameRate = flag.Duration("framerate", time.Second/60, `Cap on framerate. Provide with units, like "16.66ms"`)

	// Explicitly listing the base dir is a hack. It's needed because `go run` produces a binary in a tmp folder so we can't
	// use relative asset paths. More explanation in omustardo\gome\asset\asset.go
	baseDir = flag.String("base_dir", `C:\workspace\Go\src\github.com\omustardo\gome\demos\meshes`, "All file paths should be specified relative to this root.")
)

func init() {
	// log print with .go file and line number.
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()
	asset.Initialize(*baseDir)

	// Initialize gl constants and the glfw window. Note that this must be done before all other gl usage.
	if err := view.Initialize(*windowWidth, *windowHeight, "Model Viewer Demo"); err != nil {
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

	// Initialize singletons.
	mouse.Initialize(view.Window)
	keyboard.Initialize(view.Window)
	fps.Initialize()

	// Load standard meshes (cubes, rectangles, etc). These depend on OpenGL buffers, which depend on having an OpenGL
	// context. They must be called sometime after glfw is initialized to work.
	mesh.Initialize()
	axis.Initialize()

	// =========== Done with common initializations. From here on it's specific to this demo.

	shader.Model.SetAmbientLight(&color.NRGBA{60, 60, 60, 0}) // 3D objects don't look 3D in the default max lighting, so tone it down.

	// Load meshes.
	cubeMesh, err := asset.LoadOBJ("assets/cube.obj")
	cubeMesh.Color = &color.NRGBA{255, 25, 75, 255}
	if err != nil {
		log.Fatal(err)
	}
	vehicleMesh, err := asset.LoadDAE("assets/vehicle0.dae")
	if err != nil {
		log.Fatal(err)
	}

	capsuleMesh, err := asset.LoadOBJ("assets/capsule/capsule.obj")
	if err != nil {
		log.Fatal(err)
	}
	capsuleTexture, err := asset.LoadTexture("assets/capsule/capsule0.jpg")
	if err != nil {
		log.Fatal(err)
	}
	capsuleMesh.SetTexture(capsuleTexture)

	// Create models (meshes in world space)
	models := []*model.Model{
		// Cube
		{
			Tag:  "Built in Mesh", // Tag is *only* for human readable output/debugging.
			Mesh: mesh.NewCube(cubeMesh.Color, gl.Texture{}),
			Entity: entity.Entity{
				Scale: mgl32.Vec3{100, 100, 100},
			},
		},
		// Cube from file
		{
			Tag:  "OBJ Mesh",
			Mesh: cubeMesh,
			Entity: entity.Entity{
				Scale: mgl32.Vec3{100, 100, 100},
			},
		},
		// Rect
		{
			Tag:  "Built in Mesh",
			Mesh: mesh.NewRect(&color.NRGBA{80, 50, 100, 255}, gl.Texture{}),
			Entity: entity.Entity{
				Scale: mgl32.Vec3{100, 100, 100},
			},
		},
		// Rect outline
		{
			Mesh: mesh.NewRectOutline(&color.NRGBA{255, 25, 75, 255}),
			Entity: entity.Entity{
				Position: mgl32.Vec3{},
				Scale:    mgl32.Vec3{100, 100, 0},
				Rotation: mgl32.Vec3{},
			},
		},
		// Circle
		{
			Tag:  "Built in Mesh",
			Mesh: mesh.NewCircle(&color.NRGBA{200, 50, 100, 255}, gl.Texture{}),
			Entity: entity.Entity{
				Scale: mgl32.Vec3{100, 100, 100},
			},
		},
		// DAE mesh
		{
			Tag:  "DAE Mesh",
			Mesh: vehicleMesh,
			Entity: entity.Entity{
				Rotation: mgl32.Vec3{0, 0, 0},
				// Ideally the scale of all provided meshes fits them exactly into a unit cube, so scale is easy to work with.
				// In this case the vehicle model is already reasonably large, so don't scale it as much as other models.
				Scale: mgl32.Vec3{10, 10, 10},
			},
		},
		// Capsule
		{
			Tag:  "OBJ Textured Mesh",
			Mesh: capsuleMesh,
			Entity: entity.Entity{
				Rotation: mgl32.Vec3{0, 0, 0},
				Scale:    mgl32.Vec3{100, 100, 100},
			},
		},
	}
	// Adjust model positions so they're spaced nicely
	offset := float32(0)
	for _, m := range models {
		m.Position = mgl32.Vec3{offset, 0, 0}
		offset += 200
	}

	// Player is an empty model. It has no mesh so it can't be rendered, but it can still exist in the world.
	player := &model.Model{}
	cam := &camera.TargetCamera{
		Target:       player,
		TargetOffset: mgl32.Vec3{0, 0, 500},
		Up:           mgl32.Vec3{0, 1, 0},
		Near:         0.1,
		Far:          10000,
		FOV:          math.Pi / 2.0,
	}

	rotationPerSecond := float32(math.Pi / 4)

	ticker := time.NewTicker(*frameRate)
	for !view.Window.ShouldClose() {
		glfw.PollEvents() // Reads window events, like keyboard and mouse input.
		fps.Handler.Update()
		keyboard.Handler.Update()
		mouse.Handler.Update()

		ApplyInputs(player, cam)

		// Update the rotation.
		for _, m := range models {
			m.Rotation[0] += rotationPerSecond * fps.Handler.DeltaTimeSeconds()
			m.Rotation[1] += rotationPerSecond * fps.Handler.DeltaTimeSeconds()
			m.Rotation[2] += rotationPerSecond * fps.Handler.DeltaTimeSeconds()
		}

		// Set up Model-View-Projection Matrix and send it to the shader program.
		mvMatrix := cam.ModelView()
		w, h := view.Window.GetSize()
		pMatrix := cam.ProjectionPerspective(float32(w), float32(h))
		shader.Model.SetMVPMatrix(pMatrix, mvMatrix)

		// Clear screen, then Draw everything
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		axis.DrawXYZAxes()

		for _, m := range models {
			m.Render()
		}

		// Swaps the buffer that was drawn on to be visible. The visible buffer becomes the one that gets drawn on until it's swapped again.
		view.Window.SwapBuffers()
		<-ticker.C // wait up to the framerate cap.
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
