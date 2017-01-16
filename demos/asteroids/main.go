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
	"github.com/omustardo/gome/asset"
	"github.com/omustardo/gome/camera"
	"github.com/omustardo/gome/camera/zoom"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/demos/asteroids/asteroid"
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
	windowWidth  = flag.Int("window_width", 1000, "initial window width")
	windowHeight = flag.Int("window_height", 1000, "initial window height")

	// Explicitly listing the base dir is a hack. It's needed because `go run` produces a binary in a tmp folder so we can't
	// use relative asset paths. More explanation in omustardo\gome\asset\asset.go
	baseDir = flag.String("base_dir", `C:\workspace\Go\src\github.com\omustardo\gome\demos\asteroids`, "All file paths should be specified relative to this root.")
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

	shipMesh, err := asset.LoadOBJ("assets/ship/ship.obj")
	if err != nil {
		log.Fatalf("Unable to load ship model: %v", err)
	}
	shipTexture, err := asset.LoadTexture("assets/ship/ship.jpg")
	if err != nil {
		log.Fatalf("Unable to load asteroid texture: %v", err)
	}
	shipMesh.SetTexture(shipTexture)
	playerShip := &model.Model{
		Mesh: shipMesh,
		Entity: entity.Entity{
			Position: mgl32.Vec3{0, -200, 0},
			// Rotate the model so it starts facing directly toward the positive Y axis, which is up on the user's screen.
			// Remember that these rotations are applied in the order specified by the final parameter
			// and that like a unit circle, positive values go to the "left" and negative values go to the "right".
			// X,Y,Z correspond to Roll, Pitch, and Yaw.
			Rotation: mgl32.AnglesToQuat(mgl32.DegToRad(90), mgl32.DegToRad(-90), 0, mgl32.XYZ),
			Scale:    mgl32.Vec3{4, 4, 4},
		},
	}
	cam := &camera.TargetCamera{
		Target:       playerShip,
		TargetOffset: mgl32.Vec3{0, 0, 500},
		Up:           mgl32.Vec3{0, 1, 0},
		Zoomer: zoom.NewScrollZoom(0.25, 3,
			func() float32 {
				return mouse.Handler.Scroll().Y()
			},
		),
		Near: 0.1,
		Far:  10000,
		FOV:  math.Pi / 2.0,
	}

	asteroidMesh, err := asset.LoadOBJ("assets/rock/rock1.obj")
	if err != nil {
		log.Fatalf("Unable to load asteroid model: %v", err)
	}
	asteroidTexture, err := asset.LoadTexture("assets/rock/rock1.jpg")
	if err != nil {
		log.Fatalf("Unable to load asteroid texture: %v", err)
	}
	asteroidMesh.SetTexture(asteroidTexture)
	asteroid.SetMesh(asteroidMesh)

	asteroids := []*asteroid.Asteroid{asteroid.New()}

	ticker := time.NewTicker(time.Second / 60)
	for !view.Window.ShouldClose() {
		fps.Handler.Update()
		glfw.PollEvents() // Reads window events, like keyboard and mouse input.
		// Handler.Update takes current input and stores it. This is necessary to detect things like the start of a keypress.
		keyboard.Handler.Update()
		mouse.Handler.Update()

		ApplyInputs(playerShip, cam)

		for _, a := range asteroids {
			a.Update(fps.Handler.DeltaTimeSeconds())
		}

		cam.Update()

		// Set up Model-View-Projection Matrix and send it to the shader programs.
		mvMatrix := cam.ModelView()
		w, h := view.Window.GetSize()
		pMatrix := cam.ProjectionPerspective(float32(w), float32(h))
		shader.Model.SetMVPMatrix(pMatrix, mvMatrix)

		// Clear screen, then Draw everything
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		axis.DrawXYZAxes()

		for _, a := range asteroids {
			a.Render()
		}
		playerShip.Render()

		// Swaps the buffer that was drawn on to be visible. The visible buffer becomes the one that gets drawn on until it's swapped again.
		view.Window.SwapBuffers()
		<-ticker.C // wait up to 1/60th of a second. This caps framerate to 60 FPS.
	}
}

func ApplyInputs(target *model.Model, cam camera.Camera) {
	// rotation speed in radians per second.
	rotationSpeed := mgl32.AnglesToQuat(0, 0, mgl32.DegToRad(360/3), mgl32.XYZ)

	// rotate is the direction and amount to rotate.
	// Note that, like a unit circle, radians of higher positive value are toward the "left", while negative are to the "right".
	var rotationScale float32
	if keyboard.Handler.IsKeyDown(glfw.KeyA, glfw.KeyLeft) {
		rotationScale += fps.Handler.DeltaTimeSeconds()
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyD, glfw.KeyRight) {
		rotationScale -= fps.Handler.DeltaTimeSeconds()
	}
	if rotationScale != 0 {
		target.ModifyRotationGlobalQ(util.ScaleQuatRotation(rotationSpeed, rotationScale))
	}

	var move float32
	if keyboard.Handler.IsKeyDown(glfw.KeyW, glfw.KeyUp) {
		move -= fps.Handler.DeltaTimeSeconds()
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyS, glfw.KeyDown) {
		move += fps.Handler.DeltaTimeSeconds()
	}
	moveSpeed := float32(500)
	_, _, heading := target.RotationAngles().Elem() // direction that the target is facing in radians. Ignore roll and pitch since we're constrained to one axis of rotation - on the XY plane.

	forward := mgl32.Vec3{float32(math.Cos(float64(heading))), float32(math.Sin(float64(heading))), 0}
	forward = forward.Normalize().Mul(move * moveSpeed) // direction * speed = distance
	target.ModifyCenterV(forward)                       // current position + distance vector = final location

	w, h := view.Window.GetSize()
	if mouse.Handler.LeftPressed() {
		dir := cam.ScreenToWorldCoord2D(mouse.Handler.Position(), w, h).Sub(target.Center().Vec2())
		dir = dir.Normalize().Mul(moveSpeed * fps.Handler.DeltaTimeSeconds())
		target.ModifyCenter(dir[0], dir[1], 0)
	}
}
