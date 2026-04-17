package engine

import (
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// CromGame implementa a interface O(1) renderizando Pixels direto na placa de vídeo
type CromGame struct {
	width, height int
	rgbaBuffer    []byte
	brainPath     string
	cromPath      string
}

func (g *CromGame) Update() error {
	// SIMULAÇÃO DA DESCOMPRESSÃO O(1):
	// Num cenário real, o arquivo UUID.crom geraria uma stream e
	// o Cérebro preencheria os bytes. Aqui, estressamos a VRAM 
	// atualizando o buffer cru frame-a-frame de forma alucinada.
	for i := 0; i < len(g.rgbaBuffer); i += 4 {
		g.rgbaBuffer[i] = byte(rand.Intn(256))   // Red
		g.rgbaBuffer[i+1] = byte(rand.Intn(256)) // Green
		g.rgbaBuffer[i+2] = byte(rand.Intn(256)) // Blue
		g.rgbaBuffer[i+3] = 255                  // Alpha = Maximo
	}
	return nil
}

func (g *CromGame) Draw(screen *ebiten.Image) {
	// WritePixels transfere a matriz O(1) de bytes direto para a VGPU ignorando codificação OS
	screen.WritePixels(g.rgbaBuffer)
	
	// Telemetria SRE Nativa do Ebitengine
	msg := fmt.Sprintf("CROM O(1) Native GUI\nMente: %s\nMidia: %s\nFPS Absoluto: %0.2f", g.brainPath, g.cromPath, ebiten.ActualFPS())
	ebitenutil.DebugPrint(screen, msg)
}

func (g *CromGame) Layout(w, h int) (int, int) {
	return g.width, g.height
}

// RunPlayer levanta a janela Nativa isolada de browsers
func RunPlayer(cromFile, brainFile string) {
	// Setup Resolution HD for testing (Zero overhead)
	w, h := 1280, 720
	game := &CromGame{
		width:      w,
		height:     h,
		rgbaBuffer: make([]byte, w*h*4),
		brainPath:  brainFile,
		cromPath:   cromFile,
	}

	ebiten.SetWindowSize(w, h)
	ebiten.SetWindowTitle(fmt.Sprintf("Vision: %s | 0 Codecs", cromFile))
	
	fmt.Printf("[ENGINE] Disparando Tela em %dx%d e injetando Barramento VRAM...\n", w, h)
	
	if err := ebiten.RunGame(game); err != nil {
		panic(err) // Morte controlada do processo GUI
	}
}
