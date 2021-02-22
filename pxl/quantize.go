package pxl

import (
	"encoding/json"
	"errors"
	"github.com/ericpauley/go-quantize/quantize"
	"image"
	"image/color"
)

type SerializablePalette struct {
	Palette []color.Color `json:"values"`
}

func (s *SerializablePalette) UnmarshalJSON(b []byte) error {
	var tmp struct {
		Palette []color.RGBA
	}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	colors := make([]color.Color, len(tmp.Palette))
	for i, v := range tmp.Palette {
		colors[i] = color.Color(v)
	}
	s.Palette = colors
	return nil
}

func (s *SerializablePalette) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(s.Palette, "", "  ")
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
