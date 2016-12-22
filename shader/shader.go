package shader

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
