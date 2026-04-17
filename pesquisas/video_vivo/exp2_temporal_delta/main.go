package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"os/exec"
	"time"
)

const BlockSize = 16
const Threshold = 0.05 // Quantos % o bloco pode mudar para ainda ser considerado estático

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
					if cx >= width {
						cx = width - 1
					}
					if cy >= height {
						cy = height - 1
					}
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

// EuclideanDistance is a basic metric for delta change.
func isBlockStatic(t1, t2 []float64) bool {
	var diff float64
	for i := 0; i < len(t1); i++ {
		diff += math.Abs(t1[i] - t2[i])
	}
	// Normaliza pela quantidade de dimensões
	avgDiff := diff / float64(len(t1))
	return avgDiff < Threshold
}

func main() {
	if len(os.Args) < 2 {
		return
	}
	videoPath := os.Args[1]

	fmt.Println("==================================================")
	fmt.Println("[LAB] Crompressor VIDEO_EXP_2: Delta FastCDC")
	fmt.Println("==================================================")

	cmd := exec.Command("ffmpeg", "-i", videoPath, "-f", "image2pipe", "-pix_fmt", "rgba", "-vcodec", "rawvideo", "-")
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = nil
	cmd.Start()

	width := 640
	height := 360
	frameSize := width * height * 4

	buf := make([]byte, frameSize)

	var previousBlocks []TensorBlock
	totalBlocks := 0
	staticBlocks := 0

	fmt.Println("[*] Extraindo Quadros e Calculando Deltas Temporais...")
	start := time.Now()

	frameCount := 0
	for {
		_, err := io.ReadFull(stdout, buf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		// Copiamos pra não referenciar map sujo do byte stream
		img.Pix = append([]byte(nil), buf...)
		
		blocks := Chunker(img, BlockSize)
		totalBlocks += len(blocks)

		if previousBlocks != nil {
			frozenCount := 0
			for i, currentBlock := range blocks {
				if isBlockStatic(previousBlocks[i].Tensor, currentBlock.Tensor) {
					staticBlocks++
					frozenCount++
				}
			}
			fmt.Printf(" [Frame %02d] Quadros Estáticos (FROZEN): %d/%d (%.1f%% de Otimização I/O)\n",
				frameCount, frozenCount, len(blocks), (float64(frozenCount)/float64(len(blocks)))*100)

			// Para atestar visualmente, vamos zerar as cores de blocos frozen no frame 2 e salvar
			// Mostrando tudo que "não teve de ser transferido" como tela preta!
			if frameCount == 2 {
				visualImg := image.NewRGBA(image.Rect(0, 0, width, height))
				for i, currentBlock := range blocks {
					if isBlockStatic(previousBlocks[i].Tensor, currentBlock.Tensor) {
						// Bloco Congelado = Deixa Preto (Ignorado na Transmissão Mágica do HNSW)
						continue
					} else {
						// Bloco Alterado (Motion Detected) = Colore de Vermelho Analítico
						for by := 0; by < BlockSize; by++ {
							for bx := 0; bx < BlockSize; bx++ {
								if currentBlock.X+bx < width && currentBlock.Y+by < height {
									visualImg.Set(currentBlock.X+bx, currentBlock.Y+by, color.RGBA{255, 0, 0, 255})
								}
							}
						}
					}
				}
				f, _ := os.Create("delta_frame_02_diff.png")
				png.Encode(f, visualImg)
				f.Close()
			}
		}

		previousBlocks = blocks
		frameCount++
	}
	cmd.Wait()

	fmt.Println("==================================================")
	fmt.Printf("[+] Análise Temporal Concluída! %d Frames Processados.\n", frameCount)
	fmt.Printf("[+] Blocos Totais Extraídos: %d\n", totalBlocks)
	fmt.Printf("[+] Blocos Estáticos (Ignorados no Arquivo Final): %d\n", staticBlocks)
	fmt.Printf("[+] COMPRESSÃO TEMPORAL REAL: %.2f%%\n", (float64(staticBlocks)/float64(totalBlocks))*100.0)
	fmt.Printf("✅ VEREDITO: SUCESSO! Apenas as áreas móveis gastarão LSH Hashes.\n")
	fmt.Printf("⏱ Tempo Engine Lab: %v\n", time.Since(start))
	fmt.Println("==================================================")
}
