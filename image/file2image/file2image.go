package file2image

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash/crc32"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/divan/txqr"
	"github.com/divan/txqr/qr"
	"github.com/inovacc/utils/v2/encoding/encoder/gob"
	"github.com/inovacc/utils/v2/image/qrcode"
	"gocv.io/x/gocv"
)

var countStart, countFragment, countEnd int

type mimeType int

func (m mimeType) String() string {
	switch m {
	case MimTypeText:
		return "text/plain"
	case MimTypeImageJPEG:
		return "image/jpeg"
	case MimTypeImagePNG:
		return "image/png"
	case MimTypeImageGIF:
		return "image/gif"
	case MimTypeImageTIFF:
		return "image/tiff"
	case MimTypeImageBMP:
		return "image/bmp"
	case MimTypeImageSVG:
		return "image/svg+xml"
	case MimTypeImageWEBP:
		return "image/webp"
	case MimTypeAudio:
		return "audio/*"
	case MimTypeVideo:
		return "video/*"
	case MimTypeApplication:
		return "application/octet-stream"
	case MimTypeMessage:
		return "message/rfc822"
	case MimTypeFont:
		return "font/*"
	case MimTypeArchive:
		return "application/x-archive"
	case MimTypeDatabase:
		return "application/x-sqlite3"
	case MimTypeConfiguration:
		return "text/x-config"
	case MimTypeDesktopEntry:
		return "application/x-desktop"
	case MimTypeDirectory:
		return "inode/directory"
	case MimTypeSocket:
		return "inode/socket"
	case MimTypeSymbolicLink:
		return "inode/symlink"
	case MimTypeWhiteout:
		return "inode/whiteout"
	case MimTypeBlockSpecial:
		return "inode/blockdevice"
	case MimTypeCharacterSpecial:
		return "inode/chardevice"
	case MimTypeFifo:
		return "inode/fifo"
	case MimTypeNamedPipe:
		return "inode/named-pipe"
	default:
		return fmt.Sprintf("unknown(%d)", int(m))
	}
}

const (
	MimTypeText mimeType = iota
	MimTypeImageJPEG
	MimTypeImagePNG
	MimTypeImageGIF
	MimTypeImageTIFF
	MimTypeImageBMP
	MimTypeImageSVG
	MimTypeImageWEBP
	MimTypeAudio
	MimTypeVideo
	MimTypeApplication
	MimTypeMessage
	MimTypeFont
	MimTypeArchive
	MimTypeDatabase
	MimTypeConfiguration
	MimTypeDesktopEntry
	MimTypeDirectory
	MimTypeSocket
	MimTypeSymbolicLink
	MimTypeWhiteout
	MimTypeBlockSpecial
	MimTypeCharacterSpecial
	MimTypeFifo
	MimTypeNamedPipe
)

type kindOp int

func (k kindOp) String() string {
	switch k {
	case KindOperationStart:
		return "START"
	case KindOperationFragment:
		return "FRAGMENT"
	case KindOperationEnd:
		return "END"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", int(k))
	}
}

const (
	KindOperationStart kindOp = iota + 1
	KindOperationFragment
	KindOperationEnd
)

type chunk struct {
	data []byte
	crc  uint32
	idx  int
	kind kindOp
}

func splitFile(filename string, size int) (chan chunk, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	dataChan := make(chan chunk, 1)
	go func() {
		defer func() {
			close(dataChan)
			if err := file.Close(); err != nil {
				log.Println("Error closing file:", err)
			}
		}()

		buf := make([]byte, size)
		idx := 0
		for {
			n, err := file.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Println("Error reading file:", err)
				}
				break
			}
			dataChan <- chunk{
				data: buf[:n],
				crc:  crc32.ChecksumIEEE(buf[:n]),
				idx:  idx,
				kind: getKind(idx),
			}
			idx++
		}

		dataChan <- chunk{
			data: []byte(""),
			crc:  0,
			idx:  idx,
			kind: KindOperationEnd,
		}
	}()

	return dataChan, nil
}

func hash(data []byte) []byte {
	s := sha256.Sum256(data)
	return s[:]
}

func getKind(idx int) kindOp {
	switch idx {
	case 0:
		return KindOperationStart
	default:
		return KindOperationFragment
	}
}

func getMimeType(filename string) mimeType {
	switch ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), ".")); ext {
	case "txt":
		return MimTypeText
	case "jpg", "jpeg":
		return MimTypeImageJPEG
	case "png":
		return MimTypeImagePNG
	case "gif":
		return MimTypeImageGIF
	case "tiff":
		return MimTypeImageTIFF
	case "bmp":
		return MimTypeImageBMP
	case "svg":
		return MimTypeImageSVG
	case "webp":
		return MimTypeImageWEBP
	case "mp3", "wav", "ogg":
		return MimTypeAudio
	case "mp4", "avi", "mov":
		return MimTypeVideo
	case "zip", "rar", "tar", "gz":
		return MimTypeArchive
	case "db", "sqlite":
		return MimTypeDatabase
	case "json", "yml", "ini":
		return MimTypeConfiguration
	default:
		return MimTypeApplication
	}
}

