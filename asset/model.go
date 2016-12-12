package asset

import (
	"bytes"
	"fmt"
	"io"

	"encoding/binary"

	"github.com/GlenKelley/go-collada"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/model"
	"github.com/omustardo/gome/util/bytecoder"
)

// Model is .
type Model struct {
	Normals, Vecs                        []mgl32.Vec3
	UVs                                  []mgl32.Vec2
	VecIndices, NormalIndices, UVIndices []float32
}

// LoadObj creates a model from a local obj file.
// Based on https://gist.github.com/davemackintosh/67959fa9dfd9018d79a4
// and https://en.wikipedia.org/wiki/Wavefront_.obj_file
// and http://www.opengl-tutorial.org/beginners-tutorials/tutorial-7-model-loading/
func LoadObj(path string) (Model, error) {
	fileData, err := loadFile(path)
	if err != nil {
		return Model{}, err
	}
	return loadObjData(fileData)
}

func loadObjData(data []byte) (Model, error) {
	reader := bytes.NewBuffer(data)

	model := Model{}

	var lineType string
	for {
		// Scan the type field.
		count, err := fmt.Fscanf(reader, "%s", &lineType)
		if count != 1 {
			return Model{}, fmt.Errorf("invalid obj format: err reading line type: %v", err)
		}

		// Check if it's the end of the file
		// and break out of the loop.
		if err != nil {
			if err == io.EOF {
				break
			}
			return Model{}, err
		}

		switch lineType {
		// VERTICES.
		case "v":
			vec := mgl32.Vec3{}
			count, err := fmt.Fscanf(reader, "%f %f %f\n", &vec[0], &vec[1], &vec[2])
			if err != nil {
				return Model{}, fmt.Errorf("invalid obj format: err reading texture vertices: %v", err)
			}
			if count != 3 {
				return Model{}, fmt.Errorf("invalid obj format: got %v values for normals. Expected 3", count)
			}
			model.Vecs = append(model.Vecs, vec)

		// NORMALS.
		case "vn":
			vec := mgl32.Vec3{}
			count, err := fmt.Fscanf(reader, "%f %f %f\n", &vec[0], &vec[1], &vec[2])
			if err != nil {
				return Model{}, fmt.Errorf("invalid obj format: err reading normals: %v", err)
			}
			if count != 3 {
				return Model{}, fmt.Errorf("invalid obj format: got %v values for normals. Expected 3", count)
			}
			model.Normals = append(model.Normals, vec)

		// TEXTURE VERTICES.
		case "vt":
			vec := mgl32.Vec2{}
			count, err := fmt.Fscanf(reader, "%f %f\n", &vec[0], &vec[1])
			if err != nil {
				return Model{}, fmt.Errorf("invalid obj format: err reading texture vertices: %v", err)
			}
			if count != 2 {
				return Model{}, fmt.Errorf("invalid obj format: got %v values for texture vertices. Expected 2", count)
			}
			model.UVs = append(model.UVs, vec)

		// INDICES.
		case "f":
			norm := make([]float32, 3)
			vec := make([]float32, 3)
			uv := make([]float32, 3)
			count, err := fmt.Fscanf(reader, "%f/%f/%f %f/%f/%f %f/%f/%f\n", &vec[0], &uv[0], &norm[0], &vec[1], &uv[1], &norm[1], &vec[2], &uv[2], &norm[2])
			if err != nil {
				return Model{}, fmt.Errorf("invalid obj format: err reading indices: %v", err)
			}
			if count != 9 {
				return Model{}, fmt.Errorf("invalid obj format: got %v values for norm,vec,uv. Expected 9", count)
			}
			model.NormalIndices = append(model.NormalIndices, norm[0], norm[1], norm[2])
			model.VecIndices = append(model.VecIndices, vec[0], vec[1], vec[2])
			model.UVIndices = append(model.UVIndices, uv[0], uv[1], uv[2])

		default:
			// Do nothing - ignore unknown fields
		}
	}

	return model, nil
}

func LoadDAE(path string) (model.Mesh, error) {
	data, err := loadFile(path)
	if err != nil {
		return model.Mesh{}, err
	}
	return loadDAEData(data)
}

func loadDAEData(data []byte) (model.Mesh, error) {
	reader := bytes.NewBuffer(data)

	doc, err := collada.LoadDocumentFromReader(reader)
	if err != nil {
		return model.Mesh{}, err
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
		return model.Mesh{}, fmt.Errorf("gl.GetError: %v", glError)
	}

	return model.Mesh{
		VertexVBO:     vertexVBO,
		NormalVBO:     normalVBO,
		TriangleCount: m_TriangleCount,
	}, nil
}
