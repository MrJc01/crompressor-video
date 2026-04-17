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
					// Previne out of bounds se a imagem não for múltiplo exato de 16
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

func matchLossless(a, b image.Image) bool {
	bounds := a.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := a.At(x, y).RGBA()
			r2, g2, b2, _ := b.At(x, y).RGBA()
			if r1>>8 != r2>>8 || g1>>8 != g2>>8 || b1>>8 != b2>>8 {
				fmt.Printf("Lossless falha no X:%d Y:%d -> Original:[%d %d %d], Reconst:[%d %d %d]\n",
					x, y, r1>>8, g1>>8, b1>>8, r2>>8, g2>>8, b2>>8)
				return false
			}
		}
	}
	return true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <caminho_da_imagem>")
		return
	}
	imgPath := os.Args[1]

	fmt.Println("==================================================")
	fmt.Println("[LAB] Crompressor CROM_EXP_1: Chunker de Imagem REAL")
	fmt.Println("==================================================")

	file, err := os.Open(imgPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	originalImg, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}
	bounds := originalImg.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	fmt.Printf("[+] Imagem Carregada: %s (%dx%d)\n", imgPath, w, h)

	start := time.Now()
	blocks := Chunker(originalImg, BlockSize)
	fmt.Printf("[+] Extraídos %d Blocos Independentes Multi-Dimensional.\n", len(blocks))

	reconstructedImg := UnChunker(blocks, w, h, BlockSize)
	saveImage(reconstructedImg, "reconstrutiva_real_test.png")
	fmt.Println("[+] Imagem Costurada -> reconstrutiva_real_test.png")

	lossless := matchLossless(originalImg, reconstructedImg)
	if lossless {
		fmt.Printf("✅ VEREDITO: SUCESSO! (100%% Lossless na Imagem Real)\n")
	} else {
		fmt.Printf("❌ VEREDITO: FALHA.\n")
	}
	fmt.Printf("⏱ Tempo Lab Execução: %v\n", time.Since(start))
	fmt.Println("==================================================")
}
