// demo showing use of gome to set up a window, but to draw in it with completely custom shaders.
// Mandelbrot shader based on:
// https://blog.mayflower.de/4584-Playing-around-with-pixel-shaders-in-WebGL.html
package main

import (
	"github.com/goxjs/gl"
	"log"
	"time"
	"github.com/omustardo/gome/util/glutil"
	shaderutil "github.com/goxjs/gl/glutil"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/view"
	"github.com/goxjs/glfw"
	"github.com/omustardo/gome"
	"flag"
)

// vertSource is the vertex shader. All of the logic is done in the fragment shader, so this just passes along a vertex.
const vertSource = `
attribute vec3 pos;
void main()
{
    gl_Position = vec4(pos, 1.0);
}
`

// fragSource is the fragment shader. It uses gl_FragCoord and the screen width and height to estimate where the
// current fragment is, and then calculates what the mandelbrot set looks like at that point.
const fragSource = `
#ifdef GL_ES
precision highp float;
#endif

#define NUM_STEPS   50
#define ZOOM_FACTOR 2.0
#define X_OFFSET    0.5
uniform float uWidth;
uniform float uHeight;

void main() {
	vec2 z;
	float x,y;
	int steps;
	float normalizedX = (gl_FragCoord.x - uWidth/2.0) / uWidth * ZOOM_FACTOR * (uWidth/uHeight) - X_OFFSET;
	float normalizedY = (gl_FragCoord.y - uHeight/2.0) / uHeight * ZOOM_FACTOR;

	z.x = normalizedX;
	z.y = normalizedY;

	for (int i=0;i<NUM_STEPS;i++) {
		steps = i;

		x = (z.x * z.x - z.y * z.y) + normalizedX;
		y = (z.y * z.x + z.x * z.y) + normalizedY;

		if((x * x + y * y) > 4.0) {
			break;
		}

		z.x = x;
		z.y = y;
	}

	if (steps == NUM_STEPS-1) {
		gl_FragColor = vec4(1.0, 0.0, 0.0, 1.0);
	} else {
		gl_FragColor = vec4(0.0, 0.0, 0.0, 1.0);
	}
}`

func main() {
	flag.Parse()
	terminate := gome.Initialize("Mandelbrot", 500, 500, "")
	defer terminate()

	program, err := shaderutil.CreateProgram(vertSource, fragSource)
	if err != nil {
		log.Fatal(err)
	}
	gl.ValidateProgram(program)
	if gl.GetProgrami(program, gl.VALIDATE_STATUS) != gl.TRUE {
		log.Fatalf("shader: gl validate status: %s", gl.GetProgramInfoLog(program))
	}
	gl.UseProgram(program)
	VertexPositionAttrib := gl.GetAttribLocation(program, "pos")
	widthUniform := gl.GetUniformLocation(program, "uWidth" )
	heightUniform := gl.GetUniformLocation(program, "uHeight")

	// Set uniforms in the shader to the current size of the window.
	w, h := view.Window.GetSize()
	gl.Uniform1f(widthUniform, float32(w))
	gl.Uniform1f(heightUniform, float32(h))

	// setDimensions updates the shader with the current screen dimensions.
	// It caches values so if there's no change, then the shader doesn't need to be updated.
	// Calls to the GPU are more expensive than standard functions, so should be avoided if possible.
	setDimensions := func() {
		currentW, currentH := view.Window.GetSize()
		if currentW != w {
			w = currentW
			gl.Uniform1f(widthUniform, float32(w))
		}
		if currentH != h {
			h = currentH
			gl.Uniform1f(heightUniform, float32(h))
		}
	}

	// OpenGL uses a coordinate system from [-1,-1] to [1,1]. We want to draw the mandelbrot set over the entire screen
	// so we provide only four vertices (one in each corner). They will be drawn as two triangles.
	quad := []mgl32.Vec3{
		{-1, -1, 0},
		{1, -1, 0},
		{-1, 1, 0},
		{1, 1, 0},
	}
	vertexVBO := glutil.LoadBufferVec3(quad)

	ticker := time.NewTicker(time.Second / 60)
	for !view.Window.ShouldClose() {
		setDimensions()

		// PollEvents reads window events, like keyboard and mouse input. Without this, the window can still be drawn to
		// but it be unresponsive to any interaction, including basic things like resizing.
		glfw.PollEvents()

		// This draws the mandelbrot image every time we go through the loop. Barring inconsistencies due to floating point
		// error, it's the same image every frame so this is quite wasteful. Ideally it should only be drawn when the
		// size of the screen changes.
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.BindBuffer(gl.ARRAY_BUFFER, 	vertexVBO)
		gl.EnableVertexAttribArray(VertexPositionAttrib) // TODO: Can these VertexAttribArrays be enabled a single time in shader initialization and then just always used?
		gl.VertexAttribPointer(VertexPositionAttrib, 3, gl.FLOAT, false, 0, 0)
		gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
		view.Window.SwapBuffers()
		<-ticker.C
	}
}

