package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/goxjs/glfw"
	"github.com/omustardo/gome/camera"
	"github.com/omustardo/gome/camera/zoom"
	"github.com/omustardo/gome/input/keyboard"
	"github.com/omustardo/gome/input/mouse"
	"github.com/omustardo/gome/shader"
	"github.com/omustardo/gome/shape"
	"github.com/omustardo/gome/util"
	"github.com/omustardo/gome/util/fps"
	"github.com/omustardo/gome/view"
)

var (
	windowWidth    = flag.Int("window_width", 1000, "initial window width")
	windowHeight   = flag.Int("window_height", 1000, "initial window height")
	screenshotPath = flag.String("screenshot_dir", `C:\Users\Omar\Desktop\screenshots\`, "Folder to save screenshots in. Name is the timestamp of when they are taken.")
)

const (
	gametick  = time.Second / 3
	framerate = time.Second / 60
)

func init() {
	// log print with .go file and line number.
	log.SetFlags(log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	// Initialize gl constants and the glfw window.
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

	// Load standard meshes.
	shape.LoadModels()

	// =========== Done with common initializations. From here on it's specific to this demo. TODO: Move above stuff into a nice wrapper.

	player := &shape.Rect{
		X: 0, Y: 0,
		Width:  100,
		Height: 100,
		R:      0.8, G: 0.1, B: 0.3, A: 1,
		Angle: 0,
	}
	cam := camera.NewTargetCamera(
		player,
		zoom.NewScrollZoom(0.25, 3,
			func() float32 { return mouse.Handler.Scroll().Y() },
		),
	)

	miscCircles := []*shape.Circle{
		{
			P:      mgl32.Vec3{100, 200, 0},
			Radius: 20,
			R:      0.2, G: 0.7, B: 0.5, A: 1,
		},
		{
			P:      mgl32.Vec3{-200, -100, 0},
			Radius: 15,
			R:      0.4, G: 0.9, B: 0.1, A: 1,
		},
		{
			P:      mgl32.Vec3{0, 50, 0},
			Radius: 35,
			R:      1, G: 0.5, B: 0.2, A: 1,
		},
	}

	orbitingRects := []*shape.OrbitingRect{
		shape.NewOrbitingRect(
			shape.Rect{
				Width:  100,
				Height: 100,
				R:      0.3, G: 0.1, B: 0.9, A: 1,
				Angle: 0,
			},
			mgl32.Vec2{250, 380}, // Center of the orbit
			350,                  // Orbit radius // TODO: Allow elliptical orbits.
			nil,
			5000, // Time to make a full revolution (all the way around the orbit)
			5000, // Time to make a full rotation (turn fully around itself, i.e. 1 day)
		),
		shape.NewOrbitingRect(
			shape.Rect{
				Width:  80,
				Height: 55,
				R:      0.1, G: 0.4, B: 0.9, A: 1,
				Angle: 0,
			},
			mgl32.Vec2{-400, -30}, // Center of the orbit
			900, // Orbit radius // TODO: Allow elliptical orbits.
			nil,
			10000, // Time to make a full revolution (all the way around the orbit)
			5000,  // Time to make a full rotation (turn fully around itself, i.e. 1 day)
		),
		shape.NewOrbitingRect(
			shape.Rect{
				Width:  256,
				Height: 256,
				R:      0.8, G: 0.1, B: 0.2, A: 1,
				Angle: 0,
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
			shape.Rect{
				Width:  128,
				Height: 128,
				R:      0.4, G: 0.4, B: 0.6, A: 1,
				Angle: 0,
			},
			mgl32.Vec2{0, 0}, // Center of the orbit
			400,              // Orbit radius // TODO: Allow elliptical orbits.
			orbitingRects[0],
			-2000, // Time to make a full revolution (all the way around the orbit)
			-500,  // Time to make a full rotation (turn fully around itself, i.e. 1 day)
		),
	)

	// Generate parallax rectangles.
	parallaxObjects := shape.GenParallaxRects(cam, 500, 8, 5, 0.1, 0.2)                                // Near
	parallaxObjects = append(parallaxObjects, shape.GenParallaxRects(cam, 300, 5, 3.5, 0.35, 0.5)...)  // Med
	parallaxObjects = append(parallaxObjects, shape.GenParallaxRects(cam, 200, 2, 0.5, 0.75, 0.85)...) // Far
	parallaxObjects = append(parallaxObjects, shape.GenParallaxRects(cam, 100, 1, 0.1, 0.9, 0.95)...)  // Distant
	// Put the parallax info in buffers on the GPU. TODO: Consider using a single interleaved buffer. Stride and offset are annoying though, and I don't think a few extra buffers matter.
	parallaxPositionBuffer, parallaxTranslationBuffer, parallaxTranslationRatioBuffer, parallaxAngleBuffer, parallaxScaleBuffer, parallaxColorBuffer := shape.GetParallaxBuffers(parallaxObjects)

	ticker := time.NewTicker(framerate)
	gameTicker := time.NewTicker(gametick)
	debugLogTicker := time.NewTicker(time.Second)
	for !view.Window.ShouldClose() {
		fps.Handler.Update()
		for _, r := range orbitingRects {
			r.Update()
		}
		// Handle Input
		ApplyInputs(player, cam)

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
		pMatrix := cam.Projection(float32(w), float32(h))
		shader.Basic.SetMVPMatrix(pMatrix, mvMatrix)
		shader.Parallax.SetMVPMatrix(pMatrix, mvMatrix)

		// Clear screen, then Draw everything
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // TODO: Some cool graphical effects result from not clearing the screen.
		shape.DrawXYZAxes()
		for _, c := range miscCircles {
			c.Draw()
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

		player.Draw()

		// Debug logging - limited to once every X seconds to avoid spam.
		select {
		case _, ok := <-debugLogTicker.C:
			if ok {
				// log.Println("location:", cam.Position())
				// if mouseHandler.LeftPressed() {
				// 	 log.Println("detected mouse press at", mouseHandler.Position)
				// }
				// log.Println("zoom%:", cam.GetCurrentZoomPercent())
				// log.Println("mouse screen->world:", mouseHandler.Position, cam.ScreenToWorldCoord2D(mouseHandler.Position, WindowSize))
			}
		default:
		}

		// Swaps the buffer that was drawn on to be visible. The visible buffer becomes the one that gets drawn on until it's swapped again.
		view.Window.SwapBuffers()
		// Handler.Update takes current input and stores it. This is necessary to detect things like the start of a keypress.
		// It's important to do the update for inputs here before PollEvents. Doing these calls at the top of the game loop
		// is equivalent to doing them immediately after PollEvents, and would result in the current input state being
		// skipped, because it would immediately be stored as the previous state.
		keyboard.Handler.Update()
		mouse.Handler.Update()
		glfw.PollEvents() // Reads window events, like keyboard and mouse input.
		<-ticker.C        // wait up to 1/60th of a second. This caps framerate to 60 FPS.
	}
}

func ApplyInputs(player shape.Shape, cam camera.Camera) {
	var move mgl32.Vec2
	if keyboard.Handler.IsKeyDown(glfw.KeyA, glfw.KeyLeft) {
		move[0] = -1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyD, glfw.KeyRight) {
		move[0] = 1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyW, glfw.KeyUp) {
		move[1] = 1
	}
	if keyboard.Handler.IsKeyDown(glfw.KeyS, glfw.KeyDown) {
		move[1] = -1
	}
	move = move.Normalize().Mul(10)
	player.ModifyCenter(move[0], move[1])

	w, h := view.Window.GetSize()
	if keyboard.Handler.JustPressed(glfw.KeySpace) {
		util.SaveScreenshot(w, h, filepath.Join(*screenshotPath, fmt.Sprintf("%d.png", util.GetTimeMillis())))
	}
	if mouse.Handler.LeftPressed() {
		move = cam.ScreenToWorldCoord2D(mouse.Handler.Position(), w, h).Sub(player.Center().Vec2())

		move = move.Normalize().Mul(10)
		player.ModifyCenter(move[0], move[1])
	}
	if mouse.Handler.RightPressed() {
	}
}
