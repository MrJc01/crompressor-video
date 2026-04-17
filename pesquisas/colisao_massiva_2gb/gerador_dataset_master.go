package main

import (
	"fmt"
	"os"
	"os/exec"
)

func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func main() {
	fmt.Println("[+] Inicializando Ingestor CROM CC0...")
	os.MkdirAll("../dados/raw", 0755)

	// ============================================
	// VÍDEO (Sintel Trailer CC0)
	// ============================================
	fmt.Println("\n[...] Baixando Vídeo OpenSource (Sintel)")
	if _, err := os.Stat("../dados/raw/sintel.mp4"); os.IsNotExist(err) {
		run("wget", "-qO", "../dados/raw/sintel.mp4", "https://download.blender.org/durian/trailer/sintel_trailer-1080p.mp4")
	}
	
	// Gera loops até atingir 500MB (limitado por frames) e 100MB Unseen
	run("ffmpeg", "-y", "-stream_loop", "15", "-i", "../dados/raw/sintel.mp4", "-f", "rawvideo", "-pix_fmt", "rgb24", "-fs", "500M", "../dados/ouro_video_treino.avi")
	run("ffmpeg", "-y", "-stream_loop", "5", "-i", "../dados/raw/sintel.mp4", "-f", "rawvideo", "-pix_fmt", "rgb24", "-fs", "100M", "../dados/ouro_video_unseen.avi")

	// ============================================
	// ÁUDIO (Música Clássica Beethoven/Mozart CC-PD)
	// ============================================
	fmt.Println("\n[...] Baixando Áudio OpenSource (Beethoven)")
	if _, err := os.Stat("../dados/raw/audio.mp3"); os.IsNotExist(err) {
		run("wget", "-U", "Mozilla/5.0", "-qO", "../dados/raw/audio.mp3", "https://upload.wikimedia.org/wikipedia/commons/1/14/Beethoven_-_Symphony_No._5_in_C_minor_-_I._Allegro_con_brio.ogg")
	}
	
	run("ffmpeg", "-y", "-stream_loop", "30", "-i", "../dados/raw/audio.mp3", "-f", "f32le", "-acodec", "pcm_f32le", "-ac", "1", "-fs", "500M", "../dados/ouro_audio_treino.wav")
	run("ffmpeg", "-y", "-stream_loop", "10", "-i", "../dados/raw/audio.mp3", "-f", "f32le", "-acodec", "pcm_f32le", "-ac", "1", "-fs", "100M", "../dados/ouro_audio_unseen.wav")

	// ============================================
	// IMAGEM (Github Raw Image OpenSource)
	// ============================================
	fmt.Println("\n[...] Baixando Imagem OpenSource (Cat HD)")
	if _, err := os.Stat("../dados/raw/cat.jpg"); os.IsNotExist(err) {
		run("wget", "-U", "Mozilla/5.0", "-qO", "../dados/raw/cat.jpg", "https://raw.githubusercontent.com/image-rs/image/master/tests/images/jpg/progressive/cat.jpg")
	}
	
	run("ffmpeg", "-y", "-stream_loop", "1000", "-i", "../dados/raw/cat.jpg", "-f", "rawvideo", "-pix_fmt", "rgb24", "-fs", "500M", "../dados/ouro_imagem_treino.bmp")
	run("ffmpeg", "-y", "-stream_loop", "200", "-i", "../dados/raw/cat.jpg", "-f", "rawvideo", "-pix_fmt", "rgb24", "-fs", "100M", "../dados/ouro_imagem_unseen.bmp")

	// ============================================
	// ALIEN/TEXTO (Dados Binários Randômicos)
	// ============================================
	fmt.Println("\n[...] Gerando Mídia Alien/Ruído Puro Unseen")
	run("dd", "if=/dev/urandom", "of=../dados/alien_ruido_treino.raw", "bs=1M", "count=500")
	run("dd", "if=/dev/urandom", "of=../dados/alien_ruido_unseen.raw", "bs=1M", "count=100")

	fmt.Println("\n[✔] Ingestão CC0 2.4GB Finalizada em Disco Local!")
}
