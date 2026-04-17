package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"
)

// =========================================================
// 1. O CÉREBRO AGNÓSTICO (ACEITA QUALQUER DIMENSÃO 1D/2D/3D)
// =========================================================
type AgnosticCodebook struct {
	Memory map[uint64][]float64
}

func NewAgnosticCodebook() *AgnosticCodebook {
	return &AgnosticCodebook{Memory: make(map[uint64][]float64)}
}

func (c *AgnosticCodebook) Insert(tensor []float64) uint64 {
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

// Retorna se o bloco deu Match Perfeito (Trained Data) ou Erro (Unseen Data)
func (c *AgnosticCodebook) TestMatch(target []float64, threshold float64) bool {
	bestDist := math.MaxFloat64
	for _, dbTensor := range c.Memory {
		if len(dbTensor) != len(target) {
			continue // Ignora arrays de tamanho diferente
		}
		var diff float64
		for i := 0; i < len(target); i++ {
			diff += math.Abs(target[i] - dbTensor[i])
		}
		avgDiff := diff / float64(len(target))
		if avgDiff < bestDist {
			bestDist = avgDiff
		}
	}
	return bestDist <= threshold
}

// =========================================================
// 2. PIPE EXTRATOR (ROTEADOR DE VÍDEO vs ÁUDIO)
// =========================================================

// Extrai amplitudes sonoras (Áudio/1D) num array Flat global
func ExtractAudioPipe(path string) []float64 {
	// Puxa PCM S16LE
	cmd := exec.Command("ffmpeg", "-i", path, "-f", "s16le", "-acodec", "pcm_s16le", "-ac", "1", "-")
	out, _ := cmd.Output()
	
	// Cada sample possui 2 bytes (16-bit)
	samples := make([]float64, len(out)/2)
	for i := 0; i < len(samples); i++ {
		// Le os bits int16 e converte p float64 normalizado
		raw := int16(binary.LittleEndian.Uint16(out[i*2 : i*2+2]))
		samples[i] = float64(raw) / 32768.0 
	}
	return samples
}

// Extrai frames sequecialmente brutas (Vídeo/2D)
func ExtractVideoPipe(path string) []float64 {
	cmd := exec.Command("ffmpeg", "-i", path, "-f", "image2pipe", "-pix_fmt", "rgba", "-vcodec", "rawvideo", "-")
	out, _ := cmd.Output()
	
	samples := make([]float64, len(out))
	for i := 0; i < len(out); i++ {
		samples[i] = float64(out[i]) / 255.0
	}
	return samples
}

// Corta Array Supremo em Blocos de Tamanho Fixo N
func ChunkFlatArray(samples []float64, chunkSize int) [][]float64 {
	blocks := make([][]float64, 0, len(samples)/chunkSize)
	for i := 0; i+chunkSize <= len(samples); i += chunkSize {
		blocks = append(blocks, samples[i:i+chunkSize])
	}
	return blocks
}

// =========================================================
// 3. O ORQUESTRADOR ANALÍTICO
// =========================================================
func executeTrial(name string, trainFile string, testFile string, isAudio bool, chunkSize int) string {
	var trainFlat, testFlat []float64
	
	if isAudio {
		trainFlat = ExtractAudioPipe(trainFile)
		testFlat = ExtractAudioPipe(testFile)
	} else {
		trainFlat = ExtractVideoPipe(trainFile)
		testFlat = ExtractVideoPipe(testFile)
	}

	trainBlocks := ChunkFlatArray(trainFlat, chunkSize)
	testBlocks := ChunkFlatArray(testFlat, chunkSize)

	// Treinando o Cérebro Exclusivo
	brain := NewAgnosticCodebook()
	for _, b := range trainBlocks {
		brain.Insert(b)
	}

	// Testando Conteúdo Desconhecido (Alien) no mesmo Cérebro
	matchCount := 0
	fallbackXorCount := 0
	
	threshold := 0.05 // Pareamento de 95% exatidão

	for _, b := range testBlocks {
		if brain.TestMatch(b, threshold) {
			matchCount++
		} else {
			fallbackXorCount++
		}
	}

	total := len(testBlocks)
	matchRate := (float64(matchCount) / float64(total)) * 100.0
	fallbackRate := (float64(fallbackXorCount) / float64(total)) * 100.0

	report := fmt.Sprintf("### Ensaio: %s\n", name)
	report += fmt.Sprintf("- **Dicionário Treinado Acumulou:** %d Hashes Base.\n", len(brain.Memory))
	report += fmt.Sprintf("- **Área de Superfície do Arquivo Testado:** %d Blocos.\n", total)
	report += fmt.Sprintf("- **Taxa de Sobrevivência do Cérebro:** %.2f%% (Blocos com Match O(1)).\n", matchRate)
	report += fmt.Sprintf("- **Taxa de Fallback Delta/Xor (Unseen):** %.2f%%\n\n", fallbackRate)
	return report
}

func main() {
	fmt.Println("[MASTER BENCH] Iniciando Funil Universal de Mídias Reais...")
	start := time.Now()

	var buf bytes.Buffer
	buf.WriteString("# Relatório Executivo CROM: Mídias Universais (Trained vs Untrained Context)\n")
	buf.WriteString("Este documento prova a capacidade vetorial da Engine sobre dimensões Acústicas (1D) e Visuais (2D) simulando Oubliette Delta em conteúdos jamais vistos no dataset primário.\n\n")

	// ENSAIO 1: ÁUDIO
	// Chunk de Áudio comum na literatura acadêmica VQ-VAE é de 256/512 frames temporais
	fmt.Println(" > Processando Laboratório de Áudio...")
	repAudio := executeTrial("1D Acoustic Wave (WAV PCM)", 
		"../dados_estaveis/audio_ouro_voz.wav", 
		"../dados_estaveis/audio_alien_ruido.wav", 
		true, 512)
	buf.WriteString(repAudio)

	// ENSAIO 2: VÍDEO
	// Chunk 16x16 RGB = 768 floats (Padrão Crompressor Anterior)
	fmt.Println(" > Processando Laboratório de Vídeo...")
	repVideo := executeTrial("2D Temporal Canvas (MP4 RGBA)", 
		"../dados_estaveis/test_video.mp4", 
		"../dados_estaveis/alien_video.mp4", 
		false, 768)
	buf.WriteString(repVideo)

	elapsed := time.Since(start)
	buf.WriteString(fmt.Sprintf("---\n*Gerado nativamente via Go SRE Engine em %v*", elapsed))

	os.MkdirAll("../relatorios", 0755)
	os.WriteFile("../relatorios/multimidias_unificadas.md", buf.Bytes(), 0644)

	fmt.Printf("[+] Relatório Final Exportado p/ researches/relatorios/multimidias_unificadas.md em %v!\n", elapsed)
}
