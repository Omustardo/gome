package asset

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

func LoadFont(path string) (*truetype.Font, error) {
	data, err := LoadFile(path)
	if err != nil {
		return nil, err
	}
	return freetype.ParseFont(data)
}
