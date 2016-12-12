package shader

import (
	"errors"
	"fmt"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/goxjs/gl/glutil"
)

const (
	modelVertexSource = `//#version 120 // OpenGL 2.1.
//#version 100 // WebGL.

attribute vec3 aVertexPosition;
attribute vec3 aNormal;

uniform mat4 uTranslationMatrix;
uniform mat4 uRotationMatrix;
uniform mat4 uScaleMatrix;

uniform mat4 uMVMatrix;
uniform mat4 uPMatrix;
uniform mat4 uNormalMatrix;

varying vec3 vLighting;

void main() {
	// Position
	vec4 worldPosition = uTranslationMatrix * uRotationMatrix * uScaleMatrix * vec4(aVertexPosition, 1.0);
	gl_Position = uPMatrix * uMVMatrix * worldPosition;

	// Lighting
	vec3 ambientLight = vec3(0.3, 0.3, 0.3);
	vec3 directionalLight = vec3(0.5, 0.5, 0.5); // Color of the directional light.
	vec3 directionalVector = vec3(0, 0, 1); // Light goes from negative to positive Z.

	vec4 worldNormal = uTranslationMatrix * uRotationMatrix * vec4(aNormal, 1.0); // Put normals into world space. No need to scale since they stay as unit vectors.
	vec4 transformedNormal = uNormalMatrix * worldNormal; // Need to adjust normals a bit more. See: http://web.archive.org/web/20120228095346/http://www.arcsynthesis.org/gltut/Illumination/Tut09%20Normal%20Transformation.html

	float directional = max(dot(transformedNormal.xyz, directionalVector), 0.0);
	vLighting = ambientLight + (directionalLight * directional);
}
`
	modelFragmentSource = `//#version 120 // OpenGL 2.1.
//#version 100 // WebGL.

#ifdef GL_ES
precision lowp float;
#endif

varying vec3 vLighting;

void main(void) {
	gl_FragColor = vec4(vLighting, 1.0);
}
`
)

type model struct {
	Program gl.Program

	translationMatrixUniform gl.Uniform
	rotationMatrixUniform    gl.Uniform
	scaleMatrixUniform       gl.Uniform

	mvMatrixUniform     gl.Uniform
	pMatrixUniform      gl.Uniform
	normalMatrixUniform gl.Uniform

	VertexPositionAttrib gl.Attrib
	NormalAttrib         gl.Attrib
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
	gl.UseProgram(program)

	// Get gl "names" of variables in the shader program.
	// https://www.opengl.org/sdk/docs/man/html/glUniform.xhtml
	Model = &model{
		Program: program,

		pMatrixUniform:      gl.GetUniformLocation(program, "uPMatrix"),
		mvMatrixUniform:     gl.GetUniformLocation(program, "uMVMatrix"),
		normalMatrixUniform: gl.GetUniformLocation(program, "uNormalMatrix"),

		translationMatrixUniform: gl.GetUniformLocation(program, "uTranslationMatrix"),
		rotationMatrixUniform:    gl.GetUniformLocation(program, "uRotationMatrix"),
		scaleMatrixUniform:       gl.GetUniformLocation(program, "uScaleMatrix"),

		VertexPositionAttrib: gl.GetAttribLocation(program, "aVertexPosition"),
		NormalAttrib:         gl.GetAttribLocation(program, "aNormal"),
	}
	return nil
}

func (s *model) SetDefaults() {
	gl.UseProgram(s.Program)
	s.SetTranslationMatrix(0, 0, 0)
	s.SetRotationMatrix(0, 0, 0)
	s.SetScaleMatrix(1, 1, 1)
}

func (s *model) SetMVPMatrix(pMatrix, mvMatrix mgl32.Mat4) {
	gl.UseProgram(s.Program)
	gl.UniformMatrix4fv(s.pMatrixUniform, pMatrix[:])
	gl.UniformMatrix4fv(s.mvMatrixUniform, mvMatrix[:])
	normalMatrix := mvMatrix.Inv().Transpose()
	gl.UniformMatrix4fv(s.normalMatrixUniform, normalMatrix[:])
}

func (s *model) SetTranslationMatrix(x, y, z float32) {
	gl.UseProgram(s.Program)
	translateMatrix := mgl32.Translate3D(x, y, z)
	gl.UniformMatrix4fv(s.translationMatrixUniform, translateMatrix[:])
}

func (s *model) SetRotationMatrix(x, y, z float32) {
	gl.UseProgram(s.Program)
	rotationMatrix := mgl32.Rotate3DX(x).Mul3(mgl32.Rotate3DY(y)).Mul3(mgl32.Rotate3DZ(z)).Mat4() // TODO: Use quaternions.
	gl.UniformMatrix4fv(s.rotationMatrixUniform, rotationMatrix[:])
}

func (s *model) SetScaleMatrix(x, y, z float32) {
	gl.UseProgram(s.Program)
	scaleMatrix := mgl32.Scale3D(x, y, z)
	gl.UniformMatrix4fv(s.scaleMatrixUniform, scaleMatrix[:])
}
