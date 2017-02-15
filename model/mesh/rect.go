package mesh

import (
	"image/color"

	"github.com/goxjs/gl"
	"github.com/omustardo/gome/util/glutil"
)

func NewRect(col *color.NRGBA, texture gl.Texture) Mesh {
	r := rect
	r.Color = col
	r.SetTexture(texture)
	return r
}

func NewRectOutline(col *color.NRGBA) Mesh {
	r := wireframeRect
	r.Color = col
	return r
}

func initializeRect() Mesh {
	// Store basic rectangle vertices in a buffer.
	lower, upper := float32(-0.5), float32(0.5)
	vertexVBO := glutil.LoadBufferFloat32([]float32{
		lower, lower, 0,
		upper, lower, 0,
		upper, upper, 0,
		lower, upper, 0,
	})

	indexBuffer := glutil.LoadIndexBuffer([]uint16{0, 1, 2, 0, 2, 3})

	// Normals for a 2D object extend perpendicular to the plane it lives on.
	normalVBO := glutil.LoadBufferFloat32([]float32{
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
	})

	textureCoordBuffer := glutil.LoadBufferFloat32([]float32{
		0.0, 0.0,
		1.0, 0.0,
		1.0, 1.0,
		0.0, 1.0,
	})

	return NewMesh(vertexVBO, indexBuffer, normalVBO, gl.TRIANGLES, 6, nil, gl.Texture{}, textureCoordBuffer)
}

func initializeWireframeRect() Mesh {
	lower, upper := float32(-0.5), float32(0.5)
	vbo := glutil.LoadBufferFloat32([]float32{
		lower, lower, 0,
		upper, lower, 0,
		upper, upper, 0,
		lower, upper, 0,
	})

	return NewMesh(vbo, gl.Buffer{}, gl.Buffer{}, gl.LINE_LOOP, 4, nil, gl.Texture{}, gl.Buffer{})
}
