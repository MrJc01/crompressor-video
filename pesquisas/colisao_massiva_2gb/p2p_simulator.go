package main

import (
	"encoding/binary"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"math"
	"net"
	"os"
	"os/exec"
)

type P2PBrain struct {
	Memory map[uint64][]float64
}

func (c *P2PBrain) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewDecoder(f).Decode(&c.Memory)
}

func (c *P2PBrain) MatchForced(target []float64) uint64 {
	bestDist := math.MaxFloat64
	var bestID uint64
	for id, dbTensor := range c.Memory {
		limit := len(target)
		if len(dbTensor) < limit {
			limit = len(dbTensor)
		}
		var diff float64
		for i := 0; i < limit; i++ {
			diff += math.Abs(target[i] - dbTensor[i])
		}
		avg := diff / float64(limit)
		if avg < bestDist {
			bestDist = avg
			bestID = id
		}
	}
	return bestID
}

func runServer(brainPath string, videoPath string) {
	fmt.Println("[EMISSOR A] Carregando Cérebro Híbrido...")
	var brain P2PBrain
	brain.Load(brainPath)

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	fmt.Println("[EMISSOR A] Aguardando conexão local na porta 9000...")
	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	fmt.Println("[EMISSOR A] Receptor conectado! Iniciando Neuro-Streaming...")

	cmd := exec.Command("ffmpeg", "-i", videoPath, "-f", "rawvideo", "-pix_fmt", "rgb24", "-vcodec", "rawvideo", "-")
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	chunkSize := 768
	buf := make([]byte, chunkSize)
	totalHashes := 0
	
	for {
		n, err := io.ReadFull(stdout, buf)
		if n == 0 || (err != nil && err != io.ErrUnexpectedEOF && err != io.EOF) {
			break
		}

		samples := make([]float64, n)
		for i := 0; i < n; i++ {
			samples[i] = float64(buf[i]) / 255.0
		}

		// Neural Match
		id := brain.MatchForced(samples)

		// Transport Payload via TCP
		idBuf := make([]byte, 8)
		binary.LittleEndian.PutUint64(idBuf, id)
		conn.Write(idBuf)
		totalHashes++
	}

	fmt.Printf("[EMISSOR A] Streaming Finalizado! %d UUIDs enviados.\n", totalHashes)
	cmd.Process.Kill()
	cmd.Wait()
}

func runClient(brainPath string, outPath string) {
	fmt.Println("[RECEPTOR B] Carregando Cérebro Híbrido...")
	var brain P2PBrain
	brain.Load(brainPath)

	fmt.Println("[RECEPTOR B] Conectando ao Emissor...")
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cmd := exec.Command("ffmpeg", "-f", "rawvideo", "-pix_fmt", "rgb24", "-s", "1920x1080", "-r", "24", "-i", "-", "-c:v", "libx264", "-pix_fmt", "yuv420p", "-y", outPath)
	stdin, _ := cmd.StdinPipe()
	cmd.Start()

	idBuf := make([]byte, 8)
	totalRecebidos := 0
	
	for {
		_, err := io.ReadFull(conn, idBuf)
		if err != nil {
			break // EOF ou Fim da stream
		}

		id := binary.LittleEndian.Uint64(idBuf)
		chunk, exists := brain.Memory[id]
		
		if !exists {
			continue
		}

		byteChunk := make([]byte, len(chunk))
		for i, f := range chunk {
			val := int(f * 255.0)
			if val < 0 {
				val = 0
			}
			if val > 255 {
				val = 255
			}
			byteChunk[i] = byte(val)
		}
		
		// Escreve array reconstruído no tubo
		stdin.Write(byteChunk)
		totalRecebidos++
	}
	
	fmt.Printf("[RECEPTOR B] %d UUIDs recebidos. Reconstruindo arquivo de vídeo...\n", totalRecebidos)
	stdin.Close()
	cmd.Wait()
	fmt.Println("[RECEPTOR B] Arquivo CROM H264 MP4/AVI reconstruído e finalizado.")
}

func main() {
	mode := flag.String("mode", "client", "Modo do Simulador: server ou client")
	flag.Parse()

	if *mode == "server" {
		runServer("../dados/cerebro_HIBRIDO.gob", "../dados/ouro_video_unseen.avi")
	} else {
		runClient("../dados/cerebro_HIBRIDO.gob", "../dados/saida_reconstruida.avi") // Mudamos se quiser .mp4
	}
}
