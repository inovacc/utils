package glitch

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/inovacc/utils/v2/encoding/compression"
)

const (
	brightnessThreshold = 0x8000
	pixelSize           = 4
	width               = 1920
	height              = 1080
)

type FileMetadata struct {
	Name string
	Date int64
	Hash [32]byte
}

type Glitch struct {
	comp *compression.Compress
}

func NewGlitch() *Glitch {
	return &Glitch{
		comp: compression.NewCompress(compression.TypeGzip),
	}
}

func (g *Glitch) ImagesToBlob(imagesPath, outputDir string) (*FileMetadata, error) {
	if err := ensureDir(outputDir); err != nil {
		return nil, err
	}

	binaryStr, err := g.extractBinaryString(imagesPath)
	if err != nil {
		return nil, err
	}

	data, err := g.decodeBinaryPayload(binaryStr)
	if err != nil {
		return nil, err
	}

	meta, fileData, err := g.parseAndDecompress(data, g.comp)
	if err != nil {
		return nil, err
	}

	if err := g.validateHash(fileData, meta.Hash); err != nil {
		return nil, err
	}

	if err := os.WriteFile(filepath.Join(outputDir, meta.Name), fileData, 0644); err != nil {
		return meta, err
	}
	return meta, nil
}

func (g *Glitch) BlobToImages(filename, outputDir string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	meta := FileMetadata{
		Name: filepath.Base(filename),
		Date: time.Now().Unix(),
		Hash: sha256.Sum256(data),
	}

	metaBuf := new(bytes.Buffer)
	if err := binary.Write(metaBuf, binary.BigEndian, int32(len(meta.Name))); err != nil {
		return err
	}

	if _, err := metaBuf.WriteString(meta.Name); err != nil {
		return err
	}

	if err := binary.Write(metaBuf, binary.BigEndian, meta.Date); err != nil {
		return err
	}

	if _, err := metaBuf.Write(meta.Hash[:]); err != nil {
		return err
	}

	compressed, err := g.comp.Compress(data)
	if err != nil {
		return err
	}

	fullPayload := append(metaBuf.Bytes(), compressed...)

	bitLen := uint32(len(fullPayload) * 8)
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, bitLen)
	prefixed := append(lenBuf, fullPayload...)

	var buf bytes.Buffer
	for _, b := range prefixed {
		buf.WriteString(fmt.Sprintf("%08b", b))
	}
	binaryStr := buf.String()

	pixelsPerImage := (width / pixelSize) * (height / pixelSize)
	numImages := int(math.Ceil(float64(len(binaryStr)) / float64(pixelsPerImage)))

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, numImages)

	for i := 0; i < numImages; i++ {
		i := i
		wg.Add(1)
		go g.worker(i, binaryStr, outputDir, pixelsPerImage, errCh, &wg)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return <-errCh
	}
	return nil
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

func (g *Glitch) worker(i int, binaryStr, outputDir string, pixelsPerImage int, errCh chan error, wg *sync.WaitGroup) {
	defer wg.Done()

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

	outPath := filepath.Join(outputDir, fmt.Sprintf("frame_%04d.png", i))
	outFile, err := os.Create(outPath)
	if err != nil {
		errCh <- err
		return
	}
	err = png.Encode(outFile, img)
	if cerr := outFile.Close(); cerr != nil && err == nil {
		err = cerr
	}
	if err != nil {
		errCh <- err
	}
}

func ensureDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func (g *Glitch) extractBinaryString(imagesPath string) (string, error) {
	files, err := filepath.Glob(imagesPath)
	if err != nil || len(files) == 0 {
		return "", errors.New("no image frames found")
	}
	sort.Strings(files)

	var buf bytes.Buffer
	for _, file := range files {
		img, err := g.safeOpen(file)
		if err != nil {
			return "", err
		}
		for y := 0; y < height; y += pixelSize {
			for x := 0; x < width; x += pixelSize {
				r, g, b, _ := img.At(x, y).RGBA()
				avg := (r + g + b) / 3
				if avg < brightnessThreshold {
					buf.WriteByte('1')
				} else {
					buf.WriteByte('0')
				}
			}
		}
	}
	return g.trimPaddedBits(buf.String())
}

func (g *Glitch) trimPaddedBits(binaryStr string) (string, error) {
	if len(binaryStr) < 32 {
		return "", errors.New("binary string too short to contain length header")
	}

	var totalBits uint32
	if _, err := fmt.Sscanf(binaryStr[:32], "%032b", &totalBits); err != nil {
		return "", fmt.Errorf("failed to parse bit length prefix: %w", err)
	}

	if len(binaryStr) < int(32+totalBits) {
		return "", fmt.Errorf("binary string is shorter than declared bit length: expected %d bits, got %d", totalBits, len(binaryStr)-32)
	}

	trimmed := binaryStr[32 : 32+totalBits]
	if len(trimmed)%8 != 0 {
		return "", errors.New("trimmed binary string is not byte aligned")
	}
	return trimmed, nil
}

func (g *Glitch) decodeBinaryPayload(binaryStr string) ([]byte, error) {
	data := make([]byte, len(binaryStr)/8)
	for i := range data {
		if _, err := fmt.Sscanf(binaryStr[i*8:(i+1)*8], "%08b", &data[i]); err != nil {
			return nil, fmt.Errorf("binary decode failed at index %d: %w", i, err)
		}
	}
	return data, nil
}

func (g *Glitch) parseAndDecompress(data []byte, comp *compression.Compress) (*FileMetadata, []byte, error) {
	meta := &FileMetadata{}
	r := bytes.NewReader(data)

	var nameLen int32
	if err := binary.Read(r, binary.BigEndian, &nameLen); err != nil {
		return meta, nil, fmt.Errorf("failed to read name length: %w", err)
	}

	nameBytes := make([]byte, nameLen)
	if _, err := io.ReadFull(r, nameBytes); err != nil {
		return meta, nil, fmt.Errorf("failed to read name: %w", err)
	}
	meta.Name = string(nameBytes)

	if err := binary.Read(r, binary.BigEndian, &meta.Date); err != nil {
		return meta, nil, fmt.Errorf("failed to read timestamp: %w", err)
	}

	if _, err := io.ReadFull(r, meta.Hash[:]); err != nil {
		return meta, nil, fmt.Errorf("failed to read hash: %w", err)
	}

	compressedData, err := io.ReadAll(r)
	if err != nil {
		return meta, nil, fmt.Errorf("failed to read file content: %w", err)
	}

	decompressed, err := comp.Decompress(compressedData)
	if err != nil {
		return meta, nil, fmt.Errorf("failed to decompress: %w", err)
	}

	return meta, decompressed, nil
}

func (g *Glitch) validateHash(data []byte, expected [32]byte) error {
	actual := sha256.Sum256(data)
	if actual != expected {
		return fmt.Errorf("hash mismatch: file may be corrupted")
	}
	return nil
}
