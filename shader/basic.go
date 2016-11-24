package shader

import (
	"errors"
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/goxjs/gl/glutil"
)

const (
	basicVertexSource = `
//#version 120 // OpenGL 2.1.
//#version 100 // WebGL.
attribute vec3 aVertexPosition;

uniform mat4 uTranslationMatrix;
uniform mat4 uRotationMatrix;
uniform mat4 uScaleMatrix;

uniform mat4 uMVMatrix; // Model-View (transforms the input vertex to the camera's view of the world)
uniform mat4 uPMatrix;  // Projection (transforms camera's view into screen space)

void main() {
	vec4 worldPosition = uTranslationMatrix * uRotationMatrix * uScaleMatrix * vec4(aVertexPosition, 1.0);
	gl_Position = uPMatrix * uMVMatrix * worldPosition;
}
`
	basicFragmentSource = `
//#version 120 // OpenGL 2.1.
//#version 100 // WebGL.
#ifdef GL_ES
precision highp float; // set floating point precision. Required for WebGL.
#endif
uniform vec4 uColor;
void main() {
	gl_FragColor = uColor;
}
`
)

type basic struct {
	Program gl.Program

	translationMatrixUniform gl.Uniform
	rotationMatrixUniform    gl.Uniform
	scaleMatrixUniform       gl.Uniform

	mvMatrixUniform gl.Uniform
	pMatrixUniform  gl.Uniform
	colorUniform    gl.Uniform

	VertexPositionAttrib gl.Attrib
	ParallaxRatioAttrib  gl.Attrib
}

func setupBasicShader() error {
	if Basic != nil {
		return errors.New("Basic Shader already initialized")
	}

	program, err := glutil.CreateProgram(basicVertexSource, basicFragmentSource)
	if err != nil {
		return err
	}
	gl.ValidateProgram(program)
	if gl.GetProgrami(program, gl.VALIDATE_STATUS) != gl.TRUE {
		return fmt.Errorf("basic shader: gl validate status: %s", gl.GetProgramInfoLog(program))
	}
	gl.UseProgram(program)

	// Get gl "names" of variables in the shader program.
	// https://www.opengl.org/sdk/docs/man/html/glUniform.xhtml
	Basic = &basic{
		Program: program,

		translationMatrixUniform: gl.GetUniformLocation(program, "uTranslationMatrix"),
		rotationMatrixUniform:    gl.GetUniformLocation(program, "uRotationMatrix"),
		scaleMatrixUniform:       gl.GetUniformLocation(program, "uScaleMatrix"),
		mvMatrixUniform:          gl.GetUniformLocation(program, "uMVMatrix"),
		pMatrixUniform:           gl.GetUniformLocation(program, "uPMatrix"),
		colorUniform:             gl.GetUniformLocation(program, "uColor"),
		VertexPositionAttrib:     gl.GetAttribLocation(program, "aVertexPosition"),
	}
	return nil
}

func (s *basic) SetDefaults() {
	gl.UseProgram(s.Program)
	s.SetColor(1, 0.1, 1, 1) // Default to a bright purple.
	s.SetTranslationMatrix(0, 0, 0)
	s.SetRotationMatrix(0, 0, 0)
	s.SetScaleMatrix(1, 1, 1)
}

func (s *basic) SetMVPMatrix(pMatrix, mvMatrix mgl32.Mat4) {
	gl.UseProgram(s.Program)
	gl.UniformMatrix4fv(s.pMatrixUniform, pMatrix[:])
	gl.UniformMatrix4fv(s.mvMatrixUniform, mvMatrix[:])
}

func (s *basic) SetTranslationMatrix(x, y, z float32) {
	gl.UseProgram(s.Program)
	translateMatrix := mgl32.Translate3D(x, y, z)
	gl.UniformMatrix4fv(s.translationMatrixUniform, translateMatrix[:])
}

func (s *basic) SetRotationMatrix2D(z float32) {
	gl.UseProgram(s.Program)
	rotationMatrix := mgl32.Rotate3DZ(z).Mat4() // TODO: Use quaternions.
	gl.UniformMatrix4fv(s.rotationMatrixUniform, rotationMatrix[:])
}

func (s *basic) SetRotationMatrix(x, y, z float32) {
	gl.UseProgram(s.Program)
	rotationMatrix := mgl32.Rotate3DX(x).Mul3(mgl32.Rotate3DY(y)).Mul3(mgl32.Rotate3DZ(z)).Mat4() // TODO: Use quaternions.
	gl.UniformMatrix4fv(s.rotationMatrixUniform, rotationMatrix[:])
}

func (s *basic) SetScaleMatrix(x, y, z float32) {
	gl.UseProgram(s.Program)
	scaleMatrix := mgl32.Scale3D(x, y, z)
	gl.UniformMatrix4fv(s.scaleMatrixUniform, scaleMatrix[:])
}

func (s *basic) SetColor(r, g, b, a float32) {
	gl.UseProgram(s.Program)
	// OpenGL is supposed to clamp automatically, but I haven't found the GL ES documentation that actually states that.
	clamp := func(x float32) float32 {
		if x > 1 {
			return 1
		}
		if x < 0 {
			return 0
		}
		return x
	}
	gl.Uniform4f(s.colorUniform, clamp(r), clamp(g), clamp(b), clamp(a))
}
