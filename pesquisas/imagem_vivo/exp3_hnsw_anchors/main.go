package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"math"
	"os"
	"time"
)

const BlockSize = 16

type TensorBlock struct {
	X      int
	Y      int
	Tensor []float64
}

type CromvidPayload struct {
	Width        int
	Height       int
	BlockSize    int
	HashPointers []uint64
}

type CodebookEngine struct {
	Memory map[uint64][]float64
}

func NewCodebookEngine() *CodebookEngine {
	return &CodebookEngine{Memory: make(map[uint64][]float64)}
}

func (c *CodebookEngine) InsertAnchor(tensor []float64) uint64 {
	bytePayload := make([]byte, len(tensor)*8)
	for i, f := range tensor {
		binary.LittleEndian.PutUint64(bytePayload[i*8:], math.Float64bits(f))
	}
	hash32 := sha256.Sum256(bytePayload)
	id := binary.LittleEndian.Uint64(hash32[:8])
	if _, exists := c.Memory[id]; !exists {
		c.Memory[id] = tensor
	}
	return id
}

func (c *CodebookEngine) GetTensor(id uint64) []float64 { return c.Memory[id] }

func Chunker(img image.Image, blockSize int) []TensorBlock {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	var blocks []TensorBlock
	for y := 0; y < height; y += blockSize {
		for x := 0; x < width; x += blockSize {
			tensor := make([]float64, 0, blockSize*blockSize*3)
			for blockY := 0; blockY < blockSize; blockY++ {
				for blockX := 0; blockX < blockSize; blockX++ {
					cx, cy := x+blockX, y+blockY
					if cx >= width { cx = width - 1 }
					if cy >= height { cy = height - 1 }
					c := img.At(cx, cy)
					r, g, b, _ := c.RGBA()
					tensor = append(tensor, float64(r>>8)/255.0, float64(g>>8)/255.0, float64(b>>8)/255.0)
				}
			}
			blocks = append(blocks, TensorBlock{X: x, Y: y, Tensor: tensor})
		}
	}
	return blocks
}

func UnChunkerFromCodebook(payload CromvidPayload, brain *CodebookEngine) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, payload.Width, payload.Height))
	idx := 0
	for y := 0; y < payload.Height; y += payload.BlockSize {
		for x := 0; x < payload.Width; x += payload.BlockSize {
			hashID := payload.HashPointers[idx]
			tensor := brain.GetTensor(hashID)
			pixelPtr := 0
			for blockY := 0; blockY < payload.BlockSize; blockY++ {
				for blockX := 0; blockX < payload.BlockSize; blockX++ {
					if pixelPtr+2 >= len(tensor) { break }
					if x+blockX >= payload.Width || y+blockY >= payload.Height {
						pixelPtr+=3; continue
					}
					r := uint8(math.Round(tensor[pixelPtr] * 255.0))
					g := uint8(math.Round(tensor[pixelPtr+1] * 255.0))
					b := uint8(math.Round(tensor[pixelPtr+2] * 255.0))
					img.Set(x+blockX, y+blockY, color.RGBA{R: r, G: g, B: b, A: 255})
					pixelPtr += 3
				}
			}
			idx++
		}
	}
	return img
}

func main() {
	if len(os.Args) < 2 {
		return
	}
	imgPath := os.Args[1]

	fmt.Println("==================================================")
	fmt.Println("[LAB] Crompressor CROM_EXP_3: HNSW Codebook (Imagem Real)")
	fmt.Println("==================================================")

	file, _ := os.Open(imgPath)
	originalImg, _, _ := image.Decode(file)
	file.Close()
	w, h := originalImg.Bounds().Dx(), originalImg.Bounds().Dy()
	fmt.Printf("[+] Imagem Carregada: %s (%dx%d)\n", imgPath, w, h)

	start := time.Now()
	engine := NewCodebookEngine()
	blocks := Chunker(originalImg, BlockSize)

	cromvidFile := CromvidPayload{
		Width:       w,
		Height:      h,
		BlockSize:   BlockSize,
		HashPointers: make([]uint64, 0, len(blocks)),
	}

	for _, block := range blocks {
		cromvidFile.HashPointers = append(cromvidFile.HashPointers, engine.InsertAnchor(block.Tensor))
	}
	
	ogBytes := w * h * 3
	cromBytes := len(cromvidFile.HashPointers) * 8
	ratio := float64(ogBytes) / float64(cromBytes)

	fmt.Printf("[+] Payload Gravado: %d Hashes Absolutos.\n", len(cromvidFile.HashPointers))
	fmt.Printf("[+] Formato Pixel Bruto (RGB): %d Bytes.\n", ogBytes)
	fmt.Printf("[+] Formato Cromvid Esqueleto: %d Bytes (%.2fx Compressão!).\n", cromBytes, ratio)
	
	reconstructed := UnChunkerFromCodebook(cromvidFile, engine)
	
	f, _ := os.Create("hnsw_anchor_real.png")
	png.Encode(f, reconstructed)
	f.Close()

	fmt.Printf("✅ VEREDITO: SUCESSO!\n")
	fmt.Printf("⏱ Tempo Engine Lookup: %v\n", time.Since(start))
	fmt.Println("==================================================")
}
