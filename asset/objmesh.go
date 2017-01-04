package asset

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"

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

	// The raw per-point data that defines a mesh.
	var (
		verts, normals []mgl32.Vec3
		textureCoords  []mgl32.Vec2
	)

	// Indices are used by the OBJ file format to declare full triangles via the 'f'ace tag.
	// Except for the basic vertex indices, the read in indices are converted back to the values that they reference.
	var vertIndices, textureCoordIndices, normalIndices []uint16

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
			textureCoords = append(textureCoords, vec)

		// FACES.
		case "f":
			// Input expected to be integer indices that refer to data read into the 'v','vt', and 'vn' fields (1 based indexing).
			// Subtract 1 as they are read in to match standard 0 based indexing.
			var vec, uv, norm [3]uint16

			var count, expectedCount int
			switch {
			case strings.Contains(line, "//"):
				count, err = fmt.Sscanf(line, "%d//%d %d//%d %d//%d", &vec[0], &norm[0], &vec[1], &norm[1], &vec[2], &norm[2])
				vertIndices = append(vertIndices, vec[0]-1, vec[1]-1, vec[2]-1)
				normalIndices = append(normalIndices, norm[0]-1, norm[1]-1, norm[2]-1)
				expectedCount = 6
			case strings.Count(line, "/") == 3:
				count, err = fmt.Sscanf(line, "%d/%d %d/%d %d/%d", &vec[0], &uv[0], &vec[1], &uv[1], &vec[2], &uv[2])
				vertIndices = append(vertIndices, vec[0]-1, vec[1]-1, vec[2]-1)
				textureCoordIndices = append(textureCoordIndices, uv[0]-1, uv[1]-1, uv[2]-1)
				expectedCount = 6
			case strings.Count(line, "/") == 6:
				count, err = fmt.Sscanf(line, "%d/%d/%d %d/%d/%d %d/%d/%d", &vec[0], &uv[0], &norm[0], &vec[1], &uv[1], &norm[1], &vec[2], &uv[2], &norm[2])
				vertIndices = append(vertIndices, vec[0]-1, vec[1]-1, vec[2]-1)
				textureCoordIndices = append(textureCoordIndices, uv[0]-1, uv[1]-1, uv[2]-1)
				normalIndices = append(normalIndices, norm[0]-1, norm[1]-1, norm[2]-1)
				expectedCount = 9
			default:
				return mesh.Mesh{}, fmt.Errorf("at line #%d, error reading indices: %v", lineNum, err)
			}
			if err != nil {
				return mesh.Mesh{}, fmt.Errorf("at line #%d, error reading indices: %v", lineNum, err)
			}
			if count != expectedCount {
				return mesh.Mesh{}, fmt.Errorf("at line #%d, got %d values for vec,uv,norm. Expected %d", lineNum, count, expectedCount)
			}

		// COMMENT
		case "#":
		// Do nothing
		case "g":
		// TODO: Support groups
		default:
			// Do nothing - ignore unknown fields (
		}
	}

	// TODO: Split everything below into another function - all the data processing.

	if vertIndices != nil {
		if normalIndices != nil && len(vertIndices) != len(normalIndices) {
			return mesh.Mesh{}, fmt.Errorf("read in vertex and normal indices, but counts don't match: %d vs %d", len(vertIndices), len(normalIndices))
		}
		if textureCoordIndices != nil && len(vertIndices) != len(textureCoordIndices) {
			return mesh.Mesh{}, fmt.Errorf("read in vertex and texture coord indices, but counts don't match: %d vs %d", len(vertIndices), len(textureCoordIndices))
		}
	}

	var vertexBuffer, vertexIndexBuffer, textureCoordBuffer, normalBuffer gl.Buffer

	vertexBuffer = gl.CreateBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
	gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Vec3(binary.LittleEndian, verts...), gl.STATIC_DRAW)

	if len(vertIndices) > 0 {
		vertexIndexBuffer = gl.CreateBuffer()
		gl.BindBuffer(gl.ARRAY_BUFFER, vertexIndexBuffer)
		gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Uint16(binary.LittleEndian, vertIndices...), gl.STATIC_DRAW)
	}

	switch {
	case normalIndices != nil: // Using index buffers - dereference the normal indices and put the values in the buffer.
		normalBuffer = gl.CreateBuffer()
		gl.BindBuffer(gl.ARRAY_BUFFER, normalBuffer)
		normalValues := make([]mgl32.Vec3, len(normalIndices))
		for i, index := range normalIndices {
			if int(index) >= len(normals) {
				return mesh.Mesh{}, fmt.Errorf("unexpected Normal index %d, out of range of the provided %d normals", index+1, len(normals))
			}
			normalValues[i] = normals[index]
		}
		gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Vec3(binary.LittleEndian, normalValues...), gl.STATIC_DRAW)
	case normals != nil: // Basic case - store the values that were read in directly into the buffer.
		normalBuffer = gl.CreateBuffer()
		gl.BindBuffer(gl.ARRAY_BUFFER, normalBuffer)
		gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Vec3(binary.LittleEndian, normals...), gl.STATIC_DRAW)
	default:
		// Nothing to be done - return an uninitialized buffer which must be handled before mesh is rendered.
	}

	switch {
	case textureCoordIndices != nil: // Using index buffers - dereference the texture coordinates and put actual coords in the buffer.
		textureCoordBuffer = gl.CreateBuffer()
		gl.BindBuffer(gl.ARRAY_BUFFER, textureCoordBuffer)
		textureCoordValues := make([]mgl32.Vec2, len(textureCoordIndices))
		for i, index := range textureCoordIndices {
			if int(index) >= len(textureCoords) {
				return mesh.Mesh{}, fmt.Errorf("unexpected Texture Coordinate index %d, out of range of the provided %d texture coordinates", index+1, len(textureCoords))
			}
			textureCoordValues[i] = textureCoords[index]
		}
		gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Vec2(binary.LittleEndian, textureCoordValues...), gl.STATIC_DRAW)
	case textureCoords != nil: // Basic case - store the values that were read in directly into the buffer.
		textureCoordBuffer = gl.CreateBuffer()
		gl.BindBuffer(gl.ARRAY_BUFFER, textureCoordBuffer)
		gl.BufferData(gl.ARRAY_BUFFER, bytecoder.Vec2(binary.LittleEndian, textureCoords...), gl.STATIC_DRAW)
	default:
		// Nothing to be done - return an uninitialized buffer which must be handled before mesh is rendered.
	}

	if glError := gl.GetError(); glError != 0 {
		return mesh.Mesh{}, fmt.Errorf("gl.GetError: %v", glError)
	}

	log.Printf("Vertices: %v\n", verts)
	log.Printf("Normals: %v\n", normals)
	log.Printf("Vertex Indices: %v\n", vertIndices)

	itemCount := len(verts) / 3 // 3 vertices per triangle. // @@@@ TODO: Make sure this works.
	if vertIndices != nil {
		itemCount = len(vertIndices)
	}
	return mesh.NewMesh(vertexBuffer, vertexIndexBuffer, normalBuffer, gl.TRIANGLES, itemCount, nil, gl.Texture{}, textureCoordBuffer), nil
}
