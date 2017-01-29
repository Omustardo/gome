package shader

import "github.com/goxjs/gl"

// Online live shader editor: http://shdr.bkcore.com/
// gman's explanation is great: http://stackoverflow.com/questions/30364213/shaders-in-webgl-vs-opengl
// GLSL (GL Shading Language) Reference: http://www.shaderific.com/glsl/   Particularly the qualifiers section.

// Note that normally a shader starts with a line like:
//#version 120 // OpenGL 2.1.
// or:
//#version 100 // WebGL.
// But since these shaders must work for both desktop and webgl we leave them off and expect those to be the defaults.
// It's a bit risky, but probably fine.

var (
	Parallax *parallax
	Model    *model

	// activeProgram is the current active gl program.
	// Keeping track of this locally allows calls to gl.UseProgram to be avoided if the given program is already active.
	activeProgram gl.Program
)

func Initialize() error {
	errs := make(chan error, 10)
	errs <- setupParallaxShader()
	errs <- setupModelShader()
	close(errs)
	for err := range errs {
		if err != nil {
			return err
		}
	}
	Parallax.SetDefaults()
	Model.SetDefaults()
	return nil
}

// UseProgram sets the provided shader program as being active - meaning it will be used in rendering calls.
//
// Note that the gl program can be set directly by calling gl.UseProgram, rather than this method.
// Generally a UseProgram call is expensive, but by wrapping it in this method we can keep track of the currently active
// program so the GPU doesn't need to be interacted with if there's no change.
// Given this, it's recommended not to use gl.UseProgram and to use this method instead.
func UseProgram(p gl.Program) {
	if p != activeProgram {
		gl.UseProgram(p)
		activeProgram = p
	}
}
