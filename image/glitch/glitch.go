package glitch

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"sort"

	"github.com/inovacc/utils/v2/encoding/compression"
)

const (
	pixelSize = 4
	width     = 1920
	height    = 1080
)

type Glitch struct {
	comp *compression.Compress
}

func NewGlitch() *Glitch {
	return &Glitch{
		comp: compression.NewCompress(compression.TypeGzip),
	}
}

func (g *Glitch) BlobToImages(filePath, outputDir string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	compressed, err := g.comp.Compress(data)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	for _, b := range compressed {
		buf.WriteString(fmt.Sprintf("%08b", b))
	}
	binaryStr := buf.String()

	pixelsPerImage := (width / pixelSize) * (height / pixelSize)
	numImages := int(math.Ceil(float64(len(binaryStr)) / float64(pixelsPerImage)))

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	for i := 0; i < numImages; i++ {
		start := i * pixelsPerImage
		end := start + pixelsPerImage
		if end > len(binaryStr) {
			end = len(binaryStr)
		}
		subStr := binaryStr[start:end]
		img := image.NewRGBA(image.Rect(0, 0, width, height))

		idx := 0
		for y := 0; y < height; y += pixelSize {
			for x := 0; x < width; x += pixelSize {
				if idx >= len(subStr) {
					break
				}
				c := color.White
				if subStr[idx] == '1' {
					c = color.Black
				}
				for dy := 0; dy < pixelSize; dy++ {
					for dx := 0; dx < pixelSize; dx++ {
						img.Set(x+dx, y+dy, c)
					}
				}
				idx++
			}
		}

		outFile, err := os.Create(filepath.Join(outputDir, fmt.Sprintf("frame_%04d.png", i)))
		if err != nil {
			return err
		}
		err = png.Encode(outFile, img)
		if cerr := outFile.Close(); cerr != nil && err == nil {
			err = cerr
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *Glitch) ImagesToBlob(imagesPath string, outputPath string) error {
	imageFiles, err := filepath.Glob(imagesPath)
	if err != nil || len(imageFiles) == 0 {
		return errors.New("no image frames found")
	}

	sort.Strings(imageFiles)
	var buf bytes.Buffer

	for _, file := range imageFiles {
		imgFile, err := g.safeOpen(file)
		if err != nil {
			return err
		}

		for y := 0; y < height; y += pixelSize {
			for x := 0; x < width; x += pixelSize {
				clr := imgFile.At(x, y)
				r, g, b, _ := clr.RGBA()
				avg := (r + g + b) / 3
				if avg < 0x8000 {
					buf.WriteByte('1')
				} else {
					buf.WriteByte('0')
				}
			}
		}
	}
	binaryStr := buf.String()

	if len(binaryStr)%8 != 0 {
		binaryStr = binaryStr[:len(binaryStr)-(len(binaryStr)%8)] // trim padding
	}

	data := make([]byte, len(binaryStr)/8)
	for i := 0; i < len(data); i++ {
		slice := binaryStr[i*8 : (i+1)*8]
		var b byte
		if _, err := fmt.Sscanf(slice, "%08b", &b); err != nil {
			return err
		}
		data[i] = b
	}

	decompressed, err := g.comp.Decompress(data)
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, decompressed, 0644)
}

func (g *Glitch) safeOpen(file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}(f)

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("image decode failed for %s: %w", file, err)
	}
	return img, nil
}
