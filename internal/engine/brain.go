package engine

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"io"
	"math"
	"os"
	"os/exec"
)

type AgnosticBrain struct {
	Memory map[uint64][]float64
}

// Learn absorve um tensor e gera um UUID baseado no hash estatístico exato
func (c *AgnosticBrain) Learn(tensor []float64) uint64 {
	bytePayload := make([]byte, len(tensor)*8)
	for i, f := range tensor {
		binary.LittleEndian.PutUint64(bytePayload[i*8:], math.Float64bits(f))
	}
	hash32 := sha256.Sum256(bytePayload)
	id := binary.LittleEndian.Uint64(hash32[:8]) // Usa os primeiros bytes do Sha256

	if c.Memory == nil {
		c.Memory = make(map[uint64][]float64)
	}

	if _, exists := c.Memory[id]; !exists {
		// Congela a Memória Real
		clone := make([]float64, len(tensor))
		copy(clone, tensor)
		c.Memory[id] = clone
	}
	return id
}

// MatchForced encontra o Hash UUID para dados "Unseen" (Nunca Vistos) minimizando o MSE Injetado
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

func (c *AgnosticBrain) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(&c.Memory) // Must pass pointer if we load pointer? Gob is tricky with maps.
}

func (c *AgnosticBrain) Load(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewDecoder(f).Decode(&c.Memory)
}

// ======================================
// OS/Media Extractors
// ======================================

func ProcessFlatMedia(cmd *exec.Cmd, chunkSize int, byteMultiplier int, maxElements int, process func([]float64)) {
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
		if n == 0 || (err != nil && err != io.ErrUnexpectedEOF && err != io.EOF) {
			break
		}
		
		var samples []float64
		if byteMultiplier == 1 {
			samples = make([]float64, n)
			for i := 0; i < n; i++ {
				samples[i] = float64(buf[i]) / 255.0
			}
		} else if byteMultiplier == 4 { // Raw Float32
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

func ProcessFlatVideo(path string, chunkSize int, maxElements int, process func([]float64)) {
    // Para simplificar a POC como era antes para ler um video cru (usando FFMPEG inves de CAT)
	cmd := exec.Command("ffmpeg", "-i", path, "-f", "rawvideo", "-pix_fmt", "gray", "-")
	ProcessFlatMedia(cmd, chunkSize, 1, maxElements, process)
}
