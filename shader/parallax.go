package shader

import (
	"errors"
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/goxjs/gl/glutil"
)

const (
	parallaxVertexSource = `
attribute vec3 aPosition; // Rectangle vertices. These are the same for every rectangle.. bad.
attribute vec2 aTranslation;
attribute float aTranslationRatio; // Parallax objects move a percentage of the camera's position.
attribute float aAngle; // in radians
attribute vec2 aScale;
attribute vec3 uColor;

varying vec3 vColor; // varying lets it get passed to the fragment shader.

uniform vec2 uCameraPosition;
uniform mat4 uMVMatrix; // Model-View (transforms the input vertex to the camera's view of the world)
uniform mat4 uPMatrix;  // Projection (transforms camera's view into screen space)

// http://www.neilmendoza.com/glsl-rotation-about-an-arbitrary-axis/
mat4 rotationMatrix(float angle) {
		vec3 axis = vec3(0,0,1); // This must be normalized
    float s = sin(angle);
    float c = cos(angle);
    float oc = 1.0 - c;

    return mat4(oc * axis.x * axis.x + c,           oc * axis.x * axis.y - axis.z * s,  oc * axis.z * axis.x + axis.y * s,  0.0,
                oc * axis.x * axis.y + axis.z * s,  oc * axis.y * axis.y + c,           oc * axis.y * axis.z - axis.x * s,  0.0,
                oc * axis.z * axis.x - axis.y * s,  oc * axis.y * axis.z + axis.x * s,  oc * axis.z * axis.z + c,           0.0,
                0.0,                                0.0,                                0.0,                                1.0);
}

void main() {
	vec4 worldPosition = vec4(aPosition.xy, -10, 1.0);
	worldPosition.x *= aScale.x;
	worldPosition.y *= aScale.y;
	worldPosition = rotationMatrix(aAngle) * worldPosition;
	worldPosition.x += (aTranslation.x + aTranslationRatio * uCameraPosition.x);
	worldPosition.y += (aTranslation.y + aTranslationRatio * uCameraPosition.y);

	gl_Position = uPMatrix * uMVMatrix * worldPosition;
	vColor = uColor;
}
`
	parallaxFragmentSource = `
#ifdef GL_ES
precision highp float; // set floating point precision. Required for WebGL.
#endif
varying vec3 vColor;
void main() {
	gl_FragColor = vec4(vColor, 1.0);
}
`
)

type parallax struct {
	Program gl.Program

	cameraPositionUniform gl.Uniform
	mvMatrixUniform       gl.Uniform
	pMatrixUniform        gl.Uniform

	PositionAttrib         gl.Attrib
	TranslationAttrib      gl.Attrib
	TranslationRatioAttrib gl.Attrib
	AngleAttrib            gl.Attrib
	ScaleAttrib            gl.Attrib
	ColorAttrib            gl.Attrib
}

func setupParallaxShader() error {
	if Parallax != nil {
		return errors.New("Parallax Shader already initialized")
	}

	program, err := glutil.CreateProgram(parallaxVertexSource, parallaxFragmentSource)
	if err != nil {
		return err
	}
	gl.ValidateProgram(program)
	if gl.GetProgrami(program, gl.VALIDATE_STATUS) != gl.TRUE {
		return fmt.Errorf("parallax shader: gl validate status: %s", gl.GetProgramInfoLog(program))
	}
	UseProgram(program)

	// Get gl "names" of variables in the shader program.
	// https://www.opengl.org/sdk/docs/man/html/glUniform.xhtml
	Parallax = &parallax{
		Program: program,

		cameraPositionUniform: gl.GetUniformLocation(program, "uCameraPosition"),
		mvMatrixUniform:       gl.GetUniformLocation(program, "uMVMatrix"),
		pMatrixUniform:        gl.GetUniformLocation(program, "uPMatrix"),

		PositionAttrib:         gl.GetAttribLocation(program, "aPosition"),
		TranslationAttrib:      gl.GetAttribLocation(program, "aTranslation"),
		TranslationRatioAttrib: gl.GetAttribLocation(program, "aTranslationRatio"),
		AngleAttrib:            gl.GetAttribLocation(program, "aAngle"),
		ScaleAttrib:            gl.GetAttribLocation(program, "aScale"),
		ColorAttrib:            gl.GetAttribLocation(program, "uColor"),
	}
	return nil
}

func (s *parallax) SetDefaults() {
	UseProgram(s.Program) // TODO: If all of these UseProgram calls are expensive, make a global shader.Current = shader.Basic.Program, and it's easy to do CPU only checks.
}

func (s *parallax) SetCameraPosition(pos mgl32.Vec3) {
	UseProgram(s.Program)
	gl.Uniform2f(s.cameraPositionUniform, pos.X(), pos.Y())
}

func (s *parallax) SetMVPMatrix(pMatrix, mvMatrix mgl32.Mat4) {
	UseProgram(s.Program)
	gl.UniformMatrix4fv(s.pMatrixUniform, pMatrix[:])
	gl.UniformMatrix4fv(s.mvMatrixUniform, mvMatrix[:])
}
