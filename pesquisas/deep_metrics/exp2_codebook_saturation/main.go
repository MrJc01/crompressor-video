package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

type CodebookEngine struct {
	Memory map[uint64][]float64
}

func (c *CodebookEngine) InsertAnchor(tensor []float64) uint64 {
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

func generateRandomBlock(size int) []float64 {
	b := make([]float64, size)
	for i := range b {
		b[i] = rand.Float64()
	}
	return b
}

func main() {
	fmt.Println("==================================================")
	fmt.Println("[DEEP PROFILING] CROM_EXP_6.2: Codebook Saturation OOM & PProf")
	fmt.Println("==================================================")

	fCpu, _ := os.Create("cpu.prof")
	pprof.StartCPUProfile(fCpu)
	defer pprof.StopCPUProfile()

	brain := &CodebookEngine{Memory: make(map[uint64][]float64)}
	
	iterations := 250000 // Injeta 250 mil blocos independentes massivos (Equivale a decodificar filme 4K)
	blockSize := 16 * 16 * 3

	start := time.Now()
	fmt.Printf("[*] Iniciando injeção bruta LSH de %d blocos randômicos...\n", iterations)

	for i := 0; i < iterations; i++ {
		b := generateRandomBlock(blockSize)
		brain.InsertAnchor(b)
		if i%50000 == 0 && i != 0 {
			fmt.Printf("   > Pressionando RAM: %d Âncoras Registradas...\n", i)
		}
	}

	elapsed := time.Since(start)

	fMem, _ := os.Create("mem.prof")
	pprof.WriteHeapProfile(fMem)
	fMem.Close()

	fmt.Printf("[+] Bombardeio Encerrado. Profiling Escrito para 'cpu.prof' e 'mem.prof'\n")
	fmt.Printf("[+] Densidade Final do Cérebro: %d Índices Únicos (Colisões Sha-256 mitigadas implicitamente).\n", len(brain.Memory))
	fmt.Printf("⏱ Tempo Bruteforce: %v\n", elapsed)
	fmt.Println("==================================================")
}
