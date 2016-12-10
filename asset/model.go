package asset

import (
	"bytes"
	"fmt"
	"io"

	"github.com/go-gl/mathgl/mgl32"
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
