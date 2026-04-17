package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"os/exec"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <caminho_do_video>")
		return
	}
	videoPath := os.Args[1]

	fmt.Println("==================================================")
	fmt.Println("[LAB] Crompressor VIDEO_EXP_1: FFmpeg V-RAM Pipe")
	fmt.Println("==================================================")

	// Lógica O(1) SRE - Decodificar Frames Raw Memory via Stdout (Sem gravar PNGs no Disco)
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-f", "image2pipe", "-pix_fmt", "rgba", "-vcodec", "rawvideo", "-")
	
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(fmt.Sprintf("Erro ao abrir pipe: %v", err))
	}
	// Redirecionar stderr ao vazio para nao sujar terminal
	cmd.Stderr = nil

	if err := cmd.Start(); err != nil {
		panic(fmt.Sprintf("Erro ao iniciar FFmpeg: %v", err))
	}

	start := time.Now()
	
	// Precisamos fixar a resolução matematicamente para construir o array de bytes.
	// Em produção, leríamos via API da probe ou cabeçalho do video, no lab sabemos fixo.
	width := 640
	height := 360
	frameSize := width * height * 4 // RGBA = 4 bytes per pixel

	buf := make([]byte, frameSize)
	frameCount := 0

	fmt.Println("[*] Inicializando Absorção de Tubo Raw (Pipeline de Frames)...")
	
	var firstFrame image.Image

	for {
		// ReadFull trava até encher os exatos frameSize Bytes (1 Frame inteiro capturado)
		_, err := io.ReadFull(stdout, buf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			fmt.Printf("Erro de leitura: %v\n", err)
			break
		}

		frameCount++
		
		// Opcional: Provar O(1) Instanciação
		// Alocamos num container nativo de imagem o ByteStream
		if frameCount == 1 {
			img := image.NewRGBA(image.Rect(0, 0, width, height))
			img.Pix = append([]byte(nil), buf...)
			firstFrame = img
		}
	}

	cmd.Wait()

	fmt.Printf("[+] Extração Concluída. Total de Frames Capturados na Memória RAM: %d\n", frameCount)
	
	if firstFrame != nil {
		f, _ := os.Create("frame_0_amostra.png")
		png.Encode(f, firstFrame)
		f.Close()
		fmt.Printf("[+] Primeiro frame renderizado visualmente em disco -> frame_0_amostra.png\n")
	}

	// 10 Quadros num video de 2 seg = 20 Frames lógicos
	if frameCount >= 10 {
		fmt.Printf("✅ VEREDITO: SUCESSO!\n")
		fmt.Printf("   A interface CGO/Pipe está sugando matrizes do video sub-simbolicamente!\n")
	} else {
		fmt.Printf("❌ VEREDITO: FALHA PARCIAL.\n")
	}
	fmt.Printf("⏱ Tempo Engine Lab: %v\n", time.Since(start))
	fmt.Println("==================================================")
}
