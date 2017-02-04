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
	modelVertexSource = `
attribute vec3 aVertexPosition;
attribute vec3 aNormal;
attribute vec2 aTextureCoord;

uniform mat4 uTranslationMatrix;
uniform mat4 uRotationMatrix;
uniform mat4 uScaleMatrix;

uniform mat4 uMVMatrix;
uniform mat4 uPMatrix;

uniform vec3 uDiffuseLightDirection;
uniform vec3 uDiffuseLightColor;

varying vec3 vLighting;
varying vec2 vTextureCoord;

void main() {
	// === Texture ===
	vTextureCoord = aTextureCoord;

	// === Position ===
	vec4 worldPosition = uTranslationMatrix * uRotationMatrix * uScaleMatrix * vec4(aVertexPosition, 1.0);
	gl_Position = uPMatrix * uMVMatrix * worldPosition;

	// === Lighting ===
	// Put normals into world space. Don't translate as that would affect the direction of the normal.
	vec4 worldNormal = normalize(uRotationMatrix * uScaleMatrix * vec4(aNormal, 1.0));
	// TODO: If lighting starts looking odd with meshes that get scaled, look at: http://web.archive.org/web/20120228095346/http://www.arcsynthesis.org/gltut/Illumination/Tut09%20Normal%20Transformation.html

	float intensity = max(dot(worldNormal.xyz, uDiffuseLightDirection.xyz), 0.0);
	vLighting = uDiffuseLightColor * intensity;
}
`
	modelFragmentSource = `
#ifdef GL_ES
precision lowp float;
#endif

uniform sampler2D uSampler;
uniform vec4 uColor;
uniform vec3 uAmbientLight;

varying vec3 vLighting;
varying vec2 vTextureCoord;

void main(void) {
	// gl_FragColor = vec4(uAmbientLight + vLighting, 1.0) * uColor;
	gl_FragColor = texture2D(uSampler, vTextureCoord) * vec4(uAmbientLight + vLighting, 1.0) * uColor;
}
`
)

type model struct {
	Program gl.Program

	translationMatrixUniform gl.Uniform
	rotationMatrixUniform    gl.Uniform
	scaleMatrixUniform       gl.Uniform

	mvMatrixUniform gl.Uniform
	pMatrixUniform  gl.Uniform

	colorUniform        gl.Uniform
	ambientLightUniform gl.Uniform

	diffuseLightDirectionUniform gl.Uniform
	diffuseLightColorUniform     gl.Uniform

	VertexPositionAttrib gl.Attrib
	NormalAttrib         gl.Attrib

	samplerUniform     gl.Uniform
	TextureCoordAttrib gl.Attrib
}

func setupModelShader() error {
	if Model != nil {
		return errors.New("Model Shader already initialized")
	}

	program, err := glutil.CreateProgram(modelVertexSource, modelFragmentSource)
	if err != nil {
		return err
	}
	gl.ValidateProgram(program)
	if gl.GetProgrami(program, gl.VALIDATE_STATUS) != gl.TRUE {
		return fmt.Errorf("basic shader: gl validate status: %s", gl.GetProgramInfoLog(program))
	}
	UseProgram(program)

	// Get gl "names" of variables in the shader program.
	// https://www.opengl.org/sdk/docs/man/html/glUniform.xhtml
	Model = &model{
		Program: program,

		pMatrixUniform:  gl.GetUniformLocation(program, "uPMatrix"),
		mvMatrixUniform: gl.GetUniformLocation(program, "uMVMatrix"),

		translationMatrixUniform: gl.GetUniformLocation(program, "uTranslationMatrix"),
		rotationMatrixUniform:    gl.GetUniformLocation(program, "uRotationMatrix"),
		scaleMatrixUniform:       gl.GetUniformLocation(program, "uScaleMatrix"),

		colorUniform:        gl.GetUniformLocation(program, "uColor"),
		ambientLightUniform: gl.GetUniformLocation(program, "uAmbientLight"),

		diffuseLightDirectionUniform: gl.GetUniformLocation(program, "uDiffuseLightDirection"),
		diffuseLightColorUniform:     gl.GetUniformLocation(program, "uDiffuseLightColor"),

		VertexPositionAttrib: gl.GetAttribLocation(program, "aVertexPosition"),
		NormalAttrib:         gl.GetAttribLocation(program, "aNormal"),

		samplerUniform:     gl.GetUniformLocation(program, "uSampler"),
		TextureCoordAttrib: gl.GetAttribLocation(program, "aTextureCoord"),
	}
	return nil
}

