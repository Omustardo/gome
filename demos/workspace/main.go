package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"path/filepath"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"github.com/omustardo/gome"
	"github.com/omustardo/gome/asset"
	"github.com/omustardo/gome/camera"
	"github.com/omustardo/gome/core/entity"
	"github.com/omustardo/gome/demos/workspace/shape"
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
	windowWidth    = flag.Int("window_width", 1000, "initial window width")
	windowHeight   = flag.Int("window_height", 1000, "initial window height")
	screenshotPath = flag.String("screenshot_dir", `C:\Users\Omar\Desktop\screenshots\`, "Folder to save screenshots in. Name is the timestamp of when they are taken.")

	frameRate    = flag.Duration("framerate", time.Second/60, `Cap on framerate. Provide with units, like "16.66ms"`)
	gametickRate = flag.Duration("gametick_rate", time.Second/3, `How often to calculate major game actions. Provide with units, like "200ms"`)
	debugLogRate = flag.Duration("debug_log_rate", time.Second, `How often to do periodic debug logging. Provide with units, like "5s"`)

	// Explicitly listing the base dir is a hack. It's needed because `go run` produces a binary in a tmp folder so we can't
	// use relative asset paths. More explanation in omustardo\gome\asset\asset.go
	baseDir = flag.String("base_dir", `C:\workspace\Go\src\github.com\omustardo\gome\demos\workspace`, "All file paths should be specified relative to this root.")
)

func main() {
	flag.Parse()

	terminate := gome.Initialize("Graphics Demo", *windowWidth, *windowHeight, *baseDir)
	defer terminate()

	shader.Model.SetAmbientLight(&color.NRGBA{255, 255, 255, 0}) // 2D objects look good in max lighting.
	// shader.Model.SetAmbientLight(&color.NRGBA{60, 60, 60, 0}) // 3D objects don't look 3D in max lighting, so tone it down.

	// =========== Done with common initializations. From here on it's specific to this demo. TODO: Move above stuff into a nice wrapper. Split anything that doesn't depend on OpenGL/main thread into a goroutine.

	// Load Assets
	tex, err := asset.LoadTexture("assets/sample_texture.png")
	if err != nil {
		log.Fatalf("error loading texture: %v", err)
	}
	shipMesh, err := asset.LoadDAE("assets/vehicle0.dae")
	if err != nil {
		log.Fatal(err)
	}

	player := &model.Model{
		Mesh: mesh.NewRectOutline(&color.NRGBA{255, 25, 75, 255}),
		Entity: entity.Entity{
			Position: mgl32.Vec3{},
			Scale:    mgl32.Vec3{100, 100},
			Rotation: mgl32.QuatIdent(),
		},
	}

	//cam := &camera.TargetCamera{
	//	Target:       player,
	//	TargetOffset: mgl32.Vec3{0, 0, 500},
	//	Up:           mgl32.Vec3{0, 1, 0},
	//	Zoomer: zoom.NewScrollZoom(0.25, 3,
	//		func() float32 {
	//			return mouse.Handler.Scroll().Y()
	//		},
	//	),
	//	Near: 0.1,
	//	Far:  10000,
	//	FOV:  math.Pi / 2.0,
	//}

	cam := camera.NewFreeCamera()
	cam.SetPosition(0, 0, 500)

	miscRect := &model.Model{
		Mesh: mesh.NewRect(&color.NRGBA{180, 110, 111, 255}, gl.Texture{}),
		Entity: entity.Entity{
			Position: mgl32.Vec3{-150, 100},
			Scale:    mgl32.Vec3{100, 100},
			Rotation: mgl32.QuatIdent(),
		},
	}

	miscCircles := []model.Model{
		{
			Mesh: mesh.NewCircle(&color.NRGBA{100, 225, 25, 255}, gl.Texture{}),
			Entity: entity.Entity{
				Position: mgl32.Vec3{100, 200, 0},
				Scale:    mgl32.Vec3{20, 20, 0},
				Rotation: mgl32.QuatIdent(),
			},
		},
		{
			Mesh: mesh.NewCircle(nil, tex),
			Entity: entity.Entity{
				Position: mgl32.Vec3{-200, -100, 0},
				Scale:    mgl32.Vec3{15, 15, 0},
				Rotation: mgl32.QuatIdent(),
			},
		},
		{
			Mesh: mesh.NewCircleOutline(&color.NRGBA{255, 125, 50, 255}),
			Entity: entity.Entity{
				Position: mgl32.Vec3{0, 50, 0},
				Scale:    mgl32.Vec3{35, 35, 0},
				Rotation: mgl32.QuatIdent(),
			},
		},
	}

	orbitingRects := []*shape.OrbitingRect{
		shape.NewOrbitingRect(
			entity.Entity{Scale: mgl32.Vec3{100, 100}},
			&color.NRGBA{75, 25, 225, 255},
			mgl32.Vec2{250, 380}, // Center of the orbit
			350,                  // Orbit radius // TODO: Allow elliptical orbits.
			nil,
			5000, // Time to make a full revolution (all the way around the orbit)
			5000, // Time to make a full rotation (turn fully around itself, i.e. 1 day)
		),
		shape.NewOrbitingRect(
			entity.Entity{Scale: mgl32.Vec3{100, 100}},
			&color.NRGBA{25, 100, 225, 255},
			mgl32.Vec2{-400, -30}, // Center of the orbit
			900, // Orbit radius // TODO: Allow elliptical orbits.
			nil,
			10000, // Time to make a full revolution (all the way around the orbit)
			5000,  // Time to make a full rotation (turn fully around itself, i.e. 1 day)
		),
		shape.NewOrbitingRect(
			entity.Entity{Scale: mgl32.Vec3{256, 256}},
			&color.NRGBA{200, 25, 50, 255},
			mgl32.Vec2{-1500, 800}, // Center of the orbit
			800, // Orbit radius // TODO: Allow elliptical orbits.
			player,
			200000, // Time to make a full revolution (all the way around the orbit)
			2000,   // Time to make a full rotation (turn fully around itself, i.e. 1 day)
		),
	}
	orbitingRects = append(orbitingRects,
		shape.NewOrbitingRect(
			entity.Entity{Scale: mgl32.Vec3{128, 128}},
			&color.NRGBA{100, 100, 150, 255},
			mgl32.Vec2{0, 0}, // Center of the orbit
			400,              // Orbit radius // TODO: Allow elliptical orbits.
			orbitingRects[0],
			-2000, // Time to make a full revolution (all the way around the orbit)
			-500,  // Time to make a full rotation (turn fully around itself, i.e. 1 day)
		),
	)

	// Generate parallax rectangles.
	parallaxObjects := shape.GenParallaxRects(cam, 5000, 8, 5, 0.1, 0.2)                                // Near
	parallaxObjects = append(parallaxObjects, shape.GenParallaxRects(cam, 3000, 5, 3.5, 0.35, 0.5)...)  // Med
	parallaxObjects = append(parallaxObjects, shape.GenParallaxRects(cam, 2000, 2, 0.5, 0.75, 0.85)...) // Far
	parallaxObjects = append(parallaxObjects, shape.GenParallaxRects(cam, 1000, 1, 0.1, 0.9, 0.95)...)  // Distant
	// Put the parallax info in buffers on the GPU. TODO: Consider using a single interleaved buffer. Stride and offset are annoying though, and I don't think a few extra buffers matter.
	parallaxPositionBuffer, parallaxTranslationBuffer, parallaxTranslationRatioBuffer, parallaxAngleBuffer, parallaxScaleBuffer, parallaxColorBuffer := shape.GetParallaxBuffers(parallaxObjects)

	texturedRect := model.Model{
		Mesh: mesh.NewRect(nil, tex),
		Entity: entity.Entity{
			Position: mgl32.Vec3{0, -512, 1},
			Scale:    mgl32.Vec3{256, 256},
		},
	}

	texturedCube := model.Model{
		Tag:  "test",
		Mesh: mesh.NewCube(nil, tex),
		Entity: entity.Entity{
			Position: mgl32.Vec3{0, 0, 0},
			Scale:    mgl32.Vec3{32, 32, 32},
			Rotation: mgl32.QuatIdent(),
		},
	}
	rotationPerSecond := float32(math.Pi / 4)

	shipModel := &model.Model{
		Mesh: shipMesh,
		Entity: entity.Entity{
			Position: mgl32.Vec3{0, 100, 0},
			Scale:    mgl32.Vec3{5, 5, 5},
			Rotation: mgl32.QuatIdent(),
		},
	}
	shipModel.Mesh.Color = &color.NRGBA{155, 155, 155, 255}

	ticker := time.NewTicker(*frameRate)
	gameTicker := time.NewTicker(*gametickRate)
	debugLogTicker := time.NewTicker(*debugLogRate)
	for !view.Window.ShouldClose() {
		fps.Handler.Update()
		glfw.PollEvents() // Reads window events, like keyboard and mouse input.
		// Handler.Update takes current input and stores it. This is necessary to detect things like the start of a keypress.
		keyboard.Handler.Update()
		mouse.Handler.Update()

		// Handle Input
		ApplyInputs(player)

		// Update the cube's X and Z rotation.
		r := rotationPerSecond * float32((*frameRate).Seconds())
		texturedCube.ModifyRotationGlobal(mgl32.Vec3{r, 0, r})
		shipModel.ModifyRotationGlobal(mgl32.Vec3{r, 0, r})
		shipModel.ModifyRotationGlobal(mgl32.Vec3{r, 0, r})
		miscCircles[0].ModifyRotationGlobal(mgl32.Vec3{r, 0.8 * r, 1.3 * r})

		for _, r := range orbitingRects {
			r.Update()
		}

		// Run game logic
		select {
		case <-gameTicker.C:
			// do stuff with game logic on ticks to minimize expensive calculations.
		default:
		}
		cam.Update(fps.Handler.DeltaTime())

		// Set up Model-View-Projection Matrix and send it to the shader programs.
		mvMatrix := cam.ModelView()
		w, h := view.Window.GetSize()
		pMatrix := cam.ProjectionPerspective(float32(w), float32(h))
		shader.Parallax.SetMVPMatrix(pMatrix, mvMatrix)
		shader.Model.SetMVPMatrix(pMatrix, mvMatrix)

		// Clear screen, then Draw everything
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // TODO: Some cool graphical effects result from not clearing the screen.
		model.RenderXYZAxes()
		for _, c := range miscCircles {
			c.Render()
		}

		miscRect.Render()

		// Draw parallax objects
		// Old inefficient way of drawing the rectangles one by one:
		//for _, r := range parallaxObjects {
		//	r.DrawFilled()
		//}
		// New batched method:
		shape.DrawParallaxBuffers(6*len(parallaxObjects) /* vertices in total */, cam.Position,
			parallaxPositionBuffer, parallaxTranslationBuffer, parallaxTranslationRatioBuffer,
			parallaxAngleBuffer, parallaxScaleBuffer, parallaxColorBuffer)

		for _, r := range orbitingRects {
			r.DrawOrbit()
		}
		for _, r := range orbitingRects {
			r.Render()
		}

		texturedRect.Render()
		// texturedRect.RenderWireframe()

		texturedCube.Render()

		shipModel.Render()

		player.Render()

		// Debug logging - limited to once every X seconds to avoid spam.
		select {
		case <-debugLogTicker.C:
			// log.Println("location:", cam.Position())
			// if mouse.Handler.LeftPressed() {
			//  log.Println("detected mouse press at", mouse.Handler.Position())
			// }
			log.Println(fps.Handler.FPS(), "fps")
			// log.Println(fps.Handler.DeltaTime(), "delta time")
			// log.Println("zoom%:", cam.GetCurrentZoomPercent())

			//w, h := view.Window.GetSize()
			//log.Println("mouse screen->world:", mouse.Handler.Position(), cam.ScreenToWorldCoord2D(mouse.Handler.Position(), w, h))
		default:
		}

		// Swaps the buffer that was drawn on to be visible. The visible buffer becomes the one that gets drawn on until it's swapped again.
		view.Window.SwapBuffers()
		<-ticker.C // wait up to 1/60th of a second. This caps framerate to 60 FPS.
	}
}

func ApplyInputs(player *model.Model) {
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
	playerSpeed := float32(500)
	move = move.Normalize().Mul(playerSpeed * fps.Handler.DeltaTimeSeconds())
	player.ModifyPosition(move[0], move[1], 0)

	w, h := view.Window.GetSize()
	if keyboard.Handler.JustPressed(glfw.KeySpace) {
		util.SaveScreenshot(w, h, filepath.Join(*screenshotPath, fmt.Sprintf("%d.png", util.GetTimeMillis())))
	}
	if mouse.Handler.LeftPressed() {
		move = mgl32.Vec2{
			mouse.Handler.Position().X() - float32(w)/2,
			-(mouse.Handler.Position().Y() - float32(h)/2),
		}

		move = move.Normalize().Mul(playerSpeed * fps.Handler.DeltaTimeSeconds())
		player.ModifyPosition(move[0], move[1], 0)
	}
	if mouse.Handler.RightPressed() {

	}
}
