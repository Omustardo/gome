package asset

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/GlenKelley/go-collada"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/model/mesh"
	"github.com/omustardo/gome/util/bytecoder"
)

func LoadDAE(path string) (mesh.Mesh, error) {
	data, err := loadFile(path)
	if err != nil {
		return mesh.Mesh{}, err
	}
	return loadDAEData(data)
}

func loadDAEData(data []byte) (mesh.Mesh, error) {
	reader := bytes.NewBuffer(data)

	doc, err := collada.LoadDocumentFromReader(reader)
	if err != nil {
		return mesh.Mesh{}, err
	}

	var m_TriangleCount int
	// Calculate the total triangle and line counts.
	for _, geometry := range doc.LibraryGeometries[0].Geometry {
		for _, triangle := range geometry.Mesh.Triangles {
			m_TriangleCount += triangle.HasCount.Count
		}
	}

	vertices := make([]float32, 3*3*m_TriangleCount)
	normals := make([]float32, 3*3*m_TriangleCount)

	nTriangleNumber := 0
	for _, geometry := range doc.LibraryGeometries[0].Geometry {
		if len(geometry.Mesh.Triangles) == 0 {
			continue
		}

		// HACK. 0 seems to be position, 1 is normal, but need to not hardcode this.
		pVertexData := geometry.Mesh.Source[0].FloatArray.F32()
		pNormalData := geometry.Mesh.Source[1].FloatArray.F32()

		unsharedCount := len(geometry.Mesh.Vertices.Input)

		for _, triangles := range geometry.Mesh.Triangles {
			sharedIndicies := triangles.HasP.P.I()
			sharedCount := len(triangles.HasSharedInput.Input)

			for i := 0; i < triangles.HasCount.Count; i++ {
				offset := 0 // HACK. 0 seems to be position, 1 is normal, but need to not hardcode this.
				vertices[3*3*nTriangleNumber+0] = pVertexData[3*sharedIndicies[(3*i+0)*sharedCount+offset]+0]
				vertices[3*3*nTriangleNumber+1] = pVertexData[3*sharedIndicies[(3*i+0)*sharedCount+offset]+1]
				vertices[3*3*nTriangleNumber+2] = pVertexData[3*sharedIndicies[(3*i+0)*sharedCount+offset]+2]
				vertices[3*3*nTriangleNumber+3] = pVertexData[3*sharedIndicies[(3*i+1)*sharedCount+offset]+0]
				vertices[3*3*nTriangleNumber+4] = pVertexData[3*sharedIndicies[(3*i+1)*sharedCount+offset]+1]
				vertices[3*3*nTriangleNumber+5] = pVertexData[3*sharedIndicies[(3*i+1)*sharedCount+offset]+2]
				vertices[3*3*nTriangleNumber+6] = pVertexData[3*sharedIndicies[(3*i+2)*sharedCount+offset]+0]
				vertices[3*3*nTriangleNumber+7] = pVertexData[3*sharedIndicies[(3*i+2)*sharedCount+offset]+1]
				vertices[3*3*nTriangleNumber+8] = pVertexData[3*sharedIndicies[(3*i+2)*sharedCount+offset]+2]

				if unsharedCount*sharedCount == 2 {
					offset = sharedCount - 1 // HACK. 0 seems to be position, 1 is normal, but need to not hardcode this.
					normals[3*3*nTriangleNumber+0] = pNormalData[3*sharedIndicies[(3*i+0)*sharedCount+offset]+0]
					normals[3*3*nTriangleNumber+1] = pNormalData[3*sharedIndicies[(3*i+0)*sharedCount+offset]+1]
					normals[3*3*nTriangleNumber+2] = pNormalData[3*sharedIndicies[(3*i+0)*sharedCount+offset]+2]
					normals[3*3*nTriangleNumber+3] = pNormalData[3*sharedIndicies[(3*i+1)*sharedCount+offset]+0]
					normals[3*3*nTriangleNumber+4] = pNormalData[3*sharedIndicies[(3*i+1)*sharedCount+offset]+1]
					normals[3*3*nTriangleNumber+5] = pNormalData[3*sharedIndicies[(3*i+1)*sharedCount+offset]+2]
					normals[3*3*nTriangleNumber+6] = pNormalData[3*sharedIndicies[(3*i+2)*sharedCount+offset]+0]
					normals[3*3*nTriangleNumber+7] = pNormalData[3*sharedIndicies[(3*i+2)*sharedCount+offset]+1]
					normals[3*3*nTriangleNumber+8] = pNormalData[3*sharedIndicies[(3*i+2)*sharedCount+offset]+2]
				}

				nTriangleNumber++
			}
		}
	}

	vertexVBO := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexVBO)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Float32(binary.LittleEndian, vertices...), gl.STATIC_DRAW)

	normalVBO := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, normalVBO)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Float32(binary.LittleEndian, normals...), gl.STATIC_DRAW)

	if glError := gl.GetError(); glError != 0 {
		return mesh.Mesh{}, fmt.Errorf("gl.GetError: %v", glError)
	}

	return mesh.NewMesh(vertexVBO, gl.Buffer{}, normalVBO, gl.TRIANGLES, 3*m_TriangleCount, nil, gl.Texture{}, gl.Buffer{}), nil
}
