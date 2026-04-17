package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

// Calcula o Mean Squared Error puro
func MSE(img1, img2 image.Image) (float64, error) {
	bounds := img1.Bounds()
	if bounds != img2.Bounds() {
		return 0, fmt.Errorf("As imagens precisam ter a mesma dimensão")
	}

	var sumSq float64
	totalPixels := float64(bounds.Dx() * bounds.Dy())

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()

			// Converte para 8-bits
			c1r, c1g, c1b := float64(r1>>8), float64(g1>>8), float64(b1>>8)
			c2r, c2g, c2b := float64(r2>>8), float64(g2>>8), float64(b2>>8)

			sumSq += math.Pow(c1r-c2r, 2)
			sumSq += math.Pow(c1g-c2g, 2)
			sumSq += math.Pow(c1b-c2b, 2)
		}
	}

	// Média sobre todos os canais (RGB = 3 canais)
	mse := sumSq / (totalPixels * 3)
	return mse, nil
}

// Calcula o PSNR (Peak Signal to Noise Ratio)
func PSNR(mse float64) float64 {
	if mse == 0 {
		return math.Inf(1) // Imagens identicas, Infinito
	}
	MaxValid := 255.0
	return 10 * math.Log10((MaxValid*MaxValid)/mse)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: go run main.go <imagem_original> <imagem_compressada>")
		return
	}

	file1, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	img1, _, _ := image.Decode(file1)
	file1.Close()

	file2, err := os.Open(os.Args[2])
	if err != nil {
		panic(err)
	}
	img2, _, _ := image.Decode(file2)
	file2.Close()

	mse, err := MSE(img1, img2)
	if err != nil {
		panic(err)
	}

	psnr := PSNR(mse)

	fmt.Println("==================================================")
	fmt.Println("[DEEP PROFILING] CROM_EXP_6.1: PSNR Metrics Engine")
	fmt.Println("==================================================")
	fmt.Printf("[+] Original: %s\n", os.Args[1])
	fmt.Printf("[+] Restringida: %s\n", os.Args[2])
	fmt.Printf("--------------------------------------------------\n")
	fmt.Printf("[>] MSE Absoluto (Ruído Calculado)  : %.6f\n", mse)
	fmt.Printf("[>] PSNR Score (Sinal Pico dB)      : %.2f dB\n", psnr)

	if math.IsInf(psnr, 1) || psnr >= 40.0 {
		fmt.Println("\n✅ APROVADO GRADE ESTÚDIO: A Imagem excede 40dB. Degradação perfeitamente imperceptível para o Olho Humano/Netflix Std.")
	} else if psnr >= 30.0 {
		fmt.Println("\n⚠️ QUALIDADE DE BROADCAST: Qualidade aceitável entre 30dB e 40dB. Artefatos de quantização são ligeiramente perceptíveis.")
	} else {
		fmt.Println("\n❌ REPROVADO VISUALMENTE: PSNR Abaixo de 30dB. Destruição estrutural inaceitável. Modifique O Dicionário!")
	}
	fmt.Println("==================================================")
}
