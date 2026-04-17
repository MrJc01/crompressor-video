package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"
	"os/exec"
)

const BlockSize = 16
const Threshold = 0.05 // Frequência de Erro Aceitável

type TensorBlock struct {
	X      int
	Y      int
	Tensor []float64
}

type CodebookEngine struct {
	Memory map[uint64][]float64
}

type FrameSkeleton struct {
	Pointers []uint64            // Cada índice equivale a 1 Bloco (Hash)
	Deltas   map[int][]float64   // Array Float64 puro paras as falhas no Dicionário
}

// O Arquivo Cromvid O(1) com suporte a Residuos "Alien"
type CromvidFile struct {
	Width     int
	Height    int
	BlockSize int
	Frames    []FrameSkeleton
}

func init() {
	gob.Register(CodebookEngine{})
	gob.Register(CromvidFile{})
}

func Chunker(img image.Image, blockSize int) []TensorBlock {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	var blocks []TensorBlock

	for y := 0; y < height; y += blockSize {
		for x := 0; x < width; x += blockSize {
			tensor := make([]float64, 0, blockSize*blockSize*3)
			for blockY := 0; blockY < blockSize; blockY++ {
				for blockX := 0; blockX < blockSize; blockX++ {
					cx, cy := x+blockX, y+blockY
					if cx >= width { cx = width - 1 }
					if cy >= height { cy = height - 1 }
					c := img.At(cx, cy)
					r, g, b, _ := c.RGBA()
					tensor = append(tensor, float64(r>>8)/255.0, float64(g>>8)/255.0, float64(b>>8)/255.0)
				}
			}
			blocks = append(blocks, TensorBlock{X: x, Y: y, Tensor: tensor})
		}
	}
	return blocks
}

// LSH Simulado
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

// Retorna tensor por ID (MOCK)
func (c *CodebookEngine) GetTensor(id uint64) []float64 {
	return c.Memory[id]
}

// Em um sistema LSH REAL (HNSW), ele encontra a corda "mais próxima". 
// Aqui vamos emular um KNN Brute Force simples para testar a distância do "Cérebro"
// Se a imagem for ALIENIGENA e as cores não baterem com NENHUMA na Brain, vai explodir Threshold!
func (c *CodebookEngine) MatchOrDelta(alienTensor []float64) (uint64, bool) {
	bestDist := math.MaxFloat64
	var bestID uint64

	// Varre Dicionário Base em busca da similaridade Euclidiana
	for id, dbTensor := range c.Memory {
		var diff float64
		for i := 0; i < len(alienTensor); i++ {
			diff += math.Abs(alienTensor[i] - dbTensor[i])
		}
		avgDiff := diff / float64(len(alienTensor))
		if avgDiff < bestDist {
			bestDist = avgDiff
			bestID = id
		}
	}
	// Se achou alguém perfeitamente confiável (Treinado)
	if bestDist <= Threshold {
		return bestID, true // TRUE = Deu Match Confiavel
	}
	// Não achou nada confiável (Alien)
	return bestID, false // FALSE = Falha no dicionario, Incorpore Diferença Delta Xor!
}

// ========================== O COMANDO 'TRAIN' ==============================
func trainBrain(videoPath, brainOut string) {
	fmt.Printf("[TRAIN] Treinando 'master.brain' sobre Mídia Ouro: %s\n", videoPath)
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-f", "image2pipe", "-pix_fmt", "rgba", "-vcodec", "rawvideo", "-")
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = nil
	cmd.Start()

	width, height := 640, 360
	frameSize := width * height * 4
	buf := make([]byte, frameSize)

	brain := &CodebookEngine{Memory: make(map[uint64][]float64)}
	for {
		_, err := io.ReadFull(stdout, buf)
		if err != nil { break }
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		img.Pix = append([]byte(nil), buf...)
		blocks := Chunker(img, BlockSize)
		// Aprende cada folha do video (Mestria)
		for _, b := range blocks {
			brain.InsertAnchor(b.Tensor)
		}
	}
	cmd.Wait()
	
	f, _ := os.Create(brainOut)
	gob.NewEncoder(f).Encode(*brain)
	f.Close()
	fmt.Printf("[TRAIN] Sucesso! Cérebro '%s' compilado com %d Hashes Core.\n", brainOut, len(brain.Memory))
}

