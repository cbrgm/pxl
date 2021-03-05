package pxl

import (
	"golang.org/x/image/draw"
	"image"
	"image/color"
	"image/color/palette"
)

const scalingFactor = 0.01

const (
	Bit8   = 8 * scalingFactor
	Bit16  = 16 * scalingFactor
	Bit32  = 32 * scalingFactor
	Bit64  = 64 * scalingFactor
	Bit128 = 128 * scalingFactor
)

type Pixelizer struct {
	palette     color.Palette
	granularity float64
	maxWidth    int
	maxHeight   int
}

type Config struct {
	Palette     color.Palette
	Granularity int
	MaxWidth    int
	MaxHeight   int
}

func New() *Pixelizer {
	return &Pixelizer{
		palette:     palette.Plan9,
		granularity: Bit8,
		maxHeight:   0,
		maxWidth:    0,
	}
}

func NewFromConfig(c *Config) *Pixelizer {
	px := New()
	px.SetColors(c.Palette)
	px.SetMaxImageSize(c.MaxWidth, c.MaxHeight)
	px.SetGranularity(c.Granularity)
	return px
}

// SetGranularity sets the Pixelizer Granularity level.
// s must be a positive int value in range 0 - 128.
func (p *Pixelizer) SetGranularity(s int) {
	if isValidGranularity(s) {
		p.granularity = granularity(s)
	} else {
		p.granularity = Bit8
	}
}

func isValidGranularity(s int) bool {
	return s > 0 && s <= 128
}

func granularity(s int) float64 {
	return float64(s) * scalingFactor
}

func (p *Pixelizer) SetColors(cp color.Palette) {
	if cp == nil || len(cp) == 0 {
		return
	}
	p.palette = cp
}

func (p *Pixelizer) SetMaxImageSize(maxWidth, maxHeight int) {
	if isValidImageSize(maxWidth, maxHeight) {
		p.maxWidth = maxWidth
		p.maxHeight = maxHeight
	}
}

func isValidImageSize(width, height int) bool {
	if width <= 0 || height <= 0 {
		return false
	}
	return true
}

// Convert creates a pixelized version for a given img
func (p *Pixelizer) Convert(img image.Image, colorize bool) image.Image {
	bounds := img.Bounds()

	imgWidth, imgHeight := resizeImageBounds(p.maxWidth, p.maxHeight, bounds.Dx(), bounds.Dy())

	scaledWidth := float64(imgWidth) * p.granularity
	scaledHeight := float64(imgHeight) * p.granularity

	destination := image.Rect(0, 0, int(scaledWidth), int(scaledHeight))
	pixelImg := scaleTo(img, destination, draw.NearestNeighbor)

	var result image.Image
	{
		if colorize {
			bounds = pixelImg.Bounds()
			coloredPixelImg := image.NewPaletted(bounds, p.palette)
			draw.Draw(coloredPixelImg, coloredPixelImg.Rect, pixelImg, bounds.Min, draw.Over)
			result = coloredPixelImg
		} else {
			result = pixelImg
		}
	}
	destination = image.Rect(0, 0, imgWidth, imgHeight)
	return scaleTo(result, destination, draw.NearestNeighbor)
}

// scaleTo scales src image to the given rect size using granularity as a scaling method
func scaleTo(src image.Image, rect image.Rectangle, scale draw.Scaler) image.Image {
	result := image.NewRGBA(rect)
	scale.Scale(result, rect, src, src.Bounds(), draw.Over, nil)
	return result
}

// resizeImageBounds resizes the image.Rect bounds proportionally to MaxWidth or MaxHeight.
//
// If MaxWidth and MaxHeight is lesser or equal to 0, the original dx and dy values are returned.
// If MaxWidth is greater than dx, MaxWidth will be used for resizing.
// If MaxHeight is greater  than dy, MaxHeight will be used for resizing and has priority over the first condition.
func resizeImageBounds(maxWidth, maxHeight, dx, dy int) (int, int) {
	var ratio = 1.0
	if !isValidImageSize(maxWidth, maxHeight) {
		return dx, dy
	}
	if dx > maxWidth {
		ratio = float64(maxWidth) / float64(dx)
	}
	if dy > maxHeight {
		ratio = float64(maxHeight) / float64(dy)
	}
	x := float64(dx) * ratio
	y := float64(dy) * ratio
	return int(x), int(y)
}
