package asset

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"strings"

	"github.com/GlenKelley/go-collada"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/goxjs/gl"
	"github.com/omustardo/gome/model/mesh"
	"github.com/omustardo/gome/util/bytecoder"
)

// LoadOBJ creates a mesh from an obj file.
// Based on https://gist.github.com/davemackintosh/67959fa9dfd9018d79a4
// and https://en.wikipedia.org/wiki/Wavefront_.obj_file
// and http://www.opengl-tutorial.org/beginners-tutorials/tutorial-7-model-loading/
func LoadOBJ(path string) (mesh.Mesh, error) {
	fileData, err := loadFile(path)
	if err != nil {
		return mesh.Mesh{}, err
	}
	out, err := loadOBJData(fileData)
	if err != nil {
		return mesh.Mesh{}, fmt.Errorf("Error loading %s: %v", path, err)
	}
	return out, nil
}

func loadOBJData(data []byte) (mesh.Mesh, error) {
	lines := strings.Split(string(data), "\n")

	var (
		verts, normals                        []mgl32.Vec3
		uvs                                   []mgl32.Vec2
		normalIndices, vertIndices, uvIndices []float32
	)

	for lineNum, line := range lines {
		lineNum++ // numbering is for debug printing, and humans think of files as starting with line 1.

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		// Scan the type field.
		var lineType string
		count, err := fmt.Sscanf(line, "%s", &lineType)
		if err != nil {
			return mesh.Mesh{}, err
		}
		if count != 1 {
			return mesh.Mesh{}, fmt.Errorf("at line #%d, unable to get line type: %v", lineNum, err)
		}
		// Trim off the text that has been read.
		line = strings.TrimSpace(line[len(lineType):])

		switch lineType {
		// VERTICES.
		case "v":
			vec := mgl32.Vec3{}
			count, err := fmt.Sscanf(line, "%f %f %f", &vec[0], &vec[1], &vec[2])
			if err != nil {
				return mesh.Mesh{}, fmt.Errorf("at line #%d, error reading vertices: %v", lineNum, err)
			}
			if count != 3 {
				return mesh.Mesh{}, fmt.Errorf("at line #%d, got %d values for vertices. Expected 3", lineNum, count)
			}
			verts = append(verts, vec)

		// NORMALS.
		case "vn":
			vec := mgl32.Vec3{}
			count, err := fmt.Sscanf(line, "%f %f %f", &vec[0], &vec[1], &vec[2])
			if err != nil {
				return mesh.Mesh{}, fmt.Errorf("at line #%d, error reading normals: %v", lineNum, err)
			}
			if count != 3 {
				return mesh.Mesh{}, fmt.Errorf("at line #%d, got %d values for normals. Expected 3", lineNum, count)
			}
			normals = append(normals, vec)

		// TEXTURE VERTICES.
		case "vt":
			vec := mgl32.Vec2{}
			count, err := fmt.Sscanf(line, "%f %f", &vec[0], &vec[1])
			if err != nil {
				return mesh.Mesh{}, fmt.Errorf("at line #%d, error reading texture vertices: %v", lineNum, err)
			}
			if count != 2 {
				return mesh.Mesh{}, fmt.Errorf("at line #%d, got %v values for texture vertices. Expected 2", lineNum, count)
			}
			uvs = append(uvs, vec)

		// INDICES.
		case "f":
			norm := make([]float32, 3)
			vec := make([]float32, 3)
			uv := make([]float32, 3)
			count, err := fmt.Sscanf(line, "%f/%f/%f %f/%f/%f %f/%f/%f", &vec[0], &uv[0], &norm[0], &vec[1], &uv[1], &norm[1], &vec[2], &uv[2], &norm[2])
			if err != nil {
				return mesh.Mesh{}, fmt.Errorf("at line #%d, error reading indices: %v", lineNum, err)
			}
			if count != 9 {
				return mesh.Mesh{}, fmt.Errorf("at line #%d, got %d values for norm,vec,uv. Expected 9", lineNum, count)
			}
			normalIndices = append(normalIndices, norm[0], norm[1], norm[2])
			vertIndices = append(vertIndices, vec[0], vec[1], vec[2])
			uvIndices = append(uvIndices, uv[0], uv[1], uv[2])

		// COMMENT
		case "#":
			// Do nothing
		default:
			// Do nothing - ignore unknown fields (
		}
	}

	vertexVBO := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexVBO)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Vec3(binary.LittleEndian, verts...), gl.STATIC_DRAW)

	normalVBO := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, normalVBO)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Vec3(binary.LittleEndian, normals...), gl.STATIC_DRAW)

	vertexIndexBuffer := gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexIndexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Float32(binary.LittleEndian, vertIndices...), gl.STATIC_DRAW)

	if glError := gl.GetError(); glError != 0 {
		return mesh.Mesh{}, fmt.Errorf("gl.GetError: %v", glError)
	}

	// TODO: Index and texture buffers are currently unused.

	itemCount := len(verts) / 9 // 3 vertices per point, 3 points per triangle
	return mesh.NewMesh(vertexVBO, gl.Buffer{}, normalVBO, gl.TRIANGLES, itemCount, nil, gl.Texture{}, gl.Buffer{}), nil

	//itemCount := len(vertIndices)
	//return mesh.NewMesh(vertexVBO, vertexIndexBuffer, normalVBO, gl.TRIANGLES, itemCount, nil, gl.Texture{}, gl.Buffer{}), nil
}

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
