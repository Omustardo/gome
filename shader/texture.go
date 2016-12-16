package shader

import (
	"errors"
	"fmt"

	"image/color"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/goxjs/gl/glutil"
)

const (
	textureVertexSource = `
//#version 120 // OpenGL 2.1.
//#version 100 // WebGL.
attribute vec3 aVertexPosition;
attribute vec2 aTextureCoord;

uniform mat4 uTranslationMatrix;
uniform mat4 uRotationMatrix;
uniform mat4 uScaleMatrix;

uniform mat4 uMVMatrix; // Model-View (transforms the input vertex to the camera's view of the world)
uniform mat4 uPMatrix;  // Projection (transforms camera's view into screen space)

varying vec2 vTextureCoord;

void main() {
	vec4 worldPosition = uTranslationMatrix * uRotationMatrix * uScaleMatrix * vec4(aVertexPosition, 1.0);
	gl_Position = uPMatrix * uMVMatrix * worldPosition;
	vTextureCoord = aTextureCoord;
}
`
	textureFragmentSource = `
//#version 120 // OpenGL 2.1.
//#version 100 // WebGL.
#ifdef GL_ES
precision highp float; // set floating point precision. Required for WebGL.
#endif

uniform sampler2D uSampler;
uniform vec4 uColor;

varying vec2 vTextureCoord;

void main() {
	// gl_FragColor = uColor * texture2D(uSampler, vec2(vTextureCoord.s, vTextureCoord.t)); // TODO: Modify by color value, if it's provided.
	gl_FragColor = texture2D(uSampler, vTextureCoord); // vec2(vTextureCoord.s, vTextureCoord.t));
}
`
)

type texture struct {
	Program gl.Program

	translationMatrixUniform gl.Uniform
	rotationMatrixUniform    gl.Uniform
	scaleMatrixUniform       gl.Uniform
	samplerUniform           gl.Uniform

	mvMatrixUniform gl.Uniform
	pMatrixUniform  gl.Uniform
	colorUniform    gl.Uniform

	VertexPositionAttrib gl.Attrib
	TextureCoordAttrib   gl.Attrib
}

func setupTextureShader() error {
	if Texture != nil {
		return errors.New("Texture Shader already initialized")
	}

	program, err := glutil.CreateProgram(textureVertexSource, textureFragmentSource)
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
	Texture = &texture{
		Program: program,

		translationMatrixUniform: gl.GetUniformLocation(program, "uTranslationMatrix"),
		rotationMatrixUniform:    gl.GetUniformLocation(program, "uRotationMatrix"),
		scaleMatrixUniform:       gl.GetUniformLocation(program, "uScaleMatrix"),
		samplerUniform:           gl.GetUniformLocation(program, "uSampler"),
		mvMatrixUniform:          gl.GetUniformLocation(program, "uMVMatrix"),
		pMatrixUniform:           gl.GetUniformLocation(program, "uPMatrix"),
		colorUniform:             gl.GetUniformLocation(program, "uColor"),
		VertexPositionAttrib:     gl.GetAttribLocation(program, "aVertexPosition"),
		TextureCoordAttrib:       gl.GetAttribLocation(program, "aTextureCoord"),
	}
	return nil
}

func (s *texture) SetDefaults() {
	gl.UseProgram(s.Program)
	s.SetColor(&color.NRGBA{255, 25, 255, 255}) // Default to a bright purple.
	s.SetTranslationMatrix(0, 0, 0)
	s.SetRotationMatrix(0, 0, 0)
	s.SetScaleMatrix(1, 1, 1)
}

func (s *texture) SetMVPMatrix(pMatrix, mvMatrix mgl32.Mat4) {
	gl.UseProgram(s.Program)
	gl.UniformMatrix4fv(s.pMatrixUniform, pMatrix[:])
	gl.UniformMatrix4fv(s.mvMatrixUniform, mvMatrix[:])
}

func (s *texture) SetTranslationMatrix(x, y, z float32) {
	gl.UseProgram(s.Program)
	translateMatrix := mgl32.Translate3D(x, y, z)
	gl.UniformMatrix4fv(s.translationMatrixUniform, translateMatrix[:])
}

func (s *texture) SetRotationMatrix2D(z float32) {
	gl.UseProgram(s.Program)
	rotationMatrix := mgl32.Rotate3DZ(z).Mat4() // TODO: Use quaternions.
	gl.UniformMatrix4fv(s.rotationMatrixUniform, rotationMatrix[:])
}

func (s *texture) SetRotationMatrix(x, y, z float32) {
	gl.UseProgram(s.Program)
	rotationMatrix := mgl32.Rotate3DX(x).Mul3(mgl32.Rotate3DY(y)).Mul3(mgl32.Rotate3DZ(z)).Mat4() // TODO: Use quaternions.
	gl.UniformMatrix4fv(s.rotationMatrixUniform, rotationMatrix[:])
}

func (s *texture) SetScaleMatrix(x, y, z float32) {
	gl.UseProgram(s.Program)
	scaleMatrix := mgl32.Scale3D(x, y, z)
	gl.UniformMatrix4fv(s.scaleMatrixUniform, scaleMatrix[:])
}

func (s *texture) SetTextureSampler(texture gl.Texture) {
	gl.ActiveTexture(gl.TEXTURE0) // Determines where the BindTexture calls get bound. Necessary if using multiple textures at once. Good habit to get into using regardless.
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.Uniform1i(s.samplerUniform, 0) // Unless you are using multiple textures, the second parameter to glUniform1i should always be 0. http://stackoverflow.com/questions/14022274/hardcoding-glsl-texture-sampler2d

	// If using multiple textures:
	//// Attach Texture 0
	//glActiveTexture(GL_TEXTURE0);
	//glBindTexture(GL_TEXTURE_2D, _texture0);
	//glUniform1i(_uSampler0, 0);
	//
	//// Attach Texture 1
	//glActiveTexture(GL_TEXTURE1);
	//glBindTexture(GL_TEXTURE_2D, _texture1);
	//glUniform1i(_uSampler1, 1);
}

func (s *texture) SetColor(color *color.NRGBA) {
	gl.UseProgram(s.Program)
	if color == nil {
		gl.Uniform4f(s.colorUniform, 1, 1, 1, 1)
	}
	gl.Uniform4f(s.colorUniform, float32(color.R)/255.0, float32(color.G)/255.0, float32(color.B)/255.0, float32(color.A)/255.0)
}
