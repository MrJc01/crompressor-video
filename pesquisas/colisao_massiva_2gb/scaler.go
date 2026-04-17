package main

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func runFFmpeg(args ...string) error {
	cmd := exec.Command("ffmpeg", args...)
	return cmd.Run()
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("[CROSS-DOMAIN] Fase 8.2: Transcodificador Escalador SRE")
	fmt.Println("==================================================")

	files, err := filepath.Glob("../dados/*.*")
	if err != nil {
		fmt.Println("Erro lendo diretório:", err)
		return
	}

	start := time.Now()

	for _, file := range files {
		ext := filepath.Ext(file)
		base := strings.TrimSuffix(filepath.Base(file), ext)

		// Extrai apenas 3 segundos/3 quadros pra provar o scaler de forma rápida
		fmt.Printf("[+] Processando Escalas para: %s\n", base)

		if ext == ".avi" {
			// VÍDEO
			runFFmpeg("-y", "-i", file, "-t", "3", "-s", "426x240", "-c:v", "rawvideo", fmt.Sprintf("../dados/%s_pequeno%s", base, ext))
			runFFmpeg("-y", "-i", file, "-t", "3", "-s", "1280x720", "-c:v", "rawvideo", fmt.Sprintf("../dados/%s_medio%s", base, ext))
			runFFmpeg("-y", "-i", file, "-t", "3", "-s", "1920x1080", "-c:v", "rawvideo", fmt.Sprintf("../dados/%s_alto%s", base, ext))
		} else if ext == ".wav" {
			// ÁUDIO
			runFFmpeg("-y", "-i", file, "-t", "3", "-ar", "11025", "-c:a", "pcm_f32le", fmt.Sprintf("../dados/%s_pequeno%s", base, ext))
			runFFmpeg("-y", "-i", file, "-t", "3", "-ar", "44100", "-c:a", "pcm_f32le", fmt.Sprintf("../dados/%s_medio%s", base, ext))
			runFFmpeg("-y", "-i", file, "-t", "3", "-ar", "96000", "-c:a", "pcm_f32le", fmt.Sprintf("../dados/%s_alto%s", base, ext))
		} else if ext == ".bmp" {
			// IMAGEM
			runFFmpeg("-y", "-i", file, "-vf", "scale=640:480", "-c:v", "rawvideo", fmt.Sprintf("../dados/%s_pequeno%s", base, ext))
			runFFmpeg("-y", "-i", file, "-vf", "scale=1920:1080", "-c:v", "rawvideo", fmt.Sprintf("../dados/%s_medio%s", base, ext))
			runFFmpeg("-y", "-i", file, "-vf", "scale=3840:2160", "-c:v", "rawvideo", fmt.Sprintf("../dados/%s_alto%s", base, ext))
		}
	}

	fmt.Printf("\n✅ GERAÇÃO DE PERFILAMENTO PEQUENO, MÉDIO E ALTO CUMPRIDO (Amostragens Rápidas SRE).\n")
	fmt.Printf("⏱ Tempo Engine Transcode: %v\n", time.Since(start))
	fmt.Println("==================================================")
}