type Metadata struct {
	Kind  kindOp   `qrf:"kind"`
	Name  string   `qrf:"name"`
	Hash  []byte   `qrf:"hash"`
	Crc   uint32   `qrf:"crc"`
	Index string   `qrf:"index"`
	Mime  mimeType `qrf:"mime"`
}

type File2Image struct {
	Data chan []byte
}

func NewChunks(filename string, size int) (*File2Image, error) {
	dataChan, err := splitFile(filename, size)
	if err != nil {
		return nil, err
	}

	stream := &File2Image{
		Data: make(chan []byte),
	}

	go func() {
		defer close(stream.Data)
		seen := make(map[string]bool)

		for obj := range dataChan {
			chunkHash := sha256.Sum256(obj.data)
			hashKey := fmt.Sprintf("%x", chunkHash)
			if seen[hashKey] {
				log.Printf("🔁 Chunk duplicado ignorado (idx %d)", obj.idx)
				continue
			}
			seen[hashKey] = true

			meta := &Metadata{
				Crc:   obj.crc,
				Index: fmt.Sprintf("%06X", obj.idx),
				Mime:  getMimeType(filename),
				Kind:  obj.kind,
			}

			// Start operation
			if meta.Kind == KindOperationStart {
				meta.Name = filepath.Base(filename)
				meta.Hash = hash(obj.data)
			}

			// Fragment and End cleanup
			if meta.Kind != KindOperationStart {
				meta.Name = ""
				meta.Hash = nil
			}

			encodeGob, err := gob.EncodeGob(meta)
			if err != nil {
				fmt.Println("Error encoding metadata:", err)
				continue
			}
			stream.Data <- encodeGob
		}
	}()

	return stream, nil
}

func (q *File2Image) GenerateFramesLive() error {
	window := gocv.NewWindow("QR Dinámico")
	defer window.Close()

	qr := qrcode.NewQrcode()

	log.Println("Preparando transmisión...")

	for obj := range q.Data {
		meta := &Metadata{}
		if err := gob.DecodeGob(obj, meta); err != nil {
			log.Printf("Error decoding metadata: %v", err)
			continue
		}

		switch meta.Kind {
		case KindOperationStart:
			countStart++
		case KindOperationFragment:
			countFragment++
		case KindOperationEnd:
			countEnd++
		}

		log.Printf("Metadata: Kind=%v, Index=%s, CRC=%08X, Name=%s, Mime=%v",
			meta.Kind, meta.Index, meta.Crc, meta.Name, meta.Mime)

		if err := qr.GenerateRaw(obj); err != nil {
			return fmt.Errorf("QR generate: %v", err)
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, qr.Image(700)); err != nil {
			return fmt.Errorf("png encode: %v", err)
		}

		mat, err := gocv.IMDecode(buf.Bytes(), gocv.IMReadGrayScale)
		if err != nil {
			return fmt.Errorf("decode to mat: %v", err)
		}

		if err := window.IMShow(mat); err != nil {
			return err
		}

		if window.WaitKey(100) == 113 { // tecla 'q'
			mat.Close()
			break
		}
		mat.Close()
	}

	log.Println("======= RESUMEN DE FRAMES =======")
	log.Printf("START    : %d", countStart)
	log.Printf("FRAGMENT : %d", countFragment)
	log.Printf("END      : %d", countEnd)
	log.Println("================================")

	summary := fmt.Sprintf("RESUMEN DE FRAMES\nSTART:%d\nFRAGMENT:%d\nEND:%d",
		countStart, countFragment, countEnd)

	qrFinal := qrcode.NewQrcode()
	if err := qrFinal.Generate(summary); err != nil {
		return fmt.Errorf("QR resumen: %v", err)
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, qrFinal.Image(700)); err != nil {
		return fmt.Errorf("png resumen: %v", err)
	}

	mat, err := gocv.IMDecode(buf.Bytes(), gocv.IMReadGrayScale)
	if err != nil {
		return fmt.Errorf("decode resumen: %v", err)
	}

	log.Println("Mostrando resumen final...")
	window.IMShow(mat)
	window.WaitKey(0)
	mat.Close()

	return nil
}

func ShowTxqrSequence(frames [][]byte) error {
	window := gocv.NewWindow("QR Dinámico TXQR")
	defer window.Close()

	encoder := txqr.NewEncoder(512) // ajusta tamaño del chunk
	for _, frame := range frames {
		hexData := fmt.Sprintf("%x", frame)
		qrChunks, err := encoder.Encode(hexData)
		if err != nil {
			return fmt.Errorf("txqr encode: %v", err)
		}

		for _, chunk := range qrChunks {
			img, err := qr.Encode(chunk, 500, qr.Medium)
			if err != nil {
				return fmt.Errorf("txqr encode img: %v", err)
			}

			var buf bytes.Buffer
			if err := png.Encode(&buf, img); err != nil {
				return fmt.Errorf("png encode: %v", err)
			}

			mat, err := gocv.IMDecode(buf.Bytes(), gocv.IMReadGrayScale)
			if err != nil {
				return fmt.Errorf("decode to mat: %v", err)
			}

			window.IMShow(mat)
			if window.WaitKey(120) == 113 {
				return nil
			}
			mat.Close()
		}
	}
	return nil
}
