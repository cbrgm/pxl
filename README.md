# Pxl

A small utility written in Go to generate pixelart from PNG images.

## Installation

Use the Makefile in this repository (`make release`) or compile pixelize.go (`go build ./pixelize.go`)

## Usage

`Pxl` takes a PNG image as input and converts it to a pixelized version using different options.

```
Usage: pixelize [command] [options...] <file.png>
	Command: convert
		-o Path to the output file. Default is result.png
		-c Path to the colors.json file
		-l Level of granularity. Default is 8 Bit
		-w Maximum width for image resizing
		-h Maximum height for image resizing
	Command: colors
		-o Path to the output file. Default is colors.json
		-c Number of colors to extract into palette
```

or use `pxl` as go module

```go
import "github.com/cbrgm/pxl/pxl"
// ...
px := pxl.New()

px.SetGranularity(8)
px.SetMaxImageSize(1024, 800)
px.SetColors(palette.Websafe)

usePalette := true
pixelArt := px.Convert(img, usePalette)
// ...
```