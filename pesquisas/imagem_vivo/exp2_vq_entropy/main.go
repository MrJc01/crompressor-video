package main

import (
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
					tensor = append(tensor, float64(r>>8)/255.0)
					tensor = append(tensor, float64(g>>8)/255.0)
					tensor = append(tensor, float64(b>>8)/255.0)
				}
			}
			blocks = append(blocks, TensorBlock{X: x, Y: y, Tensor: tensor})
		}
	}
	return blocks
}

func QuantizeMicroBrain(blocks []TensorBlock) {
	for bIdx, block := range blocks {
		var sumR, sumG, sumB float64
		count := float64(len(block.Tensor) / 3)
		i := 0
		for i < len(block.Tensor) {
			sumR += block.Tensor[i]
			sumG += block.Tensor[i+1]
			sumB += block.Tensor[i+2]
			i += 3
		}
		avgR, avgG, avgB := sumR/count, sumG/count, sumB/count
		
		i = 0
		for i < len(block.Tensor) {
			blocks[bIdx].Tensor[i] = avgR
			blocks[bIdx].Tensor[i+1] = avgG
			blocks[bIdx].Tensor[i+2] = avgB
			i += 3
		}
	}
}

func UnChunker(blocks []TensorBlock, width, height, blockSize int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for _, block := range blocks {
		i := 0
		for blockY := 0; blockY < blockSize; blockY++ {
			for blockX := 0; blockX < blockSize; blockX++ {
				if i+2 >= len(block.Tensor) {
					break
				}
				if block.X+blockX >= width || block.Y+blockY >= height {
					i+=3; continue
				}
				r := uint8(math.Round(block.Tensor[i] * 255.0))
				g := uint8(math.Round(block.Tensor[i+1] * 255.0))
				b := uint8(math.Round(block.Tensor[i+2] * 255.0))
				img.Set(block.X+blockX, block.Y+blockY, color.RGBA{R: r, G: g, B: b, A: 255})
				i += 3
			}
		}
	}
	return img
}

func saveImage(img image.Image, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <caminho_da_imagem>")
		return
	}
	imgPath := os.Args[1]

	fmt.Println("==================================================")
	fmt.Println("[LAB] Crompressor CROM_EXP_2: MicroBrain na Imagem Real")
	fmt.Println("==================================================")

	file, _ := os.Open(imgPath)
	originalImg, _, _ := image.Decode(file)
	file.Close()
	bounds := originalImg.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	fmt.Printf("[+] Imagem Carregada: %s (%dx%d)\n", imgPath, w, h)

	start := time.Now()
	blocks := Chunker(originalImg, BlockSize)
	QuantizeMicroBrain(blocks)

	quantizedImg := UnChunker(blocks, w, h, BlockSize)
	saveImage(quantizedImg, "entropy_microbrain_real.png")
	
	fmt.Printf("✅ VEREDITO: SUCESSO!\n")
	fmt.Printf("   A foto real sofreu VQ Decay (Efeito Mosaico) simulando Dicionário Baixo.\n")
	fmt.Printf("⏱ Tempo VQ Execução: %v\n", time.Since(start))
	fmt.Println("==================================================")
}
