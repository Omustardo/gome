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
	"github.com/omustardo/gome/demos/asteroids/asteroid"
	"github.com/omustardo/gome/demos/asteroids/bullet"
	"github.com/omustardo/gome/demos/asteroids/player"
	"github.com/omustardo/gome/input/keyboard"
	"github.com/omustardo/gome/input/mouse"
	"github.com/omustardo/gome/model/mesh"
	"github.com/omustardo/gome/shader"
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

	shipMesh, err := asset.LoadOBJ("assets/ship/ship.obj", asset.OBJOpts{Normalize: true})
	if err != nil {
		log.Fatalf("Unable to load ship model: %v", err)
	}
	shipTexture, err := asset.LoadTexture("assets/ship/ship.jpg")
	if err != nil {
		log.Fatalf("Unable to load asteroid texture: %v", err)
	}
	shipMesh.SetTexture(shipTexture)
	ship := player.New(shipMesh)

	cam := &camera.TargetCamera{
		Target:       ship,
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

	asteroidMesh, err := asset.LoadOBJ("assets/rock/rock1.obj", asset.OBJOpts{Normalize: true, Center: &mgl32.Vec3{0.5, 0.5, 0.5}})
	if err != nil {
		log.Fatalf("Unable to load asteroid model: %v", err)
	}
	asteroidTexture, err := asset.LoadTexture("assets/rock/rock1.jpg")
	if err != nil {
		log.Fatalf("Unable to load asteroid texture: %v", err)
	}
	asteroidMesh.SetTexture(asteroidTexture)
	asteroid.SetMesh(asteroidMesh)

	bullets := []*bullet.Bullet{}
	asteroids := []*asteroid.Asteroid{}
	for i := 0; i < 5; i++ {
		asteroids = append(asteroids, asteroid.New())
	}

	ticker := time.NewTicker(time.Second / 60)
	for !view.Window.ShouldClose() {
		fps.Handler.Update()
		glfw.PollEvents() // Reads window events, like keyboard and mouse input.
		// Handler.Update takes current input and stores it. This is necessary to detect things like the start of a keypress.
		keyboard.Handler.Update()
		mouse.Handler.Update()

		ship.Move(keyboard.Handler.IsKeyDown(glfw.KeyW, glfw.KeyUp), keyboard.Handler.IsKeyDown(glfw.KeyS, glfw.KeyDown), fps.Handler.DeltaTimeSeconds())
		ship.Rotate(keyboard.Handler.IsKeyDown(glfw.KeyA, glfw.KeyLeft), keyboard.Handler.IsKeyDown(glfw.KeyD, glfw.KeyRight), fps.Handler.DeltaTimeSeconds())

		if keyboard.Handler.JustPressed(glfw.KeySpace) {
			if b := ship.FireWeapon(); b != nil {
				bullets = append(bullets, b)
			}
		}

		// TODO: Add "win" condition
		if len(asteroids) == 0 {
			log.Println("Winner!")
			return
		}

		for _, a := range asteroids {
			a.Update(fps.Handler.DeltaTimeSeconds())
		}
		for _, b := range bullets {
			b.Update(fps.Handler.DeltaTimeSeconds())
		}

		asteroidsToAdd := []*asteroid.Asteroid{}
		asteroidsToRemove := make(map[int]bool)
		bulletsToRemove := make(map[int]bool)
		for i, a := range asteroids {
			// if an asteroid is within range of a bullet, split it and destroy the bullet.
			for j, b := range bullets {
				if !bulletsToRemove[j] && a.Position.Sub(b.Position).Len() <= a.Scale.X()+b.Scale.X() {
					bulletsToRemove[j] = true
					asteroidsToRemove[i] = true
					a1, a2 := a.Split()
					if a1 != nil && a2 != nil {
						asteroidsToAdd = append(asteroidsToAdd, a1, a2)
					}
				}
			}
		}
		if asteroidsToRemove != nil {
			temp := []*asteroid.Asteroid{}
			for i := range asteroids {
				if !asteroidsToRemove[i] {
					temp = append(temp, asteroids[i])
				}
			}
			asteroids = temp
		}
		if bulletsToRemove != nil {
			temp := []*bullet.Bullet{}
			for i := range bullets {
				if !bulletsToRemove[i] && bullets[i].LifespanSeconds > 0 {
					temp = append(temp, bullets[i])
				}
			}
			bullets = temp
		}
		asteroids = append(asteroids, asteroidsToAdd...)

		// TODO: If an asteroid is within range of player, game over.
		for _, a := range asteroids {
			if a.Position.Sub(ship.Position).Len() <= a.Scale.X()+ship.Scale.X() {
				log.Println("Death")
				// return
			}
		}

		// If an asteroid is too far from the origin, reverse its velocity
		for _, a := range asteroids {
			if a.Position.Len() > 3000 {
				a.Velocity = a.Velocity.Mul(-1)
			}
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
			a.RenderDebugSphere()
		}
		for _, b := range bullets {
			b.Render()
		}
		ship.Render()

		// Swaps the buffer that was drawn on to be visible. The visible buffer becomes the one that gets drawn on until it's swapped again.
		view.Window.SwapBuffers()
		<-ticker.C // wait up to 1/60th of a second. This caps framerate to 60 FPS.
	}
}
