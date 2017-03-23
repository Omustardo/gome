package main

import (
	"encoding/binary"
	"flag"
	"image/color"
	"log"
	"math"
	"os"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"github.com/omustardo/bytecoder"
	"github.com/omustardo/gome"
	"github.com/omustardo/gome/camera"
	"github.com/omustardo/gome/camera/zoom"
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
	baseDir = flag.String("base_dir", `C:\workspace\Go\src\github.com\omustardo\gome\demos\animated_hexagons`, "All file paths should be specified relative to this root.")
)

func init() {
	// log print with .go file and line number.
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	flag.Parse()
	terminate := gome.Initialize("Animation Demo", *windowWidth, *windowHeight, *baseDir)
	defer terminate()

	shader.Model.SetAmbientLight(&color.NRGBA{60, 60, 60, 0}) // 3D objects don't look 3D in the default max lighting, so tone it down.

	initializeHexMesh()
	initializeHexWireframe()

	scale := float32(100)
	// Load meshes.
	models := []model.Model{
		{
			Mesh: hexMesh,
			Entity: entity.Entity{
				Rotation: mgl32.QuatIdent(),
				Scale:    mgl32.Vec3{scale, scale, scale},
			},
		},
		//{ // The wireframe adds a nice emphasis, but also a really bad flickering effect since it's drawn so close to the hexagon mesh.
		//	Mesh: columnWireframeMesh,
		//	Entity: entity.Entity{
		//		Rotation: mgl32.QuatIdent(),
		//		Scale:    mgl32.Vec3{scale * 1.02, scale * 1.02, scale * 1.02},
		//	},
		//},
	}

	player := model.Model{
		Mesh:   mesh.NewCube(&color.NRGBA{0, 255, 255, 255}, gl.Texture{}),
		Entity: entity.Default(),
	}
	player.Position[0] = 0

	cam := &camera.TargetCamera{
		Camera: camera.Camera{
			Entity: entity.Default(),
			Near:   0.1,
			Far:    10000,
			FOV:    math.Pi / 4.0,
		},
		Target:       &player,
		TargetOffset: mgl32.Vec3{-500, -500, 1000},
		Zoomer: zoom.NewScrollZoom(0.1, 3,
			func() float32 {
				return mouse.Handler.Scroll().Y()
			},
		),
	}
	cam.ModifyRotationLocal(mgl32.Vec3{math.Pi / 4, 0, 0})

	ticker := time.NewTicker(*frameRate)
	for !view.Window.ShouldClose() {
		glfw.PollEvents() // Reads window events, like keyboard and mouse input.
		fps.Handler.Update()
		keyboard.Handler.Update()
		mouse.Handler.Update()

		ApplyInputs(&player)

		// Set up Model-View-Projection Matrix and send it to the shader program.
		mvMatrix := cam.ModelView()
		w, h := view.Window.GetSize()
		pMatrix := cam.ProjectionPerspective(float32(w), float32(h))
		shader.Model.SetMVPMatrix(pMatrix, mvMatrix)

		cam.Update(fps.Handler.DeltaTime())
		// Clear screen, then Draw everything
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		model.RenderXYZAxes()

		count := 40
		for row := 0; row < count/2; row++ {
			for col := 0; col < count*2; col++ {
				xoffset := 2 * 0.75 * scale * float32(col)
				yoffset := 2 * 0.866 * scale * float32(row)
				if col%2 == 0 {
					yoffset += 0.866 * scale
				}
				for _, m := range models {
					m.Position = mgl32.Vec3{xoffset, yoffset, 50 * float32(math.Sin(100*float64(col)))}
					// Add time based position to animate
					loopDurationMillis := float64(5000)
					m.Position[2] *= float32(row) * float32(math.Sin(float64(util.GetTimeMillis())*math.Pi/loopDurationMillis))
					m.ModifyRotationLocal(mgl32.Vec3{0, float32(math.Sin(math.Pi / 3 * float64(int64(row)*util.GetTimeMillis()) * math.Pi / loopDurationMillis)), 0})
					m.Render()
				}
			}
		}
		player.Render()

		// Swaps the buffer that was drawn on to be visible. The visible buffer becomes the one that gets drawn on until it's swapped again.
		view.Window.SwapBuffers()
		<-ticker.C // wait up to the framerate cap.
	}
}