func (s *model) SetDefaults() {
	UseProgram(s.Program)
	s.SetColor(nil)
	s.SetDiffuse(mgl32.Vec3{0, 0, 0}, &color.NRGBA{})
	s.SetTranslationMatrix(0, 0, 0)
	s.SetRotationMatrix(0, 0, 0)
	s.SetScaleMatrix(1, 1, 1)
}

func (s *model) SetColor(color *color.NRGBA) {
	UseProgram(s.Program)
	if color == nil {
		gl.Uniform4f(s.colorUniform, 1, 1, 1, 1)
		return
	}
	gl.Uniform4f(s.colorUniform, float32(color.R)/255.0, float32(color.G)/255.0, float32(color.B)/255.0, float32(color.A)/255.0)
}

// SetDiffuse sets directional lighting. If color is nil, it's set to white. Alpha value is ignored.
func (s *model) SetDiffuse(dir mgl32.Vec3, color *color.NRGBA) {
	UseProgram(s.Program)
	if color == nil {
		gl.Uniform3f(s.diffuseLightColorUniform, 1, 1, 1)
	} else {
		gl.Uniform3f(s.diffuseLightColorUniform, float32(color.R)/255.0, float32(color.G)/255.0, float32(color.B)/255.0)
	}
	dir = dir.Normalize()
	gl.Uniform3f(s.diffuseLightDirectionUniform, dir[0], dir[1], dir[2])
}

// Alpha value is ignored since it doesn't make sense. Also why is there no color.RGB?
func (s *model) SetAmbientLight(color *color.NRGBA) {
	UseProgram(s.Program)
	if color == nil {
		gl.Uniform3f(s.ambientLightUniform, 0, 0, 0)
		return
	}
	gl.Uniform3f(s.ambientLightUniform, float32(color.R)/255.0, float32(color.G)/255.0, float32(color.B)/255.0)
}

func (s *model) SetTexture(texture gl.Texture) {
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

func (s *model) SetMVPMatrix(pMatrix, mvMatrix mgl32.Mat4) {
	UseProgram(s.Program)
	gl.UniformMatrix4fv(s.pMatrixUniform, pMatrix[:])
	gl.UniformMatrix4fv(s.mvMatrixUniform, mvMatrix[:])
}

func (s *model) SetTranslationMatrix(x, y, z float32) {
	UseProgram(s.Program)
	translateMatrix := mgl32.Translate3D(x, y, z)
	gl.UniformMatrix4fv(s.translationMatrixUniform, translateMatrix[:])
}

// SetRotationMatrix takes rotation about the X, Y, and Z axes and applies them in the same XYZ order.
func (s *model) SetRotationMatrix(x, y, z float32) {
	UseProgram(s.Program)
	rotationMatrix := mgl32.Rotate3DX(x).Mul3(mgl32.Rotate3DY(y)).Mul3(mgl32.Rotate3DZ(z)).Mat4()
	gl.UniformMatrix4fv(s.rotationMatrixUniform, rotationMatrix[:])
}

func (s *model) SetRotationMatrixQ(q mgl32.Quat) {
	UseProgram(s.Program)
	rotationMatrix := q.Mat4()
	gl.UniformMatrix4fv(s.rotationMatrixUniform, rotationMatrix[:])
}

func (s *model) SetScaleMatrix(x, y, z float32) {
	UseProgram(s.Program)
	scaleMatrix := mgl32.Scale3D(x, y, z)
	gl.UniformMatrix4fv(s.scaleMatrixUniform, scaleMatrix[:])
}
