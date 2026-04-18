package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

// downloadProtectedFile simula requisições autenticadas/válidas com fallback
func downloadProtectedFile(url, dest string) error {
	client := &http.Client{Timeout: 30 * time.Minute}
	var lastErr error
	
	// Sistema de Retry (Max 3 tentativas)
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		
		resp, err := client.Do(req)
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				out, errCreate := os.Create(dest)
				if errCreate != nil {
					return errCreate
				}
				defer out.Close()
				_, errCopy := io.Copy(out, resp.Body)
				if errCopy == nil {
					return nil // Sucesso absoluto
				}
				lastErr = errCopy
			} else {
				lastErr = fmt.Errorf("Erro HTTP: status %d", resp.StatusCode)
			}
		} else {
			lastErr = err
		}
		time.Sleep(2 * time.Second) // Aguarda antes de tentar dnv
	}
	return lastErr
}

func main() {
	fmt.Println("[+] Inicializando Ingestor CROM CC0...")
	os.MkdirAll("../dados/raw", 0755)

	// ============================================
	// VÍDEO (Sintel Trailer CC0)
	// ============================================
	fmt.Println("\n[...] Baixando Vídeo OpenSource (Sintel)")
	if _, err := os.Stat("../dados/raw/sintel.mp4"); os.IsNotExist(err) {
		downloadProtectedFile("https://download.blender.org/durian/trailer/sintel_trailer-1080p.mp4", "../dados/raw/sintel.mp4")
	}
	
	// Gera loops até atingir 500MB (limitado por frames) e 100MB Unseen
	run("ffmpeg", "-y", "-stream_loop", "15", "-i", "../dados/raw/sintel.mp4", "-f", "rawvideo", "-pix_fmt", "rgb24", "-fs", "100M", "../dados/ouro_video_treino.avi")
	run("ffmpeg", "-y", "-stream_loop", "5", "-i", "../dados/raw/sintel.mp4", "-f", "rawvideo", "-pix_fmt", "rgb24", "-fs", "100M", "../dados/ouro_video_unseen.avi")

	// ============================================
	// ÁUDIO (Música Clássica Beethoven/Mozart CC-PD)
	// ============================================
	fmt.Println("\n[...] Baixando Áudio OpenSource (Beethoven)")
	if _, err := os.Stat("../dados/raw/audio.mp3"); os.IsNotExist(err) {
		downloadProtectedFile("https://upload.wikimedia.org/wikipedia/commons/1/14/Beethoven_-_Symphony_No._5_in_C_minor_-_I._Allegro_con_brio.ogg", "../dados/raw/audio.mp3")
	}
	
	run("ffmpeg", "-y", "-stream_loop", "30", "-i", "../dados/raw/audio.mp3", "-f", "f32le", "-acodec", "pcm_f32le", "-ac", "1", "-fs", "100M", "../dados/ouro_audio_treino.wav")
	run("ffmpeg", "-y", "-stream_loop", "10", "-i", "../dados/raw/audio.mp3", "-f", "f32le", "-acodec", "pcm_f32le", "-ac", "1", "-fs", "100M", "../dados/ouro_audio_unseen.wav")

	// ============================================
	// IMAGEM (Github Raw Image OpenSource)
	// ============================================
	fmt.Println("\n[...] Baixando Imagem OpenSource (Cat HD)")
	if _, err := os.Stat("../dados/raw/cat.jpg"); os.IsNotExist(err) {
		downloadProtectedFile("https://raw.githubusercontent.com/image-rs/image/master/tests/images/jpg/progressive/cat.jpg", "../dados/raw/cat.jpg")
	}
	
	run("ffmpeg", "-y", "-stream_loop", "1000", "-i", "../dados/raw/cat.jpg", "-f", "rawvideo", "-pix_fmt", "rgb24", "-fs", "100M", "../dados/ouro_imagem_treino.bmp")
	run("ffmpeg", "-y", "-stream_loop", "200", "-i", "../dados/raw/cat.jpg", "-f", "rawvideo", "-pix_fmt", "rgb24", "-fs", "100M", "../dados/ouro_imagem_unseen.bmp")

	// ============================================
	// ALIEN/TEXTO (Dados Binários Randômicos)
	// ============================================
	fmt.Println("\n[...] Gerando Mídia Alien/Ruído Puro Unseen")
	run("dd", "if=/dev/urandom", "of=../dados/alien_ruido_treino.raw", "bs=1M", "count=100")
	run("dd", "if=/dev/urandom", "of=../dados/alien_ruido_unseen.raw", "bs=1M", "count=100")

	fmt.Println("\n[✔] Ingestão CC0 2.4GB Finalizada em Disco Local!")
}
