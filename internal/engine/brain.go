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
	Memory map[uint64][]uint8
}

// Learn absorve um tensor e gera um UUID baseado no hash estatístico exato
func (c *AgnosticBrain) Learn(tensor []uint8) uint64 {
	bytePayload := make([]byte, len(tensor))
	for i, v := range tensor {
		bytePayload[i] = v
	}
	hash32 := sha256.Sum256(bytePayload)
	id := binary.LittleEndian.Uint64(hash32[:8]) // Usa os primeiros bytes do Sha256

	if c.Memory == nil {
		c.Memory = make(map[uint64][]uint8)
	}

	if _, exists := c.Memory[id]; !exists {
		// Congela a Memória Real
		clone := make([]uint8, len(tensor))
		copy(clone, tensor)
		c.Memory[id] = clone
	}
	return id
}

// MatchForced encontra o Hash UUID para dados "Unseen" (Nunca Vistos) minimizando o MSE Injetado
func (c *AgnosticBrain) MatchForced(target []uint8) uint64 {
	bestDist := math.MaxFloat64
	var bestID uint64
	for id, dbTensor := range c.Memory {
		limit := len(target)
		if len(dbTensor) < limit {
			limit = len(dbTensor)
		}
		var diff uint64
		for i := 0; i < limit; i++ {
			d := int(target[i]) - int(dbTensor[i])
			if d < 0 {
				d = -d
			}
			diff += uint64(d)
		}
		avg := float64(diff) / float64(limit)
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

func ProcessFlatMedia(cmd *exec.Cmd, chunkSize int, byteMultiplier int, maxElements int, process func([]uint8)) {
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
		
		var samples []uint8
		if byteMultiplier == 1 {
			samples = make([]uint8, n)
			for i := 0; i < n; i++ {
				samples[i] = buf[i]
			}
		} else if byteMultiplier == 4 { // Raw Float32
			samples = make([]uint8, n/4)
			for i := 0; i < len(samples); i++ {
				bits := binary.LittleEndian.Uint32(buf[i*4 : i*4+4])
				samples[i] = uint8(math.Float32frombits(bits) * 255.0)
			}
		}
		
		process(samples)
		total += len(samples)
	}
	cmd.Process.Kill()
	cmd.Wait()
}

func ProcessFlatVideo(path string, chunkSize int, maxElements int, process func([]uint8)) {
    // Para simplificar a POC como era antes para ler um video cru (usando FFMPEG inves de CAT)
	// FORÇANDO SCALE=640:360 pra bater matematicamente com os Canvas/VRAM pre-fixados!
	cmd := exec.Command("ffmpeg", "-i", path, "-vf", "scale=640:360", "-f", "rawvideo", "-pix_fmt", "gray", "-")
	ProcessFlatMedia(cmd, chunkSize, 1, maxElements, process)
}
