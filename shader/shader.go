package shader

import (
	"github.com/goxjs/gl"
)

// Online live shader editor: http://shdr.bkcore.com/
// gman's explanation is great: http://stackoverflow.com/questions/30364213/shaders-in-webgl-vs-opengl
// GLSL (GL Shading Language) Reference: http://www.shaderific.com/glsl/   Particularly the qualifiers section.

var (
	Basic    *basic
	Parallax *parallax
)

func Initialize() error {
	errs := make(chan error, 10)
	errs <- setupBasicShader()
	errs <- setupParallaxShader()
	close(errs)
	for err := range errs {
		if err != nil {
			return err
		}
	}
	Basic.SetDefaults()
	Parallax.SetDefaults()
	gl.UseProgram(Basic.Program)
	return nil
}
