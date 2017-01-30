// Demo of loading and displaying various meshes.
package main

import (
	"flag"
	"log"
	"math"
	"os"
	"time"

	"image/color"

	"fmt"

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
	"github.com/omustardo/gome/util"
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
	terminate := gome.Initialize("Mesh Showcase", *windowWidth, *windowHeight, *baseDir)
	defer terminate()

	shader.Model.SetAmbientLight(&color.NRGBA{60, 60, 60, 0}) // 3D objects don't look 3D in the default max lighting, so tone it down.

	// Load meshes.
	// OBJ
	cubeMesh, err := asset.LoadOBJ("assets/cube.obj", asset.OBJOpts{})
	cubeMesh.Color = &color.NRGBA{255, 25, 75, 255}
	if err != nil {
		log.Fatal(err)
	}
	capsuleMesh, err := asset.LoadOBJ("assets/capsule/capsule.obj", asset.OBJOpts{Normalize: true, Center: &mgl32.Vec3{0.5, 0.5, 0.5}})
	if err != nil {
		log.Fatal(err)
	}
	capsuleTexture, err := asset.LoadTexture("assets/capsule/capsule0.jpg")
	if err != nil {
		log.Fatal(err)
	}
	capsuleMesh.SetTexture(capsuleTexture)
	shipMesh, err := asset.LoadOBJ("assets/ship/ship.obj", asset.OBJOpts{Normalize: true, Center: &mgl32.Vec3{0.5, 0.5, 0.5}})
	if err != nil {
		log.Fatal(err)
	}
	shipTexture, err := asset.LoadTexture("assets/ship/ship.jpg")
	if err != nil {
		log.Fatal(err)
	}
	shipMesh.SetTexture(shipTexture)
	// DAE
	vehicleMesh, err := asset.LoadDAE("assets/vehicle/vehicle0.dae")
	if err != nil {
		log.Fatal(err)
	}
	duckMesh, err := asset.LoadDAE("assets/duck/duck.dae")
	if err != nil {
		log.Fatal(err)
	}
	duckTexture, err := asset.LoadTexture("assets/duck/duck.png")
	if err != nil {
		log.Fatal(err)
	}
	duckMesh.SetTexture(duckTexture)
	fmt.Println("Done loading assets")

	// Create models (meshes in world space)
	models := []*model.Model{
		// Sphere
		{
			Tag:  "Generated Sphere. Detail 0",
			Mesh: mesh.NewSphere(0, &color.NRGBA{80, 50, 100, 255}, gl.Texture{}),
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
				Scale:    mgl32.Vec3{100, 100, 100},
			},
		},
		{
			Tag:  "Generated Sphere. Detail 1",
			Mesh: mesh.NewSphere(1, &color.NRGBA{80, 50, 150, 255}, gl.Texture{}),
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
				Scale:    mgl32.Vec3{100, 100, 100},
			},
		},
		{
			Tag:  "Generated Sphere. Detail 2",
			Mesh: mesh.NewSphere(2, &color.NRGBA{80, 60, 200, 255}, gl.Texture{}),
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
				Scale:    mgl32.Vec3{100, 100, 100},
			},
		},
		{
			Tag:  "Generated Sphere. Detail 3",
			Mesh: mesh.NewSphere(3, &color.NRGBA{100, 80, 240, 255}, gl.Texture{}),
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
				Scale:    mgl32.Vec3{100, 100, 100},
			},
		},
		{
			Tag:  "Generated Sphere. Detail 4",
			Mesh: mesh.NewSphere(4, &color.NRGBA{120, 100, 255, 255}, gl.Texture{}),
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
				Scale:    mgl32.Vec3{100, 100, 100},
			},
		},
		{
			Tag:  "Generated Icosahedron",
			Mesh: mesh.NewIcosahedron(&color.NRGBA{80, 50, 100, 255}, gl.Texture{}),
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
				Scale:    mgl32.Vec3{100, 100, 100},
			},
		},
		// Cube
		{
			Tag:  "Built in Mesh", // Tag is *only* for human readable output/debugging.
			Mesh: mesh.NewCube(cubeMesh.Color, gl.Texture{}),
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
				Scale:    mgl32.Vec3{100, 100, 100},
			},
		},
		// Cube from file
		{
			Tag:  "OBJ Mesh",
			Mesh: cubeMesh,
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
				Scale:    mgl32.Vec3{100, 100, 100},
			},
		},
		// Rect
		{
			Tag:  "Built in Mesh",
			Mesh: mesh.NewRect(&color.NRGBA{80, 50, 100, 255}, gl.Texture{}),
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
				Scale:    mgl32.Vec3{100, 100, 100},
			},
		},
		// Rect outline
		{
			Mesh: mesh.NewRectOutline(&color.NRGBA{255, 25, 75, 255}),
			Entity: entity.Entity{
				Position: mgl32.Vec3{},
				Rotation: mgl32.QuatIdent(),
				Scale:    mgl32.Vec3{100, 100, 0},
			},
		},
		// Circle
		{
			Tag:  "Built in Mesh",
			Mesh: mesh.NewCircle(&color.NRGBA{200, 50, 100, 255}, gl.Texture{}),
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
				Scale:    mgl32.Vec3{100, 100, 100},
			},
		},
		// Textured OBJ mesh
		{
			Tag:  "OBJ Textured Mesh",
			Mesh: shipMesh,
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
				Scale:    mgl32.Vec3{100, 100, 100},
			},
		},

		// DAE mesh
		{
			Tag:  "DAE Mesh",
			Mesh: vehicleMesh,
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
				// Ideally the scale of all provided meshes fits them exactly into a unit cube, so scale is easy to work with.
				// In this case the vehicle model is already reasonably large, so don't scale it as much as other models.
				Scale: mgl32.Vec3{10, 10, 10},
			},
		},
		{
			Tag:  "DAE Mesh",
			Mesh: duckMesh,
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
				// Ideally the scale of all provided meshes fits them exactly into a unit cube, so scale is easy to work with.
				// In this case the vehicle model is already reasonably large, so don't scale it as much as other models.
				Scale: mgl32.Vec3{0.5, 0.5, 0.5},
			},
		},
		// Capsule
		{
			Tag:  "OBJ Textured Mesh",
			Mesh: capsuleMesh,
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
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
		TargetOffset: mgl32.Vec3{0, 0, 1000},
		Up:           mgl32.Vec3{0, 1, 0},
		Near:         0.1,
		Far:          10000,
		FOV:          math.Pi / 4.0,
	}

	rotationPerSecond := mgl32.AnglesToQuat(float32(math.Pi/4), float32(math.Pi/4), float32(math.Pi/4), mgl32.XYZ)

	ticker := time.NewTicker(*frameRate)
	for !view.Window.ShouldClose() {
		glfw.PollEvents() // Reads window events, like keyboard and mouse input.
		fps.Handler.Update()
		keyboard.Handler.Update()
		mouse.Handler.Update()

		ApplyInputs(player, cam)

		// Update the rotation.
		for i := range models {
			models[i].ModifyRotationLocalQ(util.ScaleQuatRotation(rotationPerSecond, fps.Handler.DeltaTimeSeconds()))
		}

		// Set up Model-View-Projection Matrix and send it to the shader program.
		mvMatrix := cam.ModelView()
		w, h := view.Window.GetSize()
		pMatrix := cam.ProjectionPerspective(float32(w), float32(h))
		shader.Model.SetMVPMatrix(pMatrix, mvMatrix)

		// Clear screen, then Draw everything
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		model.RenderXYZAxes()

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
	target.ModifyPosition(move[0], move[1], 0)

	w, h := view.Window.GetSize()
	if mouse.Handler.LeftPressed() {
		move = cam.ScreenToWorldCoord2D(mouse.Handler.Position(), w, h).Sub(target.Center().Vec2())

		move = move.Normalize().Mul(moveSpeed * fps.Handler.DeltaTimeSeconds())
		target.ModifyPosition(move[0], move[1], 0)
	}
}
