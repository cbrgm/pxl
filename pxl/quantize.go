package pxl

import (
	"encoding/json"
	"errors"
	"github.com/ericpauley/go-quantize/quantize"
	"image"
	"image/color"
)

type SerializablePalette struct {
	Palette []color.Color
}

func (s *SerializablePalette) UnmarshalJSON(b []byte) error {
	var palette []color.RGBA
	err := json.Unmarshal(b, &palette)
	if err != nil {
		return err
	}
	colors := make([]color.Color, len(palette))
	for i, v := range palette {
		colors[i] = color.Color(v)
	}
	s.Palette = colors
	return nil
}

func (s *SerializablePalette) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Palette)
}

func GetColorsFromImage(img image.Image, count int) (*SerializablePalette, error) {
	if count <= 0 || img == nil {
		return nil, errors.New("colors amount or image file is invalid")
	}
	p := make([]color.Color, 0, count)
	q := quantize.MedianCutQuantizer{}
	colors := q.Quantize(p, img)
	return &SerializablePalette{
		Palette: colors,
	}, nil
}
