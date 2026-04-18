package engine

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// CromGame implementa a interface O(1) renderizando Pixels direto na placa de vídeo
type CromGame struct {
	width, height int
	rgbaBuffer    []byte
	brainPath     string
	cromPath      string
	
	brain       *AgnosticBrain
	videoHashes []uint64
	frameIndex  int // Agora frameIndex mapeia o Hash Inteiro
}

func (g *CromGame) Update() error {
	if len(g.videoHashes) == 0 {
		return nil
	}

	chunkSize := 768
	pixelsPerFrame := g.width * g.height

	hashesPerFrame := pixelsPerFrame / chunkSize
	if pixelsPerFrame%chunkSize != 0 {
		hashesPerFrame++
	}

	if g.frameIndex >= len(g.videoHashes) {
		g.frameIndex = 0 // loop video
	}

	// Costura de blocos paralela O(1) -> Extrai Nx Hashes e carimba no Buffer da VGPU
	pixelOffset := 0
	for h := 0; h < hashesPerFrame; h++ {
		cursorHash := g.frameIndex + h
		if cursorHash >= len(g.videoHashes) {
			break
		}

		hashO1 := g.videoHashes[cursorHash]
		tensor, exists := g.brain.Memory[hashO1]
		
		if exists {
			for t := 0; t < len(tensor); t++ {
				// Mapeando Array 1D Linear pro RGB do Ebiten
				bufferCursor := (pixelOffset + t) * 4
				if bufferCursor+3 >= len(g.rgbaBuffer) {
					break
				}
				
				val := tensor[t]
				g.rgbaBuffer[bufferCursor]   = val 
				g.rgbaBuffer[bufferCursor+1] = val 
				g.rgbaBuffer[bufferCursor+2] = val 
				g.rgbaBuffer[bufferCursor+3] = 255 
			}
		}
		pixelOffset += chunkSize
	}
	
	g.frameIndex += hashesPerFrame
	return nil
}

func (g *CromGame) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.rgbaBuffer)
	
	// Math p/ Telemetria Correta
	chunkSize := 768
	hashesPerFrame := (g.width * g.height) / chunkSize
	if hashesPerFrame == 0 { hashesPerFrame = 1 }
	quadroAtual := g.frameIndex / hashesPerFrame
	totalQuadros := len(g.videoHashes) / hashesPerFrame

	msg := fmt.Sprintf("CROM O(1) Native GUI\nMemória Carregada: %s\nHash Stream: %s\nFPS Absoluto: %0.2f\nQuadro Real: %d/%d", g.brainPath, g.cromPath, ebiten.ActualFPS(), quadroAtual, totalQuadros)
	ebitenutil.DebugPrint(screen, msg)
}

func (g *CromGame) Layout(w, h int) (int, int) {
	return g.width, g.height
}

// RunPlayer levanta a janela Nativa isolada de browsers
func RunPlayer(cromFile, brainFile string) {
	fmt.Printf("[ENGINE] Carregando Código Mestre O(1) para RAM: %s...\n", brainFile)
	b := &AgnosticBrain{}
	if err := b.Load(brainFile); err != nil {
		fmt.Printf("[ERRO] Falha ao carregar cérebro: %s\n", err)
		return
	}

	fmt.Printf("[ENGINE] Carregando Lista UUID (Fita CROM) para VRAM...\n")
	f, err := os.Open(cromFile)
	if err != nil {
		fmt.Printf("[ERRO] Falha ao carregar mídia CROM: %s\n", err)
		return
	}
	defer f.Close()

	var hashes []uint64
	var buf [8]byte
	for {
		_, err := io.ReadFull(f, buf[:])
		if err != nil {
			break
		}
		hashes = append(hashes, binary.LittleEndian.Uint64(buf[:]))
	}

	w, h := 640, 360 // Aspect limpo pra GrayScale O(1) Test
	game := &CromGame{
		width:       w,
		height:      h,
		rgbaBuffer:  make([]byte, w*h*4), // 4 bytes per pixel
		brainPath:   brainFile,
		cromPath:    cromFile,
		brain:       b,
		videoHashes: hashes,
	}

	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle(fmt.Sprintf("CROM Vision SRE: %s | Zero Codecs", cromFile))
	
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
