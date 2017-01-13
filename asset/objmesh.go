package asset

import (
	"fmt"
	"math"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/omustardo/gome/model/mesh"
)

// This code is originally based on https://gist.github.com/davemackintosh/67959fa9dfd9018d79a4
// and https://en.wikipedia.org/wiki/Wavefront_.obj_file
// and http://www.opengl-tutorial.org/beginners-tutorials/tutorial-7-model-loading/
//
// Unfortunately, found this list after implementing... I should probably use one of these instead:
// Other Golang OBJ loaders from https://github.com/mmchugh/gomobile-examples/issues/6
//https://github.com/go-qml/qml/blob/v1/examples/gopher/wavefront.go https://github.com/peterhellberg/wavefront/blob/master/wavefront.go
//https://github.com/Stymphalian/go.gl/blob/master/jgl/obj_filereader.go
//https://github.com/tobscher/go-three/blob/master/loaders/obj.go
//https://github.com/adam000/read-obj/tree/master/obj https://github.com/adam000/read-obj/blob/master/mtl/mtl.go
//https://github.com/udhos/negentropia/blob/master/webserv/src/negentropia/world/obj/obj.go
//https://github.com/fogleman/pt/blob/master/pt/obj.go
//https://github.com/luxengine/lux/blob/master/utils/objloader.go
//https://github.com/gmacd/obj/blob/master/obj.go
//https://github.com/gographics/goviewer/blob/master/loader/wavefront.go
//https://github.com/sf1/go3dm
//https://github.com/peterudkmaya11/lux/blob/master/utils/objloader.go

// LoadOBJ creates a mesh from an obj file.
func LoadOBJ(path string) (mesh.Mesh, error) {
	return loadOBJ(path, false)
}

// LoadOBJNormalized creates a mesh from an obj file.
// The loaded OBJ is scaled to be as large as possible while still fitting in a unit sphere.
func LoadOBJNormalized(path string) (mesh.Mesh, error) {
	return loadOBJ(path, true)
}

func loadOBJ(path string, normalize bool) (mesh.Mesh, error) {
	fileData, err := loadFile(path)
	if err != nil {
		return mesh.Mesh{}, err
	}
	verts, normals, textureCoords, err := loadOBJData(fileData)
	if err != nil {
		return mesh.Mesh{}, fmt.Errorf("Error loading %s: %v", path, err)
	}

	if normalize {
		// Normalize input vertices so the input mesh is exactly as large as it can be while still fitting in a unit sphere.
		// This makes scaling meshes relative to each other very easy to think about.
		// TODO: Consider centering meshes when resizing them to avoid empty space making them smaller than necessary.
		maxLength := float32(math.SmallestNonzeroFloat32)
		for _, v := range verts {
			if length := v.Len(); length > maxLength {
				maxLength = length
			}
		}
		for i := range verts {
			verts[i] = verts[i].Mul(1 / maxLength)
		}
	}
	return mesh.NewMeshFromArrays(verts, normals, textureCoords)
}

