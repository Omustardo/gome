package main

import (
	"flag"
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
	"github.com/omustardo/gome/demos/shape"
	"github.com/omustardo/gome/entity"
	"github.com/omustardo/gome/geom"
	"github.com/omustardo/gome/input/keyboard"
	"github.com/omustardo/gome/input/mouse"
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

	// Initialize singletons.
	mouse.Initialize(view.Window)
	keyboard.Initialize(view.Window)
	fps.Initialize()

	// Load standard meshes (cubes, rectangles, etc).
	geom.Initialize()

	// =========== Done with common initializations. From here on it's specific to this demo. TODO: Move above stuff into a nice wrapper. Split anything that doesn't depend on OpenGL/main thread into a goroutine.

	player := &geom.Rect{
		Entity: entity.Entity{
			Position: mgl32.Vec3{},
			Scale:    mgl32.Vec3{100, 100},
			Rotation: mgl32.Vec3{},
		},
		R: 0.8, G: 0.1, B: 0.3, A: 1,
	}

	cam := &camera.TargetCamera{
		Target:       player,
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

	miscCircles := []*geom.Circle{
		{
			Entity: entity.Entity{
				Position: mgl32.Vec3{100, 200, 0},
				Scale:    mgl32.Vec3{20, 20, 0},
				Rotation: mgl32.Vec3{},
			},
			R: 0.2, G: 0.7, B: 0.5, A: 1,
		},
		{
			Entity: entity.Entity{
				Position: mgl32.Vec3{-200, -100, 0},
				Scale:    mgl32.Vec3{15, 15, 0},
				Rotation: mgl32.Vec3{},
			},
			R: 0.4, G: 0.9, B: 0.1, A: 1,
		},
		{
			Entity: entity.Entity{
				Position: mgl32.Vec3{0, 50, 0},
				Scale:    mgl32.Vec3{35, 35, 0},
				Rotation: mgl32.Vec3{},
			},
			R: 1, G: 0.5, B: 0.2, A: 1,
		},
	}

	orbitingRects := []*shape.OrbitingRect{
		shape.NewOrbitingRect(
			geom.Rect{
				Entity: entity.Entity{
					Scale: mgl32.Vec3{100, 100},
				},
				R: 0.3, G: 0.1, B: 0.9, A: 1,
			},
			mgl32.Vec2{250, 380}, // Center of the orbit
			350,                  // Orbit radius // TODO: Allow elliptical orbits.
			nil,
			5000, // Time to make a full revolution (all the way around the orbit)
			5000, // Time to make a full rotation (turn fully around itself, i.e. 1 day)
		),
		shape.NewOrbitingRect(
			geom.Rect{
				Entity: entity.Entity{
					Scale: mgl32.Vec3{80, 55},
				},
				R: 0.1, G: 0.4, B: 0.9, A: 1,
			},
			mgl32.Vec2{-400, -30}, // Center of the orbit
			900, // Orbit radius // TODO: Allow elliptical orbits.
			nil,
			10000, // Time to make a full revolution (all the way around the orbit)
			5000,  // Time to make a full rotation (turn fully around itself, i.e. 1 day)
		),
		shape.NewOrbitingRect(
			geom.Rect{
				Entity: entity.Entity{
					Scale: mgl32.Vec3{256, 256},
				},
				R: 0.8, G: 0.1, B: 0.2, A: 1,
			},
			mgl32.Vec2{-1500, 800}, // Center of the orbit
			800, // Orbit radius // TODO: Allow elliptical orbits.
			player,
			200000, // Time to make a full revolution (all the way around the orbit)
			2000,   // Time to make a full rotation (turn fully around itself, i.e. 1 day)
		),
	}
	orbitingRects = append(orbitingRects,
		shape.NewOrbitingRect(
			geom.Rect{
				Entity: entity.Entity{
					Scale: mgl32.Vec3{128, 128},
				},
				R: 0.4, G: 0.4, B: 0.6, A: 1,
			},
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

	tex, err := asset.LoadTexture("assets/sample_texture.png")
	if err != nil {
		log.Fatalf("error loading texture: %v", err)
	}
	texturedRect := geom.Rect{
		Entity: entity.Entity{
			Position: mgl32.Vec3{0, -512},
			Scale:    mgl32.Vec3{256, 256},
		},
		R: 0.4, G: 0.8, B: 0.6, A: 1,
	}

	texturedCube := geom.Cube{
		Entity: entity.Entity{
			Position: mgl32.Vec3{0, 0, 0},
			Scale:    mgl32.Vec3{32, 32, 32},
			Rotation: mgl32.Vec3{},
		},
		Texture: tex,
	}
	rotationPerSecond := float32(math.Pi / 4)

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
		ApplyInputs(player, cam)

		// Update the cube's X and Z rotation.
		texturedCube.Rotation[0] += rotationPerSecond * float32((*frameRate).Seconds())
		texturedCube.Rotation[2] += rotationPerSecond * float32((*frameRate).Seconds())

		for _, r := range orbitingRects {
			r.Update()
		}

		// Run game logic
		select {
		case _, ok := <-gameTicker.C: // do stuff with game logic on ticks to minimize expensive calculations.
			if ok {
			}
		default:
		}
		cam.Update()

		// Set up Model-View-Projection Matrix and send it to the shader programs.
		mvMatrix := cam.ModelView()
		w, h := view.Window.GetSize()
		pMatrix := cam.ProjectionPerspective(float32(w), float32(h))
		shader.Basic.SetMVPMatrix(pMatrix, mvMatrix)
		shader.Parallax.SetMVPMatrix(pMatrix, mvMatrix)
		shader.Texture.SetMVPMatrix(pMatrix, mvMatrix)
		shader.UI.SetMVPMatrix(cam.ProjectionOrthographic(float32(w), float32(h)), mvMatrix)

		// Clear screen, then Draw everything
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // TODO: Some cool graphical effects result from not clearing the screen.
		util.DrawXYZAxes()
		for _, c := range miscCircles {
			c.DrawWireframe()
		}

		// Draw parallax objects
		// Old inefficient way of drawing the rectangles one by one:
		//for _, r := range parallaxObjects {
		//	r.DrawFilled()
		//}
		// New batched method:
		shape.DrawParallaxBuffers(6*len(parallaxObjects) /* vertices in total */, cam.Position(),
			parallaxPositionBuffer, parallaxTranslationBuffer, parallaxTranslationRatioBuffer,
			parallaxAngleBuffer, parallaxScaleBuffer, parallaxColorBuffer)

		for _, r := range orbitingRects {
			r.DrawOrbit()
		}
		for _, r := range orbitingRects {
			r.DrawFilled()
		}

		texturedRect.DrawTextured(tex)
		texturedRect.Draw()

		texturedCube.Draw()

		player.Draw()

		// Debug logging - limited to once every X seconds to avoid spam.
		select {
		case _, ok := <-debugLogTicker.C:
			if ok {
				// log.Println("location:", cam.Position())
				// if mouse.Handler.LeftPressed() {
				//  log.Println("detected mouse press at", mouse.Handler.Position())
				// }
				// log.Println(fps.Handler.FPS(), "fps")
				// log.Println(fps.Handler.DeltaTime(), "delta time")
				// log.Println("zoom%:", cam.GetCurrentZoomPercent())

				//w, h := view.Window.GetSize()
				//log.Println("mouse screen->world:", mouse.Handler.Position(), cam.ScreenToWorldCoord2D(mouse.Handler.Position(), w, h))
			}
		default:
		}

		// Swaps the buffer that was drawn on to be visible. The visible buffer becomes the one that gets drawn on until it's swapped again.
		view.Window.SwapBuffers()
		<-ticker.C // wait up to 1/60th of a second. This caps framerate to 60 FPS.
	}
}

func ApplyInputs(player *geom.Rect, cam camera.Camera) {
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
	player.ModifyCenter(move[0], move[1])

	w, h := view.Window.GetSize()
	if keyboard.Handler.JustPressed(glfw.KeySpace) {
		// util.SaveScreenshot(w, h, filepath.Join(*screenshotPath, fmt.Sprintf("%d.png", util.GetTimeMillis())))
	}
	if mouse.Handler.LeftPressed() {
		move = cam.ScreenToWorldCoord2D(mouse.Handler.Position(), w, h).Sub(player.Center().Vec2())

		move = move.Normalize().Mul(playerSpeed * fps.Handler.DeltaTimeSeconds())
		player.ModifyCenter(move[0], move[1])
	}
	if mouse.Handler.RightPressed() {

	}
}