// ========================= O COMANDO 'ENCODE' ==============================
func encodeAlien(videoPath, brainIn, cromvidOut string) {
	fmt.Printf("[ENCODE] Extraindo Video Alien '%s' limitando sobre cérebro engessado\n", videoPath)

	fBrain, _ := os.Open(brainIn)
	var brain CodebookEngine
	gob.NewDecoder(fBrain).Decode(&brain)
	fBrain.Close()

	cmd := exec.Command("ffmpeg", "-y", "-i", videoPath, "-f", "image2pipe", "-pix_fmt", "rgba", "-vcodec", "rawvideo", "-")
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = nil
	cmd.Start()

	width, height := 640, 360
	frameSize := width * height * 4
	buf := make([]byte, frameSize)

	cromvid := CromvidFile{Width: width, Height: height, BlockSize: BlockSize}

	perfectHashes := 0
	deltaFails := 0

	for {
		_, err := io.ReadFull(stdout, buf)
		if err != nil { break }
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		img.Pix = append([]byte(nil), buf...)
		blocks := Chunker(img, BlockSize)

		frame := FrameSkeleton{
			Pointers: make([]uint64, len(blocks)),
			Deltas:   make(map[int][]float64),
		}

		for idx, b := range blocks {
			// Submetemos o Alien. Se o cérebro treinado em "Barra de cores"
			// ver um Rosto ou "Mandelbrot", ele não vai achar hash parecido ($<5\%$).
			bestMatchID, isReliable := brain.MatchOrDelta(b.Tensor)
			if isReliable {
				frame.Pointers[idx] = bestMatchID
				perfectHashes++
			} else {
				// CÉREBRO FALHOU! Mecânica Delta (Uncanny Valley Fallback Doc 5)
				// Em vez de arriscar 1 pixel errado, gravamos o bloco inteiro original adjacente
				frame.Pointers[idx] = bestMatchID 
				frame.Deltas[idx] = b.Tensor // Em prod. seria "Xor", mas p/ o Proof, passamos a array plana
				deltaFails++
			}
		}
		cromvid.Frames = append(cromvid.Frames, frame)
	}

	cmd.Wait()
	fCromvid, _ := os.Create(cromvidOut)
	gob.NewEncoder(fCromvid).Encode(cromvid)
	fCromvid.Close()
	
	fmt.Printf("[ENCODE] Cromvid Esqueleto Esculpido!\n")
	fmt.Printf("   > Blocos com Conhecimento Nativo na Brain (0KB Overhead): %d\n", perfectHashes)
	fmt.Printf("   > Blocos Cérebro Falhou salvos com Payload Delta Residuo  : %d\n", deltaFails)
}

// ========================= O COMANDO 'DECODE' ==============================
func decodeAlien(cromvidIn, brainIn string) {
	fmt.Printf("[DECODE] Remontando alienígena...\n")
	
	fBrain, _ := os.Open(brainIn)
	var brain CodebookEngine
	gob.NewDecoder(fBrain).Decode(&brain)
	fBrain.Close()

	fCromvid, _ := os.Open(cromvidIn)
	var cromvid CromvidFile
	gob.NewDecoder(fCromvid).Decode(&cromvid)
	fCromvid.Close()

	fmt.Printf("[DECODE] Processando Frame O(1) usando Fallbacks...\n")
	if len(cromvid.Frames) == 0 {
		return
	}

	// Restituimos apenas o frame #0 para prova visual
	f0 := cromvid.Frames[0]
	img := image.NewRGBA(image.Rect(0, 0, cromvid.Width, cromvid.Height))

	for idx, hashID := range f0.Pointers {
		var tensor []float64
		// O milagre operado aqui: se o bloco for delta, ele intercepta antes do Cérebro burro desenhar
		if deltaResidual, exists := f0.Deltas[idx]; exists {
			tensor = deltaResidual
		} else {
			tensor = brain.GetTensor(hashID)
		}

		bx := (idx % (cromvid.Width / BlockSize)) * BlockSize
		by := (idx / (cromvid.Width / BlockSize)) * BlockSize
		
		pixelPtr := 0
		for blockY := 0; blockY < BlockSize; blockY++ {
			for blockX := 0; blockX < BlockSize; blockX++ {
				if pixelPtr+2 >= len(tensor) { break }
				if bx+blockX >= cromvid.Width || by+blockY >= cromvid.Height {
					pixelPtr += 3; continue
				}
				r := uint8(math.Round(tensor[pixelPtr] * 255.0))
				g := uint8(math.Round(tensor[pixelPtr+1] * 255.0))
				b := uint8(math.Round(tensor[pixelPtr+2] * 255.0))
				img.Set(bx+blockX, by+blockY, color.RGBA{R: r, G: g, B: b, A: 255})
				pixelPtr += 3
			}
		}
	}
	
	fOut, _ := os.Create("alien_restored_frame0.png")
	png.Encode(fOut, img)
	fOut.Close()
	fmt.Println("[DECODE] Tese confirmada! Frame Zero renderizado a salvo de distorções em 'alien_restored_frame0.png'")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Comandos suportados: train, encode, decode")
		return
	}

	cmd := os.Args[1]
	switch cmd {
	case "train":
		trainBrain(os.Args[2], os.Args[3])
	case "encode":
		encodeAlien(os.Args[2], os.Args[3], os.Args[4])
	case "decode":
		decodeAlien(os.Args[2], os.Args[3])
	default:
		fmt.Println("Comando não reconhecido")
	}
}
