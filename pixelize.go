package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cbrgm/pxl/pxl"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	cmdConvert = flag.NewFlagSet("convert", flag.ExitOnError)
	cmdColors  = flag.NewFlagSet("colors", flag.ExitOnError)

	cmdConvertOutputOpt = cmdConvert.String("o", "result.png", "")
	cmdConvertColorsOpt = cmdConvert.String("c", "", "")
	cmdConvertLevelOpt  = cmdConvert.Int("l", 8, "")
	cmdConvertWidthOpt  = cmdConvert.Int("w", 0, "")
	cmdConvertHeightOpt = cmdConvert.Int("h", 0, "")

	cmdColorsOutputOpt = cmdColors.String("o", "colors.json", "")
	cmdColorsColorsOpt = cmdColors.Int("c", 48, "")
)

var usage = `Usage: pixelize [command] [options...] file.png
	Command: convert
		-o Path to the output file. Default is result.png
		-c Path to the colors.json file
		-l Level of granularity. Default is 8 Bit
		-w Maximum width for image resizing
		-h Maximum height for image resizing
	Command: colors
		-o Path to the output file. Default is colors.json
		-c Number of colors to extract into palette
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	if len(os.Args) < 2 {
		usageAndExit("no command selected")
	}
	switch os.Args[1] {
	case "convert":
		cmdConvert.Parse(os.Args[2:])
		convertCmd()
	case "colors":
		cmdColors.Parse(os.Args[2:])
		colorsCmd()
	}
}

func convertCmd() {
	f := cmdConvert.Arg(1)
	o := *cmdConvertOutputOpt
	c := *cmdConvertColorsOpt
	l := *cmdConvertLevelOpt
	w := *cmdConvertWidthOpt
	h := *cmdConvertHeightOpt

	fmt.Println(f)

	img, err := loadImage(f)
	if err != nil {
		usageAndExit("unable to load image file")
	}

	px := pxl.New()
	px.SetGranularity(l)
	px.SetMaxImageSize(w, h)

	colorize := false

	if c != "" {
		colors, err := loadPalette(c)
		if err != nil {
			usageAndExit("unable to load color palette")
		}
		px.SetColors(colors.Palette)
		colorize = true
	}

	res := px.Convert(img, colorize)
	err = saveImage(res, o)
	if err != nil {
		usageAndExit("unable to save output file")
	}
}

func colorsCmd() {
	f := cmdColors.Arg(1)
	o := *cmdColorsOutputOpt
	c := *cmdColorsColorsOpt

	img, err := loadImage(f)
	if err != nil {
		usageAndExit("unable to load image file")
	}

	colors, err := pxl.GetColorsFromImage(img, c)
	if err != nil {
		usageAndExit("unable to extract colors from image file")
	}

	err = savePalette(colors, o)
	if err != nil {
		usageAndExit("unable to save output file")
	}
}

func saveImage(img image.Image, path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	imageFile, err := os.Create(absPath)
	defer imageFile.Close()
	if err != nil {
		return err
	}

	err = png.Encode(imageFile, img)
	if err != nil {
		return err
	}
	return nil
}

func loadImage(path string) (image.Image, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	imageFile, err := os.Open(absPath)
	defer imageFile.Close()
	if err != nil {
		return nil, err
	}

	imgData, imgType, err := image.Decode(imageFile)
	if err != nil || imgType != "png" {
		return nil, err
	}
	return imgData, err
}

func loadPalette(path string) (*pxl.SerializablePalette, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	p := pxl.SerializablePalette{}
	err = json.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}
	return &p, err
}

func savePalette(p *pxl.SerializablePalette, path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	f, _ := os.OpenFile(absPath, os.O_CREATE, os.ModePerm)
	defer f.Close()

	e := json.NewEncoder(f)
	return e.Encode(p)
}

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}
