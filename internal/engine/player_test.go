package engine

import (
	"os"
	"testing"
)

func TestPlayer_UpdateDraw(t *testing.T) {
	// Cria uma mook de memória e simulamos a execução Update do CromGame
	b := &AgnosticBrain{Memory: make(map[uint64][]uint8)}
	b.Memory[999] = []uint8{255, 128, 50, 20} // Chunk 4 bytes

	game := &CromGame{
		width: 2, height: 2, // 4 pixels = 16 bytes rgba
		rgbaBuffer:  make([]byte, 16),
		brainPath:   "falso.gob",
		cromPath:    "falso.crom",
		brain:       b,
		videoHashes: []uint64{999}, // Simular que a Fita possui 1 bloco de 4
		frameIndex:  0,
	}

	err := game.Update() // Dispara a construção do rgbaBuffer baseado nos hashes
	if err != nil {
		t.Fatalf("Erro no Game Update: %v", err)
	}

	if game.rgbaBuffer[0] != 255 {
		t.Errorf("Decodificador de VRAM falhou em mesclar uint8. Obtido: %d", game.rgbaBuffer[0])
	}
}

func TestRunPlayer_NoWindow(t *testing.T) {
	// Apenas para engatilhar os retornos da função RunPlayer quando o arquivo não existe (Cobertura)
	RunPlayer("fantasma.crom", "fantasma.gob")
	
	tmpGob := "fantastma2.gob"
	b := &AgnosticBrain{Memory: make(map[uint64][]uint8)}
	b.Save(tmpGob)
	defer os.Remove(tmpGob)
	RunPlayer("fantasma.crom", tmpGob)
}
