package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
)

type AgnosticBrain struct {
	Memory map[uint64][]float64
}

func (c *AgnosticBrain) Learn(tensor []float64) uint64 {
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

func (c *AgnosticBrain) MatchForced(target []float64) uint64 {
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

// Qualimetria
func CalcMSE(original, decoded []float64) float64 {
	limit := len(original)
	if len(decoded) < limit {
		limit = len(decoded)
	}
	if limit == 0 {
		return math.MaxFloat64
	}
	var sum float64
	for i := 0; i < limit; i++ {
		diff := original[i] - decoded[i]
		sum += diff * diff
	}
	return sum / float64(limit)
}

func CalcPSNR(mse float64) float64 {
	if mse == 0 {
		return 100.0
	}
	return 10 * math.Log10(1.0/mse)
}

// Process pipes
func processFlatMedia(cmd *exec.Cmd, chunkSize int, byteMultiplier int, maxElements int, process func([]float64)) {
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	defer cmd.Process.Kill()
	buf := make([]byte, chunkSize*byteMultiplier)
	total := 0
	for {
		if maxElements > 0 && total >= maxElements {
			break
		}
		n, err := io.ReadFull(stdout, buf)
		if n == 0 || err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			break
		}
		var samples []float64
		if byteMultiplier == 1 {
			samples = make([]float64, n)
			for i := 0; i < n; i++ {
				samples[i] = float64(buf[i]) / 255.0
			}
		} else if byteMultiplier == 4 {
			samples = make([]float64, n/4)
			for i := 0; i < len(samples); i++ {
				bits := binary.LittleEndian.Uint32(buf[i*4 : i*4+4])
				samples[i] = float64(math.Float32frombits(bits))
			}
		}
		process(samples)
		total += len(samples)
	}
	cmd.Process.Kill()
	cmd.Wait()
}

func processFlatVideo(path string, chunkSize int, maxElements int, process func([]float64)) {
	processFlatMedia(exec.Command("ffmpeg", "-i", path, "-f", "image2pipe", "-pix_fmt", "rgb24", "-vcodec", "rawvideo", "-"), chunkSize, 1, maxElements, process)
}

func processFlatAudio(path string, chunkSize int, maxElements int, process func([]float64)) {
	processFlatMedia(exec.Command("ffmpeg", "-i", path, "-f", "f32le", "-acodec", "pcm_f32le", "-ac", "1", "-"), chunkSize, 4, maxElements, process)
}

func processFlatRaw(path string, chunkSize int, maxElements int, process func([]float64)) {
	cmd := exec.Command("cat", path)
	processFlatMedia(cmd, chunkSize, 1, maxElements, process)
}

// Salva e Carrega
func (c *AgnosticBrain) Save(path string) error {
	f, err := os.Create(path)
	if err != nil { return err }
	defer f.Close()
	return gob.NewEncoder(f).Encode(c.Memory)
}

func (c *AgnosticBrain) Load(path string) error {
	f, err := os.Open(path)
	if err != nil { return err }
	defer f.Close()
	return gob.NewDecoder(f).Decode(&c.Memory)
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("[CROSS-DOMAIN MATRIZ 2.0] Qualimetria SRE e Mentes Isoladas.")
	fmt.Println("==================================================")
	
	chunkSize := 768
	brains := map[string]*AgnosticBrain{
		"VIDEO":   {Memory: make(map[uint64][]float64)},
		"AUDIO":   {Memory: make(map[uint64][]float64)},
		"IMAGE":   {Memory: make(map[uint64][]float64)},
		"ALIEN":   {Memory: make(map[uint64][]float64)},
		"HIBRIDO": {Memory: make(map[uint64][]float64)},
	}

	// 1. FASE DE TREINO (Aprender todos os Datasets de 500MB)
	// Limite: 50.000.000 samples por brain (~400MB de floats extraidos)
	trainLimit := 50000000 
	fmt.Println("\n[...] FASE 1: Treinamento Neural e Persistência")

	tasks := []struct{ Name, Path, Func string }{
		{"VIDEO", "../dados/ouro_video_treino.avi", "video"},
		{"AUDIO", "../dados/ouro_audio_treino.wav", "audio"},
		{"IMAGE", "../dados/ouro_imagem_treino.bmp", "video"}, // bmp processado via ffmpeg
		{"ALIEN", "../dados/alien_ruido_treino.raw", "raw"},
	}

	for _, t := range tasks {
		fmt.Printf("      -> Lendo %s e preenchendo Cérebro Isolado e Híbrido...\n", t.Name)
		processor := processFlatVideo
		if t.Func == "audio" { processor = processFlatAudio }
		if t.Func == "raw" { processor = processFlatRaw }

		if _, err := os.Stat(t.Path); err != nil {
			fmt.Printf("         [ERRO] Base %s ausente, pulando...\n", t.Path)
			continue
		}

		processor(t.Path, chunkSize, trainLimit, func(chunk []float64) {
			brains[t.Name].Learn(chunk)
			brains["HIBRIDO"].Learn(chunk)
		})
		brains[t.Name].Save(fmt.Sprintf("../dados/cerebro_%s.gob", t.Name))
	}
	brains["HIBRIDO"].Save("../dados/cerebro_HIBRIDO.gob")
	fmt.Println("[+] Todos os Cérebros (5) foram persistidos em DISCO!")

	// 2. FASE INTERLAÇADA N x N
	fmt.Println("\n[...] FASE 2: Intersecção Cross-Domain (100MB Unseen Data)")
	unseenLimit := 100000 // Teste rápido em 100k samples (~800KB decodificados) pra não demorar
	
	csvF, _ := os.Create("../dados/resultados_crossover.csv")
	defer csvF.Close()
	csvF.WriteString("Teste,Cerebro_Usado,Dominio_Entrada,MSE,PSNR\n")

	unseenTasks := []struct{ Name, Path, Func string }{
		{"VIDEO_TEST", "../dados/ouro_video_unseen.avi", "video"},
		{"AUDIO_TEST", "../dados/ouro_audio_unseen.wav", "audio"},
		{"IMAGE_TEST", "../dados/ouro_imagem_unseen.bmp", "video"},
		{"ALIEN_TEST", "../dados/alien_ruido_unseen.raw", "raw"},
	}

	for BName, brain := range brains {
		if len(brain.Memory) == 0 { continue }
		for _, UTask := range unseenTasks {
			var orig, rec []float64

			processor := processFlatVideo
			if UTask.Func == "audio" { processor = processFlatAudio }
			if UTask.Func == "raw" { processor = processFlatRaw }
			
			if _, err := os.Stat(UTask.Path); err != nil { continue }

			processor(UTask.Path, chunkSize, unseenLimit, func(chunk []float64) {
				orig = append(orig, chunk...)
				uuid := brain.MatchForced(chunk)
				rec = append(rec, brain.Memory[uuid]...)
			})

			mse := CalcMSE(orig, rec)
			psnr := CalcPSNR(mse)
			
			name := fmt.Sprintf("%s_VS_%s", BName, UTask.Name)
			fmt.Printf("   [CROSS] %s -> MSE: %.4f | PSNR: %.2fdB\n", name, mse, psnr)
			csvF.WriteString(fmt.Sprintf("%s,%s,%s,%.4f,%.2f\n", name, BName, UTask.Name, mse, psnr))
		}
	}
	fmt.Println("\n✅ AVALIAÇÃO PSNR CROSS-DOMAIN (N x N) FINALIZADA. csv gerado.")
}