func loadOBJData(data []byte) (verts, normals []mgl32.Vec3, textureCoords []mgl32.Vec2, err error) {
	lines := strings.Split(string(data), "\n")

	// Indices are used by the OBJ file format to declare full triangles via the 'f'ace tag.
	// All of these indices are converted back to the values that they reference and stored in gl buffers to be returned.
	var vertIndices, uvIndices, normalIndices []uint16

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
			return nil, nil, nil, err
		}
		if count != 1 {
			return nil, nil, nil, fmt.Errorf("at line #%d, unable to get line type: %v", lineNum, err)
		}
		// Trim off the text that has been read.
		line = strings.TrimSpace(line[len(lineType):])

		switch lineType {
		// VERTICES.
		case "v":
			vec := mgl32.Vec3{}
			count, err := fmt.Sscanf(line, "%f %f %f", &vec[0], &vec[1], &vec[2])
			if err != nil {
				return nil, nil, nil, fmt.Errorf("at line #%d, error reading vertices: %v", lineNum, err)
			}
			if count != 3 {
				return nil, nil, nil, fmt.Errorf("at line #%d, got %d values for vertices. Expected 3", lineNum, count)
			}
			verts = append(verts, vec)

		// NORMALS.
		case "vn":
			vec := mgl32.Vec3{}
			count, err := fmt.Sscanf(line, "%f %f %f", &vec[0], &vec[1], &vec[2])
			if err != nil {
				return nil, nil, nil, fmt.Errorf("at line #%d, error reading normals: %v", lineNum, err)
			}
			if count != 3 {
				return nil, nil, nil, fmt.Errorf("at line #%d, got %d values for normals. Expected 3", lineNum, count)
			}
			normals = append(normals, vec)

		// TEXTURE VERTICES.
		case "vt":
			vec := mgl32.Vec2{}
			count, err := fmt.Sscanf(line, "%f %f", &vec[0], &vec[1])
			if err != nil {
				return nil, nil, nil, fmt.Errorf("at line #%d, error reading texture vertices: %v", lineNum, err)
			}
			if count != 2 {
				return nil, nil, nil, fmt.Errorf("at line #%d, got %v values for texture vertices. Expected 2", lineNum, count)
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
				uvIndices = append(uvIndices, uv[0]-1, uv[1]-1, uv[2]-1)
				expectedCount = 6
			case strings.Count(line, "/") == 6:
				count, err = fmt.Sscanf(line, "%d/%d/%d %d/%d/%d %d/%d/%d", &vec[0], &uv[0], &norm[0], &vec[1], &uv[1], &norm[1], &vec[2], &uv[2], &norm[2])
				vertIndices = append(vertIndices, vec[0]-1, vec[1]-1, vec[2]-1)
				uvIndices = append(uvIndices, uv[0]-1, uv[1]-1, uv[2]-1)
				normalIndices = append(normalIndices, norm[0]-1, norm[1]-1, norm[2]-1)
				expectedCount = 9
			default:
				return nil, nil, nil, fmt.Errorf("at line #%d, error reading indices: %v", lineNum, err)
			}
			if err != nil {
				return nil, nil, nil, fmt.Errorf("at line #%d, error reading indices: %v", lineNum, err)
			}
			if count != expectedCount {
				return nil, nil, nil, fmt.Errorf("at line #%d, got %d values for vec,uv,norm. Expected %d", lineNum, count, expectedCount)
			}

		// COMMENT
		case "#":
		// Do nothing
		case "g":
		// TODO: Support groups
		default:
			// Do nothing - ignore unknown fields
		}
	}

	if vertIndices != nil {
		if normalIndices != nil && len(vertIndices) != len(normalIndices) {
			return nil, nil, nil, fmt.Errorf("read in vertex and normal indices, but counts don't match: %d vs %d", len(vertIndices), len(normalIndices))
		}
		if uvIndices != nil && len(vertIndices) != len(uvIndices) {
			return nil, nil, nil, fmt.Errorf("read in vertex and texture coord indices, but counts don't match: %d vs %d", len(vertIndices), len(uvIndices))
		}
	}

	// If vertices were provided with an index buffer, transform it into a list of raw vertices.
	if vertIndices != nil {
		verts, err = indicesToValues(vertIndices, verts)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	if normalIndices != nil {
		normals, err = indicesToValues(normalIndices, normals)
		if err != nil {
			return nil, nil, nil, err
		}
	}
	if uvIndices != nil {
		textureCoordValues := make([]mgl32.Vec2, len(uvIndices))
		for i, index := range uvIndices {
			if int(index) >= len(textureCoords) {
				return nil, nil, nil, fmt.Errorf("unexpected Texture Coordinate index %d, out of range of the provided %d texture coordinates", index+1, len(textureCoords))
			}
			textureCoordValues[i] = textureCoords[index]
		}
		textureCoords = textureCoordValues
	}
	return verts, normals, textureCoords, nil
}

// indicesToValues takes a list of indices and the data they reference, and returns the raw list of referenced data
// with all of the duplicate values that entails.
// Note that the indices are expected to be zero based, even though OBJ files use 1 based indexing.
func indicesToValues(indices []uint16, data []mgl32.Vec3) ([]mgl32.Vec3, error) {
	values := make([]mgl32.Vec3, len(indices))
	for i, index := range indices {
		if int(index) >= len(data) {
			return nil, fmt.Errorf("unexpected index %d, out of range of the provided %d data", index+1, len(data))
		}
		values[i] = data[index]
	}
	return values, nil
}
