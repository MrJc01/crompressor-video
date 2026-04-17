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
	frameIndex  int
}

func (g *CromGame) Update() error {
	if g.frameIndex >= len(g.videoHashes) {
		g.frameIndex = 0 // loop video
	}
	if len(g.videoHashes) == 0 {
		return nil
	}

	hashO1 := g.videoHashes[g.frameIndex]
	tensor, exists := g.brain.Memory[hashO1]

	if exists {
		// Decodificação Extrema: Tensor -> Tela
		// Tensor está em 0.0 - 1.0 ou cru (0-255). No nosso Extrator (gray) é de 0 a 1 em base 1.
		// Vamos injetar O(1) de forma agressiva
		for i := 0; i < len(g.rgbaBuffer) && i/4 < len(tensor); i += 4 {
			val := byte(tensor[i/4] * 255.0)
			g.rgbaBuffer[i]   = val // R
			g.rgbaBuffer[i+1] = val // G
			g.rgbaBuffer[i+2] = val // B
			g.rgbaBuffer[i+3] = 255 // A
		}
	}
	
	g.frameIndex++
	return nil
}

func (g *CromGame) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.rgbaBuffer)
	msg := fmt.Sprintf("CROM O(1) Native GUI\nMemória Carregada: %s\nHash Stream: %s\nFPS Absoluto: %0.2f\nQuadro: %d/%d", g.brainPath, g.cromPath, ebiten.ActualFPS(), g.frameIndex, len(g.videoHashes))
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