func ApplyInputs(target *model.Model) {
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
}

var hexMesh mesh.Mesh
var columnWireframeMesh mesh.Mesh

func initializeHexMesh() {
	// Store basic vertices in a buffer. This is a unit cube centered at the origin.
	lower, upper := -float32(1/(2*math.Sqrt(2))), float32(1/(2*math.Sqrt(2)))

	triangles := make([][3]mgl32.Vec3, 0, 6+6+2*6)

	// The hexagonal column can be created in six nearly equivalent steps. The only difference is the angle to use when
	// calculating the points.
	for i := float64(0); i < 6; i++ {
		// startSegment and endSegment are the corners of the hexagon that we need to deal with in this step.
		// For example, on the first step, we start at (1, 0) and end at (cos(pi/3), sin(pi/3)).
		pi := float64(math.Pi)
		startSegment := mgl32.Vec2{float32(math.Cos(i * pi / 3)), float32(math.Sin(i * pi / 3))}
		endSegment := mgl32.Vec2{float32(math.Cos((i + 1) * pi / 3)), float32(math.Sin((i + 1) * pi / 3))}
		top := [3]mgl32.Vec3{
			{0, 0, upper},
			startSegment.Vec3(upper),
			endSegment.Vec3(upper),
		}
		bottom := [3]mgl32.Vec3{
			{0, 0, lower},
			endSegment.Vec3(lower),
			startSegment.Vec3(lower),
		}
		side1 := [3]mgl32.Vec3{
			startSegment.Vec3(upper),
			startSegment.Vec3(lower),
			endSegment.Vec3(lower),
		}
		side2 := [3]mgl32.Vec3{
			startSegment.Vec3(upper),
			endSegment.Vec3(lower),
			endSegment.Vec3(upper),
		}
		triangles = append(triangles, top, bottom, side1, side2)
	}

	vertices := make([]mgl32.Vec3, 0, len(triangles)*3)
	for _, tri := range triangles {
		vertices = append(vertices, tri[0], tri[1], tri[2])
	}

	vertexBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Vec3(binary.LittleEndian, vertices...), gl.STATIC_DRAW)

	normals := mesh.TriangleNormals(triangles)
	normalBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, normalBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Vec3(binary.LittleEndian, normals...), gl.STATIC_DRAW)

	hexMesh = mesh.NewMesh(vertexBuffer, gl.Buffer{}, normalBuffer, gl.TRIANGLES, len(triangles)*3, nil, gl.Texture{}, gl.Buffer{})
	hexMesh.Color = &color.NRGBA{100, 10, 10, 255}
}

func initializeHexWireframe() {
	lower, upper := -float32(1/(2*math.Sqrt(2))), float32(1/(2*math.Sqrt(2)))
	vertices := make([]mgl32.Vec3, 0, 36)

	// The hexagonal column can be created in six nearly equivalent steps. The only difference is the angle to use when
	// calculating the points.
	for i := float64(0); i < 6; i++ {
		// startSegment and endSegment are the corners of the hexagon that we need to deal with in this step.
		// For example, on the first step, we start at (1, 0) and end at (cos(pi/3), sin(pi/3)).
		pi := float64(math.Pi)
		startSegment := mgl32.Vec2{float32(math.Cos(i * pi / 3)), float32(math.Sin(i * pi / 3))}
		endSegment := mgl32.Vec2{float32(math.Cos((i + 1) * pi / 3)), float32(math.Sin((i + 1) * pi / 3))}

		vertices = append(vertices,
			startSegment.Vec3(upper), endSegment.Vec3(upper),
			startSegment.Vec3(upper), startSegment.Vec3(lower),
			startSegment.Vec3(lower), endSegment.Vec3(lower),
		)
	}

	vertexBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Vec3(binary.LittleEndian, vertices...), gl.STATIC_DRAW)

	columnWireframeMesh = mesh.NewMesh(vertexBuffer, gl.Buffer{}, gl.Buffer{}, gl.LINES, len(vertices), nil, gl.Texture{}, gl.Buffer{})
	columnWireframeMesh.Color = &color.NRGBA{255, 255, 255, 255}
}
