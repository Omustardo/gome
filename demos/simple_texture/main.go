package main

import (
	"flag"
	"image/color"
	"log"
	"math"
	"os"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"github.com/omustardo/gome"
	"github.com/omustardo/gome/asset"
	"github.com/omustardo/gome/camera"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/input/keyboard"
	"github.com/omustardo/gome/input/mouse"
	"github.com/omustardo/gome/model"
	"github.com/omustardo/gome/model/mesh"
	"github.com/omustardo/gome/shader"
	"github.com/omustardo/gome/util/fps"
	"github.com/omustardo/gome/view"
)

var (
	windowWidth  = flag.Int("window_width", 1000, "initial window width")
	windowHeight = flag.Int("window_height", 1000, "initial window height")

	// Explicitly listing the base dir is a hack. It's needed because `go run` produces a binary in a tmp folder so we can't
	// use relative asset paths. More explanation in omustardo\gome\asset\asset.go
	baseDir = flag.String("base_dir", `C:\workspace\Go\src\github.com\omustardo\gome\demos\simple_texture`, "All file paths should be specified relative to this root.")
)

func init() {
	// log print with .go file and line number.
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()
	terminate := gome.Initialize("Texture Demo", *windowWidth, *windowHeight, *baseDir)
	defer terminate()

	shader.Model.SetAmbientLight(&color.NRGBA{255, 255, 255, 0})

	tex, err := asset.LoadTexture("assets/ship.jpg")
	if err != nil {
		log.Fatalln(err)
	}

	target := &model.Model{
		Mesh: mesh.NewRect(&color.NRGBA{255, 255, 255, 255}, tex),
		Entity: entity.Entity{
			Position: mgl32.Vec3{},
			Scale:    mgl32.Vec3{300, 300, 0},
			Rotation: mgl32.QuatIdent(),
		},
	}

	cam := &camera.TargetCamera{
		Target:       target,
		TargetOffset: mgl32.Vec3{0, 0, 500},
		Up:           mgl32.Vec3{0, 1, 0},
		Near:         0.1,
		Far:          10000,
		FOV:          math.Pi / 4.0,
	}

	ticker := time.NewTicker(time.Second / 60)
	for !view.Window.ShouldClose() {
		fps.Handler.Update()
		glfw.PollEvents() // Reads window events, like keyboard and mouse input.
		// Handler.Update takes current input and stores it. This is necessary to detect things like the start of a keypress.
		keyboard.Handler.Update()
		mouse.Handler.Update()

		ApplyInputs(target, cam)

		cam.Update()

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
	target.ModifyPosition(move[0], move[1], 0)

	w, h := view.Window.GetSize()
	if mouse.Handler.LeftPressed() {
		move = cam.ScreenToWorldCoord2D(mouse.Handler.Position(), w, h).Sub(target.Position.Vec2())

		move = move.Normalize().Mul(moveSpeed * fps.Handler.DeltaTimeSeconds())
		target.ModifyPosition(move[0], move[1], 0)
	}
}
